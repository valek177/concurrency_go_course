package main

import (
	"bufio"
	"fmt"
	"os"

	"go.uber.org/zap"

	"concurrency_go_course/internal/compute"
	"concurrency_go_course/internal/service"
	"concurrency_go_course/internal/storage"
	"concurrency_go_course/pkg/logger"
)

// LogLevel is log level for logging
const LogLevel = "debug"

func main() {
	logger.InitLogger(LogLevel)

	storage := storage.NewEngine()

	requestParser := compute.NewRequestParser()
	compute := compute.NewCompute(requestParser)

	service := service.NewService(storage, compute)

	logger.Debug("App started")

	fmt.Println("Waiting for commands:")

	reader := bufio.NewReader(os.Stdin)
	for {
		query, err := reader.ReadString('\n')
		if err != nil {
			logger.Error("Failed to read query", zap.Error(err))
			continue
		}

		res, err := service.Handle(query)
		if err != nil {
			logger.Error("Failed to handle query", zap.String("query", query),
				zap.Error(err))
			continue
		}

		fmt.Println(res)
	}
}
