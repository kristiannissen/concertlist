// internal/adapters/router.go
package adapters

import (
	"encoding/json"
	"net/http"

	"github.com/go-resty/resty/v2"
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
	mux.HandleFunc("GET /api/scrape/trigger", func(w http.ResponseWriter, r *http.Request) {
		reg := NewScraperRegistry(logger)
		client := resty.New()
		client.SetAuthToken(r.Header.Get("x-vercel-oidc-token"))

		var failed []string
		for venue := range reg {
			body, _ := json.Marshal(map[string]string{"venue": venue})

			resp, err := client.R().SetBody(body).Post("https://arn1.vercel-queue.com/api/v3/topic/venue-scrape")
			if err != nil || resp.IsError() {
				logger.Error("failed to enqueue venue", zap.String("venue", venue), zap.Error(err))
				failed = append(failed, venue)
				continue // one bad publish shouldn't block the rest
			}
			logger.Info("enqueued venue for scraping", zap.String("venue", venue))
		}

		w.Header().Set("Content-Type", "application/json")
		if len(failed) > 0 {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]any{"status": "partial failure", "failed": failed})
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
	})

	return mux
}

// QueueConsumer is the skeleton handler for the "musicevent" Vercel Queues
// topic (the same topic Richter.Scrape publishes to). It's wired up as a
// queue/v2beta trigger in vercel.json, bound to its own serverless function
// (api/queue-consumer/index.go) rather than the shared mux in NewRouter — Vercel
// Queues triggers apply per-function, and NewRouter's function
// (api/index.go) is also the public entry point for /api/health and
// /api/musicevent/richter, so it can't be reused here without air-gapping
// those routes too.
//
// For now this just logs that a message arrived; fill in real message
// parsing/handling once the consumer's request/payload shape is confirmed.
func QueueConsumer(w http.ResponseWriter, r *http.Request) {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	logger.Info("received queue message", zap.String("path", r.URL.Path))

	w.WriteHeader(http.StatusOK)
}
