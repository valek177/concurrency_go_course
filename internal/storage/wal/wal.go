package wal

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"go.uber.org/zap"

	"concurrency_go_course/internal/compute"
	"concurrency_go_course/internal/config"
	"concurrency_go_course/internal/filesystem"
	"concurrency_go_course/pkg/logger"
	"concurrency_go_course/pkg/parser"
)

const (
	defaultFlushingBatchSize    = 100
	defaultFlushingBatchTimeout = "10ms"
	defaultMaxSegmentSize       = "10MB"
)

// Settings is a WAL settings struct
type Settings struct {
	MaxSegmentSize       int
	FlushingBatchSize    int
	FlushingBatchTimeout time.Duration
	DataDirectory        string
}

// WAL is a write ahead log struct
type WAL struct {
	settings *Settings

	logsManager LogsManager

	mutexBuffer sync.Mutex
	buffer      []Request

	bufferCh chan []Request

	writeStatus <-chan error
}

// New creates new WAL
func New(cfg *config.WALCfg) (*WAL, error) {
	if cfg == nil {
		return nil, fmt.Errorf("unable to create WAL: cfg is empty")
	}

	settings, err := walSettings(cfg)
	if err != nil {
		return nil, err
	}

	logsManager, err := getLogsManager(settings)
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(settings.DataDirectory); err != nil {
		if os.IsNotExist(err) {
			err := os.MkdirAll(settings.DataDirectory, os.ModePerm) //nolint:gosec
			if err != nil {
				return nil, fmt.Errorf("mkdir error: %w", err)
			}
		} else {
			return nil, fmt.Errorf("unable to get dir info: %w", err)
		}
	}

	return &WAL{
		settings:    settings,
		mutexBuffer: sync.Mutex{},
		buffer:      make([]Request, 0),
		bufferCh:    make(chan []Request, 1),
		logsManager: logsManager,
	}, nil
}

// Start initializes WAL
func (w *WAL) Start(ctx context.Context) {
	logger.Info("Starting WAL with settings",
		zap.String("flushing_timeout", w.settings.FlushingBatchTimeout.String()),
		zap.Int("flushing_batch_size", w.settings.FlushingBatchSize),
		zap.Int("max_segment_size", w.settings.MaxSegmentSize),
	)

	go func() {
		ticker := time.NewTicker(w.settings.FlushingBatchTimeout)
		defer ticker.Stop()

		for {
			select {
			case <-ctx.Done():
				w.flushBatch()
				logger.Debug("Batch was flushed by ctx")
				return
			default:
			}

			select {
			case <-ctx.Done():
				w.flushBatch()
				logger.Debug("Batch was flushed by ctx")
				return
			case batch := <-w.bufferCh:
				w.logsManager.Write(batch)
				ticker.Reset(w.settings.FlushingBatchTimeout * time.Second)
				logger.Debug("Batch was flushed by buffer")
			case <-ticker.C:
				w.flushBatch()
				logger.Debug("Batch was flushed by timeout")
			}
		}
	}()
}

// Recover recover from files
func (w *WAL) Recover() ([]Request, error) {
	return w.logsManager.ReadAll()
}

// Set sets new value
func (w *WAL) Set(key, value string) error {
	w.push(compute.CommandSet, []string{key, value})

	return <-w.writeStatus
}

// Del deletes key
func (w *WAL) Del(key string) error {
	w.push(compute.CommandDelete, []string{key})

	return <-w.writeStatus
}

func (w *WAL) push(cmd string, args []string) {
	request := NewRequest(cmd, args)

	w.mutexBuffer.Lock()
	w.buffer = append(w.buffer, request)
	if len(w.buffer) == w.settings.FlushingBatchSize {
		w.bufferCh <- w.buffer
		w.buffer = nil
	}
	w.mutexBuffer.Unlock()

	w.writeStatus = request.doneStatus
}

func (w *WAL) flushBatch() {
	var batch []Request

	w.mutexBuffer.Lock()
	batch = w.buffer
	w.buffer = nil
	w.mutexBuffer.Unlock()

	if len(batch) != 0 {
		w.logsManager.Write(batch)
	}
}

func walSettings(cfg *config.WALCfg) (*Settings, error) {
	segmentSize, err := parser.ParseSize(defaultMaxSegmentSize)
	if err != nil {
		return nil, err
	}

	timeout, err := time.ParseDuration(defaultFlushingBatchTimeout)
	if err != nil {
		return nil, err
	}

	settings := Settings{
		MaxSegmentSize:       segmentSize,
		FlushingBatchTimeout: timeout,
		FlushingBatchSize:    defaultFlushingBatchSize,
		DataDirectory:        cfg.WalConfig.DataDirectory,
	}

	segmentSize, err = parser.ParseSize(cfg.WalConfig.MaxSegmentSize)
	if err == nil && segmentSize != 0 {
		settings.MaxSegmentSize = segmentSize
	}

	if cfg.WalConfig.FlushingBatchSize != 0 {
		settings.FlushingBatchSize = cfg.WalConfig.FlushingBatchSize
	}

	batchTimeout, err := time.ParseDuration(cfg.WalConfig.FlushingBatchTimeout)
	if err == nil && batchTimeout != 0 {
		settings.FlushingBatchTimeout = batchTimeout
	}

	return &settings, nil
}

func getLogsManager(settings *Settings) (LogsManager, error) {
	fileLib := filesystem.NewFileLib()

	segment := filesystem.NewSegment(settings.DataDirectory,
		settings.MaxSegmentSize, fileLib)

	logsManager, err := NewLogsManager(segment)
	if err != nil {
		return nil, err
	}

	return logsManager, nil
}
