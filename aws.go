package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

type SubscriptionInfo struct {
	TopicArn        string
	SubscriptionArn string
	Endpoint        string
	Owner           string
}

var knownAccounts = map[string]string{
	"571653956102": "SNSAccount", // Replace with your actual account numbers and names
	"471112726481": "LambdaAccount",
}

func listSubscriptions(accountNumber string) {
	// Load AWS config
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion("us-east-1"),
	)
	if err != nil {
		log.Fatalf("unable to load SDK config: %v", err)
	}

	// Create STS client to verify current identity
	stsClient := sts.NewFromConfig(cfg)
	identity, err := stsClient.GetCallerIdentity(context.TODO(), &sts.GetCallerIdentityInput{})
	if err != nil {
		log.Fatalf("unable to get caller identity: %v", err)
	}

	if *identity.Account != accountNumber {
		log.Fatalf("Current credentials are for account %s, but requested account is %s", *identity.Account, accountNumber)
	}

	// Create SNS client
	snsClient := sns.NewFromConfig(cfg)

	// List all SNS topics
	topics, err := snsClient.ListTopics(context.TODO(), &sns.ListTopicsInput{})
	if err != nil {
		log.Fatalf("unable to list topics: %v", err)
	}

	crossAccountSubs := make([]SubscriptionInfo, 0)

	// For each topic, list subscriptions
	for _, topic := range topics.Topics {
		subs, err := snsClient.ListSubscriptionsByTopic(context.TODO(), &sns.ListSubscriptionsByTopicInput{
			TopicArn: topic.TopicArn,
		})
		if err != nil {
			log.Printf("unable to list subscriptions for topic %s: %v", *topic.TopicArn, err)
			continue
		}

		// Check each subscription
		for _, sub := range subs.Subscriptions {
			if *sub.Protocol == "lambda" {
				// Extract account ID from Lambda ARN
				lambdaArn := *sub.Endpoint
				parts := strings.Split(lambdaArn, ":")
				if len(parts) >= 6 {
					subAccountID := parts[4]
					if subAccountID != accountNumber {
						crossAccountSubs = append(crossAccountSubs, SubscriptionInfo{
							TopicArn:        *sub.TopicArn,
							SubscriptionArn: *sub.SubscriptionArn,
							Endpoint:        *sub.Endpoint,
							Owner:           getAccountName(subAccountID),
						})
					}
				}
			}
		}
	}

	// Print results
	if len(crossAccountSubs) == 0 {
		fmt.Println("No cross-account Lambda subscriptions found")
		return
	}

	fmt.Printf("\nFound %d cross-account Lambda subscriptions:\n\n", len(crossAccountSubs))
	for _, sub := range crossAccountSubs {
		fmt.Printf("Topic: %s\n", sub.TopicArn)
		fmt.Printf("Subscription: %s\n", sub.SubscriptionArn)
		fmt.Printf("Lambda Function: %s\n", sub.Endpoint)
		fmt.Printf("Owner Account: %s\n", sub.Owner)
		fmt.Println("---")
	}
}

func getAccountName(accountID string) string {
	if name, ok := knownAccounts[accountID]; ok {
		return fmt.Sprintf("%s (%s)", name, accountID)
	}
	return accountID
}
