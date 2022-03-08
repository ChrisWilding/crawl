package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

type page struct {
	url   string
	links []string
}

func get(url string) page {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	links, _ := parse(resp.Body)

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

func crawl(url string) []page {
	var pages []page

	seen := make(map[string]struct{})
	todo := make(map[string]struct{})
	todo[url] = struct{}{}
	next := make(map[string]struct{})

	for {
		for url := range todo {
			if _, ok := seen[url]; ok {
				continue
			}
			seen[url] = struct{}{}
			page := get(url)
			pages = append(pages, page)

			links := filter(page.links, url)
			links = filterSameDomain(links, url)

			for _, link := range links {
				next[link] = struct{}{}
			}
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
