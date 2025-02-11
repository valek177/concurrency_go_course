package replication

import (
	"bytes"
	"encoding/gob"
	"fmt"
)

// SlaveRequest is a struct for request from slave node
type SlaveRequest struct {
	LastSegmentName string
}

// NewRequest returns new slave request
func NewRequest(lastSegmentName string) SlaveRequest {
	return SlaveRequest{
		LastSegmentName: lastSegmentName,
	}
}

// MasterResponse is a struct for response from master node
type MasterResponse struct {
	Succeed     bool
	SegmentName string
	SegmentData []byte
}

// NewMasterResponse returns new master response
func NewMasterResponse(succeed bool, segmentName string, segmentData []byte) MasterResponse {
	return MasterResponse{
		Succeed:     succeed,
		SegmentName: segmentName,
		SegmentData: segmentData,
	}
}

// EncodeResponse encodes master response
func EncodeResponse(response *MasterResponse) ([]byte, error) {
	if response == nil {
		return nil, fmt.Errorf("failed to encode object: nil object")
	}
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(response); err != nil {
		return nil, fmt.Errorf("failed to encode object: %w", err)
	}
	return buffer.Bytes(), nil
}

// EncodeSlaveRequest encodes slave request
func EncodeSlaveRequest(request *SlaveRequest) ([]byte, error) {
	if request == nil {
		return nil, fmt.Errorf("failed to encode object: nil object")
	}
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(request); err != nil {
		return nil, fmt.Errorf("failed to encode object: %w", err)
	}
	return buffer.Bytes(), nil
}

// DecodeResponse decodes master response
func DecodeResponse(response *MasterResponse, data []byte) error {
	if response == nil {
		return fmt.Errorf("failed to decode object: nil object")
	}
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	if err := decoder.Decode(response); err != nil {
		return fmt.Errorf("failed to decode object: %w", err)
	}
	return nil
}

// DecodeSlaveRequest decodes slave request
func DecodeSlaveRequest(request *SlaveRequest, data []byte) error {
	if request == nil {
		return fmt.Errorf("failed to decode object: nil object")
	}
	buffer := bytes.NewBuffer(data)
	decoder := gob.NewDecoder(buffer)
	if err := decoder.Decode(request); err != nil {
		return fmt.Errorf("failed to decode object: %w", err)
	}
	return nil
}
