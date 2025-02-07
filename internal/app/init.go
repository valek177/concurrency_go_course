package app

import (
	"fmt"

	"concurrency_go_course/internal/compute"
	"concurrency_go_course/internal/config"
	"concurrency_go_course/internal/database"
	"concurrency_go_course/internal/storage"
	"concurrency_go_course/internal/storage/wal"
	"concurrency_go_course/pkg/logger"
)

// Init initializes new database and wal service
func Init(cfg *config.Config, walCfg *config.WALCfg) (database.Database, *wal.WAL, error) {
	var err error

	if cfg == nil {
		return nil, nil, fmt.Errorf("config is empty")
	}

	var walObj *wal.WAL

	if walCfg == nil || walCfg.WalConfig == nil {
		logger.Debug("WAL config is empty, WAL is not used")
	} else {
		walObj, err = wal.New(walCfg)
		if err != nil {
			return nil, nil, fmt.Errorf("unable to create new WAL: %v", err)
		}
	}

	engine := storage.NewEngine()

	storage, err := storage.New(engine, walObj)
	if err != nil {
		return nil, nil, fmt.Errorf("unable to init storage: %v", err)
	}

	requestParser := compute.NewRequestParser()
	compute := compute.NewCompute(requestParser)

	db := database.NewDatabase(storage, compute)

	return db, walObj, nil
}
