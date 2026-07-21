package handler

import (
	"net/http"

	"github.com/kristiannissen/concertlist/gateway"
)

func Handler(w http.ResponseWriter, r *http.Request) {
	gateway.EventExtractConsumer(w, r)
}
