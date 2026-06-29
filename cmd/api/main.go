// Package main provides the entry point for Vercel deployment.
// This file is required by Vercel's Go runtime.
package main

import (
	"log"
	"net/http"
	"os"
)

// main is the entry point for Vercel's Go runtime.
// It must listen on the PORT environment variable.
func main() {
	router := adapters.NewRouter()

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Printf("Starting server on :%s", port)
	if err := http.ListenAndServe(":"+port, router); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
