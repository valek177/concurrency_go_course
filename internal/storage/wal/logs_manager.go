package wal

import (
	"bytes"
	"errors"
	"fmt"

	fs "concurrency_go_course/internal/filesystem"
	"concurrency_go_course/pkg/logger"
)

// LogsManager is interface for manager
type LogsManager interface {
	Write(requests []Request)
	ReadAll() ([]Request, error)
}

// LogsManager is a struct for logs manager
type logsmanager struct {
	segment fs.Segment
}

// NewLogsManager returns new logs manager
func NewLogsManager(segment fs.Segment) (LogsManager, error) {
	if segment == nil {
		return nil, errors.New("segment is invalid")
	}

	return &logsmanager{segment: segment}, nil
}

// Write writes requests
func (l *logsmanager) Write(requests []Request) {
	var buffer bytes.Buffer
	for _, req := range requests {
		if err := req.Encode(&buffer); err != nil {
			logger.ErrorWithMsg("failed to encode requests", err)
			l.acknowledgeWrite(requests, err)
			return
		}
	}

	err := l.segment.Write(buffer.Bytes())
	if err != nil {
		logger.ErrorWithMsg("failed to write request data:", err)
	}

	l.acknowledgeWrite(requests, err)
}

// ReadAll reads all requests
func (l *logsmanager) ReadAll() ([]Request, error) {
	segmentsData, err := l.segment.ReadAll()
	if err != nil {
		return nil, fmt.Errorf("failed to read segments: %w", err)
	}

	var requests []Request
	for _, data := range segmentsData {
		requests, err = l.readSegment(requests, data)
		if err != nil {
			return nil, fmt.Errorf("failed to read segments: %w", err)
		}
	}

	logger.Debug("WAL requests was readed")

	return requests, nil
}

func (l *logsmanager) readSegment(requests []Request, data []byte) ([]Request, error) {
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

func (l *logsmanager) acknowledgeWrite(requests []Request, err error) {
	for _, req := range requests {
		req.doneStatus <- err
		close(req.doneStatus)
	}
}
