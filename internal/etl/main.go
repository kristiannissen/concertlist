package main

import (
	"concertlist/internal/etl/boilerplate"
	"concertlist/internal/etl/inators"
	"log"
	"time"

	"github.com/gocolly/colly"
)

const (
	cachedir = "./cache"
)

/**
* TODO: Create colly configuration
 * Pass colly as argument
*/
func main() {
	// Create a channel
	/*
		ok := make(chan string)

		// Perfor async
		go func() {
			ok <- richtergladsaxe.Hello()
		}()

		msg := <-ok
		fmt.Println(msg)
	*/
	coll := colly.NewCollector(
		colly.CacheDir(cachedir),
		colly.Async(true),
		colly.UserAgent("Mozilla/5.0 (iPhone; CPU iPhone OS 12_2 like Mac OS X)"),
	)

	coll.OnError(func(r *colly.Response, err error) {
		log.Fatal("Error ", err)
	})

	bookstore := boilerplate.Resource{URL: "https://books.toscrape.com/", Delay: 5 * time.Second}
	bookstore.New(coll)

	inators.Runner(bookstore)
}
