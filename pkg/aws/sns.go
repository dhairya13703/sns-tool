package aws

import (
	"context"
	"fmt"
	"strings"

	"sns-tool/pkg/models"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sns"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

// ClientConfig holds configuration for creating a new SNS client
type ClientConfig struct {
	Region    string
	AccountID string
	Profile   string
}

// SNSClient wraps the AWS SNS client with additional functionality
type SNSClient struct {
	snsClient *sns.Client
	stsClient *sts.Client
	accountID string
}

// NewSNSClient creates a new SNS client with the specified configuration and verifies the account
func NewSNSClient(cfg ClientConfig) (*SNSClient, error) {
	// Create the AWS SDK configuration options
	var opts []func(*config.LoadOptions) error

	// Add region
	opts = append(opts, config.WithRegion(cfg.Region))

	// Add profile if specified
	if cfg.Profile != "" {
		opts = append(opts, config.WithSharedConfigProfile(cfg.Profile))
	}

	// Load AWS config with options
	awsCfg, err := config.LoadDefaultConfig(context.TODO(), opts...)
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config: %v", err)
	}

	// Create clients
	stsClient := sts.NewFromConfig(awsCfg)
	snsClient := sns.NewFromConfig(awsCfg)

	// Verify account
	identity, err := stsClient.GetCallerIdentity(context.TODO(), &sts.GetCallerIdentityInput{})
	if err != nil {
		return nil, fmt.Errorf("unable to get caller identity: %v", err)
	}

	if *identity.Account != cfg.AccountID {
		return nil, fmt.Errorf("current credentials are for account %s, but requested account is %s",
			*identity.Account, cfg.AccountID)
	}

	return &SNSClient{
		snsClient: snsClient,
		stsClient: stsClient,
		accountID: cfg.AccountID,
	}, nil
}

// ListCrossAccountSubscriptions returns all Lambda subscriptions from different accounts
func (c *SNSClient) ListCrossAccountSubscriptions() ([]models.SubscriptionInfo, error) {
	crossAccountSubs := make([]models.SubscriptionInfo, 0)

	// List all SNS topics
	var nextToken *string
	for {
		topics, err := c.snsClient.ListTopics(context.TODO(), &sns.ListTopicsInput{
			NextToken: nextToken,
		})
		if err != nil {
			return nil, fmt.Errorf("unable to list topics: %v", err)
		}

		// Process each topic
		for _, topic := range topics.Topics {
			// List subscriptions for this topic with pagination
			var subNextToken *string
			for {
				subs, err := c.snsClient.ListSubscriptionsByTopic(context.TODO(), &sns.ListSubscriptionsByTopicInput{
					TopicArn:  topic.TopicArn,
					NextToken: subNextToken,
				})
				if err != nil {
					fmt.Printf("Warning: unable to list subscriptions for topic %s: %v\n", *topic.TopicArn, err)
					break
				}

				// Check each subscription
				for _, sub := range subs.Subscriptions {
					if *sub.Protocol == "lambda" {
						// Extract account ID from Lambda ARN
						lambdaArn := *sub.Endpoint
						parts := strings.Split(lambdaArn, ":")
						if len(parts) >= 6 {
							subAccountID := parts[4]
							if subAccountID != c.accountID {
								crossAccountSubs = append(crossAccountSubs, models.SubscriptionInfo{
									TopicArn:        *sub.TopicArn,
									SubscriptionArn: *sub.SubscriptionArn,
									Endpoint:        *sub.Endpoint,
									Owner:           models.GetAccountName(subAccountID),
								})
							}
						}
					}
				}

				// Check if there are more subscriptions
				if subs.NextToken == nil {
					break
				}
				subNextToken = subs.NextToken
			}
		}

		// Check if there are more topics
		if topics.NextToken == nil {
			break
		}
		nextToken = topics.NextToken
	}

	return crossAccountSubs, nil
}
