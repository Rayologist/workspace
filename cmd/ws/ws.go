package main

import (
	"os"

	"workspace/internal/cli"
	"workspace/internal/cmd/root"
)

func main() {
	r := cli.NewRuntime()

	rootCmd := root.New(r)
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
