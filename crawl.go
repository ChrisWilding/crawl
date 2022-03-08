package main

import (
	"log"
	"net/http"
)

func get(url string) []string {
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	links, _ := parse(resp.Body)
	return links
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
