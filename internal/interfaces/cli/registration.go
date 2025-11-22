package cli

import (
	"sync"

	"github.com/spf13/cobra"
	"github.com/vertikon/mcp-core-inventory/internal/interfaces/cli/analytics"
	"github.com/vertikon/mcp-core-inventory/internal/interfaces/cli/ci"
	"github.com/vertikon/mcp-core-inventory/pkg/logger"
)

var registerOnce sync.Once

// RegisterCommands registers all CLI commands to the root command
// This centralizes command registration to avoid init() conflicts
func RegisterCommands(rootCmd *cobra.Command) {
	registerOnce.Do(func() {
		logger.Info("Registering CLI commands")

		// Register main command groups
		rootCmd.AddCommand(analytics.AnalyticsCmd)
		rootCmd.AddCommand(ci.CICmd)

		// Register individual commands
		rootCmd.AddCommand(generateCmd)
		rootCmd.AddCommand(monitorCmd)
		rootCmd.AddCommand(stateCmd)
		rootCmd.AddCommand(templateCmd)
		rootCmd.AddCommand(versionCmd)
		rootCmd.AddCommand(aiCmd)

		// Register subcommands and flags
		configureFlags()

		logger.Info("All CLI commands registered successfully")
	})
}

// configureFlags configures flags for all commands
func configureFlags() {
	// Generate command flags
	generateCmd.Flags().StringP("template", "t", "", "Template ID to use")
	generateCmd.Flags().StringP("output", "o", ".", "Output directory")
	generateCmd.MarkFlagRequired("template")

	// Template command flags
	templateCmd.AddCommand(templateListCmd)
	templateCmd.AddCommand(templateCreateCmd)
	templateCreateCmd.Flags().StringP("name", "n", "", "Template name")
	templateCreateCmd.MarkFlagRequired("name")

	// AI command flags
	aiCmd.AddCommand(aiChatCmd)
	aiChatCmd.Flags().StringP("message", "m", "", "Message to send")
	aiChatCmd.MarkFlagRequired("message")
}
