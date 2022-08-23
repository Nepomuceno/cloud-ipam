package cmd

import (
	"context"
	"os"
	"strings"

	"github.com/nepomuceno/cloud-ipam/ipconfig"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var client *ipconfig.IpConfigClient

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "cloud-ipam",
	Short: "cloud-ipam is a tool to manage IP addresses in cloud environments",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
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

func GetRootCmd() *cobra.Command {
	return rootCmd
}

func init() {
	rootCmd.SetOutput(os.Stdout)
	rootCmd.PersistentFlags().StringP("config", "c", "", "config file")
	rootCmd.PersistentFlags().String("log-file", ".cloud-ipam.log", "log file")
	rootCmd.PersistentFlags().String("log-level", "INFO", "log level")
	rootCmd.PersistentFlags().String("table-name", "cloudIpam", "ipam storage table name")
	rootCmd.PersistentFlags().StringP("storage-account-name", "s", "", "storage account name")

	_ = viper.BindPFlag("config", rootCmd.PersistentFlags().Lookup("config"))
	_ = viper.BindPFlag("log-file", rootCmd.PersistentFlags().Lookup("log-file"))
	_ = viper.BindPFlag("log-level", rootCmd.PersistentFlags().Lookup("log-level"))
	_ = viper.BindPFlag("table-name", rootCmd.PersistentFlags().Lookup("table-name"))
	_ = viper.BindPFlag("storage-account-name", rootCmd.PersistentFlags().Lookup("storage-account-name"))

	viper.SetConfigName(".cloud-ipam.yaml")
	viper.SetEnvPrefix("CLOUD_IPAM")
	replacer := strings.NewReplacer("-", "_")
	viper.SetEnvKeyReplacer(replacer)
	viper.AutomaticEnv()
}
