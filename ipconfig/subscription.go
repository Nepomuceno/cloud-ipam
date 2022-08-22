package ipconfig

import (
	"encoding/json"
	"fmt"
	"math"
	"net"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/data/aztables"
	"github.com/mepomuceno/cloud-ipam/model"
)

func (client *IpConfigClient) AddOrUpdateSubscription(subscription model.SubscriptionDefinition) error {
	entity, err := json.Marshal(subscription)
	if err != nil {
		return err
	}
	_, err = client.UpsertEntity(client.ctx, entity, nil)
	return err
}

func (client *IpConfigClient) DeleteSubscription(environmnetID, id string) error {
	_, err := client.DeleteEntity(client.ctx, environmnetID, id, nil)
	return err
}

func (client *IpConfigClient) GetSubscription(environmentID, id string) (model.SubscriptionDefinition, error) {
	result := model.SubscriptionDefinition{}
	entity, err := client.GetEntity(client.ctx, environmentID, id, nil)
	if err != nil {
		return result, err
	}
	err = json.Unmarshal(entity.Value, &result)
	return result, err
}

func (client *IpConfigClient) ListSubscriptions(environmentID string) ([]model.SubscriptionDefinition, error) {
	result := []model.SubscriptionDefinition{}
	pages := client.NewListEntitiesPager(&aztables.ListEntitiesOptions{
		Filter: to.Ptr("PartitionKey eq '" + environmentID + "'"),
	})
	for pages.More() {
		pageResp, err := pages.NextPage(client.ctx)
		if err != nil {
			return result, err
		}
		for _, entity := range pageResp.Entities {
			var subscription model.SubscriptionDefinition
			err = json.Unmarshal(entity, &subscription)
			if err != nil {
				return result, err
			}
			result = append(result, subscription)
		}
	}
	return result, nil
}

func (client *IpConfigClient) GetNextAvailableIpRange(environmentID string, subnetRangeSize int) (*net.IPNet, error) {
	netEnvironment, err := client.GetEnvironment(environmentID)
	if err != nil {
		return nil, err
	}
	ipranges := strings.Split(netEnvironment.IPRanges, ",")
	usedIpRanges, err := client.GetUsedRangesForEnvironment(environmentID)
	if err != nil {
		return nil, err
	}
	for _, iprange := range ipranges {
		_, mainRange, err := net.ParseCIDR(iprange)
		if err != nil {
			return nil, err
		}
		ipmask := net.CIDRMask(subnetRangeSize, 32)
		ipnetResult := net.IPNet{IP: mainRange.IP, Mask: ipmask}
		for ok := ipPresentinAny(usedIpRanges, ipnetResult.IP); ok && mainRange.Contains(ipnetResult.IP); ok = ipPresentinAny(usedIpRanges, ipnetResult.IP) {
			newIp := nextIP(ipnetResult.IP, uint(math.Pow(2, float64(32-subnetRangeSize))))
			ipnetResult.IP = newIp
		}
		if mainRange.Contains(ipnetResult.IP) {
			return &ipnetResult, nil
		}

	}
	return nil, fmt.Errorf("no available ip range found")
}

func ipPresentinAny(ipRanges []*net.IPNet, ip net.IP) bool {
	for _, ipRange := range ipRanges {
		if ipRange.Contains(ip) {
			return true
		}
	}
	return false
}

func (client *IpConfigClient) GetUsedRangesForEnvironment(environmentID string) ([]*net.IPNet, error) {
	result := make([]*net.IPNet, 0)
	pages := client.NewListEntitiesPager(&aztables.ListEntitiesOptions{
		Filter: to.Ptr(fmt.Sprintf("PartitionKey eq '%s'", environmentID)),
	})
	for pages.More() {
		pageResp, err := pages.NextPage(client.ctx)
		if err != nil {
			return result, err
		}
		for _, entity := range pageResp.Entities {
			var definition model.SubscriptionDefinition
			err = json.Unmarshal(entity, &definition)
			if err != nil {
				return result, err
			}
			ipRanges := strings.Split(definition.IPRanges, ",")
			for _, ipRange := range ipRanges {
				_, ipnet, err := net.ParseCIDR(ipRange)
				if err != nil {
					return result, err
				}
				result = append(result, ipnet)
			}
		}
	}
	return result, nil
}

func nextIP(ip net.IP, inc uint) net.IP {
	i := ip.To4()
	v := uint(i[0])<<24 + uint(i[1])<<16 + uint(i[2])<<8 + uint(i[3])
	v += inc
	v3 := byte(v & 0xFF)
	v2 := byte((v >> 8) & 0xFF)
	v1 := byte((v >> 16) & 0xFF)
	v0 := byte((v >> 24) & 0xFF)
	return net.IPv4(v0, v1, v2, v3)
}
