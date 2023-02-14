package commands

import (
	"v2ray-database-api/config"
	"v2ray-database-api/log"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "show the version",
	RunE:  printVersion,
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

func printVersion(_ *cobra.Command, _ []string) error {
	log.Info("Application:         ", config.ApplicationName)
	log.Info("Version:             ", config.Version)
	return nil
}
