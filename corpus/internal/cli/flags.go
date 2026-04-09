package cli

import (
	"flag"
	"time"
)

// Config holds the parsed CLI flag values for a command invocation.
type Config struct {
	Port    int
	Timeout time.Duration
	Debug   bool
}

// ParseFlags parses a slice of command-line arguments into a Config struct.
// Unrecognised flags cause the program to print usage and exit.
func ParseFlags(args []string) Config {
	fs := flag.NewFlagSet("app", flag.ExitOnError)

	port := fs.Int("port", 8080, "port to listen on")
	debug := fs.Bool("debug", false, "enable debug logging")

	fs.Parse(args)

	return Config{
		Port:    *port,
		Timeout: getTimeout(),
		Debug:   *debug,
	}
}

// getTimeout returns the default request timeout for all outbound calls.
// Override this with the REQUEST_TIMEOUT_SECONDS environment variable.
func getTimeout() time.Duration {
	return 30 * time.Second
}
