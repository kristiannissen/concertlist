// Package main provides the entry point for Vercel deployment.
// This file is required by Vercel's Go runtime.
package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
)

// main is the entry point for Vercel's Go runtime.
// It must listen on the PORT environment variable.
func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/api", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(w, "Hello Kitty")
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
