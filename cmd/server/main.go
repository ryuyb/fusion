package main

import (
	"log"

	"github.com/ryuyb/fusion/internal/cli"
)

func main() {
	if err := cli.Execute(); err != nil {
		log.Fatalf("Failed to execute command: %v", err)
	}
}
