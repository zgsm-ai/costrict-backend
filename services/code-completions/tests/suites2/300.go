package main

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"
	"sync"
	"time"
)

type WebCrawler struct {
	visited map[string]bool
	queue   []string
	mu      sync.Mutex
	client  *http.Client
	results []string
}

func NewWebCrawler() *WebCrawler {
	return &WebCrawler{
		visited: make(map[string]bool),
		queue:   make([]string, 0),
		client:  &http.Client{Timeout: 10 * time.Second},
		results: make([]string, 0),
	}
}

func (wc *WebCrawler) Start(seedURL string) {
	wc.mu.Lock()
	wc.queue = append(wc.queue, seedURL)
	wc.mu.Unlock()
	
	for len(wc.queue) > 0 && len(wc.results) < 10 {
		wc.mu.Lock()
		url := wc.queue[0]
		wc.queue = wc.queue[1:]
		wc.visited[url] = true
		wc.mu.Unlock()
		
		title, err := wc.crawlPage(url)
		if err != nil {
			fmt.Printf("Error crawling %s: %v\n", url, err)
			continue
		}
		
		wc.mu.Lock()
		wc.results = append(wc.results, title)
		wc.mu.Unlock()
		
		time.Sleep(1 * time.Second)
	}
}

func (wc *WebCrawler) crawlPage(url string) (string, error) {
	resp, err := wc.client.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP error: %s", resp.Status)
	}
	
	// Extract title
	titleRegex := regexp.MustCompile(`<title[^>]*>(.*?)</title>`)
	titleMatches := titleRegex.FindStringSubmatch(resp.Status)
	if len(titleMatches) > 1 {
		return titleMatches[1], nil
	}
	
	// Extract links
	linkRegex := regexp.MustCompile(`<a[^>]*href=["\']([^"\']*)["\'][^>]*>`)
	linkMatches := linkRegex.FindAllStringSubmatch(resp.Status, -1)
	for _, match := range linkMatches {
		if len(match) > 1 {
			link := match[1]
			if strings.HasPrefix(link, "http") && !wc.isVisited(link) {
				wc.mu.Lock()
				<｜fim▁hole｜>
				wc.mu.Unlock()
			}
		}
	}
	
	return "Untitled", nil
}

func (wc *WebCrawler) isVisited(url string) bool {
	wc.mu.Lock()
	defer wc.mu.Unlock()
	return wc.visited[url]
}

func (wc *WebCrawler) GetResults() []string {
	wc.mu.Lock()
	defer wc.mu.Unlock()
	return wc.results
}

func main() {
	crawler := NewWebCrawler()
	
	fmt.Println("Starting web crawler...")
	crawler.Start("https://example.com")
	
	results := crawler.GetResults()
	fmt.Printf("Crawled %d pages\n", len(results))
	
	for i, title := range results {
		fmt.Printf("Page %d: %s\n", i+1, title)
	}
}