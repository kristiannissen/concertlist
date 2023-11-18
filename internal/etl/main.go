package main

import (
	"fmt"
	"internal/etl/boilerplate"
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

	fmt.Println("Hello Kitty")
}
