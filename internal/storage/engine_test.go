package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetEngine(t *testing.T) {
	t.Parallel()

	engine := NewEngine(4)
	engine.Set("key1", "a")

	tests := map[string]struct {
		key           string
		expectedValue string
	}{
		"get existing key": {
			key:           "key1",
			expectedValue: "a",
		},
		"get non existing key": {
			key:           "key2",
			expectedValue: "",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			value, _ := engine.Get(test.key)
			assert.Equal(t, value, test.expectedValue)
		})
	}
}

func TestSetEngine(t *testing.T) {
	t.Parallel()

	engine := NewEngine(4)

	tests := map[string]struct {
		key   string
		value string
	}{
		"set new key": {
			key:   "key2",
			value: "b",
		},
		"set existing key": {
			key:   "key1",
			value: "c",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			engine.Set(test.key, test.value)
			value, _ := engine.Get(test.key)
			assert.Equal(t, value, test.value)
		})
	}
}

func TestDeleteEngine(t *testing.T) {
	t.Parallel()

	engine := NewEngine(4)
	engine.Set("key1", "a")
	engine.Set("key2", "a")

	tests := map[string]struct {
		key   string
		value string
	}{
		"delete existing key": {
			key: "key1",
		},
		"delete non existing key": {
			key: "key2",
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			engine.Delete(test.key)
			value, _ := engine.Get(test.key)
			assert.Equal(t, value, "")
		})
	}
}
