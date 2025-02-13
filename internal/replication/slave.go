package replication

import (
	"bytes"
	"context"
	"fmt"
	"path"
	"time"

	"concurrency_go_course/internal/config"
	"concurrency_go_course/internal/filesystem"
	"concurrency_go_course/internal/network"
	"concurrency_go_course/internal/storage/wal"
	"concurrency_go_course/pkg/logger"

	"go.uber.org/zap"
)

// Slave is struct for slave replication
type Slave struct {
	masterAddress string
	connection    *network.TCPClient
	syncInterval  time.Duration
	walDirectory  string
	stream        chan []wal.Request
	fileLib       filesystem.FileLib
}

// NewReplicationClient returns new replication client
func NewReplicationClient(
	cfg *config.Config, walCfg *config.WALCfg,
) (*Slave, error) {
	if cfg == nil || cfg.Replication == nil {
		return nil, fmt.Errorf("config is empty")
	}

	if walCfg == nil || walCfg.WalConfig == nil {
		return nil, fmt.Errorf("WAL config is empty")
	}

	connection, err := network.NewClient(cfg.Replication.MasterAddress)
	if err != nil {
		return nil, fmt.Errorf("connection create error: %w", err)
	}

	return &Slave{
		connection:    connection,
		masterAddress: cfg.Replication.MasterAddress,
		syncInterval:  cfg.Replication.SyncInterval,
		walDirectory:  walCfg.WalConfig.DataDirectory,
		stream:        make(chan []wal.Request),
		fileLib:       filesystem.NewFileLib(),
	}, nil
}

// Start starts slave
func (s *Slave) Start(ctx context.Context) {
	logger.Debug("replication client was started",
		zap.String("sync_interval", s.syncInterval.String()))
	ticker := time.NewTicker(s.syncInterval)
	defer func() {
		ticker.Stop()
		s.connection.Close()
	}()

	for {
		select {
		case <-ctx.Done():
			logger.Debug("replication client stopping")
			s.connection.Close()
			return
		default:
		}

		select {
		case <-ticker.C:
			s.syncWithMaster()

		case <-ctx.Done():
			logger.Debug("replication client stopping")
			s.connection.Close()
			return
		}
	}
}

// ReplicationStream returns replication stream channel
func (s *Slave) ReplicationStream() chan []wal.Request {
	return s.stream
}

// IsMaster returns flag
func (s *Slave) IsMaster() bool {
	return false
}

func (s *Slave) syncWithMaster() {
	lastSegmentName, err := s.fileLib.SegmentLast(s.walDirectory)
	if err != nil {
		logger.ErrorWithMsg("unable to sync on slave:", err)
	}
	req := SlaveRequest{LastSegmentName: lastSegmentName}

	data, err := EncodeSlaveRequest(&req)
	if err != nil {
		logger.ErrorWithMsg("unable to encode request", err)
		return
	}

	resp, err := s.connection.Send(data)
	if err != nil {
		logger.ErrorWithMsg("unable to connect with master", err)
		return
	}

	response := &MasterResponse{}
	err = DecodeResponse(response, resp)
	if err != nil {
		logger.ErrorWithMsg("unable to decode response", err)
		return
	}

	err = s.saveSegment(response.SegmentName, response.SegmentData)
	if err != nil {
		logger.ErrorWithMsg("unable to save segment", err)
		return
	}

	err = s.applyDataToEngine(response.SegmentData)
	if err != nil {
		logger.ErrorWithMsg("unable to apply data to engine", err)
		return
	}
}

func (s *Slave) saveSegment(name string, data []byte) error {
	if name == "" {
		return nil
	}
	filename := path.Join(s.walDirectory, name)
	segmentFile, err := s.fileLib.CreateFile(filename)
	if err != nil {
		return err
	}

	if _, err = s.fileLib.WriteFile(segmentFile, data); err != nil {
		return err
	}

	return nil
}

func (s *Slave) applyDataToEngine(segmentData []byte) error {
	if len(segmentData) == 0 {
		return nil
	}

	var queries []wal.Request
	buffer := bytes.NewBuffer(segmentData)
	for buffer.Len() > 0 {
		var request wal.Request
		if err := request.Decode(buffer); err != nil {
			return fmt.Errorf("unable to parse request data: %w", err)
		}

		queries = append(queries, request)
	}

	s.stream <- queries
	return nil
}
