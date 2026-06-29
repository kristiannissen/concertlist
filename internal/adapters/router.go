// Package adapters contains HTTP handlers and routing implementations.
package adapters

import (
	"net/http"
)

// NewRouter creates and configures the main router for the application.
// This router is used by both Vercel (api/index.go) and local development (cmd/api/main.go).
func NewRouter() *http.ServeMux {
	mux := http.NewServeMux()

	// TODO: Register your routes here
	// Example: mux.HandleFunc("GET /api/concerts", handleGetConcerts)
	// Example: mux.HandleFunc("POST /api/concerts", handleCreateConcert)

	// Health check endpoint
	mux.HandleFunc("GET /api/health", handleHealthCheck)

	return mux
}

// handleHealthCheck is a simple health check endpoint.
func handleHealthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status": "ok"}`)) //nolint:errcheck
}
