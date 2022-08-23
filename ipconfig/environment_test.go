package ipconfig_test

import (
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/data/aztables"
	"github.com/google/uuid"
	"github.com/nepomuceno/cloud-ipam/ipconfig"
	"github.com/nepomuceno/cloud-ipam/model"
)

func TestAddEnvironment(t *testing.T) {
	client, err := getTestClient("ganepomutfstate")
	if err != nil {
		t.Fatal(err)
	}
	defer client.RemoveTable()
	environmentId := uuid.New().String()
	environmentName := uuid.New().String()
	environment := model.EnvironmentDefinition{
		Entity: aztables.Entity{
			PartitionKey: ipconfig.ENVIRONMENT_PARTITION_KEY,
			RowKey:       environmentId,
		},
		IPRanges: "10.0.0.0/8",
		Name:     environmentName,
	}
	err = client.AddOrUpdateEnvironment(environment)
	if err != nil {
		t.Fatal(err)
	}
	environments, err := client.ListEnvironments()
	if err != nil {
		t.Fatal(err)
	}
	if len(environments) != 1 {
		t.Fatalf("expected 1 environment, got %d", len(environments))
	}
	if environments[0].Name != environmentName {
		t.Fatalf("expected environment name %s, got %s", environmentName, environments[0].Name)
	}
	if environments[0].IPRanges != environment.IPRanges {
		t.Fatalf("expected environment ip ranges %s, got %s", environment.IPRanges, environments[0].IPRanges)
	}
	if environments[0].RowKey != environmentId {
		t.Fatalf("expected environment id %s, got %s", environmentId, environments[0].RowKey)
	}
	if environments[0].PartitionKey != environment.PartitionKey {
		t.Fatalf("expected environment partition key %s, got %s", environment.PartitionKey, environments[0].PartitionKey)
	}
}
