package app

import (
	"fmt"

	"concurrency_go_course/internal/compute"
	"concurrency_go_course/internal/config"
	"concurrency_go_course/internal/service"
	"concurrency_go_course/internal/storage"
)

// Init initializes new database service
func Init(cfg *config.Config) (service.Service, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config is empty")
	}

	storage := storage.NewEngine()

	requestParser := compute.NewRequestParser()
	compute := compute.NewCompute(requestParser)

	service := service.NewService(storage, compute)

	return service, nil
}
