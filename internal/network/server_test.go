package network

import (
	"fmt"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"concurrency_go_course/internal/config"
	mstorage "concurrency_go_course/internal/storage/mock"
	"concurrency_go_course/pkg/logger"
)

func TestNewServerNil(t *testing.T) {
	t.Parallel()

	logger.MockLogger()

	tests := []struct {
		name          string
		cfg           *config.ServerConfig
		resultServer  *TCPServer
		expectedError error
	}{
		{
			name:          "New server without config",
			cfg:           nil,
			resultServer:  nil,
			expectedError: fmt.Errorf("config is empty"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server, err := NewServer(tt.cfg)
			assert.Nil(t, server)
			assert.Equal(t, tt.expectedError, err)
		})
	}
}

func TestNewServer(t *testing.T) {
	t.Parallel()

	logger.MockLogger()

	cfg := config.ServerConfig{
		Engine: &config.EngineConfig{
			Type: "in_memory",
		},
		Network: &config.NetworkConfig{
			Address:        "127.0.0.1:7777",
			MaxConnections: 100,
			MaxMessageSize: "4KB",
			IdleTimeout:    "5m",
		},
		Logging: &config.LoggingConfig{
			Level:  "info",
			Output: "log/output.log",
		},
	}

	tests := []struct {
		name          string
		cfg           *config.ServerConfig
		resultServer  *TCPServer
		expectedError error
	}{
		{
			name: "New server without config",
			cfg:  &cfg,
			resultServer: &TCPServer{
				cfg: &cfg,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server, err := NewServer(tt.cfg)
			assert.Nil(t, err)
			assert.Equal(t, tt.resultServer.cfg, server.cfg)
		})
	}
}

func TestRun(t *testing.T) {
	t.Parallel()

	logger.MockLogger()

	ctrl := gomock.NewController(t)
	defer t.Cleanup(ctrl.Finish)

	mockStorage := mstorage.NewMockStorage(ctrl)
	service := &MockService{
		storage: mockStorage,
	}

	addr := "127.0.0.1:5555"

	cfg := config.ServerConfig{
		Engine: &config.EngineConfig{
			Type: "in_memory",
		},
		Network: &config.NetworkConfig{
			Address:        addr,
			MaxConnections: 100,
			MaxMessageSize: "4KB",
			IdleTimeout:    "5m",
		},
		Logging: &config.LoggingConfig{
			Level:  "info",
			Output: "log/output.log",
		},
	}

	server, err := NewServer(&cfg)
	server.dbService = service
	if err != nil {
		t.Errorf("want nil error; got %+v", err)
	}

	time.Sleep(100 * time.Millisecond)

	go server.Run()

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()

		conn, err := net.Dial("tcp", addr)
		if err != nil {
			t.Errorf("want nil error; got %+v", err)
		}

		_, err = conn.Write([]byte("first"))
		if err != nil {
			t.Errorf("want nil error; got %+v", err)
		}

		buffer := make([]byte, 1024)
		size, err := conn.Read(buffer)
		if err != nil {
			t.Errorf("want nil error; got %+v", err)
		}

		assert.Equal(t, "hello first", string(buffer[:size]))
	}()

	go func() {
		defer wg.Done()

		conn, err := net.Dial("tcp", addr)
		if err != nil {
			t.Errorf("want nil error; got %+v", err)
		}

		_, err = conn.Write([]byte("second"))
		if err != nil {
			t.Errorf("want nil error; got %+v", err)
		}

		buffer := make([]byte, 1024)
		size, err := conn.Read(buffer)
		if err != nil {
			t.Errorf("want nil error; got %+v", err)
		}

		assert.Equal(t, "hello second", string(buffer[:size]))
	}()

	wg.Wait()

	if err := server.listener.Close(); err != nil {
		t.Errorf("unable to close listener %s", err.Error())
	}
}
