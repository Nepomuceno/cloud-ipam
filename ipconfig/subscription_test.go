package ipconfig_test

import (
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/data/aztables"
	"github.com/google/uuid"
	"github.com/mepomuceno/cloud-ipam/ipconfig"
	"github.com/mepomuceno/cloud-ipam/model"
)

func TestAddSubscription(t *testing.T) {
	client, environmentId := setupSubscriptionTest(t)
	defer client.RemoveTable()
	subscriptionId := uuid.New().String()
	subscriptionName := uuid.New().String()
	subscription := model.SubscriptionDefinition{
		Entity: aztables.Entity{
			PartitionKey: environmentId,
			RowKey:       subscriptionId,
		},
		IPRanges: "10.0.0.0/24",
		Name:     subscriptionName,
	}
	err := client.AddOrUpdateSubscription(subscription)
	if err != nil {
		t.Fatal(err)
	}
	subscriptions, err := client.ListSubscriptions(environmentId)
	if err != nil {
		t.Fatal(err)
	}
	if len(subscriptions) != 1 {
		t.Fatalf("expected 1 subscription, got %d", len(subscriptions))
	}
	if subscriptions[0].Name != subscriptionName {
		t.Fatalf("expected subscription name %s, got %s", subscriptionName, subscriptions[0].Name)
	}
	if subscriptions[0].IPRanges != subscription.IPRanges {
		t.Fatalf("expected subscription ip ranges %s, got %s", subscription.IPRanges, subscriptions[0].IPRanges)
	}
	if subscriptions[0].RowKey != subscriptionId {
		t.Fatalf("expected subscription id %s, got %s", subscriptionId, subscriptions[0].RowKey)
	}
	if subscriptions[0].PartitionKey != subscription.PartitionKey {
		t.Fatalf("expected subscription partition key %s, got %s", subscription.PartitionKey, subscriptions[0].PartitionKey)
	}
}

func TestNextIpRange(t *testing.T) {
	client, environmentId := setupSubscriptionTest(t)
	defer client.RemoveTable()
	subscriptionId := uuid.New().String()
	subscriptionName := uuid.New().String()
	subscription := model.SubscriptionDefinition{
		Entity: aztables.Entity{
			PartitionKey: environmentId,
			RowKey:       subscriptionId,
		},
		IPRanges: "10.0.0.0/24",
		Name:     subscriptionName,
	}
	err := client.AddOrUpdateSubscription(subscription)
	if err != nil {
		t.Fatal(err)
	}
	subscriptionId = uuid.New().String()
	subscriptionName = uuid.New().String()
	subscription = model.SubscriptionDefinition{
		Entity: aztables.Entity{
			PartitionKey: environmentId,
			RowKey:       subscriptionId,
		},
		IPRanges: "10.0.1.0/24",
		Name:     subscriptionName,
	}
	err = client.AddOrUpdateSubscription(subscription)
	if err != nil {
		t.Fatal(err)
	}
	ip, err := client.GetNextAvailableIpRange(environmentId, 24)
	if err != nil {
		t.Fatal(err)
	}
	if ip.String() != "10.0.2.0/24" {
		t.Fatalf("expected ip range %s, got %s", "10.0.2.0/24", ip.String())
	}
}

func TestNextIpRangeFirstIp(t *testing.T) {
	client, environmentId := setupSubscriptionTest(t)
	defer client.RemoveTable()
	ip, err := client.GetNextAvailableIpRange(environmentId, 24)
	if err != nil {
		t.Fatal(err)
	}
	if ip.String() != "10.0.0.0/24" {
		t.Fatalf("expected ip range %s, got %s", "10.0.0.0/24", ip.String())
	}
}

func setupSubscriptionTest(t *testing.T) (*ipconfig.IpConfigClient, string) {
	client, err := getTestClient("ganepomutfstate")
	if err != nil {
		t.Fatal(err)
	}
	environmentId := uuid.New().String()
	environmentName := uuid.New().String()
	environment := model.EnvironmentDefinition{
		Entity: aztables.Entity{
			PartitionKey: ipconfig.ENVIRONMENT_PARTITION_KEY,
			RowKey:       environmentId,
		},
		IPRanges: "10.0.0.0/16",
		Name:     environmentName,
	}
	err = client.AddOrUpdateEnvironment(environment)
	if err != nil {
		t.Fatal(err)
	}
	return client, environmentId
}
