package service

import (
	"fmt"

	"concurrency_go_course/internal/compute"
	"concurrency_go_course/pkg/logger"

	"go.uber.org/zap"
)

func (s *serv) Handle(request string) (string, error) {
	query, err := s.compute.Handle(request)
	if err != nil {
		fmt.Printf("Parsing request error: %v", err.Error())

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
		s.storage.Set(query.Args[0], query.Args[1])

		logger.Debug("Key with value was saved",
			zap.String("key", query.Args[0]), zap.String("value", query.Args[1]))

		return "", nil
	case compute.CommandDelete:
		s.storage.Delete(query.Args[0])

		logger.Debug("Key was deleted", zap.String("key", query.Args[0]))

		return "", nil
	}

	return "", nil
}
