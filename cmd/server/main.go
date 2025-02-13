package main

import (
	"context"
	"flag"
	"log"
	"os/signal"
	"sync"
	"syscall"

	"concurrency_go_course/internal/app"
	"concurrency_go_course/internal/config"
	"concurrency_go_course/internal/network"
	"concurrency_go_course/internal/replication"
	"concurrency_go_course/pkg/logger"
)

var configPathMaster = "config.yaml"

func main() {
	configPath := flag.String("config-path", configPathMaster, "path to config file")
	flag.Parse()

	cfg, err := config.NewConfig(*configPath)
	if err != nil {
		log.Fatal("unable to start server: unable to read cfg")
	}

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	defer cancel()

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
	if wal != nil && (cfg.Replication == nil ||
		cfg.Replication != nil &&
			cfg.Replication.ReplicaType == replication.ReplicaTypeMaster) {
		wg.Add(1)
		go func() {
			defer func() {
				wg.Done()
			}()

			logger.Debug("starting WAL")
			wal.Start(ctx)
		}()
	}

	if cfg.Replication != nil {
		logger.Debug("starting replication")
		if repl.Master != nil {
			wg.Add(1)
			go func() {
				defer wg.Done()

				repl.Master.Start(ctx)
			}()
		} else if repl.Slave != nil {
			wg.Add(1)
			go func() {
				defer wg.Done()

				repl.Slave.Start(ctx)
			}()
		}
	}

	server, err := network.NewServer(cfg, cfg.Network.Address)
	if err != nil {
		log.Fatal("unable to start server")
	}

	server.Run(ctx, func(_ context.Context, s []byte) []byte {
		response, err := db.Handle(string(s) + "\n")
		if err != nil {
			logger.ErrorWithMsg("unable to handle query:", err)
			response = err.Error()
		}
		return []byte(response)
	})

	wg.Wait()
}
