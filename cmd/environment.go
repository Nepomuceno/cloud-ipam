package cmd

import (
	"fmt"
	"net"
	"strings"

	"github.com/Azure/azure-sdk-for-go/sdk/data/aztables"
	"github.com/nepomuceno/cloud-ipam/ipconfig"
	"github.com/nepomuceno/cloud-ipam/model"
	"github.com/spf13/cobra"
)

var environmentCmd = &cobra.Command{
	Use:   "env",
	Short: "Manage cloud ipam environments",
}

var listEnvironmentsCmd = &cobra.Command{
	Use:   "list",
	Short: "List environments",
	RunE:  listEnvironments,
}

var addEnvironmentCmd = &cobra.Command{
	Use:   "add",
	Short: "Add or update environment",
	RunE:  addOrUpdateEnvironment,
}

var deleteEnvironmentCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete environment",
	RunE:  deleteEnvironment,
}

func listEnvironments(cmd *cobra.Command, args []string) error {
	environment, err := client.ListEnvironments()
	if err != nil {
		return err
	}
	for _, env := range environment {
		cmd.Printf("%s, %s, %s\n", env.RowKey, env.Name, env.IPRanges)
	}
	return nil
}

func addOrUpdateEnvironment(cmd *cobra.Command, args []string) error {
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
	if id == "" {
		id = strings.ReplaceAll(strings.ToLower(name), " ", "-")
	}

	if id == ipconfig.ENVIRONMENT_PARTITION_KEY {
		return fmt.Errorf("%s is a reserved id", ipconfig.ENVIRONMENT_PARTITION_KEY)
	}

	environment := model.EnvironmentDefinition{
		Entity: aztables.Entity{
			PartitionKey: ipconfig.ENVIRONMENT_PARTITION_KEY,
			RowKey:       id,
		},
		Name:     name,
		IPRanges: ipRange.String(),
	}

	err = client.AddOrUpdateEnvironment(environment)
	if err != nil {
		return err
	}
	cmd.Printf("Environment %s added or updated\n", id)
	return nil
}

func deleteEnvironment(cmd *cobra.Command, args []string) error {
	id, err := cmd.Flags().GetString("id")
	if err != nil {
		return err
	}
	err = client.DeleteEnvironment(id)
	if err != nil {
		return err
	}
	cmd.Printf("Environment %s deleted\n", id)
	return nil
}

func init() {
	rootCmd.AddCommand(environmentCmd)
	_, defaultIPCIDR, err := net.ParseCIDR("10.0.0.0/8")
	if err != nil {
		panic(err)
	}

	addEnvironmentCmd.Flags().IPNet("ip-range", *defaultIPCIDR, "IP range to be used by the environment (default 10.0.0.0/8")
	addEnvironmentCmd.Flags().StringP("name", "n", "", "Name of the environment")
	addEnvironmentCmd.Flags().StringP("id", "i", "", "ID of the environment (be default the same as name lowercase and dashes for spaces)")
	err = addEnvironmentCmd.MarkFlagRequired("name")
	if err != nil {
		panic(err)
	}

	deleteEnvironmentCmd.Flags().StringP("id", "i", "", "ID of the environment")
	err = deleteEnvironmentCmd.MarkFlagRequired("id")
	if err != nil {
		panic(err)
	}

	environmentCmd.AddCommand(listEnvironmentsCmd)
	environmentCmd.AddCommand(addEnvironmentCmd)
	environmentCmd.AddCommand(deleteEnvironmentCmd)

}
