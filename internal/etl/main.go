package main

import (
	"concertlist/internal/etl/boilerplate"
	"concertlist/internal/etl/inators"
)

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
	bookstore := boilerplate.Resource{URL: "https://books.toscrape.com/"}

	inators.Runner(bookstore)
}
