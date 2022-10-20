# Goodreads Quote Scraping

Here is a simple scraper written in Go using (Colly)[https://github.com/gocolly/colly] to scrape quotes from (Goodreads)[https://www.goodreads.com/quotes/].

# Usage

Included is an executable `colly_goodreads_scraper` (built on 64 bit linux):

```bash
$ ./colly_goodreads_scraper life
...
```

Alternatively you can install directly:
```bash
$ go install github.com/nikulpatel3141/colly_goodreads_scraper@latest
$ colly_goodreads_scraper humor # assuming GOPATH is in your PATH
...
```

... or run/build it:
```bash 
$ git clone https://github.com/nikulpatel3141/colly_goodreads_scraper
$ cd colly_goodreads_scraper
$ go run main.go wisdom
...
```

