package app

import (
	"context"
	"fmt"

	"concurrency_go_course/internal/compute"
	"concurrency_go_course/internal/config"
	"concurrency_go_course/internal/service"
	"concurrency_go_course/internal/storage"
	"concurrency_go_course/internal/storage/wal"
	"concurrency_go_course/pkg/logger"
)

// Init initializes new database service
func Init(ctx context.Context, cfg *config.Config, walCfg *config.WALCfg) (service.Service, error) {
	var err error

	if cfg == nil {
		return nil, fmt.Errorf("config is empty")
	}

	walObj := &wal.WAL{}

	if walCfg == nil {
		logger.Debug("WAL config is empty, WAL is not used")
	} else {
		walObj, err = wal.New(walCfg)
		if err != nil {
			return nil, fmt.Errorf("unable to create new WAL: %v", err)
		}
	}

	engine := storage.NewEngine()

	storage, err := storage.New(engine, walObj)
	if err != nil {
		return nil, fmt.Errorf("unable to init storage: %v", err)
	}

	requestParser := compute.NewRequestParser()
	compute := compute.NewCompute(requestParser)

	service := service.NewService(storage, compute)

	if walObj != nil {
		go walObj.Start(ctx)
	}

	return service, nil
}
