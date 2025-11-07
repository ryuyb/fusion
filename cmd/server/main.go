package main

import (
	"log"

	"github.com/ryuyb/fusion/internal/cli"
)

var (
	// Version is the application version, set via ldflags during build
	Version = "dev"
	// BuildTime is the build timestamp, set via ldflags during build
	BuildTime = "unknown"
	// GitCommit is the git commit hash, set via ldflags during build
	GitCommit = "unknown"
	// GoVersion is the Go version used to build, set via ldflags during build
	GoVersion = "unknown"
)

func main() {
	cli.SetVersionInfo(Version, BuildTime, GitCommit, GoVersion)
	if err := cli.Execute(); err != nil {
		log.Fatalf("Failed to execute command: %v", err)
	}
}
