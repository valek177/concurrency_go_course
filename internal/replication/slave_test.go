package replication

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"concurrency_go_course/internal/config"
	"concurrency_go_course/pkg/logger"
)

func TestNewSlaveErr(t *testing.T) {
	t.Parallel()

	logger.MockLogger()

	cfgWithoutReplicaAddr := &config.Config{
		Replication: &config.ReplicationConfig{
			ReplicaType:   "slave",
			MasterAddress: "",
		},
	}

	cfgWithReplicaAddr := &config.Config{
		Replication: &config.ReplicationConfig{
			ReplicaType:   "slave",
			MasterAddress: "127.0.0.1:9996",
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
			name:          "New replica client without cfg",
			cfg:           nil,
			walCfg:        nil,
			expectedError: fmt.Errorf("config is empty"),
		},
		{
			name:          "New replica client without WAL config",
			cfg:           cfgWithReplicaAddr,
			walCfg:        nil,
			expectedError: fmt.Errorf("WAL config is empty"),
		},
		{
			name:          "New replica client without address",
			cfg:           cfgWithoutReplicaAddr,
			walCfg:        walCfg,
			expectedError: fmt.Errorf("connection create error: dial tcp: missing address"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := NewReplicationClient(tt.cfg, tt.walCfg)
			assert.Nil(t, client)
			assert.Equal(t, tt.expectedError.Error(), err.Error())
		})
	}
}
