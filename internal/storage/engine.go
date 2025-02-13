package storage

import (
	"hash/fnv"
	"sync"
)

// Engine is interface for engine
type Engine interface {
	Get(key string) (string, bool)
	Set(key string, value string)
	Delete(key string)
}

type engine struct {
	parts []*HashTable
}

const defaultKeyCount = 8

// NewEngine returns new engine
func NewEngine(partsNumber int) Engine {
	engine := &engine{
		parts: make([]*HashTable, partsNumber),
	}

	for i := 0; i < partsNumber; i++ {
		engine.parts[i] = &HashTable{
			mutex: sync.RWMutex{},
			data:  make(map[string]string, defaultKeyCount),
		}
	}
	return engine
}

// Get returns value
func (e *engine) Get(key string) (string, bool) {
	hash := getHash(key, len(e.parts))

	part := e.parts[hash]

	value, ok := part.Get(key)
	if !ok {
		return "", false
	}

	return value, ok
}

// Set sets new value for key
func (e *engine) Set(key string, value string) {
	hash := getHash(key, len(e.parts))
	part := e.parts[hash]

	part.Set(key, value)
}

// Delete deletes key-value pair
func (e *engine) Delete(key string) {
	hash := getHash(key, len(e.parts))
	part := e.parts[hash]

	part.Del(key)
}

func getHash(key string, partsCount int) int {
	hash := fnv.New32a()

	_, _ = hash.Write([]byte(key))
	return int(hash.Sum32()) % partsCount
}
