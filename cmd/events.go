package cmd

import (
	"github.com/spf13/cobra"
)

var eventsCmd = &cobra.Command{
	Use:   "events",
	Short: "Manage CloudWatch Event integration",
	Long: `Manage CloudWatch Event integration

The events command provides subcommands for working with CloudWatch Events (scheduled tasks, etc.)`,
}

var ruleName string

func init() {
	rootCmd.AddCommand(eventsCmd)
}
