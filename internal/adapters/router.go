// internal/adapters/router.go
package adapters

import (
	"context"
	"encoding/json"
	"net/http"
	"os"
	"sync"

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
		reg := NewScraperRegistry(logger, r.Header.Get("x-vercel-oidc-token"))
		client := resty.New()
		client.SetAuthToken(r.Header.Get("x-vercel-oidc-token"))

		var failed []string
		for venue := range reg {
			body, _ := json.Marshal(map[string]string{"venue": venue})

			resp, err := client.R().
				SetHeader("Vqs-Deployment-Id", os.Getenv("VERCEL_DEPLOYMENT_ID")).
				SetBody(body).
				Post("https://arn1.vercel-queue.com/api/v3/topic/venue-scrape")
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

// QueueScrapeConsumer is the handler for the "venue-scrape" Vercel Queues topic.
// It's wired up as a queue/v2beta trigger in vercel.json, bound to its own
// serverless function (api/queue-consumer/index.go) rather than the shared
// mux in NewRouter — Vercel Queues triggers apply per-function, and
// NewRouter's function (api/index.go) is also the public entry point for
// /api/health and /api/scrape/trigger, so it can't be reused here without
// air-gapping those routes too.
func EventScrapeConsumer(w http.ResponseWriter, r *http.Request) {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	var msg struct {
		Venue string `json:"venue"`
	}
	json.NewDecoder(r.Body).Decode(&msg)
	logger.Info("incoming message", zap.String("venue", msg.Venue))

	scraper, ok := NewScraperRegistry(logger, r.Header.Get("x-vercel-oidc-token"))[msg.Venue]
	if !ok {
		logger.Error("unknown venue", zap.String("venue", msg.Venue))
		w.WriteHeader(http.StatusOK)
		return
	}

	ctx := context.WithValue(r.Context(), "x-vercel-oidc-token", r.Header.Get("x-vercel-oidc-token"))

	wg := &sync.WaitGroup{}

	if err := scraper.Scrape(ctx, wg); err != nil {
		logger.Error("scrape failed", zap.String("venue", msg.Venue), zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	wg.Wait()

	logger.Info("scrape complete", zap.String("venue", msg.Venue))
	w.WriteHeader(http.StatusOK)
}

func EventExtractConsumer(w http.ResponseWriter, r *http.Request) {
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	var msg struct {
		Venue string `json:"venue"`
		URL   string `json:"url"`
	}
	json.NewDecoder(r.Body).Decode(&msg)
	logger.Info("incoming message", zap.String("venue", msg.Venue))

	scraper, ok := NewScraperRegistry(logger, r.Header.Get("x-vercel-oidc-token"))[msg.Venue]
	if !ok {
		logger.Error("unknown venue", zap.String("venue", msg.Venue))
		w.WriteHeader(http.StatusOK)
		return
	}

	wg := &sync.WaitGroup{}

	if err := scraper.Extract(r.Context(), wg, msg.URL); err != nil {
		logger.Error("extract failed", zap.String("URL", msg.URL), zap.String("venue", msg.Venue), zap.Error(err))
		w.WriteHeader(http.StatusInternalServerError)
	}
	wg.Wait()

	logger.Info("extract complete", zap.String("URL", msg.URL))
	w.WriteHeader(http.StatusOK)
}
