package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/ChrisWilding/crawl/link"
)

type page struct {
	url   string
	links []string
}

func get(url string) page {
	client := http.Client{
		Timeout: 15 * time.Second,
	}
	resp, err := client.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	links, _ := link.Parse(resp.Body)

	return page{
		url:   url,
		links: links,
	}
}

// filter returns a slice of all http(s) and relative links
// having filtered and removed any mailto, tel, app links or fragments
func filter(links []string, url string) []string {
	var filtered []string
	for _, l := range links {
		switch {
		case strings.HasPrefix(l, "/"):
			filtered = append(filtered, url+l)
		case strings.HasPrefix(l, "http"):
			filtered = append(filtered, l)
		}
	}
	return filtered
}

func filterSameDomain(links []string, url string) []string {
	var filtered []string
	for _, l := range links {
		if strings.HasPrefix(l, url) {
			filtered = append(filtered, l)
		}
	}
	return filtered
}

func filterSeen(links map[string]struct{}, seen map[string]struct{}) []string {
	var queue []string
	for link := range links {
		if _, ok := seen[link]; ok {
			continue
		}
		seen[link] = struct{}{}
		queue = append(queue, link)
	}
	return queue
}

func crawl(url string) []page {
	var pages []page
	var mu sync.Mutex

	seen := make(map[string]struct{})
	todo := make(map[string]struct{})
	todo[url] = struct{}{}
	next := make(map[string]struct{})

	for {
		queue := filterSeen(todo, seen)

		c := make(chan page, len(queue))
		for _, url := range queue {
			go func(url string) {
				page := get(url)
				c <- page
			}(url)
		}
		for i := 0; i < len(queue); i++ {
			page := <-c
			mu.Lock()
			pages = append(pages, page)
			links := filter(page.links, url)
			links = filterSameDomain(links, url)
			for _, link := range links {
				next[link] = struct{}{}
			}
			mu.Unlock()
		}

		todo = next
		next = make(map[string]struct{})

		if len(todo) == 0 {
			break
		}
	}

	return pages
}

func printPage(p page, w io.Writer) {
	fmt.Fprintf(w, "Page: %s\n", p.url)
	for _, l := range p.links {
		fmt.Fprintln(w, l)
	}
	fmt.Fprintln(w)
}
