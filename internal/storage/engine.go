package storage

import "sync"

// Engine is struct for key value data
type Engine struct {
	m    sync.Mutex
	data map[string]string
}

// NewEngine returns new engine
func NewEngine() *Engine {
	return &Engine{
		data: make(map[string]string),
	}
}

// Get returns value
func (e *Engine) Get(key string) (string, bool) {
	e.m.Lock()
	defer e.m.Unlock()
	value, ok := e.data[key]

	return value, ok
}

// Set sets new value for key
func (e *Engine) Set(key string, value string) {
	e.m.Lock()
	defer e.m.Unlock()
	e.data[key] = value
}

// Delete deletes key-value pair
func (e *Engine) Delete(key string) {
	e.m.Lock()
	defer e.m.Unlock()
	delete(e.data, key)
}
