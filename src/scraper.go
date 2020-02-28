package main

import (
	"fmt"

	// run go get -u github.com/gocolly/colly if you get an error for this package
	// if you are in your $GOPATH, it should find it. If not manually add it to the directory it is looking

	"github.com/gocolly/colly"

	//ALSO, if you haven't figured it out yet... run go build in the src directory, then to execute run src.exe
)

type pageInfo struct {
	StatusCode int
	Links      map[string]int
	Headings   map[string]int
 }

func main() {
	// Instantiate default collector
	c := colly.NewCollector(
		// MaxDepth is 2, so only the links on the scraped page
		// and links on those pages are visited
		colly.MaxDepth(1),
		colly.Async(true),

		// colly.AllowedDomains("espn.com"),
	)

	   // We add Headings here
	   p := &pageInfo{Links: make(map[string]int), Headings: make(map[string]int)}


	// Limit the maximum parallelism to 2
	// This is necessary if the goroutines are dynamically
	// created to control the limit of simultaneous requests.
	
	// Parallelism can be controlled also by spawning fixed
	// number of go routines.
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 2,
	})

   // count links and map to struct
   c.OnHTML("a[href]", func(e *colly.HTMLElement) {
	link := e.Request.AbsoluteURL(e.Attr("href"))
	if link != "" {
		// fmt.Println(link)
	   	p.Links[link]++

	}
 	})

	// count headings and map to Struct
	// h2 works for FoxNews titles
	c.OnHTML("h2", func(e *colly.HTMLElement) {
		// We are looping through the h2 elements and then getting the text of the a element
		heading := e.ChildText("a")
		if heading != "" {
			// fmt.Println(p.Headings[heading])

			fmt.Println(heading)
		   	p.Headings[heading]++
		}
	 })

	// On every a element which has href attribute call callback
	// c.OnHTML("a[href]", func(e *colly.HTMLElement) {
	// 	link := e.Attr("href")
	// 	// h1 := e.Attr("title")
	// 	fmt.Println(e.Text) 		//print text from page
	// 	// Print link
	// 	// fmt.Println(link)		//print links from page
	// 	// text := e.Attr("h1")

	// 	// fmt.Println(text)
	// 	// Visit link found on page on a new thread
	// 	e.Request.Visit(link)

	// })


	// Start scraping on.... en.wikipedia.org, foxnews.com, cnn.com, etc...

	// c.Visit("https://en.wikipedia.org/")

	c.Visit("https://www.foxnews.com/")

	// c.Visit("https://www.cnn.com/")

	// c.Visit("https://www.espn.com/")

	// Wait until threads are finished
	c.Wait()
}
