// Package main provides the entry point for Vercel deployment.
// This file is required by Vercel's Go runtime even though we use api/index.go as the actual handler.
package main

import (
	"net/http"

	"github.com/kristiannissen/concertlist/internal/adapters"
)

// This main function is required by Vercel but delegates to the same router used by api/index.go.
func main() {
	router := adapters.NewRouter()
	http.ListenAndServe(":8080", router) //nolint:errcheck
}
