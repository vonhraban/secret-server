package main

import (
	"github.com/vonhraban/secret-server/app/http"
	"github.com/vonhraban/secret-server/persistence"
	"github.com/vonhraban/secret-server/secret"
	"github.com/vonhraban/secret-server/core/log"
	"github.com/sirupsen/logrus"
)

func main() {
	// TODO! Move to a log factory
	log := log.NewLogrusLogger(logrus.New())

	clock := &secret.TimeClock{}
	vault := persistence.NewInMemoryVault(clock)
	
	httpService := http.New(vault, clock, log)
	
	httpService.Serve()
}
