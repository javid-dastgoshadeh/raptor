package cmd

import (
	"fmt"
	"raptor/config"

	"github.com/spf13/cobra"
)

// versionCmd represents the versionCmd command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "App version",
	Long:  `Print app version in console screen`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("App version: " + env.GetString("version"))
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
