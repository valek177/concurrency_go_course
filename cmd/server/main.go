package main

import (
	"log"

	"concurrency_go_course/internal/config"
	"concurrency_go_course/internal/network"
)

var configPath = "config.yaml"

func main() {
	cfg, err := config.NewServerConfig(configPath)
	if err != nil {
		log.Fatal("unable to start server: unable to read cfg")
	}

	server, err := network.NewServer(cfg)
	if err != nil {
		log.Fatal("unable to start server")
	}

	defer func() {
		_ = server.Close()
	}()

	server.Run()
}
