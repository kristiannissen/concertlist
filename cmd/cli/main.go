// Package main is the general entry point for th CLI application
// This file is located at /cmd/cli/main.go
package main

import (
	"fmt"

	"go.uber.org/zap"
)

func main() {
	//
	logger, _ := zap.NewDevelopment()
	//
	defer logger.Sync()

	//
	logger.Info("Running")
	fmt.Println("Hello Kitty")
}
