package main

import (
	"fmt"
	"net/url"
	"sync"
)

type config struct {
	pages    map[string]PageData
	baseURL  *url.URL
	mu       *sync.Mutex
	sema     chan struct{}
	wg       *sync.WaitGroup
	maxPages int
}

func (c *config) crawlPage(rawCurrentURL string) {
	// Aquire semaphore lock
	c.sema <- struct{}{}
	defer func() {
		// Release the lock
		<-c.sema
		c.wg.Done()
	}()

	currentURL, err := url.Parse(rawCurrentURL)
	if err != nil || currentURL.Host != c.baseURL.Host {
		return
	}

	normalizedCurrent, err := normalizeURL(rawCurrentURL)
	if err != nil {
		return
	}

	c.mu.Lock()
	if len(c.pages) >= c.maxPages {
		c.mu.Unlock()
		return
	}
	// If it already exists return
	if _, exists := c.pages[normalizedCurrent]; exists {
		c.mu.Unlock()
		return
	}
	// set current page as extracted
	c.pages[normalizedCurrent] = PageData{}
	c.mu.Unlock()

	html, err := getHTML(rawCurrentURL)
	if err != nil {
		return
	}

	pageData, errs := extractPageData(html, currentURL.String(), c.baseURL)
	if len(errs) > 0 {
		fmt.Println("found some errors: ", errs)
	}

	c.mu.Lock()
	c.pages[normalizedCurrent] = pageData
	c.mu.Unlock()

	for _, link := range pageData.OutgoingLinks {
		c.wg.Add(1)
		go func(l string) {
			c.crawlPage(l)
		}(link)
	}

	for _, link := range pageData.ImageURLs {
		c.wg.Add(1)
		go func(l string) {
			c.crawlPage(l)
		}(link)
	}
}
