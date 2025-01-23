package main

import (
	"log"

	"concurrency_go_course/internal/config"
	"concurrency_go_course/internal/network"
)

func main() {
	cfg, err := config.NewServerConfig()
	if err != nil {
		log.Fatal("unable to start server: unable to read cfg")
	}

	server, err := network.New(cfg)
	if err != nil {
		log.Fatal("unable to start server")
	}

	defer server.Close()

	server.Run()
}
