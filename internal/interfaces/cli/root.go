package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/vertikon/mcp-core-inventory/pkg/logger"
	"go.uber.org/zap"
)

var (
	// Version is set at build time
	Version = "dev"
	// BuildDate is set at build time
	BuildDate = "unknown"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "hulk",
	Short: "mcp-core-inventory CLI - Generate and manage MCP projects",
	Long: `mcp-core-inventory is a powerful CLI tool for generating and managing
Model Context Protocol (MCP) projects.`,
	Version: fmt.Sprintf("%s (built %s)", Version, BuildDate),
}

// Execute adds all child commands to the root command and sets flags appropriately
func Execute() {
	RegisterCommands(rootCmd)
	if err := rootCmd.Execute(); err != nil {
		logger.Error("Command execution failed", zap.Error(err))
		os.Exit(1)
	}
}

// init initializes CLI - MOVED TO registration.go to avoid conflicts
// All command registration is now centralized in registration.go

// GetRootCmd returns root command (for testing)
func GetRootCmd() *cobra.Command {
	return rootCmd
}
