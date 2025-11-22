package ci

import "github.com/spf13/cobra"

// CICmd aggregates CI-related commands.
var CICmd = &cobra.Command{
	Use:   "ci",
	Short: "Continuous integration helpers",
	Long:  "Commands that integrate MCP projects with CI/CD workflows.",
}
