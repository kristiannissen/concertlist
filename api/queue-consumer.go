// Package handler is the Vercel entry point for the Vercel Queues consumer
// function. It's invoked exclusively by Vercel's queue infrastructure via
// the queue/v2beta trigger declared for this file in vercel.json — it has
// no public URL and cannot be reached over the internet.
package handler

import (
	"net/http"

	"github.com/kristiannissen/concertlist/gateway"
)

// Handler is the Vercel entry point for queue message delivery.
// This function must be public for Vercel to detect it.
func Handler(w http.ResponseWriter, r *http.Request) {
	gateway.QueueConsumer(w, r)
}
