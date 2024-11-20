package root

import (
	"sns-tool/cmd/list"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "sns-tool",
	Short: "AWS SNS Cross-Account Subscription Checker",
	Long: `A CLI tool to check cross-account SNS subscriptions across different AWS accounts.
This tool helps identify Lambda subscriptions that are connected across different AWS accounts.`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.AddCommand(list.NewListCmd())
}
