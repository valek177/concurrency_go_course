package replication

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

type SlaveRequest struct {
	LastSegmentName string
}

func NewRequest(lastSegmentName string) SlaveRequest {
	return SlaveRequest{
		LastSegmentName: lastSegmentName,
	}
}

type MasterResponse struct {
	Succeed     bool
	SegmentName string
	SegmentData []byte
}

func NewMasterResponse(succeed bool, segmentName string, segmentData []byte) MasterResponse {
	return MasterResponse{
		Succeed:     succeed,
		SegmentName: segmentName,
		SegmentData: segmentData,
	}
}

func EncodeResponse(response *MasterResponse) ([]byte, error) {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(response); err != nil {
		return nil, fmt.Errorf("failed to encode object: %w", err)
	}
	return buffer.Bytes(), nil
}

func EncodeSlaveRequest(request *SlaveRequest) ([]byte, error) {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(request); err != nil {
		return nil, fmt.Errorf("failed to encode object: %w", err)
	}
	return buffer.Bytes(), nil
}

func DecodeResponse(response *MasterResponse, data []byte) error {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	if err := decoder.Decode(response); err != nil {
		return fmt.Errorf("failed to decode object: %w", err)
	}
	return nil
}

func DecodeSlaveRequest(request *SlaveRequest, data []byte) error {
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	if err := decoder.Decode(request); err != nil {
		return fmt.Errorf("failed to decode object: %w", err)
	}
	return nil
}
