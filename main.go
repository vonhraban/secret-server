package main

import "github.com/vonhraban/secret-server/app/http"

func main() {
	// router
	httpService := http.New()
	httpService.Serve()
}
