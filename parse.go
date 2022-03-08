package main

import (
	"io"

	"golang.org/x/net/html"
)

func parse(r io.Reader) ([]string, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return nil, err
	}
	links := extractLinks(doc)
	return links, nil
}

func extractLinks(n *html.Node) []string {
	if n.Type == html.ElementNode && n.Data == "a" {
		return []string{extractHrefAttr(n)}
	}
	var hrefs []string
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		hrefs = append(hrefs, extractLinks(c)...)
	}
	return hrefs
}

func extractHrefAttr(n *html.Node) string {
	for _, attr := range n.Attr {
		if attr.Key == "href" {
			return attr.Val
		}
	}
	return ""
}
