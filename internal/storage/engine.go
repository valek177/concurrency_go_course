package storage

import "sync"

// Engine is struct for key value data
type EngineObj struct {
	m    sync.RWMutex
	data map[string]string
}

// NewEngine returns new engine
func NewEngine() *EngineObj {
	return &EngineObj{
		data: make(map[string]string),
	}
}

// Get returns value
func (e *EngineObj) Get(key string) (string, bool) {
	e.m.Lock()
	defer e.m.Unlock()
	value, ok := e.data[key]

	return value, ok
}

// Set sets new value for key
func (e *EngineObj) Set(key string, value string) {
	e.m.Lock()
	defer e.m.Unlock()
	e.data[key] = value
}

// Delete deletes key-value pair
func (e *EngineObj) Delete(key string) {
	e.m.Lock()
	defer e.m.Unlock()
	delete(e.data, key)
}
