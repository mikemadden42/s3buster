// Package main implements s3buster, a simple S3 bucket recon tool that
// resolves candidate bucket hostnames and attempts to list their contents.
package main

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
)

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
	file, err := os.Open("words.txt")
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		host := scanner.Text() + ".s3-us-west-1.amazonaws.com"
		fmt.Println(host)
		lookup(host)
		list("http://" + host)
	}

	if err = scanner.Err(); err != nil {
		log.Fatal(err)
	}
}

func lookup(host string) {
	hosts, err := net.LookupHost(host)
	if err != nil {
		log.Println(err)
	}

	for _, host := range hosts {
		fmt.Println(host)
	}
}

func list(url string) {
	response, err := http.Get(url) //nolint:gosec // URL is built from the wordlist; this is the tool's intended behavior.
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

	var results ListBucketResult
	if err = xml.Unmarshal(responseData, &results); err != nil {
		log.Println(err)
		return
	}

	for _, result := range results.Contents {
		fmt.Printf("key: %s, size: %s\n", result.Key, result.Size)
	}
}
