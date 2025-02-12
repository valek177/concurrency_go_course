package network

import (
	"context"
	"fmt"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"concurrency_go_course/internal/config"
	"concurrency_go_course/pkg/logger"
)

func TestNewServerNil(t *testing.T) {
	t.Parallel()

	logger.MockLogger()

	serverAddr := "127.0.0.1:7777"

	cfg := &config.Config{
		Engine: &config.EngineConfig{
			Type:             "in_memory",
			PartitionsNumber: 256,
		},
		Network: &config.NetworkConfig{
			Address:        serverAddr,
			MaxConnections: 100,
			MaxMessageSize: "4KB",
			IdleTimeout:    "5m",
		},
		Logging: &config.LoggingConfig{
			Level:  "info",
			Output: "log/output.log",
		},
		Replication: &config.ReplicationConfig{
			ReplicaType: "master",
		},
	}

	tests := []struct {
		name    string
		cfg     *config.Config
		address string

		resultServer  *TCPServer
		expectedError error
	}{
		{
			name:          "New server without config",
			cfg:           nil,
			address:       serverAddr,
			resultServer:  nil,
			expectedError: fmt.Errorf("config is empty"),
		},
		{
			name:          "New server without address",
			cfg:           cfg,
			address:       "",
			resultServer:  nil,
			expectedError: fmt.Errorf("address is empty"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server, err := NewServer(tt.cfg, tt.address)
			assert.Nil(t, server)
			assert.Equal(t, tt.expectedError, err)
		})
	}
}

func TestNewServer(t *testing.T) {
	t.Parallel()

	logger.MockLogger()

	serverAddr := "127.0.0.1:8888"

	cfg := &config.Config{
		Engine: &config.EngineConfig{
			Type: "in_memory",
		},
		Network: &config.NetworkConfig{
			Address:        serverAddr,
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
		cfg           *config.Config
		resultServer  *TCPServer
		expectedError error
	}{
		{
			name: "New server with config",
			cfg:  cfg,
			resultServer: &TCPServer{
				address: serverAddr,
				cfg:     cfg,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server, err := NewServer(tt.cfg, serverAddr)
			assert.Nil(t, err)
			assert.Equal(t, tt.resultServer.cfg, server.cfg)
		})
	}
}

func TestRun(t *testing.T) {
	t.Parallel()

	logger.MockLogger()

	ctx := context.Background()

	db := NewMockDatabase()

	addr := "127.0.0.1:5555"

	cfg := config.Config{
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

	server, err := NewServer(&cfg, addr)
	if err != nil {
		t.Errorf("want nil error; got %+v", err)
	}

	time.Sleep(100 * time.Millisecond)

	go server.Run(ctx, func(_ context.Context, s []byte) []byte {
		response, err := db.Handle(string(s))
		if err != nil {
			response = err.Error()
		}
		return []byte(response)
	})

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
