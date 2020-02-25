package main

import (
	"fmt"

	"github.com/gocolly/colly"
)

func main() {
	// Instantiate default collector
	c := colly.NewCollector(
		// MaxDepth is 2, so only the links on the scraped page
		// and links on those pages are visited
		colly.MaxDepth(1),
		colly.Async(true),

		// colly.AllowedDomains("espn.com"),
	)

	// Limit the maximum parallelism to 2
	// This is necessary if the goroutines are dynamically
	// created to control the limit of simultaneous requests.
	//
	// Parallelism can be controlled also by spawning fixed
	// number of go routines.
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 2,
	})

	// On every a element which has href attribute call callback
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		// Print link
		fmt.Println(link)
		// text := e.Attr("h1")

		// fmt.Println(text)
		// Visit link found on page on a new thread
		e.Request.Visit(link)

	})

	// Trying to print h1 tags from the pages...
	// c.OnHTML("a[h1]", func(e *colly.HTMLElement) {
	// 	h1 := e.Attr("h1")
	// 	// Print link
	// 	fmt.Println(h1)

	// })

	// Start scraping on https://en.wikipedia.org
	c.Visit("https://en.wikipedia.org/")

	// c.Visit("https://www.foxnews.com/")

	// c.Visit("https://www.cnn.com/")

	// c.Visit("https://www.espn.com/")
	// Wait until threads are finished
	c.Wait()
}
