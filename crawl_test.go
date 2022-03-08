package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCrawl(t *testing.T) {
	crawl("")
	// TODO
}

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
	expected := []string{
		"https://www.iana.org/domains/example",
		"/domains/example",
		"mailto:example@example.com",
	}

	svr := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, exampleHTMLWith3Links)
	}))
	defer svr.Close()

	actual := get(svr.URL)

	assert.ElementsMatch(t, actual, expected)
}
