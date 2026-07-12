// Package main is the general entry point for th CLI application
// This file is located at /cmd/cli/main.go
package main

import (
	"context"
	"sync"

	"github.com/kristiannissen/concertlist/internal/adapters/scrapers"
	"github.com/kristiannissen/concertlist/internal/ports"
	"go.uber.org/zap"
)

func main() {
	//
	logger, _ := zap.NewDevelopment()
	//
	defer logger.Sync()

	//
	logger.Info("Running")

	// Wire up the RicAx adapter behind the ports.Scraper interface.
	var scraper ports.Scraper = &scrapers.RicAx{
		URL: "https://richter-gladsaxe.dk/", // TODO: replace with the real target venue URL
		Log: logger,
	}

	ctx := context.Background()

	wg := &sync.WaitGroup{}

	// RicAx drives the scrape itself via colly using r.URL; data/contentType
	// aren't used by this adapter, so they're passed empty.
	err := scraper.Scrape(ctx, wg)
	//
	wg.Wait()
	if err != nil {
		logger.Error("scrape failed", zap.Error(err))
		return
	} else {
		logger.Info("Scrape done")
	}

}
