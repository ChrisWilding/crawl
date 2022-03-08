package main

import (
	"flag"
	"fmt"
)

func main() {
	urlFlag := flag.String("url", "https://www.example.com", "the url to crawl")
	limitFlag := flag.Int("limit", 100, "limit to the number of levels of links to follow")
	flag.Parse()
	pages := crawl(*urlFlag, *limitFlag)
	for _, p := range pages {
		fmt.Printf("%v", p)
	}
}
