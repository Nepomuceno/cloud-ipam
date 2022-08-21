package model

import "github.com/Azure/azure-sdk-for-go/sdk/data/aztables"

type SubscriptionDefinition struct {
	aztables.Entity        // PartitionKey: <environment-id> | RowKey: <SubscriptionId>
	IpRanges        string `json:"ipRanges"` // ipRanges is a comma separated list of IP ranges
	Name            string `json:"name"`     // name is the name of the subscription
}
