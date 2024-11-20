package list

import (
	"encoding/csv"
	"fmt"
	"os"
	"time"

	"github.com/dhairya13703/sns-tool/pkg/models"

	"github.com/dhairya13703/sns-tool/pkg/aws"

	"github.com/spf13/cobra"
)

type ListOptions struct {
	AccountID  string
	Region     string
	Profile    string
	ExportCSV  string
	OutputFile string
}

func NewListCmd() *cobra.Command {
	opts := &ListOptions{}

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List cross-account SNS subscriptions",
		Long: `List all SNS topics that have Lambda subscriptions from different AWS accounts.
Examples:
  sns-tool list -a 123456789012
  sns-tool list -a 123456789012 -p dev
  sns-tool list -a 123456789012 -p prod -r us-west-2
  sns-tool list -a 123456789012 -e csv -o subscriptions.csv`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runList(opts)
		},
	}

	// Add flags with shorthand
	cmd.Flags().StringVarP(&opts.AccountID, "account", "a", "", "AWS account number to check")
	cmd.Flags().StringVarP(&opts.Region, "region", "r", "us-east-1", "AWS region to check")
	cmd.Flags().StringVarP(&opts.Profile, "profile", "p", "", "AWS profile to use")
	cmd.Flags().StringVarP(&opts.ExportCSV, "export", "e", "", "Export format (csv)")
	cmd.Flags().StringVarP(&opts.OutputFile, "output", "o", "", "Output file name")

	// Mark required flags
	cmd.MarkFlagRequired("account")

	return cmd
}

func runList(opts *ListOptions) error {
	snsClient, err := aws.NewSNSClient(aws.ClientConfig{
		Region:    opts.Region,
		AccountID: opts.AccountID,
		Profile:   opts.Profile,
	})
	if err != nil {
		return fmt.Errorf("failed to create SNS client: %v", err)
	}

	subscriptions, err := snsClient.ListCrossAccountSubscriptions()
	if err != nil {
		return fmt.Errorf("failed to list subscriptions: %v", err)
	}

	if len(subscriptions) == 0 {
		fmt.Println("No cross-account Lambda subscriptions found")
		return nil
	}

	// Handle CSV export if requested
	if opts.ExportCSV == "csv" {
		return exportToCSV(subscriptions, opts.OutputFile)
	}

	// Default console output
	fmt.Printf("\nFound %d cross-account Lambda subscriptions:\n\n", len(subscriptions))
	for _, sub := range subscriptions {
		fmt.Printf("Topic: %s\n", sub.TopicArn)
		fmt.Printf("Subscription: %s\n", sub.SubscriptionArn)
		fmt.Printf("Lambda Function: %s\n", sub.Endpoint)
		fmt.Printf("Owner Account: %s\n", sub.Owner)
		fmt.Println("---")
	}

	return nil
}

func exportToCSV(subscriptions []models.SubscriptionInfo, outputFile string) error {
	// If no output file specified, generate a default name
	if outputFile == "" {
		timestamp := time.Now().Format("20060102-150405")
		outputFile = fmt.Sprintf("sns-subscriptions-%s.csv", timestamp)
	}

	// Create or open the file
	file, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("failed to create output file: %v", err)
	}
	defer file.Close()

	// Create CSV writer
	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write header
	header := []string{"Topic ARN", "Subscription ARN", "Lambda Function", "Owner Account"}
	if err := writer.Write(header); err != nil {
		return fmt.Errorf("failed to write CSV header: %v", err)
	}

	// Write data
	for _, sub := range subscriptions {
		record := []string{
			sub.TopicArn,
			sub.SubscriptionArn,
			sub.Endpoint,
			sub.Owner,
		}
		if err := writer.Write(record); err != nil {
			return fmt.Errorf("failed to write CSV record: %v", err)
		}
	}

	fmt.Printf("Successfully exported to %s\n", outputFile)
	return nil
}
