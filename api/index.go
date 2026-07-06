// Package handler provides the Vercel entry point for the concertlist API.
package handler

import (
	"net/http"
	"strings"

	"github.com/kristiannissen/concertlist/pkg/router"
)

// Handler is the Vercel entry point that routes all /api/* requests.
// This function must be public for Vercel to detect it.
func Handler(w http.ResponseWriter, r *http.Request) {
	// Strip the /api prefix from the request path
	path := r.URL.Path
	if strings.HasPrefix(path, "/api") {
		r.URL.Path = strings.TrimPrefix(path, "/api")
	}

	// Initialize the router from pkg/router
	mux := router.NewRouter()
	mux.ServeHTTP(w, r)
}
