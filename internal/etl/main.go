package main

import (
	"concertlist/internal/etl/richtergladsaxe"
	"fmt"
)

func main() {
	// Create a channel
	ok := make(chan string)

	// Perfor async
	go func() {
		ok <- richtergladsaxe.Hello()
	}()

	msg := <-ok
	fmt.Println(msg)
}
