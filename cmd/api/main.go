// Package main provides the entry point for Vercel deployment.
// This file is required by Vercel's Go runtime.
package main

import (
	"log"
	"net/http"
	"os"

	"github.com/kristiannissen/concertlist/internal/adapters"
)

// main is the entry point for Vercel's Go runtime.
// It must listen on the PORT environment variable.
func main() {
	mux := adapters.NewRouter()

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
