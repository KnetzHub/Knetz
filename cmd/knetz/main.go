package main

import (
	"fmt"
	"os"

	"github.com/knetz-io/knetz/cmd/knetz/commands"
)

var (
	// Version information (injected at build time)
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

func main() {
	// Set version info
	commands.Version = version
	commands.Commit = commit
	commands.Date = date

	if err := commands.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

