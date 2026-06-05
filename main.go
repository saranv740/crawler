package main

import (
	"fmt"
	"net/url"
	"os"
)

func crawlPage(rawBaseURL string, rawCurrentURL string, pages map[string]int) {
	baseURL, err := url.Parse(rawBaseURL)
	if err != nil {
		return
	}

	currentURL, err := url.Parse(rawCurrentURL)
	if err != nil {
		return
	}

	if currentURL.Host != baseURL.Host {
		return
	}

	normalizedCurrent, err := normalizeURL(rawCurrentURL)
	if err != nil {
		return
	}

	// If it's already exists just increment count
	if pages[normalizedCurrent] > 0 {
		pages[normalizedCurrent]++
		return
	}

	html, err := getHTML(rawCurrentURL)
	if err != nil {
		return
	}

	pageData, errs := extractPageData(html, baseURL)
	if len(errs) > 0 {
		fmt.Println("found some errors: ", errs)
	}

	// set current page as extracted
	pages[normalizedCurrent] = 1

	for _, link := range pageData.OutgoingLinks {
		crawlPage(rawBaseURL, link, pages)
	}

	for _, link := range pageData.OutgoingLinks {
		crawlPage(rawBaseURL, link, pages)
	}
}

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		fmt.Println("no url provided")
		os.Exit(1)
	}

	if len(args) > 1 {
		fmt.Println("too many arguments provided")
		os.Exit(1)
	}

	rawURL := args[0]

	pages := make(map[string]int)
	crawlPage(rawURL, rawURL, pages)
	for k, v := range pages {
		fmt.Println(k, ":", v)
	}
}
