package env

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// SetConfigFile ...
func SetConfigFile(file string) {
	viper.SetConfigFile(file)
}

// SetConfigType ...
func SetConfigType(typeName string) {
	viper.SetConfigType(typeName)
}

// AddConfigPath ...
func AddConfigPath(path string) {
	viper.AddConfigPath(path)
}

// SetConfigName ...
func SetConfigName(name string) {
	viper.SetConfigName(name)
}

// Get ...
func Get(key string) interface{} {
	return viper.Get(key)
}

// GetString ...
func GetString(key string) string {
	return viper.GetString(key)
}

// GetInt ...
func GetInt(key string) int {
	return viper.GetInt(key)
}

// GetBool ...
func GetBool(key string) bool {
	return viper.GetBool(key)
}

// GetStringMap ...
func GetStringMap(key string) map[string]interface{} {
	return viper.GetStringMap(key)
}

// GetStringMapString ...
func GetStringMapString(key string) map[string]string {
	return viper.GetStringMapString(key)
}

// GetStringSlice ...
func GetStringSlice(key string) []string {
	return viper.GetStringSlice(key)
}

// IsSet ...
func IsSet(key string) bool {
	return viper.IsSet(key)
}

// GetBasePath ...
func GetBasePath() string {
	ex, err := os.Executable()
	if err != nil {
		panic(err)
	}
	exPath := filepath.Dir(ex)

	return exPath
}

// BindPFlag ...
func BindPFlag(key string, flag *pflag.Flag) error {
	return viper.BindPFlag(key, flag)
}

// Init ...
func Init(filename string) (string, error) {
	SetConfigType("json")

	if filename != "" {
		SetConfigFile(filename)
	} else {
		SetConfigName("config")
		AddConfigPath(GetBasePath())
	}

	viper.SetEnvPrefix("AS")
	viper.AutomaticEnv()

	// If a config file is found, read it in.
	//if err := viper.ReadInConfig(); err == nil {
	//	logger.Info("Using config file:" + viper.ConfigFileUsed())
	//}
	return viper.ConfigFileUsed(), viper.ReadInConfig()
}

// HTTPApiAddress  return the API host
func HTTPApiAddress() string {
	return fmt.Sprintf("%s://%s:%d",
		GetString("services.http.scheme"),
		GetString("services.http.address"),
		GetInt("services.http.port"))
}
