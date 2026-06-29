// Package handler provides the Vercel entry point for the concertlist API.
package handler

import (
	"net/http"

	"github.com/kristiannissen/concertlist/internal/adapters"
)

// Handler is the Vercel entry point that routes all /api/* requests.
// This function must be public for Vercel to detect it.
func Handler(w http.ResponseWriter, r *http.Request) {
	// Initialize the router from internal adapters
	router := adapters.NewRouter()
	router.ServeHTTP(w, r)
}
