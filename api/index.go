// Package handler provides the Vercel entry point for the concertlist API.
package handler

import (
	"fmt"
	"net/http"
)

// Handler is the Vercel entry point that handles /api requests.
// This function must be public for Vercel to detect it.
func Handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprint(w, "Hello Kitty")
}
