package replication

import (
	"bytes"
	"context"
	"encoding/gob"
	"fmt"
	"path"
	"time"

	"concurrency_go_course/internal/config"
	"concurrency_go_course/internal/filesystem"
	"concurrency_go_course/internal/network"
	"concurrency_go_course/internal/storage/wal"
	"concurrency_go_course/pkg/logger"
)

type Slave struct {
	masterAddress string
	connection    *network.TCPClient
	syncInterval  time.Duration
	walDirectory  string
	stream        chan []wal.Request
	fileLib       filesystem.FileLib
}

func NewReplicationClient(
	cfg *config.Config, walCfg *config.WALCfg,
) (*Slave, error) {
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

func (s *Slave) Start(ctx context.Context) error {
	logger.Debug("replication client was started")
	ticker := time.NewTicker(s.syncInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			s.connection.Close()
			return nil
		default:
		}

		select {
		case <-ticker.C:
			s.syncWithMaster()

		case <-ctx.Done():
			s.connection.Close()
			return nil
		}
	}
}

func (s *Slave) ReplicationStream() chan []wal.Request {
	fmt.Println("slave ", s)
	return s.stream
}

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

	fmt.Println("segment name ", response.SegmentName)

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

	fmt.Println("get segment data ", segmentData)

	var queries []wal.Request
	buffer := bytes.NewBuffer(segmentData)
	fmt.Println("get segment data 11", buffer)
	decoder := gob.NewDecoder(buffer)
	if err := decoder.Decode(&queries); err != nil {
		return fmt.Errorf("unable to decode data: %w", err)
	}

	s.stream <- queries
	return nil
}
