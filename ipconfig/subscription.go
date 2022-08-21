package ipconfig

import (
	"encoding/json"
	"fmt"
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
	_, err = client.UpdateEntity(client.ctx, entity, nil)
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

func (client *IpConfigClient) NextAvailableIpRange(environmentID string, subnetRangeSize int) (string, error) {
	netEnvironment, err := client.GetEnvironment(environmentID)
	if err != nil {
		return "", err
	}
	ipranges := strings.Split(netEnvironment.IpRanges, ",")
	usedIpRanges, err := client.GetUsedRangesForEnvironment(environmentID)
	if err != nil {
		return "", err
	}
	for _, iprange := range ipranges {
		startIp, _, err := net.ParseCIDR(iprange)
		if err != nil {
			return "", err
		}
		ipmask := net.CIDRMask(subnetRangeSize, 32)
		ipnet := net.IPNet{IP: startIp, Mask: ipmask}

		for _, ok := usedIpRanges[ipnet.String()]; !ok; {
			newIp := nextIP(startIp, uint(32-subnetRangeSize))
			ipnet.IP = newIp
			fmt.Printf("Trying %s\n", ipnet.String())
		}

	}
	return "", nil
}

func (client *IpConfigClient) GetUsedRangesForEnvironment(environmentID string) (map[string]interface{}, error) {
	result := make(map[string]interface{})
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
			ipRanges := strings.Split(definition.IpRanges, ",")
			for _, ipRange := range ipRanges {
				result[ipRange] = true
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
