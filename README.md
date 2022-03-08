# Crawl

## Introduction

We'd like you to write a simple web crawler in a programming language you're familiar with. Given a starting URL, the crawler should visit each URL it finds on the same domain. It should print each URL visited, and a list of links found on that page. The crawler should be limited to one subdomain - so when you start with *https://monzo.com/*, it would crawl all pages on the monzo.com website, but not follow external links, for example to facebook.com or community.monzo.com.

We would like to see your own implementation of a web crawler. Please do not use frameworks like scrapy or go-colly which handle all the crawling behind the scenes or someone else's code. You are welcome to use libraries to handle things like HTML parsing.

## Prerequisites

1. [Go 1.17](https://go.dev/dl/)

## Setup

```sh
$ git clone git@github.com:ChrisWilding/crawl.git
$ cd crawl
```

## How To

### Test

```sh
$ go test -v ./...
```

### Run

```
$ go build .
$ ./crawl --help
Usage of ./crawl:
  -url string
        the url to crawl (default "https://www.example.com")
```

### Run with Docker

```sh
$ docker pull ghcr.io/chriswilding/crawl:latest
$ docker run --rm -ti ghcr.io/chriswilding/crawl:latest --help
Usage of /ko-app/crawl:
  -url string
        the url to crawl (default "https://www.example.com")
```
