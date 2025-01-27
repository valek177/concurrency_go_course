package main

import (
	"log"

	"concurrency_go_course/internal/app"
	"concurrency_go_course/internal/config"
	"concurrency_go_course/internal/network"
)

var configPath = "config.yaml"

func main() {
	cfg, err := config.NewConfig(configPath)
	if err != nil {
		log.Fatal("unable to start server: unable to read cfg")
	}

	dbService, err := app.Init(cfg)
	if err != nil {
		log.Fatal("unable to init app")
	}

	server, err := network.NewServer(dbService, cfg)
	if err != nil {
		log.Fatal("unable to start server")
	}

	defer func() {
		_ = server.Close()
	}()

	server.Run()
}
