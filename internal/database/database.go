package database

import (
	"fmt"

	"concurrency_go_course/internal/compute"
	"concurrency_go_course/internal/storage"
	"concurrency_go_course/pkg/logger"

	"go.uber.org/zap"
)

var resultOK = "OK"

// Database is interface for database
type Database interface {
	Handle(request string) (string, error)
}

type database struct {
	storage storage.Storage
	compute compute.Compute
}

// NewDatabase returns new database
func NewDatabase(
	storage storage.Storage,
	compute compute.Compute,
) Database {
	return &database{
		storage: storage,
		compute: compute,
	}
}

// Handle handles request
func (s *database) Handle(request string) (string, error) {
	query, err := s.compute.Handle(request)
	if err != nil {
		fmt.Printf("Parsing request error: %v\n", err.Error())

		return "", err
	}

	switch query.Command {
	case compute.CommandGet:
		v, ok := s.storage.Get(query.Args[0])
		if !ok {
			logger.Error("get error: value not found")
			fmt.Printf("Value by key %s not found", query.Args[0])

			return "", fmt.Errorf("value not found")
		}

		logger.Debug("Value for key was found",
			zap.String("key", query.Args[0]), zap.String("value", v))

		return v, nil
	case compute.CommandSet:
		_ = s.storage.Set(query.Args[0], query.Args[1])

		logger.Debug("Key with value was saved",
			zap.String("key", query.Args[0]), zap.String("value", query.Args[1]))

		return resultOK, nil
	case compute.CommandDelete:
		_ = s.storage.Del(query.Args[0])

		logger.Debug("Key was deleted", zap.String("key", query.Args[0]))

		return resultOK, nil
	}

	return "", fmt.Errorf("unknown command: %s", query.Command)
}
