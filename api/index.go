// Package handler provides the Vercel entry point for the concertlist API.
package handler

import (
	"fmt"
	"net/http"
)

// Handler is the Vercel entry point that routes all /api/* requests.
// This function must be public for Vercel to detect it.
//
// Deliberately minimal right now: no internal/ imports, no router, just a
// fixed response. This is to confirm the basic Vercel Go function deploy
// path works before reintroducing the app's routing logic.
func Handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Hello Kitty")
}
