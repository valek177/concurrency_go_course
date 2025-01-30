package wal

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"concurrency_go_course/internal/compute"
	"concurrency_go_course/internal/config"
	"concurrency_go_course/internal/filesystem"
	"concurrency_go_course/pkg/logger"
	"concurrency_go_course/pkg/parser"

	"go.uber.org/zap"
)

const (
	defaultFlushingBatchSize    = 100
	defaultFlushingBatchTimeout = "10ms"
	defaultMaxSegmentSize       = "10MB"
)

type Settings struct {
	MaxSegmentSize       int
	FlushingBatchSize    int
	FlushingBatchTimeout time.Duration
	DataDirectory        string
}

type WAL struct {
	settings Settings

	logsManager *LogsManager

	m      sync.Mutex
	buffer []Request

	bufferCh chan []Request

	writeStatus <-chan error
}

func New(cfg *config.WALCfg) (*WAL, error) {
	if cfg == nil {
		return nil, fmt.Errorf("unable to create WAL: cfg is empty")
	}

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
	}

	segmentSize, err = parser.ParseSize(cfg.WalConfig.MaxSegmentSize)
	if err != nil {
		settings.MaxSegmentSize = segmentSize
	}

	if cfg.WalConfig.FlushingBatchSize != 0 {
		settings.FlushingBatchSize = cfg.WalConfig.FlushingBatchSize
	}

	batchTimeout, err := time.ParseDuration(cfg.WalConfig.FlushingBatchTimeout)
	if err == nil && batchTimeout != 0 {
		settings.FlushingBatchTimeout = batchTimeout
	}

	segment := filesystem.NewSegment(cfg.WalConfig.DataDirectory, segmentSize)

	logsManager, err := NewLogsManager(segment)
	if err != nil {
		return nil, err
	}

	if _, err := os.Stat(cfg.WalConfig.DataDirectory); err != nil {
		if os.IsNotExist(err) {
			err := os.MkdirAll(cfg.WalConfig.DataDirectory, os.ModePerm)
			if err != nil {
				return nil, fmt.Errorf("mkdir error: %w", err)
			}
		} else {
			return nil, fmt.Errorf("unable to get dir info: %w", err)
		}
	}

	return &WAL{
		settings:    settings,
		m:           sync.Mutex{},
		buffer:      make([]Request, 0),
		bufferCh:    make(chan []Request, 1),
		logsManager: logsManager,
	}, nil
}

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
				logger.Debug("Write batch by buffer")
			case <-ticker.C:
				w.flushBatch()
				logger.Debug("Write batch by timeout")
			}
		}
	}()
}

func (w *WAL) Recover() ([]Request, error) {
	return w.logsManager.ReadAll()
}

func (w *WAL) Set(key, value string) error {
	w.push(compute.CommandSet, []string{key, value})

	return <-w.writeStatus
}

func (w *WAL) Del(key string) error {
	w.push(compute.CommandDelete, []string{key})

	return <-w.writeStatus
}

func (w *WAL) push(cmd string, args []string) {
	request := NewRequest(cmd, args)

	w.m.Lock()
	w.buffer = append(w.buffer, request)
	if len(w.buffer) == w.settings.FlushingBatchSize {
		w.bufferCh <- w.buffer
		w.buffer = nil
	}
	w.m.Unlock()

	w.writeStatus = request.doneStatus
}

func (w *WAL) flushBatch() {
	var batch []Request

	w.m.Lock()
	batch = w.buffer
	w.buffer = nil
	w.m.Unlock()

	if len(batch) != 0 {
		w.logsManager.Write(batch)
	}
}
