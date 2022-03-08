package main

import (
	"flag"
)

func main() {
	urlFlag := flag.String("url", "https://www.example.com", "the url to crawl")
	flag.Parse()
	crawl(*urlFlag)
}
