// Package router provides the HTTP routing configuration.
package router

import (
	"net/http"

	"github.com/kristiannissen/concertlist/pkg/handler"
)

// NewRouter creates and configures the main ServeMux router.
func NewRouter() *http.ServeMux {
	mux := http.NewServeMux()

	// Health check endpoint. In Go's net/http enhanced routing (1.22+), a
	// pattern ending in "/" already matches every path under it, so this
	// single registration acts as the catch-all for all paths.
	// The events endpoint will be added when we wire up the dependencies.
	mux.HandleFunc("GET /", handler.HelloHandler)

	return mux
}
