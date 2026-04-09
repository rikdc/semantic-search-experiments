package cli

import (
	"fmt"
	"os"
)

// Execute is the entry point for the CLI. It parses os.Args, dispatches to the
// appropriate subcommand handler, and exits with a non-zero status on failure.
func Execute() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "usage: app <command> [flags]")
		fmt.Fprintln(os.Stderr, "commands: serve, migrate, version")
		os.Exit(1)
	}

	switch os.Args[1] {
	case "serve":
		cfg := ParseFlags(os.Args[2:])
		fmt.Printf("starting server on :%d (timeout: %s)\n", cfg.Port, cfg.Timeout)
	case "migrate":
		fmt.Println("running database migrations...")
	case "version":
		fmt.Println("v1.0.0")
	default:
		fmt.Fprintf(os.Stderr, "unknown command: %s\n", os.Args[1])
		os.Exit(1)
	}
}
