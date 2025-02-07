package storage

import "sync"

// Engine is interface for engine
type Engine interface {
	Get(key string) (string, bool)
	Set(key string, value string)
	Delete(key string)
}

type engine struct {
	m    sync.RWMutex
	data map[string]string
}

// NewEngine returns new engine
func NewEngine() Engine {
	return &engine{
		data: make(map[string]string),
	}
}

// Get returns value
func (e *engine) Get(key string) (string, bool) {
	e.m.Lock()
	defer e.m.Unlock()
	value, ok := e.data[key]

	return value, ok
}

// Set sets new value for key
func (e *engine) Set(key string, value string) {
	e.m.Lock()
	defer e.m.Unlock()
	e.data[key] = value
}

// Delete deletes key-value pair
func (e *engine) Delete(key string) {
	e.m.Lock()
	defer e.m.Unlock()
	delete(e.data, key)
}
