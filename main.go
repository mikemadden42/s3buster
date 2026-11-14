// Package main implements s3buster, a simple S3 bucket recon tool that
// resolves candidate bucket hostnames and attempts to list their contents.
package main

import (
	"bufio"
	"bytes"
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"sync"
	"time"
)

// httpClient enforces a timeout so a single unresponsive host cannot
// stall the whole scan.
var httpClient = &http.Client{Timeout: 10 * time.Second}

type ListBucketResult struct {
	XMLName     xml.Name `xml:"ListBucketResult"`
	Text        string   `xml:",chardata"`
	Xmlns       string   `xml:"xmlns,attr"`
	Name        string   `xml:"Name"`
	Prefix      string   `xml:"Prefix"`
	Marker      string   `xml:"Marker"`
	MaxKeys     string   `xml:"MaxKeys"`
	IsTruncated string   `xml:"IsTruncated"`
	Contents    []struct {
		Text         string `xml:",chardata"`
		Key          string `xml:"Key"`
		LastModified string `xml:"LastModified"`
		ETag         string `xml:"ETag"`
		Size         string `xml:"Size"`
		StorageClass string `xml:"StorageClass"`
	} `xml:"Contents"`
}

func main() {
	wordlist := flag.String("wordlist", "words.txt", "path to the file of candidate bucket names")
	region := flag.String("region", "us-west-1", "AWS region of the S3 endpoint")
	concurrency := flag.Int("concurrency", 10, "number of hosts to scan in parallel")
	flag.Parse()

	if *concurrency < 1 {
		log.Fatal("concurrency must be at least 1")
	}

	file, err := os.Open(*wordlist)
	if err != nil {
		log.Fatal(err)
	}
	defer func() {
		if cerr := file.Close(); cerr != nil {
			log.Println(cerr)
		}
	}()

	suffix := fmt.Sprintf(".s3-%s.amazonaws.com", *region)

	hosts := make(chan string)

	// Serialize stdout so each host's output stays grouped.
	var outMu sync.Mutex
	var wg sync.WaitGroup
	for i := 0; i < *concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for host := range hosts {
				out := scan(host)
				outMu.Lock()
				_, _ = os.Stdout.Write(out)
				outMu.Unlock()
			}
		}()
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		name := scanner.Text()
		if name == "" {
			continue
		}
		hosts <- name + suffix
	}
	close(hosts)
	wg.Wait()

	if err = scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

// scan resolves and lists a single host, returning its collected output so
// callers can emit it atomically.
func scan(host string) []byte {
	var buf bytes.Buffer
	fmt.Fprintln(&buf, host)
	lookup(&buf, host)
	list(&buf, "http://"+host)
	return buf.Bytes()
}

func lookup(w io.Writer, host string) {
	hosts, err := net.LookupHost(host)
	if err != nil {
		log.Println(err)
	}

	for _, host := range hosts {
		_, _ = fmt.Fprintln(w, host)
	}
}

func list(w io.Writer, url string) {
	response, err := httpClient.Get(url) //nolint:gosec // URL is built from the wordlist; this is the tool's intended behavior.
	if err != nil {
		log.Println(err)
		return
	}
	defer func() {
		if cerr := response.Body.Close(); cerr != nil {
			log.Println(cerr)
		}
	}()

	if response.StatusCode != http.StatusOK {
		log.Printf("%s: unexpected status %s", url, response.Status)
		return
	}

	responseData, err := io.ReadAll(response.Body)
	if err != nil {
		log.Println(err)
		return
	}

	results, err := parseListBucketResult(responseData)
	if err != nil {
		log.Println(err)
		return
	}

	for _, result := range results.Contents {
		_, _ = fmt.Fprintf(w, "key: %s, size: %s\n", result.Key, result.Size)
	}
}

// parseListBucketResult unmarshals an S3 ListBucketResult XML document.
func parseListBucketResult(data []byte) (ListBucketResult, error) {
	var results ListBucketResult
	if err := xml.Unmarshal(data, &results); err != nil {
		return ListBucketResult{}, err
	}
	return results, nil
}
