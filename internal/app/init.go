package app

import (
	"fmt"

	"concurrency_go_course/internal/compute"
	"concurrency_go_course/internal/config"
	"concurrency_go_course/internal/database"
	"concurrency_go_course/internal/replication"
	"concurrency_go_course/internal/storage"
	"concurrency_go_course/internal/storage/wal"
	"concurrency_go_course/pkg/logger"
)

// Init initializes new database and wal service and other objects
func Init(cfg *config.Config, walCfg *config.WALCfg) (
	database.Database, *wal.WAL, *replication.Replication, error,
) {
	var err error
	var replicaType string

	if cfg.Replication != nil && cfg.Replication.ReplicaType != "" {
		replicaType = cfg.Replication.ReplicaType
	}

	if cfg == nil {
		return nil, nil, nil, fmt.Errorf("config is empty")
	}

	var walObj *wal.WAL

	if walCfg == nil || walCfg.WalConfig == nil {
		logger.Debug("WAL config is empty, WAL is not used")
	} else {
		walObj, err = wal.New(walCfg)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("unable to create new WAL: %v", err)
		}
	}

	repl := &replication.Replication{}

	if replicaType == replication.ReplicaTypeMaster {
		replServer, err := replication.NewReplicationServer(cfg, walCfg)
		if err != nil {
			logger.ErrorWithMsg("unable to create replication master server:", err)
		} else {
			repl.Master = replServer
		}

	} else if replicaType == replication.ReplicaTypeSlave {
		replClient, err := replication.NewReplicationClient(cfg, walCfg)
		if err != nil {
			logger.ErrorWithMsg("unable to create replication slave server:", err)
		} else {
			repl.Slave = replClient
		}
	}

	var replStream chan []wal.Request
	if repl.Slave != nil {
		replStream = repl.Slave.ReplicationStream()
	}

	engine := storage.NewEngine(cfg.Engine.PartitionsNumber)

	storage, err := storage.New(engine, walObj, replicaType, replStream)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("unable to init storage: %v", err)
	}

	requestParser := compute.NewRequestParser()
	compute := compute.NewCompute(requestParser)

	db := database.NewDatabase(storage, compute)

	return db, walObj, repl, nil
}
