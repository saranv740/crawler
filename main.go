package main

import (
	"flag"
	"fmt"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/saranv740/crawler/internal/crawler"
)

// start will begin to crawl base url with single goroutine, further pages will be extracted by their own goroutines
func start(rawURL string, maxConcurrency int, maxPages int) (map[string]crawler.PageData, error) {
	now := time.Now()
	pages := make(map[string]crawler.PageData)

	baseURL, err := url.Parse(strings.TrimSuffix(rawURL, "/"))
	if err != nil {
		return pages, fmt.Errorf("error in parsing URL")
	}

	crawlr := crawler.New(baseURL, maxConcurrency, maxPages)
	crawlr.Run(rawURL)

	fmt.Printf("concurrent crawler: crawled %d urls in %s\n", len(crawlr.Pages()), time.Since(now))
	return crawlr.Pages(), nil
}

func main() {
	var rawURL string
	var maxConcurrency int
	var maxPages int
	var output string

	flag.StringVar(&rawURL, "url", "", "the url to crawl")
	flag.IntVar(&maxConcurrency, "workers", 4, "max number of concurrent workers")
	flag.IntVar(&maxPages, "pages", 10, "max pages to collect")
	flag.StringVar(&output, "output", "output/report.json", "the file name to write the report")

	flag.Parse()

	if rawURL == "" {
		fmt.Println("URL is not provided")
		os.Exit(1)
	}

	data, err := start(rawURL, maxConcurrency, maxPages)
	if err != nil {
		fmt.Println("encountered error: ", err)
		os.Exit(1)
	}

	// Generate report
	err = writeJSONReport(data, output)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
