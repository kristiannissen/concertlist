// Package handler provides the Vercel entry point for the concertlist API.
package handler

import (
	"fmt"
	"net/http"

	"github.com/kristiannissen/concertlist/internal/adapters"
)

// Handler is the Vercel entry point that routes all /api/* requests.
// This function must be public for Vercel to detect it.
func Handler(w http.ResponseWriter, r *http.Request) {
	// For browser requests, return a simple hello message
	if r.Header.Get("Accept") == "" || r.Header.Get("User-Agent") != "" {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		fmt.Fprint(w, "hello kitty")
		return
	}

	// For API requests, use the router from internal adapters
	router := adapters.NewRouter()
	router.ServeHTTP(w, r)
}
