package storage

import (
	"fmt"

	"concurrency_go_course/internal/compute"
	"concurrency_go_course/internal/replication"
	"concurrency_go_course/internal/storage/wal"
	"concurrency_go_course/pkg/logger"

	"go.uber.org/zap"
)

// Storage is interface for storage
type Storage interface {
	Set(key, value string) error
	Get(key string) (string, bool)
	Del(key string) error
	Restore(requests []wal.Request)
}

type storage struct {
	engine            Engine
	replicationStream chan []wal.Request
	wal               *wal.WAL
	isMasterRepl      bool
}

// WAL is interface for write ahead log
type WAL interface {
	Set(string, string) error
	Del(string) error
	Recover() ([]wal.Request, error)
}

// New creates new storage
func New(engine Engine, wal *wal.WAL,
	replicationType string, replStream chan []wal.Request,
) (Storage, error) {
	if engine == nil {
		return nil, fmt.Errorf("unable to create storage: engine is empty")
	}

	stor := &storage{
		engine:            engine,
		wal:               wal,
		replicationStream: replStream,
		isMasterRepl:      replicationType == replication.ReplicaTypeMaster,
	}

	if wal != nil {
		requests, err := stor.wal.Recover()
		if err != nil {
			logger.ErrorWithMsg("unable to get requests from WAL", err)
		} else {
			stor.Restore(requests)
		}
	}

	if replStream != nil {
		go func() {
			for request := range replStream {
				logger.Debug("applying request from replication stream")
				stor.Restore(request)
			}
		}()
	}

	return stor, nil
}

// Set sets new value
func (s *storage) Set(key, value string) error {
	if !s.isMasterRepl {
		return fmt.Errorf("unable to execute set command on slave")
	}

	if s.wal != nil {
		err := s.wal.Set(key, value)
		if err != nil {
			return err
		}
	}

	s.engine.Set(key, value)
	return nil
}

// Get returns value by key
func (s *storage) Get(key string) (string, bool) {
	return s.engine.Get(key)
}

// Del deletes key
func (s *storage) Del(key string) error {
	if !s.isMasterRepl {
		return fmt.Errorf("unable to execute delete command on slave")
	}

	if s.wal != nil {
		if err := s.wal.Del(key); err != nil {
			return err
		}
	}

	s.engine.Delete(key)
	return nil
}

// Restore restores WAL settings
func (s *storage) Restore(requests []wal.Request) {
	for _, request := range requests {
		switch request.Command {
		case compute.CommandSet:
			s.engine.Set(request.Args[0], request.Args[1])
			logger.Debug("Was restored", zap.String("key", request.Args[0]),
				zap.String("value", request.Args[1]))
		case compute.CommandDelete:
			s.engine.Delete(request.Args[0])
			logger.Debug("Was deleted", zap.String("key", request.Args[0]))
		}
	}
}
