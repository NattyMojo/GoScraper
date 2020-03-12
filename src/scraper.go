package main

import (
	"fmt"

	"strconv"

	"strings"

	// run go get -u github.com/gocolly/colly if you get an error for this package
	// if you are in your $GOPATH, it should find it. If not manually add it to the directory it is looking

	"github.com/gocolly/colly"

	// used to eliminate stopwords 
	"github.com/bbalet/stopwords"

	//ALSO, if you haven't figured it out yet... run go build in the src directory, then to execute run src.exe
)

type pageInfo struct {
	StatusCode int
	Links      map[string]int
	Headings   map[string]int


 }

func main() {

	// Initialize PageInfo struct of maps
	pageDetails := &pageInfo{Links: make(map[string]int), Headings: make(map[string]int)}
	

	//Pass our struct to the Scrape function, able to modularize other processes
	scrape(*pageDetails)

	// //Print key and number of occurences ( element )
	// for key, element := range pageDetails.Headings {

	// 	fmt.Println("Key:", key, "=>", "Element:", element)
	// }

	fmt.Println("Total Keys: " + strconv.Itoa(len(pageDetails.Headings)))

	//Returns count for number of times key appears on page
	if value, ok := pageDetails.Headings["Explaining Trump's blowout numbers in primary states"]; ok {
		fmt.Println("value: ", value)
	} else {
		fmt.Println("key not found")
	}
	
	concatString := KeysString(pageDetails.Headings)

	noStopWords := stopwords.CleanString(concatString, "english", false)

	fmt.Println(noStopWords)
	// fmt.Println(concatString)
}


// Concats our keys together, then we can parse and split with Stop words
func KeysString(m map[string]int) string {
    keys := make([]string, 0, len(m))
    for k := range m {
        keys = append(keys, k)
    }
    return strings.Join(keys, " ")
}


func scrape(pageDetails pageInfo) {

		// Instantiate default collector
		c := colly.NewCollector(
			// IF MaxDepth is 2, only the links on the scraped page
			// and links on those pages are visited
			colly.MaxDepth(1),
			colly.Async(true),
	
			// colly.AllowedDomains("espn.com"),
		)
	
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
			   pageDetails.Links[link]++
	
		}
		 })
	
	
		// count headings and map to Struct
		// h2 works for FoxNews titles
		c.OnHTML("h2", func(e *colly.HTMLElement) {
			// We are looping through the h2 elements and then getting the text of the a element
			heading := e.ChildText("a")
			if heading != "" {
				// fmt.Println(p.Headings[heading])
	
				// fmt.Println( strconv.Itoa(len(pageDetails.Headings)) + " : " +  heading)
				   pageDetails.Headings[heading]++
	
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
	
	



		// START scraping on.... en.wikipedia.org, foxnews.com, cnn.com, etc...
		// c.Visit("https://en.wikipedia.org/")
		c.Visit("https://www.foxnews.com/")
		// c.Visit("https://www.cnbc.com/")
		// c.Visit("https://www.cnn.com/")
		// c.Visit("https://www.espn.com/")
	
		// Wait until threads are finished
		c.Wait()
}
