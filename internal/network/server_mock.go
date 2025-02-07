package network

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

// MockDatabase is a struct for mocking database
type MockDatabase interface {
	Handle(request string) (string, error)
	StartWAL(ctx context.Context)
}

type mockdatabase struct{}

// NewMockDatabase mocks Database
func NewMockDatabase() MockDatabase {
	var db mockdatabase
	return &db
}

// Handle is a mock function for Handle
func (m *mockdatabase) Handle(request string) (string, error) {
	if strings.Contains(request, "error") {
		return "", errors.New("unable to handle request")
	}

	return fmt.Sprintf("hello %s", request), nil
}

func (m *mockdatabase) StartWAL(_ context.Context) {}
