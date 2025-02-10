package storage

import (
	"sync"
)

// HashTable is a struct for hash table
type HashTable struct {
	mutex sync.RWMutex
	data  map[string]string
}

// NewHashTable returns new hash table
func NewHashTable() *HashTable {
	return &HashTable{
		data: make(map[string]string),
	}
}

// Set sets new key-value
func (s *HashTable) Set(key, value string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	s.data[key] = value
}

// Get returns value for key
func (s *HashTable) Get(key string) (string, bool) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	value, found := s.data[key]
	return value, found
}

// Del deletes key
func (s *HashTable) Del(key string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	delete(s.data, key)
}
