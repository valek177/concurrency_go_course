package compute

import (
	"fmt"
	"slices"
	"strings"
)

// Parser is interface for parser
type Parser interface {
	Parse(query string) (Query, error)
}

// RequestParser is struct for request parser
type RequestParser struct{}

// NewRequestParser returns new request parser
func NewRequestParser() *RequestParser {
	return &RequestParser{}
}

// Parse parses request string
func (r *RequestParser) Parse(query string) (Query, error) {
	queryFields := strings.Fields(query)

	if len(queryFields) == 0 {
		return Query{}, fmt.Errorf("invalid query length (0)")
	}

	command := queryFields[0]

	allCommands := []string{CommandGet, CommandSet, CommandDelete}
	if !slices.Contains(allCommands, command) {
		return Query{}, fmt.Errorf("invalid command %s", command)
	}

	argsLen := len(queryFields[1:])

	switch command {
	case CommandGet:
		if len(queryFields[1:]) != 1 {
			return Query{}, fmt.Errorf("for command %s expected 1 argument, got %d",
				CommandGet, argsLen)
		}
	case CommandSet:
		if len(queryFields[1:]) != 2 {
			return Query{}, fmt.Errorf("for command %s expected 2 arguments, got %d",
				CommandSet, argsLen)
		}
	case CommandDelete:
		if len(queryFields[1:]) != 1 {
			return Query{}, fmt.Errorf("for command %s expected 1 argument, got %d",
				CommandDelete, argsLen)
		}
	}

	return NewQuery(command, queryFields[1:]), nil
}
