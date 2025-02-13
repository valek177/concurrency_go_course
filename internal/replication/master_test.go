package replication

import (
	"context"
	"fmt"
	"net"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"

	"concurrency_go_course/internal/config"
	"concurrency_go_course/pkg/logger"
)

func TestNewServerErr(t *testing.T) {
	t.Parallel()

	logger.MockLogger()

	cfgWithoutReplicaAddr := &config.Config{
		Replication: &config.ReplicationConfig{
			ReplicaType:   "master",
			MasterAddress: "",
		},
	}

	cfgWithReplicaAddr := &config.Config{
		Replication: &config.ReplicationConfig{
			ReplicaType:   "master",
			MasterAddress: "127.0.0.1:9998",
		},
	}

	walCfg := &config.WALCfg{
		WalConfig: &config.WALSettings{
			FlushingBatchSize:    100,
			FlushingBatchTimeout: "10ms",
			MaxSegmentSize:       "10MB",
			DataDirectory:        "tmp",
		},
	}

	tests := []struct {
		name   string
		cfg    *config.Config
		walCfg *config.WALCfg

		expectedError error
	}{
		{
			name:          "New server without cfg",
			cfg:           nil,
			walCfg:        nil,
			expectedError: fmt.Errorf("config is empty"),
		},
		{
			name:          "New server without WAL config",
			cfg:           cfgWithReplicaAddr,
			walCfg:        nil,
			expectedError: fmt.Errorf("WAL config is empty"),
		},
		{
			name:          "New server without address",
			cfg:           cfgWithoutReplicaAddr,
			walCfg:        walCfg,
			expectedError: fmt.Errorf("address is empty"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server, err := NewReplicationServer(tt.cfg, tt.walCfg)
			assert.Nil(t, server)
			assert.Equal(t, tt.expectedError, err)
		})
	}
}

func TestNewServerOk(t *testing.T) {
	t.Parallel()

	logger.MockLogger()

	cfg := &config.Config{
		Engine: &config.EngineConfig{
			Type:             "in_memory",
			PartitionsNumber: 256,
		},
		Network: &config.NetworkConfig{
			Address:        "127.0.0.1:9997",
			MaxConnections: 100,
			MaxMessageSize: "4KB",
			IdleTimeout:    "5m",
		},
		Logging: &config.LoggingConfig{
			Level:  "info",
			Output: "log/output.log",
		},
		Replication: &config.ReplicationConfig{
			ReplicaType:   "master",
			MasterAddress: "127.0.0.1:9998",
		},
	}

	walCfg := &config.WALCfg{
		WalConfig: &config.WALSettings{
			FlushingBatchSize:    100,
			FlushingBatchTimeout: "10ms",
			MaxSegmentSize:       "10MB",
			DataDirectory:        "tmp",
		},
	}

	tests := []struct {
		name   string
		cfg    *config.Config
		walCfg *config.WALCfg

		expectedError error
	}{
		{
			name:          "New server master OK",
			cfg:           cfg,
			walCfg:        walCfg,
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server, err := NewReplicationServer(tt.cfg, tt.walCfg)
			assert.Nil(t, err)
			assert.NotNil(t, server)
		})
	}
}

func TestNewServerStart(t *testing.T) {
	logger.MockLogger()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	replMasterAddr := "127.0.0.1:9998"

	cfg := &config.Config{
		Engine: &config.EngineConfig{
			Type:             "in_memory",
			PartitionsNumber: 256,
		},
		Network: &config.NetworkConfig{
			Address:        "127.0.0.1:9997",
			MaxConnections: 100,
			MaxMessageSize: "4KB",
			IdleTimeout:    "5m",
		},
		Logging: &config.LoggingConfig{
			Level:  "info",
			Output: "log/output.log",
		},
		Replication: &config.ReplicationConfig{
			ReplicaType:   "master",
			MasterAddress: replMasterAddr,
		},
	}

	walCfg := &config.WALCfg{
		WalConfig: &config.WALSettings{
			FlushingBatchSize:    100,
			FlushingBatchTimeout: "10ms",
			MaxSegmentSize:       "10MB",
			DataDirectory:        "test_data",
		},
	}

	server, err := NewReplicationServer(cfg, walCfg)
	if err != nil {
		t.Errorf("want nil error; got %+v", err)
	}
	assert.NotNil(t, server)

	go server.Start(ctx)

	wg := sync.WaitGroup{}

	expectedResponse := MasterResponse{
		Succeed:     true,
		SegmentName: "wal_1.log",
	}

	wg.Add(1)
	go func() {
		defer wg.Done()

		conn, err := net.Dial("tcp", replMasterAddr)
		if err != nil {
			t.Errorf("want nil error; got %+v", err)
		}

		req := SlaveRequest{LastSegmentName: "wal_0.log"}

		data, err := EncodeSlaveRequest(&req)
		if err != nil {
			t.Errorf("want nil error; got %+v", err)
		}

		_, err = conn.Write(data)
		if err != nil {
			t.Errorf("want nil error; got %+v", err)
		}

		buffer := make([]byte, 1024)
		_, err = conn.Read(buffer)
		if err != nil {
			t.Errorf("want nil error; got %+v", err)
		}

		response := &MasterResponse{}
		err = DecodeResponse(response, buffer)
		if err != nil {
			t.Errorf("want nil error; got %+v", err)
		}

		assert.Equal(t, expectedResponse.Succeed, response.Succeed)
		assert.Equal(t, expectedResponse.SegmentName, response.SegmentName)
	}()

	wg.Wait()
}
