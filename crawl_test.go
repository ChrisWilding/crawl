package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

var exampleHTMLWith3Links = `
<!doctype html>
<html>
<head>
    <title>Example Domain</title>
    <meta charset="utf-8" />
    <meta http-equiv="Content-type" content="text/html; charset=utf-8" />
    <meta name="viewport" content="width=device-width, initial-scale=1" />
</head>

<body>
<div>
    <h1>Example Domain</h1>
    <p>This domain is for use in illustrative examples in documents. You may use this
    domain in literature without prior coordination or asking for permission.</p>
    <p><a href="https://www.iana.org/domains/example">More information...</a></p>
	<a href="/domains/example">More information...</a>
	<a href="mailto:example@example.com">EMail</a>
</div>
</body>
</html>
`

func TestGet(t *testing.T) {
	expectedLinks := []string{
		"https://www.iana.org/domains/example",
		"/domains/example",
		"mailto:example@example.com",
	}

	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, exampleHTMLWith3Links)
	}))
	defer svr.Close()

	page := get(svr.URL)

	assert.Equal(t, svr.URL, page.url)
	assert.ElementsMatch(t, page.links, expectedLinks)
}

func TestFilter(t *testing.T) {
	links := []string{
		"/hello",
		"mailto:example@example.com",
		"https://www.example.com/world",
		"#fragment",
		"tel:0700000000",
		"app://www.example.com",
	}

	expected := []string{
		"https://www.example.com/hello",
		"https://www.example.com/world",
	}

	actual := filter(links, "https://www.example.com")

	assert.ElementsMatch(t, actual, expected)
}

func TestFilterSameDomain(t *testing.T) {
	links := []string{
		"https://www.example.com/hello",
		"https://subdomain.example.com/world",
		"https://www.example.com/world",
		"https://www.another-example.com/world",
	}

	expected := []string{
		"https://www.example.com/hello",
		"https://www.example.com/world",
	}

	actual := filterSameDomain(links, "https://www.example.com")

	assert.ElementsMatch(t, actual, expected)
}

func TestPrintPage(t *testing.T) {
	p := page{
		url: "https://www.example.com",
		links: []string{
			"https://www.example.com/hello",
			"https://www.example.com/world",
		},
	}
	var s strings.Builder
	printPage(p, &s)
	assert.Equal(t, s.String(), `Page: https://www.example.com
https://www.example.com/hello
https://www.example.com/world

`)
}

func TestFilterSeen(t *testing.T) {
	links := map[string]struct{}{
		"a": {},
		"b": {},
		"c": {},
	}
	seen := map[string]struct{}{
		"b": {},
	}
	actual := filterSeen(links, seen)
	assert.ElementsMatch(t, actual, []string{"a", "c"})
	assert.True(t, reflect.DeepEqual(seen, map[string]struct{}{
		"a": {},
		"b": {},
		"c": {},
	}))
}
