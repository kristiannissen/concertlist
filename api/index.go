// Package handler provides the Vercel entry point for the concertlist API.
package handler

import (
	"net/http"

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
	gateway.NewRouter().ServeHTTP(w, r)
}
