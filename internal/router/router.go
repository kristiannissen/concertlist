// Package router provides the HTTP routing configuration.
package router

import (
	"net/http"

	"github.com/kristiannissen/concertlist/internal/handler"
)

// NewRouter creates and configures the main ServeMux router.
func NewRouter() *http.ServeMux {
	mux := http.NewServeMux()

	// Register the hello handler for the root path
	mux.HandleFunc("GET /", handler.HelloHandler)

	return mux
}
