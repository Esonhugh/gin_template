package cmd

import (
	"gin_template/cmd/createapp"
	"gin_template/cmd/newroute"
	"gin_template/cmd/serve"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:          "app",
	Short:        "app",
	SilenceUsage: true,
	Long:         `app`,
}

func init() {
	rootCmd.AddCommand(createapp.StartCmd)
	rootCmd.AddCommand(serve.ServerCmd)
	rootCmd.AddCommand(newroute.RouterCmd)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
