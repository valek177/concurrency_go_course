package wal

import (
	"context"
	"fmt"
	"reflect"
	"sync"
	"testing"
	"time"

	"concurrency_go_course/internal/compute"
	"concurrency_go_course/internal/config"
	"concurrency_go_course/pkg/logger"

	"github.com/stretchr/testify/assert"
)

func TestNewWAL(t *testing.T) {
	logger.MockLogger()

	tests := []struct {
		name     string
		cfg      *config.WALCfg
		settings *Settings
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
			settings: &Settings{
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
			settings: &Settings{
				FlushingBatchSize:    100,
				FlushingBatchTimeout: 100 * time.Millisecond,
				MaxSegmentSize:       1024 * 1024,
				DataDirectory:        "tmp",
			},
		},
		{
			name: "New WAL with invalid FlushingBatchTimeout",
			cfg: &config.WALCfg{
				WalConfig: &config.WALSettings{
					FlushingBatchSize:    10,
					FlushingBatchTimeout: "abcd",
					MaxSegmentSize:       "1MB",
					DataDirectory:        "tmp",
				},
			},
			settings: &Settings{
				FlushingBatchSize:    10,
				FlushingBatchTimeout: 10 * time.Millisecond,
				MaxSegmentSize:       1024 * 1024,
				DataDirectory:        "tmp",
			},
		},
		{
			name: "New WAL with invalid MaxSegmentSize",
			cfg: &config.WALCfg{
				WalConfig: &config.WALSettings{
					FlushingBatchSize:    10,
					FlushingBatchTimeout: "20ms",
					MaxSegmentSize:       "abcd",
					DataDirectory:        "tmp",
				},
			},
			settings: &Settings{
				FlushingBatchSize:    10,
				FlushingBatchTimeout: 20 * time.Millisecond,
				MaxSegmentSize:       10 * 1024 * 1024,
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

func TestWAL_Start_WriteByTimeout(t *testing.T) {
	logger.MockLogger()

	cfg := &config.WALCfg{
		WalConfig: &config.WALSettings{
			FlushingBatchSize:    100,
			FlushingBatchTimeout: "100ms",
			MaxSegmentSize:       "1MB",
			DataDirectory:        "tmp",
		},
	}

	wal, err := New(cfg)
	if err != nil {
		t.Errorf("unable to create WAL: %s", err)
	}

	start := time.Now()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wal.Start(ctx)
	err = wal.Set("key", "value")
	if err != nil {
		t.Errorf("unable to set value: %s", err)
	}

	duration := time.Since(start)
	if duration < 100*time.Millisecond || duration > 120*time.Millisecond {
		t.Errorf("wrong WAL write timeout: expected 100-120ms, got %d", duration)
	}
}

func TestWAL_Start_WriteByBatchSize(t *testing.T) {
	t.Parallel()
	logger.MockLogger()

	cfg := &config.WALCfg{
		WalConfig: &config.WALSettings{
			FlushingBatchSize:    2,
			FlushingBatchTimeout: "100ms",
			MaxSegmentSize:       "1MB",
			DataDirectory:        "tmp",
		},
	}

	wal, err := New(cfg)
	if err != nil {
		t.Errorf("unable to create WAL: %s", err)
	}

	start := time.Now()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	wal.Start(ctx)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()

		err = wal.Set("key1", "value1")
		if err != nil {
			t.Errorf("unable to set value: %s", err)
		}
	}()

	go func() {
		defer wg.Done()

		err = wal.Set("key2", "value2")
		if err != nil {
			t.Errorf("unable to set value: %s", err)
		}
	}()

	wg.Wait()

	duration := time.Since(start)
	if duration > 10*time.Millisecond {
		t.Errorf("wrong WAL write timeout: expected less than 10ms, got %d", duration)
	}
}

func TestWAL_Recover(t *testing.T) {
	t.Parallel()
	logger.MockLogger()

	cfg := &config.WALCfg{
		WalConfig: &config.WALSettings{
			FlushingBatchSize:    100,
			FlushingBatchTimeout: "100ms",
			MaxSegmentSize:       "1MB",
			DataDirectory:        "test_data",
		},
	}

	wal, err := New(cfg)
	if err != nil {
		t.Errorf("unable to create WAL: %s", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	wal.Start(ctx)

	requests, err := wal.Recover()
	if err != nil {
		t.Errorf("recover error: %s", err)
	}

	if len(requests) != 3 {
		t.Errorf("recover: got %d requests, expected 3", len(requests))
	}

	for i := range 2 {
		if requests[i].Command != compute.CommandSet {
			t.Errorf("recover error: got command %s, expected %s", requests[i].Command, compute.CommandSet)
		}
	}

	if requests[2].Command != compute.CommandDelete {
		t.Errorf("recover error: got command %s, expected %s", requests[2].Command, compute.CommandDelete)
	}

	if !reflect.DeepEqual(requests[0].Args, []string{"ozzy", "osbourne"}) {
		t.Errorf("recover error: got args %+v, expected %+v", requests[0].Args, []string{"ozzy", "osbourne"})
	}

	if !reflect.DeepEqual(requests[1].Args, []string{"lemmy", "kilmister"}) {
		t.Errorf("recover error: got args %+v, expected %+v", requests[1].Args, []string{"ozzy", "osbourne"})
	}

	if !reflect.DeepEqual(requests[2].Args, []string{"lemmy"}) {
		t.Errorf("recover error: got args %+v, expected %+v", requests[2].Args, []string{"lemmy"})
	}
}
