package cmd

import (
	"github.com/spf13/cobra"
	"raptor/config"
	"raptor/database"
	"raptor/database/seeds"
)

// seedCmd represents the seed command
var seedCmd = &cobra.Command{
	Use:   "seed",
	Short: "Seed data",
	Long:  `Seed all data in database`,
	Run: func(cmd *cobra.Command, args []string) {
		defaultDriver := env.GetString("database.default")
		seeds.Init(database.DBConfig{
			Engine:   defaultDriver,
			DbName:   env.GetString("database.drivers." + defaultDriver + ".db_name"),
			Host:     env.GetString("database.drivers." + defaultDriver + ".host"),
			Port:     env.GetString("database.drivers." + defaultDriver + ".port"),
			Username: env.GetString("database.drivers." + defaultDriver + ".username"),
			Password: env.GetString("database.drivers." + defaultDriver + ".password"),
			Log:      env.GetBool("database.drivers." + defaultDriver + ".log"),
			SslMode:  env.GetString("database.drivers." + defaultDriver + ".ssl_mode"),
		})
		seeds.Seed()
	},
}

func init() {
	//rootCmd.AddCommand(seedCmd)

	seedCmd.Flags().String("db_address", "127.0.0.1", "Database full address")
	seedCmd.Flags().Int("db_port", 5432, "Database instance port")
	seedCmd.Flags().String("db_username", "postgres", "Database username")
	seedCmd.Flags().String("db_password", "", "Database password")
	seedCmd.Flags().String("db_dbname", "oauth_otp", "Default database name")

	defaultDriver := env.GetString("database.default")

	env.BindPFlag("database.drivers."+defaultDriver+".host", seedCmd.Flags().Lookup("db_address"))
	env.BindPFlag("database.drivers."+defaultDriver+".port", seedCmd.Flags().Lookup("db_port"))
	env.BindPFlag("database.drivers."+defaultDriver+".username", seedCmd.Flags().Lookup("db_username"))
	env.BindPFlag("database.drivers."+defaultDriver+".password", seedCmd.Flags().Lookup("db_password"))
	env.BindPFlag("database.drivers."+defaultDriver+".db_name", seedCmd.Flags().Lookup("db_dbname"))
}
