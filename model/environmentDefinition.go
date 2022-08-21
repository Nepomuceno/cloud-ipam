package model

import "github.com/Azure/azure-sdk-for-go/sdk/data/aztables"

type EnvironmentDefinition struct {
	aztables.Entity        // PartitionKey: EnvironmentDefinition | RowKey: <EnvironmentId>
	IpRanges        string `json:"ipRanges"` // ipRanges is a comma separated list of IP ranges
	Name            string `json:"name"`     // name is the name of the environment
}
