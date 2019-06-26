package main

import (
	"fmt"

	"github.com/vonhraban/secret-server/app/http"
	"github.com/vonhraban/secret-server/persistence"
	"github.com/vonhraban/secret-server/secret"
	"github.com/vonhraban/secret-server/core/log"
	"github.com/vonhraban/secret-server/app/config"
)

func main() {
	cfg, err := config.NewViperConfig("config", ".", "yaml")
	if err != nil {
		panic(err)
	}

	
	logLevel, err := cfg.GetString("log.level")
	if err != nil {
		panic("Can not read log.level config value")
	}
	
	logger := log.NewLogrusLogger(logLevel)

	logger.Infof("Log level %s", logLevel)

	serverPort := requireConfigInt(logger, cfg, false, "api.port")
	apiVersion := requireConfigString(logger, cfg, false, "api.version")
	dbHost := requireConfigString(logger, cfg, false, "db.mongo.host")
	dbPort := requireConfigInt(logger, cfg, false, "db.mongo.port")
	databaseName := requireConfigString(logger, cfg, false, "db.mongo.database")
	dbUsername := requireConfigString(logger, cfg, false, "db.mongo.username")
	dbPassword := requireConfigString(logger, cfg, true, "db.mongo.password")

	clock := &secret.TimeClock{}
	vault := persistence.NewMongoVault(clock, dbHost, dbPort, databaseName, dbUsername, dbPassword)
	
	httpService := http.New(vault, clock, logger, serverPort, apiVersion)
	
	httpService.Serve()
}

func requireConfigString(logger log.Logger, cfg config.Config, hideOutput bool, key string) string {
	val, err := cfg.GetString(key)
	if err != nil {
		panic(fmt.Sprintf("Can not read %s config value", key))
	}

	if !hideOutput {
		logger.Infof("Config: %s = %s", key, val)
	}

	return val
}

func requireConfigInt(logger log.Logger, cfg config.Config, hideOutput bool, key string) int {
	val, err := cfg.GetInt(key)
	if err != nil {
		panic(fmt.Sprintf("Can not read %s config value", key))
	}

	if !hideOutput {
		logger.Infof("Config: %s = %d", key, val)
	}

	return val
}