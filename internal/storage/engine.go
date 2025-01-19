package storage

// Engine is struct for key value data
type Engine struct {
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
	value, ok := e.data[key]

	return value, ok
}

// Set sets new value for key
func (e *Engine) Set(key string, value string) {
	e.data[key] = value
}

// Delete deletes key-value pair
func (e *Engine) Delete(key string) {
	delete(e.data, key)
}
