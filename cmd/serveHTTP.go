package cmd

import (
	env "raptor/config"
	"raptor/logger"
	httpServer "raptor/server/http"
	"strconv"

	"github.com/spf13/cobra"
	"raptor/cache"
)

// serveHTTPCmd represents the serve command
var serveHTTPCmd = &cobra.Command{
	Use:   "http",
	Short: "Serve our json http server",
	Long:  `test service sum two numbers http -port=804`,
	Run: func(cmd *cobra.Command, args []string) {
		//defaultDriver := env.GetString("database.default")
		//
		//_, err := database.New(database.DBConfig{
		//	Engine:   defaultDriver,
		//	DbName:   env.GetString("database.drivers." + defaultDriver + ".db_name"),
		//	Host:     env.GetString("database.drivers." + defaultDriver + ".host"),
		//	Port:     env.GetString("database.drivers." + defaultDriver + ".port"),
		//	Username: env.GetString("database.drivers." + defaultDriver + ".username"),
		//	Password: env.GetString("database.drivers." + defaultDriver + ".password"),
		//	Log:      env.GetBool("database.drivers." + defaultDriver + ".log"),
		//	SslMode:  env.GetString("database.drivers." + defaultDriver + ".ssl_mode"),
		//})
		//if err != nil {
		//
		//	log.Fatalln("Error when connecting to database")
		//	return
		//	log.Println("Error when connecting to database")
		//
		//}
		//log.Println("Connected to database successfully")

		expireStr := env.GetString("cache.redis.expire")
		expire, _ := strconv.Atoi(expireStr)
		cache, err := cache.New(cache.Config{
			Host:   env.GetString("cache.redis.host"),
			Port:   env.GetString("cache.redis.port"),
			Expire: expire,
		})
		logger.Info("Successfully connected to redis redis in: " + env.GetString("cache.redis.host") + ":" + env.GetString("cache.redis.port"))
		if err != nil {
			panic("Error when connecting to redis")
		}
		defer cache.Conn.Close()

		httpServer.Serve(&httpServer.ServerConfig{
			Address: env.GetString("services.http.address"),
			Port:    env.GetInt("services.http.port"),
			Debug:   env.GetBool("debug"),
		})
	},
}

func init() {
	serveCmd.AddCommand(serveHTTPCmd)

	serveHTTPCmd.Flags().IntP("port", "p", 8001, "Port number for serving")
	serveHTTPCmd.Flags().StringP("address", "a", "0.0.0.0", "Host address(domain) for serving")
	serveHTTPCmd.Flags().BoolP("debug", "d", true, "Debug mode")

	env.BindPFlag("services.http.port", serveHTTPCmd.Flags().Lookup("port"))
	env.BindPFlag("services.http.address", serveHTTPCmd.Flags().Lookup("address"))
	env.BindPFlag("debug", serveHTTPCmd.Flags().Lookup("debug"))

}
