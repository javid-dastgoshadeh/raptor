package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
	"raptor/config"
	"raptor/logger"
)

var cfgFile string

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "authenticator",
	Short: "authenticate services with otp by sms and email code",
	Long:  `authenticate services with otp by sms and email code http -port=8001`,
	// has an action associated with it:
	Run: func(cmd *cobra.Command, args []string) {
		// Default action on serve
		serveCmd.Run(cmd, args)
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	//helpers.GeneratePrivateKey()
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}

func init() {
	cobra.OnInitialize(initConfig)
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file (default is $COMMAND_HOME/config.json)")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	fileName, err := env.Init(cfgFile)
	if err != nil {
		panic("config file note set")
	}
	logger.Info("Using config file:" + fileName)
}
