package compute

import (
	"fmt"

	"go.uber.org/zap"

	"concurrency_go_course/internal/storage"
	"concurrency_go_course/pkg/logger"
)

// Compute is struct for compute object
type Compute struct {
	storage       storage.Storage
	requestParser Parser
}

// NewCompute returns new compute object
func NewCompute(
	storage storage.Storage,
	requestParser Parser,
) *Compute {
	return &Compute{
		storage:       storage,
		requestParser: requestParser,
	}
}

// Handle handles requests
func (c *Compute) Handle(request string) (string, error) {
	query, err := c.requestParser.Parse(request)
	if err != nil {
		logger.Error("parsing request error", zap.Error(err))
		fmt.Printf("Parsing request error: %v", err.Error())

		return "", err
	}

	switch query.Command {
	case CommandGet:
		v, ok := c.storage.Get(query.Args[0])
		if !ok {
			logger.Error("get error: value not found")
			fmt.Printf("Value by key %s not found", query.Args[0])

			return "", fmt.Errorf("value not found")
		}

		logger.Debug("Value for key was found",
			zap.String("key", query.Args[0]), zap.String("value", v))

		return v, nil
	case CommandSet:
		c.storage.Set(query.Args[0], query.Args[1])

		logger.Debug("Key with value was saved",
			zap.String("key", query.Args[0]), zap.String("value", query.Args[1]))

		return "", nil
	case CommandDelete:
		c.storage.Delete(query.Args[0])

		logger.Debug("Key was deleted", zap.String("key", query.Args[0]))

		return "", nil
	}

	return "", nil
}
