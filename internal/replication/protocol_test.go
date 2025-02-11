package replication

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"testing"

	"concurrency_go_course/pkg/logger"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_EncodeResponseOk(t *testing.T) {
	t.Parallel()

	logger.MockLogger()

	masterResponse := &MasterResponse{
		Succeed:     true,
		SegmentName: "wal_1.log",
		SegmentData: []byte{},
	}
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(masterResponse); err != nil {
		t.Errorf("expected nil error; got %+v", err)
	}
	expected := buffer.Bytes()

	tests := []struct {
		name     string
		response *MasterResponse
		expected []byte
	}{
		{
			name:     "correct encoding",
			response: masterResponse,
			expected: expected,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bytes, err := EncodeResponse(masterResponse)
			assert.Equal(t, tt.expected, bytes)
			require.NoError(t, err)
		})
	}
}

func Test_EncodeResponseFail(t *testing.T) {
	t.Parallel()

	logger.MockLogger()

	masterResponse := &MasterResponse{
		Succeed:     true,
		SegmentName: "wal_1.log",
		SegmentData: []byte{},
	}
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	if err := encoder.Encode(masterResponse); err != nil {
		t.Errorf("expected nil error; got %+v", err)
	}
	expected := buffer.Bytes()

	tests := []struct {
		name     string
		response *MasterResponse
		expected []byte
		err      error
	}{
		{
			name:     "error of encoding object",
			response: nil,
			expected: expected,
			err:      fmt.Errorf("failed to encode object: nil object"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bytes, err := EncodeResponse(nil)
			assert.Equal(t, tt.err, err)
			assert.Nil(t, bytes)
		})
	}
}

func Test_DecodeResponseOk(t *testing.T) {
	t.Parallel()

	logger.MockLogger()

	req := &MasterResponse{
		Succeed:     true,
		SegmentName: "wal_1.log",
		SegmentData: []byte{},
	}
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	if err := encoder.Encode(req); err != nil {
		t.Errorf("expected nil error; got %+v", err)
	}

	masterResponse := &MasterResponse{}

	buffer := bytes.NewBuffer(buf.Bytes())
	decoder := gob.NewDecoder(buffer)
	if err := decoder.Decode(masterResponse); err != nil {
		t.Errorf("expected nil error; got %+v", err)
	}

	masterResponse = &MasterResponse{}

	tests := []struct {
		name     string
		response *MasterResponse
		data     []byte
		err      error
	}{
		{
			name:     "correct decoding",
			response: masterResponse,
			data:     buf.Bytes(),
			err:      nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := DecodeResponse(tt.response, tt.data)
			require.NoError(t, err)
		})
	}
}

func Test_DecodeResponseFail(t *testing.T) {
	t.Parallel()

	logger.MockLogger()

	req := &MasterResponse{
		Succeed:     true,
		SegmentName: "wal_1.log",
		SegmentData: []byte{},
	}
	var buf bytes.Buffer
	encoder := gob.NewEncoder(&buf)
	if err := encoder.Encode(req); err != nil {
		t.Errorf("expected nil error; got %+v", err)
	}

	tests := []struct {
		name     string
		response *MasterResponse
		data     []byte
		err      error
	}{
		{
			name:     "error decoding",
			response: nil,
			data:     buf.Bytes(),
			err:      fmt.Errorf("failed to decode object: nil object"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := DecodeResponse(tt.response, tt.data)
			assert.Equal(t, tt.err, err)
		})
	}
}
