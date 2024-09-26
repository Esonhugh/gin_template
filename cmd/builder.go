package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var RootCmd = &cobra.Command{
	Use:          "app",
	Short:        "app",
	SilenceUsage: true,
	Long:         `app`,
}

func init() {
}

func Execute() {
	if err := RootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
