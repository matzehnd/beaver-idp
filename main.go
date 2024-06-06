package main

import (
	"beaver/idp/api"
	"os"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		panic("no port defined")
	}
	api.Serve(port)
}
