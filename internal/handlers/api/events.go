// Package api provides HTTP handlers for the API endpoints.
package api

import (
	"encoding/json"
	"net/http"

	"github.com/kristiannissen/concertlist/internal/domain"
)

// EventsHandler handles requests to the /events endpoint.
type EventsHandler struct {
	etlService domain.ETLService
}

// NewEventsHandler creates a new EventsHandler.
func NewEventsHandler(etlService domain.ETLService) *EventsHandler {
	return &EventsHandler{
		etlService: etlService,
	}
}

// GetEvents handles GET requests to retrieve all events.
func (h *EventsHandler) GetEvents(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	events, err := h.etlService.Run()
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(map[string]string{"error": err.Error()}) //nolint:errcheck
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(events); err != nil {
		// If we can't encode the response, there's not much we can do
		// The client will receive an empty response
		return
	}
}
