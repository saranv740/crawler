package main

import (
	"errors"
	"fmt"
	"net/url"
	"sync"
)

var (
	externalHostErr = errors.New("external host")
	urlSanityErr    = errors.New("url sanitizing failed")
	pageFetchErr    = errors.New("error in fetching page")
)

type Crawler struct {
	baseURL    *url.URL
	maxWorkers int
	maxPages   int
	pages      map[string]PageData
	result     chan crawlResult
	links      chan string
}

type crawlResult struct {
	URL         string
	Data        PageData
	NewLinks    []string
	Err         error
	ParseErrors []error
}

func New(baseURL *url.URL, maxWorkers int, maxPages int) *Crawler {
	return &Crawler{
		baseURL:    baseURL,
		maxPages:   maxPages,
		maxWorkers: maxWorkers,
		pages:      make(map[string]PageData),
		// keep result small so workers will block
		result: make(chan crawlResult, 1),
		// keep links to maximum number for pages
		links: make(chan string, maxPages+1),
	}
}

func (c *Crawler) Run(rawStartURL string) {
	var wg sync.WaitGroup

	// start workers
	for i := range c.maxWorkers {
		wg.Add(1)

		go func(id int) {
			defer wg.Done()
			c.worker(id)
		}(i)
	}

	c.links <- rawStartURL
	activeLinks := 1

	for activeLinks > 0 {
		select {
		case result := <-c.result:
			activeLinks--
			if result.Err != nil {
				delete(c.pages, result.URL)
				continue
			}

			if len(result.ParseErrors) != 0 {
				fmt.Printf("encountered some parsing errors %v\n", result.ParseErrors)
			}

			c.pages[result.URL] = result.Data

			for _, link := range result.NewLinks {
				normalized, err := normalizeURL(link)
				if err != nil {
					// Skip completely broken URLs
					continue
				}

				if len(c.pages) >= c.maxPages {
					continue
				}

				if _, exists := c.pages[normalized]; !exists {
					activeLinks++
					// mark as done in advance so the link won't be processed again
					c.pages[normalized] = PageData{}
					c.links <- link
				}
			}
		}
	}

	close(c.links)
	wg.Wait()
}

func (c *Crawler) worker(id int) {
	for link := range c.links {
		normalized, err := normalizeURL(link)
		if err != nil {
			c.result <- crawlResult{URL: link, Err: urlSanityErr}
			continue
		}

		currentURL, err := url.Parse(link)
		if err != nil {
			c.result <- crawlResult{URL: normalized, Err: urlSanityErr}
			continue
		}
		if currentURL.Host != c.baseURL.Host {
			c.result <- crawlResult{URL: normalized, Err: externalHostErr}
			continue
		}

		html, err := getHTML(currentURL.String())
		if err != nil {
			c.result <- crawlResult{URL: normalized, Err: pageFetchErr}
			continue
		}

		data, errs := extractPageData(html, normalized, c.baseURL)
		c.result <- crawlResult{
			URL:         normalized,
			ParseErrors: errs,
			Data:        data,
			NewLinks:    append(data.OutgoingLinks, data.ImageURLs...),
		}
	}
}
