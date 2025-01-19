package compute

import (
	"fmt"
	"testing"

	"concurrency_go_course/pkg/logger"

	"github.com/stretchr/testify/assert"
)

func TestHandleNegCompute(t *testing.T) {
	t.Parallel()

	logger.MockLogger()

	parser := NewRequestParser()
	compute := NewCompute(parser)

	negTests := map[string]struct {
		in  string
		res Query
		err error
	}{
		"empty request": {
			in:  "",
			res: Query{},
			err: fmt.Errorf("invalid query length (0)"),
		},
		"cmd: invalid command": {
			in:  "cmd unknown",
			res: Query{},
			err: fmt.Errorf("invalid command cmd"),
		},
		"get: invalid command name case": {
			in:  "get unknown",
			res: Query{},
			err: fmt.Errorf("invalid command get"),
		},
	}

	for name, test := range negTests {
		t.Run(name, func(t *testing.T) {
			res, err := compute.Handle(test.in)
			assert.Equal(t, err, test.err)
			assert.Equal(t, res, test.res)
		})
	}
}

func TestHandlePosCompute(t *testing.T) {
	t.Parallel()

	logger.MockLogger()

	parser := NewRequestParser()
	compute := NewCompute(parser)

	posTests := map[string]struct {
		in  string
		res Query
		err error
	}{
		"GET: existing value": {
			in:  "GET key1",
			res: Query{Command: "GET", Args: []string{"key1"}},
			err: nil,
		},
		"SET: new value": {
			in:  "SET key2 value2",
			res: Query{Command: "SET", Args: []string{"key2", "value2"}},
			err: nil,
		},
		"DEL: existing value": {
			in:  "DEL key1",
			res: Query{Command: "DEL", Args: []string{"key1"}},
			err: nil,
		},
	}

	for name, test := range posTests {
		t.Run(name, func(t *testing.T) {
			res, _ := compute.Handle(test.in)
			assert.Equal(t, res, test.res)
		})
	}
}
