// Package handler provides the Vercel entry point for the concertlist API.
package handler

import (
	"net/http"
	"strings"

	"github.com/kristiannissen/concertlist/gateway"
)

// Handler is the Vercel entry point that routes all /api/* requests.
// This function must be public for Vercel to detect it.
//
// This file must not import anything under internal/ directly — Vercel's
// Go builder wraps it in a synthetic package outside this module's import
// path, which trips Go's internal-package visibility rule. See
// gateway.NewRouter for why the indirection exists.
func Handler(w http.ResponseWriter, r *http.Request) {
	// Strip the /api prefix from the request path
	path := r.URL.Path
	if strings.HasPrefix(path, "/api") {
		r.URL.Path = strings.TrimPrefix(path, "/api")
	}

	mux := gateway.NewRouter()
	mux.ServeHTTP(w, r)
}
