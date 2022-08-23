package cmd

import "github.com/spf13/cobra"

var cloudCmd = &cobra.Command{
	Use:   "cloud",
	Short: "Cloud sync functionality",
}

var cloudSyncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync cloud ipam environments",
	RunE:  syncCloud,
}

var cloudImportCmd = &cobra.Command{
	Use:   "import",
	Short: "Import cloud ipam environments",
	RunE:  importCloud,
}

func syncCloud(cmd *cobra.Command, args []string) error {
	return nil
}

func importCloud(cmd *cobra.Command, args []string) error {
	return nil
}

func init() {
	rootCmd.AddCommand(cloudCmd)

	cloudCmd.AddCommand(cloudSyncCmd)
	cloudCmd.AddCommand(cloudImportCmd)
}
