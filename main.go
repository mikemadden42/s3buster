package main

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"io/ioutil"
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

func list(host string) {
	response, err := http.Get(host)
	if err != nil {
		log.Fatal(err)
	}

	responseData, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}

	var results ListBucketResult
	err = xml.Unmarshal(responseData, &results)
	if err != nil {
		log.Fatal(err)
	}

	for _, result := range results.Contents {
		fmt.Printf("key: %s, size: %s\n", result.Key, result.Size)
	}
}
