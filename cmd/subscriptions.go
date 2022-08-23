package cmd

import (
	"encoding/json"
	"net"

	"github.com/Azure/azure-sdk-for-go/sdk/data/aztables"
	"github.com/nepomuceno/cloud-ipam/model"
	"github.com/spf13/cobra"
)

var subscriptionCmd = &cobra.Command{
	Use:   "sub",
	Short: "Manage cloud ipam subscriptions",
}

var listSubCmd = &cobra.Command{
	Use:   "list",
	Short: "List environments",
	RunE:  listSubscriptions,
}

var addSubCmd = &cobra.Command{
	Use:   "add",
	Short: "Add or update subscription",
	RunE:  addOrUpdateSubscription,
}

var deleteSubCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete subscription",
	RunE:  deleteSubscription,
}

func listSubscriptions(cmd *cobra.Command, args []string) error {
	environmentID, err := cmd.Flags().GetString("environment-id")
	if err != nil {
		return err
	}
	subscriptions, err := client.ListSubscriptions(environmentID)
	if err != nil {
		return err
	}
	content, err := json.MarshalIndent(subscriptions, "", "  ")
	if err != nil {
		return err
	}
	cmd.Printf("%s\n", content)
	return nil
}

func addOrUpdateSubscription(cmd *cobra.Command, args []string) error {

	environmentID, err := cmd.Flags().GetString("environment-id")
	if err != nil {
		return err
	}

	ipRange, err := cmd.Flags().GetIPNet("ip-range")
	if err != nil {
		return err
	}
	if ipRange.IP.To4() == nil || ipRange.IP.Equal(net.IPv4zero) {
		ipMask, err := cmd.Flags().GetIPv4Mask("ip-mask")
		if err != nil {
			return err
		}
		leading, other := ipMask.Size()
		_ = other
		nextIpRange, err := client.GetNextAvailableIpRange(environmentID, leading)
		if err != nil {
			return err
		}
		ipRange = *nextIpRange
	}
	id, err := cmd.Flags().GetString("id")
	if err != nil {
		return err
	}
	name, err := cmd.Flags().GetString("name")
	if err != nil {
		return err
	}
	if name == "" {
		name = id
	}

	subscription := model.SubscriptionDefinition{
		Entity: aztables.Entity{
			PartitionKey: environmentID,
			RowKey:       id,
		},
		Name:     name,
		IPRanges: ipRange.String(),
	}

	err = client.AddOrUpdateSubscription(subscription)
	if err != nil {
		return err
	}
	content, err := json.MarshalIndent(subscription, "", "  ")
	if err != nil {
		return err
	}
	cmd.Printf("%s\n", content)
	return nil
}

func deleteSubscription(cmd *cobra.Command, args []string) error {
	environmentID, err := cmd.Flags().GetString("environment-id")
	if err != nil {
		return err
	}
	id, err := cmd.Flags().GetString("id")
	if err != nil {
		return err
	}
	err = client.DeleteSubscription(environmentID, id)
	if err != nil {
		return err
	}
	cmd.Printf("{ id: '%s' }\n", id)
	return nil
}

func init() {
	rootCmd.AddCommand(subscriptionCmd)
	_, defaultIPCIDR, err := net.ParseCIDR("0.0.0.0/0")
	if err != nil {
		panic(err)
	}
	subscriptionCmd.PersistentFlags().StringP("environment-id", "e", "", "Environment ID for that subscription")
	err = subscriptionCmd.MarkPersistentFlagRequired("environment-id")
	if err != nil {
		panic(err)
	}

	addSubCmd.Flags().IPNet("ip-range", *defaultIPCIDR, "IP range to be used by the subscription")
	addSubCmd.Flags().StringP("name", "n", "", "Name of the subscription")
	addSubCmd.Flags().StringP("id", "i", "", "ID of the subscription")
	defaultMask := net.IPv4Mask(0xff, 0xff, 0xff, 0x00) // 255.255.255.0
	addSubCmd.Flags().IPMaskP("ip-mask", "m", defaultMask, "IP mask to be used by the subscription")
	addSubCmd.MarkFlagsMutuallyExclusive("ip-range", "ip-mask")

	err = addSubCmd.MarkFlagRequired("id")
	if err != nil {
		panic(err)
	}

	deleteSubCmd.Flags().StringP("id", "i", "", "ID of the subscription")
	err = deleteSubCmd.MarkFlagRequired("id")
	if err != nil {
		panic(err)
	}

	subscriptionCmd.AddCommand(listSubCmd)
	subscriptionCmd.AddCommand(addSubCmd)
	subscriptionCmd.AddCommand(deleteSubCmd)

}
