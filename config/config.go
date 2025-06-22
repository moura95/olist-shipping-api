package config

import (
	"os"

	"github.com/spf13/viper"
)

type Config struct {
	DBSource          string `mapstructure:"DB_SOURCE"`
	HTTPServerAddress string `mapstructure:"HTTP_SERVER_ADDRESS"`
}

func LoadConfig(path string) (config Config, err error) {
	// Set defaults
	config.HTTPServerAddress = "0.0.0.0:8080"

	viper.AddConfigPath(path)
	viper.SetConfigType("env")
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	// Try to read config file, but don't fail if it doesn't exist
	viper.ReadInConfig()
	viper.Unmarshal(&config)

	// Cloud Run PORT variable takes precedence
	if port := os.Getenv("PORT"); port != "" {
		config.HTTPServerAddress = "0.0.0.0:" + port
	}

	return config, nil
}
