package main

import (
	"log"
	"os"

	"github/moura95/olist-shipping-api/config"
	"github/moura95/olist-shipping-api/db"
	server "github/moura95/olist-shipping-api/internal"
	"github/moura95/olist-shipping-api/internal/repository"
	"go.uber.org/zap"
)

func main() {
	loadConfig, err := config.LoadConfig(".")
	if err != nil {
		log.Printf("Failed to load config: %v", err)
		os.Exit(1)
	}

	conn, err := db.ConnectPostgres(loadConfig.DBSource)
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
		os.Exit(1)
	}
	defer conn.Close()

	log.Print("connection is repository establish")

	store := repository.New(conn.DB())

	logger, err := zap.NewProduction()
	if err != nil {
		log.Printf("Failed to create logger: %v", err)
		os.Exit(1)
	}
	defer logger.Sync()
	sugar := logger.Sugar()

	sugar.Infof("Starting server on %s", loadConfig.HTTPServerAddress)

	server.RunGinServer(loadConfig, store, sugar)
}
