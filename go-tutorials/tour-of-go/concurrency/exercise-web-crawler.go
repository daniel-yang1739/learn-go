package main

import (
	"fmt"
	"sync"
)

type Fetcher interface {
	Fetch(url string) (body string, urls []string, err error)
}

type safeCounter struct {
	mu sync.Mutex
	visited map[string]bool
}

func (c *safeCounter) IsVisited(url string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()
	if _, ok := c.visited[url]; ok {
		return true
	}
	c.visited[url] = true
	return false
}

func Crawl(url string, depth int, fetcher Fetcher) {
	var wg sync.WaitGroup
	counter := safeCounter{ visited: make(map[string]bool) }

	var internalCrawl func(string, int)
	internalCrawl = func(url string, depth int) {
		defer wg.Done()
		if depth <= 0 {
			return
		}

		if counter.IsVisited(url) {
			return
		}

		body, urls, err := fetcher.Fetch(url)
		if err != nil {
			fmt.Println(err)
			return
		}
		fmt.Printf("found: %s %q\n", url, body)
		for _, u := range urls {
			wg.Add(1)
			go internalCrawl(u, depth-1)
		}
	}

	wg.Add(1)
	go internalCrawl(url, depth)
	wg.Wait()

	return
}

func main() {
	Crawl("https://golang.org/", 4, fetcher)
}

// fakeFetcher is Fetcher that returns canned results.
type fakeFetcher map[string]*fakeResult

type fakeResult struct {
	body string
	urls []string
}

func (f fakeFetcher) Fetch(url string) (string, []string, error) {
	if res, ok := f[url]; ok {
		return res.body, res.urls, nil
	}
	return "", nil, fmt.Errorf("not found: %s", url)
}

// fetcher is a populated fakeFetcher.
var fetcher = fakeFetcher{
	"https://golang.org/": &fakeResult{
		"The Go Programming Language",
		[]string{
			"https://golang.org/pkg/",
			"https://golang.org/cmd/",
		},
	},
	"https://golang.org/pkg/": &fakeResult{
		"Packages",
		[]string{
			"https://golang.org/",
			"https://golang.org/cmd/",
			"https://golang.org/pkg/fmt/",
			"https://golang.org/pkg/os/",
		},
	},
	"https://golang.org/pkg/fmt/": &fakeResult{
		"Package fmt",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
	"https://golang.org/pkg/os/": &fakeResult{
		"Package os",
		[]string{
			"https://golang.org/",
			"https://golang.org/pkg/",
		},
	},
}
