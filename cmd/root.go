package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/mepomuceno/cloud-ipam/ipconfig"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var client *ipconfig.IpConfigClient

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cloud-ipam",
	Short: "cloud-ipam is a tool to manage IP addresses in cloud environments",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if !viper.IsSet("storage-account-name") {
			return fmt.Errorf("storage-account-name is required")
		}
		storageAccountName := viper.GetString("storage-account-name")
		tableName := viper.GetString("table-name")
		ipconfigClient, err := ipconfig.GetClient(tableName, storageAccountName, context.Background())
		if err != nil {
			return err
		}
		client = ipconfigClient
		return nil
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.SetOutput(os.Stdout)
	rootCmd.PersistentFlags().StringP("config", "c", "", "config file (default is $HOME/.cloud-ipam.yaml)")
	rootCmd.PersistentFlags().String("log-file", ".cloud-ipam.log", "log file (default is $HOME/.cloud-ipam.log)")
	rootCmd.PersistentFlags().String("log-level", "INFO", "log level (default is info)")
	rootCmd.PersistentFlags().String("table-name", "cloudIpam", "ipam storage table name (default is cloudIpam)")
	rootCmd.PersistentFlags().StringP("storage-account-name", "s", "", "storage account name")

	err := rootCmd.MarkPersistentFlagRequired("storage-account-name")
	if err != nil {
		panic(err)
	}

	_ = viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))
	_ = viper.BindPFlag("log-file", rootCmd.PersistentFlags().Lookup("log-file"))
	_ = viper.BindPFlag("log-level", rootCmd.PersistentFlags().Lookup("log-level"))
	_ = viper.BindPFlag("table-name", rootCmd.PersistentFlags().Lookup("table-name"))
	_ = viper.BindPFlag("storage-account-name", rootCmd.PersistentFlags().Lookup("storage-account-name"))

	viper.SetConfigName(".cloud-ipam.yaml")
	viper.SetEnvPrefix("CLOUD_IPAM")
	viper.AutomaticEnv()

}
