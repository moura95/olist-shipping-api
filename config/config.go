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
		return
	}
	err = viper.Unmarshal(&config)
	// Cloud Run
	if port := os.Getenv("PORT"); port != "" {
		config.HTTPServerAddress = "0.0.0.0:" + port
	}

	return
}
