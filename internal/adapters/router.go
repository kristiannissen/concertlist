// internal/adapters/router.go
package adapters

import (
	"encoding/json"
	"net/http"
	"sync"

	"github.com/kristiannissen/concertlist/internal/adapters/scrapers"
	"github.com/kristiannissen/concertlist/internal/ports"
	"go.uber.org/zap"
)

func NewRouter() *http.ServeMux {
	mux := http.NewServeMux()
	logger, _ := zap.NewDevelopment()
	//
	defer logger.Sync()

	//
	mux.HandleFunc("GET /api/health", func(w http.ResponseWriter, r *http.Request) {
		// Pass to queue
		// ctx := r.Context()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{"status": "Ok"})
	})
	// Returns 201 on success. Registered as GET because Vercel cron jobs
	// always trigger via HTTP GET, never POST.
	mux.HandleFunc("GET /api/musicevent/richter", func(w http.ResponseWriter, r *http.Request) {
		// Pass to queue
		var scraper ports.Scraper = &scrapers.Richter{
			URL: "https://richter-gladsaxe.dk",
			Log: logger,
		}
		wg := &sync.WaitGroup{}
		err := scraper.Scrape(r.Context(), wg)
		wg.Wait()
		// ctx := r.Context()
		w.Header().Set("Content-Type", "application/json")

		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"status": err.Error()})
			return
		}

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	return mux
}
