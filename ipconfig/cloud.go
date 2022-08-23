package ipconfig

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/network/armnetwork"
	"github.com/nepomuceno/cloud-ipam/exp"
	"github.com/nepomuceno/cloud-ipam/model"
)

func (client *IpConfigClient) Sync(environmentIDs []string) error {
	for _, environmentID := range environmentIDs {
		subscriptions, err := client.ListSubscriptions(environmentID)
		if err != nil {
			return err
		}
		for _, subscription := range subscriptions {
			err := client.ValidateSubscription(environmentID, subscription.RowKey)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (client *IpConfigClient) getSubDefinitionFromCloud(environmentID, subscriptionID, vnetName string) (*model.SubscriptionDefinition, error) {
	result := model.SubscriptionDefinition{}
	vnetClient, err := armnetwork.NewVirtualNetworksClient(subscriptionID, client.cred, nil)
	if err != nil {
		return nil, err
	}
	ipranges := make(map[string]bool)
	vnetsPager := vnetClient.NewListAllPager(nil)
	for vnetsPager.More() {
		vnets, err := vnetsPager.NextPage(client.ctx)
		if err != nil {
			return nil, err
		}
		for _, vnet := range vnets.Value {
			if strings.EqualFold(*vnet.Name, vnetName) || strings.EqualFold("", vnetName) {
				for _, iprange := range vnet.Properties.AddressSpace.AddressPrefixes {
					ipranges[*iprange] = true
				}
			}
		}
	}
	result.IPRanges = strings.Join(exp.Keys(ipranges), ",")

	return &result, err
}

func (client *IpConfigClient) ValidateSubscription(environmentID, subscriptionID string) error {
	entity, err := client.GetEntity(client.ctx, environmentID, subscriptionID, nil)
	if err != nil {
		return err
	}
	sub := model.SubscriptionDefinition{}
	err = json.Unmarshal(entity.Value, &sub)
	if err != nil {
		return err
	}
	vnetClient, err := armnetwork.NewVirtualNetworksClient(sub.RowKey, client.cred, nil)
	if err != nil {
		return err
	}
	ipRanges := []string{}
	allVnetsPager := vnetClient.NewListAllPager(nil)
	for allVnetsPager.More() {
		vnetPage, err := allVnetsPager.NextPage(client.ctx)
		if err != nil {
			return err
		}
		for _, vnet := range vnetPage.Value {
			for _, ipRange := range vnet.Properties.AddressSpace.AddressPrefixes {
				ipRanges = append(ipRanges, *ipRange)
			}
		}
	}
	actualIPRanges := strings.Join(ipRanges, ",")
	if actualIPRanges != sub.IPRanges {
		return fmt.Errorf("ip ranges do not match for sub %s. expected: %s, actual: %s", sub.RowKey, sub.IPRanges, actualIPRanges)
	}

	return nil

}
