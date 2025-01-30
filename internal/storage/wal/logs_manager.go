package wal

import (
	"bytes"
	"errors"
	"fmt"

	fs "concurrency_go_course/internal/filesystem"
	"concurrency_go_course/pkg/logger"
)

type Manager interface {
	Write([]byte) error
	ReadAll() ([][]byte, error)
}

type LogsManager struct {
	segment *fs.Segment
}

func NewLogsManager(segment *fs.Segment) (*LogsManager, error) {
	if segment == nil {
		return nil, errors.New("segment is invalid")
	}

	return &LogsManager{segment: segment}, nil
}

func (m *LogsManager) Write(requests []Request) {
	var buffer bytes.Buffer
	for _, req := range requests {
		if err := req.Encode(&buffer); err != nil {
			logger.ErrorWithMsg("failed to encode requests", err)
			m.acknowledgeWrite(requests, err)
			return
		}
	}

	err := m.segment.Write(buffer.Bytes())
	if err != nil {
		logger.ErrorWithMsg("failed to write request data:", err)
	}

	m.acknowledgeWrite(requests, err)
}

func (m *LogsManager) ReadAll() ([]Request, error) {
	segmentsData, err := m.segment.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read segments: %w", err)
	}

	var requests []Request
	for _, data := range segmentsData {
		requests, err = m.readSegment(requests, data)
		if err != nil {
			return nil, fmt.Errorf("failed to read segments: %w", err)
		}
	}

	logger.Debug("WAL requests was readed")

	return requests, nil
}

func (l *LogsManager) readSegment(requests []Request, data []byte) ([]Request, error) {
	buffer := bytes.NewBuffer(data)
	for buffer.Len() > 0 {
		var request Request
		if err := request.Decode(buffer); err != nil {
			return nil, fmt.Errorf("failed to parse logs data: %w", err)
		}

		requests = append(requests, request)
	}

	return requests, nil
}

func (l *LogsManager) acknowledgeWrite(requests []Request, err error) {
	for _, req := range requests {
		req.doneStatus <- err
		close(req.doneStatus)
	}
}
