package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"sync/atomic"
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

var page1 = `
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
	<a href="/page/two">P2</a>
	<a href="/page/three">P3</a>
</div>
</body>
</html>
`

var page2 = `
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
	<a href="/">P1</a>
	<a href="/page/three">P3</a>
</div>
</body>
</html>
`

var page3 = `
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
	<a href="/page/two">P2</a>
	<a href="/">P1</a>
</div>
</body>
</html>
`

func TestCrawl(t *testing.T) {
	var calls uint64
	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddUint64(&calls, 1)
		switch r.URL.Path {
		case "/":
			fmt.Fprint(w, page1)
		case "/page/two":
			fmt.Fprint(w, page2)
		case "/page/three":
			fmt.Fprint(w, page3)
		default:
			t.FailNow()
		}
	}))
	defer svr.Close()

	pages := crawl(svr.URL)

	assert.Equal(t, 3, len(pages))
	assert.Equal(t, uint64(3), calls)
}
