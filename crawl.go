package main

import (
	"bytes"
	"fmt"
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

func (p page) String() string {
	b := new(bytes.Buffer)
	fmt.Fprintf(b, "Page: %s\n", p.url)
	for _, l := range p.links {
		fmt.Fprintln(b, l)
	}
	fmt.Fprintln(b)
	return b.String()
}

func get(url string) page {
	client := http.Client{
		Timeout: 15 * time.Second,
	}
	resp, err := client.Get(url)
	if err != nil {
		log.Printf("An error occured getting %s, error was %v\n", url, err)
	}
	defer resp.Body.Close()

	links, err := link.Parse(resp.Body)
	if err != nil {
		log.Printf("An error occured parsing HTML for URL %s, error was %v\n", url, err)
	}

	return page{
		url:   url,
		links: links,
	}
}

// filterIsSameDomain returns a slice of all http(s) and relative links
// having filtered and removed any mailto, tel, app links, fragments
// or links on other domains
func filterIsSameDomain(links []string, url string) []string {
	var filtered []string
	for _, l := range links {
		switch {
		case strings.HasPrefix(l, "/"):
			filtered = append(filtered, url+l)
		case strings.HasPrefix(l, "http") && strings.HasPrefix(l, url):
			filtered = append(filtered, l)
		}
	}
	return filtered
}

// filterIsUnseen returns a slice of the unseen links from links
// and adds any unseen links to seen
func filterIsUnseen(links map[string]struct{}, seen map[string]struct{}) []string {
	var unseen []string
	for link := range links {
		if _, ok := seen[link]; ok {
			continue
		}
		seen[link] = struct{}{}
		unseen = append(unseen, link)
	}
	return unseen
}

func crawl(url string, limit int) []page {
	var pages []page
	var mu sync.Mutex

	seen := make(map[string]struct{})
	todo := make(map[string]struct{})
	todo[url+"/"] = struct{}{}
	next := make(map[string]struct{})

	for i := 0; i < limit; i++ {
		queue := filterIsUnseen(todo, seen)

		c := make(chan page, len(queue))
		for _, url := range queue {
			go func(url string) {
				c <- get(url)
			}(url)
		}
		for i := 0; i < len(queue); i++ {
			page := <-c
			mu.Lock()
			pages = append(pages, page)
			for _, link := range filterIsSameDomain(page.links, url) {
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
