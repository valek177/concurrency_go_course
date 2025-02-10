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
	database.Database, *wal.WAL, replication.Replication, error,
) {
	var err error

	if cfg == nil {
		return nil, nil, nil, fmt.Errorf("config is empty")
	}

	var walObj *wal.WAL

	if walCfg == nil || walCfg.WalConfig == nil {
		logger.Debug("WAL config is empty, WAL is not used")
	} else {
		walObj, err = wal.New(walCfg, cfg.Replication.ReplicaType)
		if err != nil {
			return nil, nil, nil, fmt.Errorf("unable to create new WAL: %v", err)
		}
	}

	var repl replication.Replication

	if cfg.Replication != nil {
		if cfg.Replication.ReplicaType == replication.ReplicaTypeMaster {
			replServer, err := replication.NewReplicationServer(cfg, walCfg)
			if err != nil {
				logger.ErrorWithMsg("unable to create replication master server:",
					err)
			} else {
				repl = replServer
			}

		} else if cfg.Replication.ReplicaType == replication.ReplicaTypeSlave {
			fmt.Println("replc in slave", cfg.Replication.ReplicaType)
			replClient, err := replication.NewReplicationClient(cfg, walCfg)
			fmt.Println("err ", err, replClient)
			if err != nil {
				logger.ErrorWithMsg("unable to create replication slave server:",
					err)
			} else {
				repl = replClient
			}
		}
	}

	fmt.Println("replc 11", repl)
	var replStream chan []wal.Request
	if repl != nil && !repl.IsMaster() {
		replStream = repl.ReplicationStream()
	}

	engine := storage.NewEngine()
	fmt.Println("stream ", replStream)

	storage, err := storage.New(engine, walObj, cfg, cfg.Replication.ReplicaType,
		replStream)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("unable to init storage: %v", err)
	}

	requestParser := compute.NewRequestParser()
	compute := compute.NewCompute(requestParser)

	db := database.NewDatabase(storage, compute)

	return db, walObj, repl, nil
}
