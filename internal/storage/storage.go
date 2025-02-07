package storage

import (
	"fmt"

	"concurrency_go_course/internal/compute"
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
	engine Engine
	wal    *wal.WAL
}

// WAL is interface for write ahead log
type WAL interface {
	Set(string, string) error
	Del(string) error
	Recover() ([]wal.Request, error)
}

// New creates new storage
func New(engine Engine, wal *wal.WAL) (Storage, error) {
	if engine == nil {
		return nil, fmt.Errorf("unable to create storage: engine is empty")
	}

	if wal == nil {
		logger.Debug("WAL is not used")
	}

	stor := &storage{
		engine: engine,
		wal:    wal,
	}

	if stor.wal != nil {
		requests, err := stor.wal.Recover()
		if err != nil {
			logger.ErrorWithMsg("unable to get requests from WAL", err)
		} else {
			stor.Restore(requests)
		}
	}

	return stor, nil
}

// Set sets new value
func (s *storage) Set(key, value string) error {
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
