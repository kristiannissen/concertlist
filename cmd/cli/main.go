package main

import (
	"fmt"
	"log/slog"
	"os"
)

func main() {
	//
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, nil)))

	slog.Info("Hello", "who", "Kitty")
	fmt.Println("Hello Kitty")
}
