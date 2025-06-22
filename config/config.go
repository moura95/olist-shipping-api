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
	viper.AddConfigPath(path)
	viper.SetConfigType("env")
	viper.SetConfigFile(".env")

	viper.AutomaticEnv()

	err = viper.ReadInConfig()
	if err != nil {
		config.HTTPServerAddress = "0.0.0.0:8080"
	} else {
		err = viper.Unmarshal(&config)
		if err != nil {
			return
		}
	}

	// Cloud Run sempre define PORT, use ela se disponível
	if port := os.Getenv("PORT"); port != "" {
		config.HTTPServerAddress = "0.0.0.0:" + port
	}

	// Fallback se não tiver PORT nem HTTP_SERVER_ADDRESS
	if config.HTTPServerAddress == "" {
		config.HTTPServerAddress = "0.0.0.0:8080"
	}

	return
}
