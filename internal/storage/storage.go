package storage

import (
	"fmt"

	"concurrency_go_course/internal/compute"
	"concurrency_go_course/internal/storage/wal"
	"concurrency_go_course/pkg/logger"

	"go.uber.org/zap"
)

type Engine interface {
	Get(key string) (string, bool)
	Set(key string, value string)
	Delete(key string)
}

type Storage struct {
	engine Engine
	wal    *wal.WAL
}

type WAL interface {
	Set(string, string) error
	Del(string) error
	Recover() ([]wal.Request, error)
}

func New(engine Engine, wal *wal.WAL) (*Storage, error) {
	if engine == nil {
		return nil, fmt.Errorf("unable to create storage: engine is empty")
	}

	if wal == nil {
		logger.Debug("WAL is not used")
	}

	storage := &Storage{
		engine: engine,
		wal:    wal,
	}

	if storage.wal != nil {
		requests, err := storage.wal.Recover()
		if err != nil {
			logger.ErrorWithMsg("unable to get requests from WAL", err)
		} else {
			storage.restore(requests)
		}
	}

	return storage, nil
}

func (s *Storage) Set(key, value string) error {
	if s.wal != nil {
		if err := s.wal.Set(key, value); err != nil {
			return err
		}
	}

	s.engine.Set(key, value)
	return nil
}

func (s *Storage) Get(key string) (string, bool) {
	return s.engine.Get(key)
}

func (s *Storage) Del(key string) error {
	if s.wal != nil {
		if err := s.wal.Del(key); err != nil {
			return err
		}
	}

	s.engine.Delete(key)
	return nil
}

func (s *Storage) restore(requests []wal.Request) {
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
