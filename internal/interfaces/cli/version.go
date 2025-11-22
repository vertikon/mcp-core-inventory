package cli

import (
	"fmt"

	"github.com/spf13/cobra"
)

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print version information",
	Long:  `Print the version number and build information`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("mcp-core-inventory CLI\n")
		fmt.Printf("Version: %s\n", Version)
		fmt.Printf("Build Date: %s\n", BuildDate)
	},
}
