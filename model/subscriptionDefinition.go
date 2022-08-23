package model

import "github.com/Azure/azure-sdk-for-go/sdk/data/aztables"

type SubscriptionDefinition struct {
	aztables.Entity        // PartitionKey: <environment-id> | RowKey: <SubscriptionId>
	IPRanges        string `json:"ipRanges"` // ipRanges is a comma separated list of IP ranges
	Name            string `json:"name"`     // name is the name of the subscription
}

type SubcriptionValidationResult struct {
	Valid          bool   `json:"valid"`          // valid is true if the subscription is valid
	Error          error  `json:"error"`          // error is the error message if the subscription is invalid
	SubscriptionID string `json:"subscriptionId"` // subscriptionId is the subscription id
	EnvironmentID  string `json:"environmentId"`  // environmentId is the environment id
}
