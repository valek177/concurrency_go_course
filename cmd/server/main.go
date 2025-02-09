package main

import (
	"context"
	"flag"
	"log"
	"sync"

	"concurrency_go_course/internal/app"
	"concurrency_go_course/internal/config"
	"concurrency_go_course/internal/network"
	"concurrency_go_course/pkg/logger"
)

var (
	configPathMaster = "config.yaml"
	configPathSlave  = "config_slave.yaml"
)

func main() {
	// ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	// defer cancel()

	ctx := context.Background()

	configPath := flag.String("config-path", configPathMaster, "path to config file")
	flag.Parse()

	cfg, err := config.NewConfig(*configPath)
	if err != nil {
		log.Fatal("unable to start server: unable to read cfg")
	}

	logger.InitLogger(cfg.Logging.Level, cfg.Logging.Output)
	logger.Debug("init logger")

	walCfg, err := config.NewWALConfig(*configPath)
	if err != nil {
		logger.Info("unable to set WAL settings, WAL is disabled")
	}

	db, wal, repl, err := app.Init(cfg, walCfg)
	if err != nil {
		log.Fatal("unable to init app")
	}

	wg := sync.WaitGroup{}
	wg.Add(3)
	if wal != nil && walCfg != nil && walCfg.WalConfig != nil {
		go func() {
			defer wg.Done()

			logger.Debug("starting WAL")
			wal.Start(ctx)
		}()
	}

	if cfg.Replication != nil {
		go func() {
			defer wg.Done()

			err = repl.Start(ctx)
			if err != nil {
				logger.ErrorWithMsg("unable to start replication", err)
			}
		}()
	}

	server, err := network.NewServer(cfg, cfg.Network.Address)
	if err != nil {
		log.Fatal("unable to start server")
	}

	go func() {
		defer wg.Done()
		defer func() {
			_ = server.Close()
		}()

		server.Run(ctx, func(ctx context.Context, s []byte) []byte {
			response, err := db.Handle(string(s) + "\n")
			if err != nil {
				logger.ErrorWithMsg("unable to handle query:", err)
				response = err.Error()
			}
			return []byte(response)
		})
	}()

	wg.Wait()
}
