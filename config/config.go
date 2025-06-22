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
	viper.ReadInConfig()

	// Unmarshal to config
	viper.Unmarshal(&config)

	viper.AutomaticEnv()

	// Override with environment variables if they exist
	if dbSource := os.Getenv("DB_SOURCE"); dbSource != "" {
		config.DBSource = dbSource
	}

	if port := os.Getenv("PORT"); port != "" {
		config.HTTPServerAddress = "0.0.0.0:" + port
	}

	return config, nil
}
