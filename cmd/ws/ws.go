package main

import (
	"os"

	"workspace/internal/cmd/root"
)

func main() {
	rootCmd := root.New()
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
