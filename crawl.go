package main

import (
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

func crawl(url string) {
	// GET the page
	// Parse the HTML and extract the links
	// Print the URL
	// Print each of the links
	// Filter the links so that the slice only contains links on the same domain
	// If not already visited
	//   Crawl the link
	//   Mark as visited
}
