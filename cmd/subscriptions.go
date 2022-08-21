package cmd

import (
	"net"

	"github.com/Azure/azure-sdk-for-go/sdk/data/aztables"
	"github.com/mepomuceno/cloud-ipam/model"
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
	for _, env := range subscriptions {
		cmd.Printf("%s, %s, %s\n", env.RowKey, env.Name, env.IpRanges)
	}
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
	name, err := cmd.Flags().GetString("name")
	if err != nil {
		return err
	}
	id, err := cmd.Flags().GetString("id")
	if err != nil {
		return err
	}

	subscription := model.SubscriptionDefinition{
		Entity: aztables.Entity{
			PartitionKey: environmentID,
			RowKey:       id,
		},
		Name:     name,
		IpRanges: ipRange.String(),
	}

	err = client.AddOrUpdateSubscription(subscription)
	if err != nil {
		return err
	}
	cmd.Printf("Subcription %s added or updated\n", id)
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
	cmd.Printf("Subscription %s deleted\n", id)
	return nil
}

func init() {
	rootCmd.AddCommand(subscriptionCmd)
	_, defaultIPCIDR, err := net.ParseCIDR("0.0.0.0/0")
	if err != nil {
		panic(err)
	}
	subscriptionCmd.PersistentFlags().String("environment-id", "", "Environment ID for that subscription")
	err = subscriptionCmd.MarkPersistentFlagRequired("environment-id")
	if err != nil {
		panic(err)
	}

	addSubCmd.Flags().IPNet("ip-range", *defaultIPCIDR, "IP range to be used by the subscription")
	addSubCmd.Flags().StringP("name", "n", "", "Name of the subscription")
	addSubCmd.Flags().StringP("id", "i", "", "ID of the subscription")
	err = addSubCmd.MarkFlagRequired("id")
	if err != nil {
		panic(err)
	}
	err = addSubCmd.MarkFlagRequired("ip-range")
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
