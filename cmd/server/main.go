package main

import (
	"context"
	"log"

	"concurrency_go_course/internal/app"
	"concurrency_go_course/internal/config"
	"concurrency_go_course/internal/network"
	"concurrency_go_course/pkg/logger"
)

var configPath = "config.yaml"

func main() {
	ctx := context.Background()

	cfg, err := config.NewConfig(configPath)
	if err != nil {
		log.Fatal("unable to start server: unable to read cfg")
	}

	logger.InitLogger(cfg.Logging.Level, cfg.Logging.Output)
	logger.Debug("init logger")

	walCfg, err := config.NewWALConfig(configPath)
	if err != nil {
		logger.Info("unable to set WAL settings, WAL is disabled")
	}

	db, err := app.Init(cfg, walCfg)
	if err != nil {
		log.Fatal("unable to init app")
	}

	if walCfg != nil && walCfg.WalConfig != nil {
		logger.Debug("starting WAL")
		go db.StartWAL(ctx)
	}

	server, err := network.NewServer(db, cfg)
	if err != nil {
		log.Fatal("unable to start server")
	}

	defer func() {
		_ = server.Close()
	}()

	server.Run()
}
