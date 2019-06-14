package main

import "github.com/vonhraban/secret-server/secret_server"

func main() {
	// router
	app := secret_server.NewApp()
	app.Serve()
}
