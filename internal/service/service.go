package service

import (
	"concurrency_go_course/internal/compute"
	"concurrency_go_course/internal/storage"
)

// Service is interface for service
type Service interface {
	Handle(request string) (string, error)
}

type serv struct {
	storage storage.Storage
	compute compute.Compute
}

// NewService returns new service
func NewService(
	storage storage.Storage,
	compute compute.Compute,
) *serv {
	return &serv{
		storage: storage,
		compute: compute,
	}
}
