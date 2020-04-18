package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"sort"
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

	concatString := KeysString(pageDetails.Headings)

	//cleanString gets rid of stopwords, english lang parameter,
	//and false is for HTML tags, so we could try to use this function earliar on with our scraped content.
	noStopWords := stopwords.CleanString(concatString, "english", false)

	// fmt.Println(concatString)
	// fmt.Println(noStopWords)

	wordCountMap := wordCount(noStopWords)

	// printSortedKey(pageDetails.Headings)
	printSortedKey(wordCountMap)

	fmt.Println("Total Keys: " + strconv.Itoa(len(wordCountMap)))

	keyWord1 := "coronavirus"
	keyWord2 := "quarantine"
	keyWord3 := "trump"
	keyWord4 := "trump's"
	keyWord5 := "biden"
	keyWord6 := "covid-"

	keyWordCount(wordCountMap, keyWord1)
	keyWordCount(wordCountMap, keyWord2)
	keyWordCount(wordCountMap, keyWord3)
	keyWordCount(wordCountMap, keyWord4)
	keyWordCount(wordCountMap, keyWord5)
	keyWordCount(wordCountMap, keyWord6)

	//Possibly build the index.html file here? (This does work btw after testing it)
	f, err := os.Create("index.html")
	if err != nil {
		fmt.Println("Had error creating: index.html")
		return
	}
			//This takes care of the beginning of the file, up to where we put in the first link
	l, err := f.WriteString("<!DOCTYPE html>\n<html lang=\"en\">\n<head>\n <meta charset=\"UTF-8\">\n<link rel=\"stylesheet\" href=\"style.css\">\n</head>\n<body>\n<div>\n")
	if err != nil {
		fmt.Println("Couldn't write to file: index.html",l)
		return
	}
			//Insert first URL, all we need is to change the URL under "src"
	l, err = f.WriteString("<iframe src=\"https://www.nbcnews.com/news/us-news/texas-gov-wants-slowly-reopen-private-business-trump-says-s-n1182886\" style=\"height:100%;width:50%;position:absolute;top:0;left:0\"></iframe>\n")
	if err != nil {
		fmt.Println("Couldn't write to file: index.html")
		return
	}
			//Insert Second URL, change the URL under "src"
	l, err = f.WriteString("<iframe src=\"https://www.breitbart.com/politics/2020/04/13/dr-anthony-fauci-repeatedly-downplayed-coronavirus-threat/\" style=\"height:100%;width:50%;position:absolute;top:0;right:0\"></iframe>\n")
	if err != nil {
		fmt.Println("Couldn't write to file: index.html")
		return
	}
			//Finish it off
	l, err = f.WriteString("</div>\n</body>\n</html>")
	if err != nil {
		fmt.Println("Couldn't write to file: index.html")
		return
	}
	//Close the File
	err = f.Close()
	if err != nil {
		fmt.Println("Couldn't close file: index.html")
		return
	}


	startDash()

}

func keyWordCount(countMap map[string]int, keyWord string) {

	//returns count on our cleaned up map of words for a provided keyWord
	if value, ok := countMap[keyWord]; ok {
		fmt.Println("KeyWord: "+keyWord+", # of times found: ", value)
	} else {
		fmt.Println("Key not found.")
	}

}

func printSortedKey(countMap map[string]int) {

	keys := make([]string, 0, len(countMap))
	for k := range countMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		fmt.Println(k, countMap[k])
	}
}

func wordCount(str string) map[string]int {
	wordList := strings.Fields(str)
	counts := make(map[string]int)
	for _, word := range wordList {
		_, ok := counts[word]
		if ok {
			counts[word] += 1
		} else {
			counts[word] = 1
		}
	}
	return counts
}

// Concats our keys together, then we can parse and split with Stop words
func KeysString(m map[string]int) string {
	keys := make([]string, 0, len(m))
	for k, element := range m {
		if element == 1 {
			keys = append(keys, (k))
		} else {
			for i := 0; i < element; i++ {
				keys = append(keys, (k))
			}
		}
	}
	joined := strings.Join(keys, " ")

	return joined
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
			//fmt.Println(link)
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
	c.Visit("https://www.breitbart.com/")
	c.Visit("https://www.nbcnews.com/")
	// c.Visit("https://en.wikipedia.org/")
	//c.Visit("https://www.foxnews.com/")
	// c.Visit("https://www.cnbc.com/")
	// c.Visit("https://www.cnn.com/")
	// c.Visit("https://www.espn.com/")

	// Wait until threads are finished
	c.Wait()
}
