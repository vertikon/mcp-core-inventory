package main

import (
	"fmt"
	"os"

	"github.com/vertikon/mcp-core-inventory/pkg/logger"
)

// This executable is intentionally minimal; the real server wiring lives in the
// internal packages and will evolve alongside the domain use cases. Keeping a
// lightweight entry point avoids bit rot while still giving integrators a
// stable binary target.
func main() {
	if err := logger.Init("info", true); err != nil {
		fmt.Fprintf(os.Stderr, "failed to initialize logger: %v\n", err)
		os.Exit(1)
	}
	defer logger.Sync()

	logger.Info("mcp-core-inventory executable stub running")
}
