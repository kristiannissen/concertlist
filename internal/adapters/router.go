// internal/adapters/router.go
package adapters

import (
	"encoding/json"
	"net/http"
)

func NewRouter() *http.ServeMux {
	mux := http.NewServeMux()

	//
	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, r *http.Request) {
		// Pass to queue
		// ctx := r.Context()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "Ok"})
	})
	// Returns 201 on success
	mux.HandleFunc("POST /api/v3/topic/musicevent", func(w http.ResponseWriter, r *http.Request) {
		// Pass to queue
		// ctx := r.Context()
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"status": "Created"})
	})

	return mux
}
