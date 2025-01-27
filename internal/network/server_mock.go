package network

import (
	"errors"
	"fmt"
	"strings"

	"concurrency_go_course/internal/storage"
)

// MockService is a struct for mocking service
type MockService struct {
	storage storage.Storage
}

// Handle is a mock function for Handle
func (m *MockService) Handle(request string) (string, error) {
	if strings.Contains(request, "error") {
		return "", errors.New("unable to handle request")
	}

	return fmt.Sprintf("hello %s", request), nil
}
