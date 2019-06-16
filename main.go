package main

import (
	"github.com/vonhraban/secret-server/app/http"
	"github.com/vonhraban/secret-server/persistence"
	"github.com/vonhraban/secret-server/secret"
)

func main() {
	vault := persistence.NewInMemoryVault()
	clock := &secret.TimeClock{}
	httpService := http.New(vault, clock)
	httpService.Serve()
}
