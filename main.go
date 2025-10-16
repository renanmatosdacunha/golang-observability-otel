package main

import (
	"log"

	"github.com/renanmatosdacunha/golang-observability-otel.git/api"
)

const (
	serverAddress = "0.0.0.0:8080"
)

func main() {
	server := api.NewServer()

	err := server.Start(serverAddress)
	if err != nil {
		log.Fatal("Cannot start server", err)
	}
}
