// Package main provides the entry point for Vercel deployment.
// This file is required by Vercel's Go runtime.
package main

import (
	"net/http"
	"os"

	"github.com/kristiannissen/concertlist/gateway"
	"go.uber.org/zap"
)

// main is the entry point for Vercel's Go runtime.
// It must listen on the PORT environment variable.
func main() {
	//
	logger, _ := zap.NewDevelopment()
	defer logger.Sync()

	mux := gateway.NewRouter()

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	if err := http.ListenAndServe(":"+port, mux); err != nil {
		logger.Error("Server failed %v", zap.Error(err))
	}
}
