// Package handler provides HTTP request handlers.
package handler

import (
	"fmt"
	"net/http"
)

// HelloHandler returns a simple "Hello Kitty" response.
func HelloHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	fmt.Fprint(w, "Hello Kitty")
}
