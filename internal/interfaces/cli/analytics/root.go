package analytics

import (
	"github.com/spf13/cobra"
)

// AnalyticsCmd represents the analytics command group
var AnalyticsCmd = &cobra.Command{
	Use:   "analytics",
	Short: "Analytics and metrics commands",
	Long:  `Commands for viewing system analytics, metrics, and performance data`,
}
