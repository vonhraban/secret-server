package main

import (
	"github.com/vonhraban/secret-server/app/http"
	"github.com/vonhraban/secret-server/persistence"
	"github.com/vonhraban/secret-server/secret"
)

func main() {
	clock := &secret.TimeClock{}
	vault := persistence.NewInMemoryVault(clock)
	httpService := http.New(vault, clock)
	httpService.Serve()
}
