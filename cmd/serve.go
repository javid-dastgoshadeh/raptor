package cmd

import (
	"github.com/spf13/cobra"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serve the authenticator server",
	Long:  `test service sum two numbers http -port=804`,
	Run: func(cmd *cobra.Command, args []string) {
		serveHTTPCmd.Run(cmd, args)
	},
}

func init() {
	rootCmd.AddCommand(serveCmd)
}
