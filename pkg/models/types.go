package models

import "fmt"

// SubscriptionInfo represents an SNS subscription with cross-account information
type SubscriptionInfo struct {
	TopicArn        string
	SubscriptionArn string
	Endpoint        string
	Owner           string
}

// AccountMap represents a mapping of AWS account IDs to friendly names
var AccountMap = map[string]string{
	"571653956102": "SNSAccount", // Replace with your actual account numbers and names
	"471112726481": "LambdaAccount",
}

// GetAccountName returns the friendly name for an account ID if it exists
func GetAccountName(accountID string) string {
	if name, ok := AccountMap[accountID]; ok {
		return fmt.Sprintf("%s (%s)", name, accountID)
	}
	return accountID
}
