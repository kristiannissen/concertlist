// Package main provides the CLI entry point for testing concertlist components.
package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/kristiannissen/concertlist/pkg/adapters/etl/extractors/richter_gladsaxe"
	"github.com/kristiannissen/concertlist/pkg/domain"
)

func main() {
	// Define CLI flags
	venue := flag.String("venue", "richter_gladsaxe", "Venue extractor to test (currently only richter_gladsaxe is supported)")
	outputFormat := flag.String("format", "json", "Output format: json or text")
	flag.Parse()

	// Set up context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	var concerts []domain.Concert
	var err error

	// Test the specified extractor
	switch *venue {
	case "richter_gladsaxe":
		extractor := richter_gladsaxe.NewExtractor()
		concerts, err = extractor.Extract(ctx)
		if err != nil {
			log.Fatalf("Failed to extract from Richter Gladsaxe: %v", err)
		}
	default:
		log.Fatalf("Unsupported venue: %s. Currently only richter_gladsaxe is supported.", *venue)
	}

	// Output results based on format
	switch *outputFormat {
	case "json":
		jsonData, err := json.MarshalIndent(concerts, "", "  ")
		if err != nil {
			log.Fatalf("Failed to marshal concerts to JSON: %v", err)
		}
		fmt.Println(string(jsonData))
	case "text":
		if len(concerts) == 0 {
			fmt.Println("No concerts found.")
			return
		}
		fmt.Printf("Found %d concerts:", len(concerts))
		for i, concert := range concerts {
			fmt.Printf("%d. %s", i+1, concert.Title)
			fmt.Printf("   Date: %s", concert.Date)
			fmt.Printf("   Venue: %s", concert.Venue)
			if concert.URL != "" {
				fmt.Printf("   URL: %s", concert.URL)
			}
			if concert.Description != "" {
				fmt.Printf("   Description: %s", concert.Description)
			}
			fmt.Println()
		}
	default:
		log.Fatalf("Unsupported output format: %s. Use 'json' or 'text'.", *outputFormat)
	}
}
