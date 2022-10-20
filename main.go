package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
)

const (
	urlBase    = "https://www.goodreads.com/quotes/tag/"
	jsStart    = "//<![CDATA"
	trimChars  = " ,\"“”"
	csvDelim   = "|"
	numPages   = 35
	numThreads = 10
)

func checkErr(e error) {
	if e != nil {
		panic(e)
	}
}

func pageURL(tag string, pagenum int) string {
	return urlBase + fmt.Sprintf("%s?page=%d", tag, pagenum)
}

// Take text from queried HTML and remove all extra leading + trailing
// characters and javascript
func processData(s string) string {
	if strings.Contains(s, jsStart) {
		s = s[:strings.Index(s, jsStart)]
	}
	s = strings.Split(s, "\n")[0]
	s = strings.TrimLeft(s, trimChars)
	s = strings.TrimRight(s, trimChars)
	return s
}

// Takes a div.quote block for a single quote and process it into
// a quote|author|title string for writing into a CSV file
func quoteParser(e *colly.HTMLElement) string {
	fields := []string{
		e.ChildText("div.quoteText"),
		e.ChildText("span.authorOrTitle"),
		e.ChildText("a.authorOrTitle"),
	}
	row := ""
	for _, x := range fields {
		row += processData(x) + csvDelim
	}
	return row
}

func setupCollyCollector(outFile *os.File) *colly.Collector {
	c := colly.NewCollector(colly.Async())
	c.Limit(&colly.LimitRule{DomainGlob: "*", Parallelism: numThreads, Delay: time.Second})

	c.OnHTML("div.quote", func(e *colly.HTMLElement) {
		quote_info := quoteParser(e)
		_, err := outFile.WriteString(quote_info + "\n")
		checkErr(err)
	})
	c.OnResponse(func(r *colly.Response) {
		fmt.Println("Visited", r.Request.URL)
	})
	return c
}

func main() {
	tag := os.Args[1]
	filename := tag + "_quotes.csv"
	fmt.Printf("Will scrape %s quotes from %s and save to %s\n", tag, urlBase, filename)

	f, err := os.Create(filename)
	checkErr(err)
	defer f.Close()

	c := setupCollyCollector(f)

	fmt.Println("Scraping...")
	startTime := time.Now()
	for i := 0; i < numPages; i++ {
		c.Visit(pageURL(tag, i))
	}
	c.Wait()
	elapsedTime := time.Since(startTime).Abs().Seconds()
	fmt.Printf("Done in %.2f seconds!\n", elapsedTime)
}
