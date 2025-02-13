package replication

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"concurrency_go_course/internal/config"
	"concurrency_go_course/internal/filesystem"
	"concurrency_go_course/internal/network"
	"concurrency_go_course/pkg/logger"

	"go.uber.org/zap"
)

// Master is a struct for master node
type Master struct {
	server       *network.TCPServer
	walDirectory string
	fileLib      filesystem.FileLib
}

// TCPServer is interface for TCP server
type TCPServer interface {
	Run(context.Context, func(context.Context, []byte) []byte)
}

// IsMaster returns flag
func (m *Master) IsMaster() bool {
	return true
}

// NewReplicationServer creates new master replication server
func NewReplicationServer(cfg *config.Config, walCfg *config.WALCfg) (*Master, error) {
	if cfg == nil || cfg.Replication == nil {
		return nil, fmt.Errorf("config is empty")
	}

	if walCfg == nil || walCfg.WalConfig == nil {
		return nil, fmt.Errorf("WAL config is empty")
	}

	server, err := network.NewServer(cfg, cfg.Replication.MasterAddress)
	if err != nil {
		return nil, err
	}

	return &Master{
		server:       server,
		walDirectory: walCfg.WalConfig.DataDirectory,
		fileLib:      filesystem.NewFileLib(),
	}, nil
}

// Start starts master
func (m *Master) Start(ctx context.Context) {
	logger.Debug("replication master server was started")
	m.server.Run(ctx, func(ctx context.Context, requestData []byte) []byte {
		if ctx.Err() != nil {
			return nil
		}

		var request SlaveRequest
		if err := DecodeSlaveRequest(&request, requestData); err != nil {
			logger.Error("unable to decode replication request", zap.Error(err))
			return nil
		}

		response := m.lastSegment(request)
		responseData, err := EncodeResponse(&response)
		if err != nil {
			logger.Error("unable to encode replication response", zap.Error(err))
		}

		return responseData
	})
}

func (m *Master) lastSegment(request SlaveRequest) MasterResponse {
	var response MasterResponse

	segmentName, err := m.fileLib.SegmentNext(m.walDirectory, request.LastSegmentName)
	if err != nil {
		logger.Error("failed to find WAL segment", zap.Error(err))
		return response
	}

	if segmentName == "" {
		response.Succeed = true
		return response
	}

	filename := fmt.Sprintf("%s/%s", m.walDirectory, segmentName)
	data, err := os.ReadFile(filepath.Clean(filename))
	if err != nil {
		logger.Error("failed to read WAL segment", zap.Error(err))
		return response
	}

	response.Succeed = true
	response.SegmentData = data
	response.SegmentName = segmentName

	logger.Debug("sending response to client ",
		zap.String("name", response.SegmentName))

	return response
}
