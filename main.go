package main

import (
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
	
	log := log.NewLogrusLogger(logLevel)

	log.Infof("Log level %s", logLevel)

	serverPort, err := cfg.GetInt("api.port")
	if err != nil {
		panic("Can not read api.port config value")
	}
	log.Infof("Server port %d", serverPort)

	apiVersion, err := cfg.GetString("api.version")
	if err != nil {
		panic("Can not read api.version config value")
	}
	log.Infof("Api version %s", apiVersion)


	clock := &secret.TimeClock{}
	vault := persistence.NewInMemoryVault(clock)
	
	httpService := http.New(vault, clock, log, serverPort, apiVersion)
	
	httpService.Serve()
}
