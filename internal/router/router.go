// Package router provides the HTTP routing configuration.
package router

import (
	"net/http"

	"github.com/kristiannissen/concertlist/internal/handler"
)

// NewRouter creates and configures the main ServeMux router.
func NewRouter() *http.ServeMux {
	mux := http.NewServeMux()

	// Health check endpoint
	mux.HandleFunc("GET /", handler.HelloHandler)

	// For now, just keep the hello handler for all paths
	// The events endpoint will be added when we wire up the dependencies
	mux.HandleFunc("GET /{...}", handler.HelloHandler)

	return mux
}
