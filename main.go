package main

import (
	"github.com/joho/godotenv"
	"github.com/vonhraban/secret-server/app/config"
	"github.com/vonhraban/secret-server/app/http"
	"github.com/vonhraban/secret-server/core/log"
	"github.com/vonhraban/secret-server/persistence"
	"github.com/vonhraban/secret-server/secret"
)

func init() {
	// we want to ignore errors here , env file is not required
	_ = godotenv.Load()
}

func main() {
	cfg, err := config.New()
	if err != nil {
		panic(err)
	}

	logger := log.NewLogrusLogger(cfg.LogLevel)

	logger.Infof("Config %s", cfg.Printable())

	clock := &secret.TimeClock{}
	vault := persistence.NewMongoVault(clock, cfg.MongoHost, cfg.MongoPort, cfg.MongoDatabase, cfg.MongoUsername, cfg.MongoPassword)

	httpService := http.New(vault, clock, logger, cfg.ServerPort, cfg.ApiVersion)

	httpService.Serve()
}
