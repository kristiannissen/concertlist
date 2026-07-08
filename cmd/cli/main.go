package main

import (
	"fmt"
	"os"

	"log/slog"

	"github.com/kristiannissen/concertlist/internal/adapters/extractor"
	"github.com/kristiannissen/concertlist/internal/ports" // Importerer interfacet
	"github.com/kristiannissen/concertlist/internal/usecase/etl"
)

func main() {
	//
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, nil)))

	//
	e := extractor.NewSiteExtractor()

	var orc ports.Orchestrator = etl.NewOrchestrator(e)
	orc.RunCLI("https://richter-gladsaxe.dk/")

	fmt.Println("Hello Kitty")
}
