package cmd

import (
	"github.com/spf13/cobra"
	"raptor/config"
	"raptor/database"
	"raptor/database/migrations"
)

// migrateCmd represents the migrateCmd command
var migrateCmd = &cobra.Command{
	Use:   "migrate",
	Short: "Migrate database",
	Long:  `Migrate all changes to database`,
	Run: func(cmd *cobra.Command, args []string) {
		defaultDriver := env.GetString("database.default")

		migrations.Init(database.DBConfig{
			Engine:   defaultDriver,
			DbName:   env.GetString("database.drivers." + defaultDriver + ".db_name"),
			Host:     env.GetString("database.drivers." + defaultDriver + ".host"),
			Port:     env.GetString("database.drivers." + defaultDriver + ".port"),
			Username: env.GetString("database.drivers." + defaultDriver + ".username"),
			Password: env.GetString("database.drivers." + defaultDriver + ".password"),
			Log:      env.GetBool("database.drivers." + defaultDriver + ".log"),
			SslMode:  env.GetString("database.drivers." + defaultDriver + ".ssl_mode"),
		})

		migrations.Migrate()
	},
}

func init() {
	//rootCmd.AddCommand(migrateCmd)
	migrateCmd.Flags().String("db_address", "127.0.0.1", "Database full address")
	migrateCmd.Flags().Int("db_port", 5432, "Database instance port")
	migrateCmd.Flags().String("db_username", "postgres", "Database username")
	migrateCmd.Flags().String("db_password", "", "Database password")
	migrateCmd.Flags().String("db_dbname", "oauth-otp", "Default database name")

	defaultDriver := env.GetString("database.default")

	env.BindPFlag("database.drivers."+defaultDriver+".host", migrateCmd.Flags().Lookup("db_address"))
	env.BindPFlag("database.drivers."+defaultDriver+".port", migrateCmd.Flags().Lookup("db_port"))
	env.BindPFlag("database.drivers."+defaultDriver+".username", migrateCmd.Flags().Lookup("db_username"))
	env.BindPFlag("database.drivers."+defaultDriver+".password", migrateCmd.Flags().Lookup("db_password"))
	env.BindPFlag("database.drivers."+defaultDriver+".db_name", migrateCmd.Flags().Lookup("db_dbname"))
}
