// Package handler provides the Vercel entry point for the concertlist API.
package handler

import (
	"net/http"

	"github.com/kristiannissen/concertlist/internal/router"
)

// Handler is the Vercel entry point that routes all /api/* requests.
// This function must be public for Vercel to detect it.
func Handler(w http.ResponseWriter, r *http.Request) {
	// Initialize the router from internal/router
	mux := router.NewRouter()
	mux.ServeHTTP(w, r)
}
