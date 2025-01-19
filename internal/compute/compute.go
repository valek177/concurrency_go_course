package compute

import (
	"go.uber.org/zap"

	"concurrency_go_course/pkg/logger"
)

const (
	// CommandGet is a get command
	CommandGet = "GET"
	// CommandSet is a set command
	CommandSet = "SET"
	// CommandDelete is a delete command
	CommandDelete = "DEL"
)

// Compute is interface for compute object
type Compute interface {
	Handle(request string) (Query, error)
}

// Comp is struct for compute object
type Comp struct {
	requestParser Parser
}

// NewCompute returns new compute object
func NewCompute(
	requestParser Parser,
) *Comp {
	return &Comp{
		requestParser: requestParser,
	}
}

// Handle handles requests
func (c *Comp) Handle(request string) (Query, error) {
	query, err := c.requestParser.Parse(request)
	if err != nil {
		logger.Error("parsing request error", zap.Error(err))

		return Query{}, err
	}

	return query, nil
}
