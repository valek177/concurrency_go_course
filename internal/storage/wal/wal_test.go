package wal

import (
	"fmt"
	"testing"
	"time"

	"concurrency_go_course/internal/config"
	"concurrency_go_course/pkg/logger"

	"github.com/stretchr/testify/assert"
)

func TestNewWAL(t *testing.T) {
	t.Parallel()
	logger.MockLogger()

	tests := []struct {
		name     string
		cfg      *config.WALCfg
		settings Settings
		err      error
	}{
		{
			name: "New correct WAL",
			cfg: &config.WALCfg{
				WalConfig: &config.WALSettings{
					FlushingBatchSize:    100,
					FlushingBatchTimeout: "100ms",
					MaxSegmentSize:       "1MB",
					DataDirectory:        "tmp",
				},
			},
			settings: Settings{
				FlushingBatchSize:    100,
				FlushingBatchTimeout: 100 * time.Millisecond,
				MaxSegmentSize:       1024 * 1024,
				DataDirectory:        "tmp",
			},
		},
		{
			name: "New WAL with invalid FlushingBatchSize",
			cfg: &config.WALCfg{
				WalConfig: &config.WALSettings{
					FlushingBatchSize:    0,
					FlushingBatchTimeout: "100ms",
					MaxSegmentSize:       "1MB",
					DataDirectory:        "tmp",
				},
			},
			settings: Settings{
				FlushingBatchSize:    100,
				FlushingBatchTimeout: 100 * time.Millisecond,
				MaxSegmentSize:       1024 * 1024,
				DataDirectory:        "tmp",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wal, err := New(tt.cfg)
			assert.Nil(t, err)
			assert.Equal(t, tt.settings, wal.settings)
		})
	}
}

func TestNewWALNeg(t *testing.T) {
	t.Parallel()
	logger.MockLogger()

	tests := []struct {
		name string
		cfg  *config.WALCfg
		err  error
	}{
		{
			name: "Empty config (error)",
			cfg:  nil,
			err:  fmt.Errorf("unable to create WAL: cfg is empty"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			wal, err := New(tt.cfg)
			assert.Nil(t, wal)
			assert.Equal(t, tt.err, err)
		})
	}
}
