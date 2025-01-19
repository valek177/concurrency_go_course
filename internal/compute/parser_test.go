package compute

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseRequestParserNeg(t *testing.T) {
	t.Parallel()

	parser := NewRequestParser()

	negTests := map[string]struct {
		in    string
		query Query
		err   error
	}{
		"empty request": {
			in:    "",
			query: Query{},
			err:   fmt.Errorf("invalid query length (0)"),
		},
		"incorrect command": {
			in:    "somecmd",
			query: Query{},
			err:   fmt.Errorf("invalid command somecmd"),
		},
		"incorrect command with lower case get": {
			in:    "get",
			query: Query{},
			err:   fmt.Errorf("invalid command get"),
		},
		"incorrect command with lower case set": {
			in:    "set",
			query: Query{},
			err:   fmt.Errorf("invalid command set"),
		},
		"incorrect command with lower case del": {
			in:    "del",
			query: Query{},
			err:   fmt.Errorf("invalid command del"),
		},
		"GET: without args": {
			in:    "GET",
			query: Query{},
			err:   fmt.Errorf("for command GET expected 1 argument, got 0"),
		},
		"GET: with 2 args": {
			in:    "GET key value",
			query: Query{},
			err:   fmt.Errorf("for command GET expected 1 argument, got 2"),
		},
		"SET: without args": {
			in:    "SET",
			query: Query{},
			err:   fmt.Errorf("for command SET expected 2 arguments, got 0"),
		},
		"SET: with 1 args": {
			in:    "SET key",
			query: Query{},
			err:   fmt.Errorf("for command SET expected 2 arguments, got 1"),
		},
		"SET: with 3 args": {
			in:    "SET key key key",
			query: Query{},
			err:   fmt.Errorf("for command SET expected 2 arguments, got 3"),
		},
		"DEL: without args": {
			in:    "DEL",
			query: Query{},
			err:   fmt.Errorf("for command DEL expected 1 argument, got 0"),
		},
		"DEL: with 2 args": {
			in:    "DEL key value",
			query: Query{},
			err:   fmt.Errorf("for command DEL expected 1 argument, got 2"),
		},
	}

	for name, test := range negTests {
		t.Run(name, func(t *testing.T) {
			_, err := parser.Parse(test.in)
			assert.Equal(t, err, test.err)
		})
	}
}

func TestParseRequestParserPos(t *testing.T) {
	t.Parallel()

	parser := NewRequestParser()

	posTests := map[string]struct {
		in    string
		query Query
	}{
		"correct GET test": {
			in:    "GET key",
			query: Query{Command: "GET", Args: []string{"key"}},
		},
		"correct SET test": {
			in:    "SET key value",
			query: Query{Command: "SET", Args: []string{"key", "value"}},
		},
		"correct DEL test": {
			in:    "DEL key",
			query: Query{Command: "DEL", Args: []string{"key"}},
		},
	}

	for name, test := range posTests {
		t.Run(name, func(t *testing.T) {
			query, _ := parser.Parse(test.in)
			assert.Equal(t, query, test.query)
		})
	}
}
