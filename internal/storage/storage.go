package storage

// Storage is interface for storage object
type Storage interface {
	Get(key string) (string, bool)
	Set(key string, value string)
	Delete(key string)
}
