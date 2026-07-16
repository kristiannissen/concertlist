// Package main is the general entry point for th CLI application
// This file is located at /cmd/cli/main.go
package main

import (
	"context"
	"flag"
	"sync"

	"github.com/kristiannissen/concertlist/internal/adapters/scrapers"
	"github.com/kristiannissen/concertlist/internal/ports"
	"go.uber.org/zap"
)

func main() {
	var venue string
	flag.StringVar(&venue, "venue", "", "enter venue name")
	flag.Parse()
	//
	logger, _ := zap.NewDevelopment()
	//
	defer logger.Sync()

	//
	logger.Info("Running")

	ctx := context.Background()

	wg := &sync.WaitGroup{}

	var s ports.Scraper

	switch venue {
	case "vega":
		s = &scrapers.Vega{URL: "https://vega.dk/?view=calendar", Log: logger}
	case "richter":
		s = &scrapers.Richter{URL: "https://richter-gladsaxe.dk/", Log: logger}
	default:
		logger.Info("No venue matches")
	}
	err := s.Scrape(ctx, wg)
	//
	wg.Wait()
	if err != nil {
		logger.Error("scrape failed", zap.Error(err))
	} else {
		logger.Info("Scrape done")
	}
}
