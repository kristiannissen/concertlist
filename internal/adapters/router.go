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

	return mux
}
