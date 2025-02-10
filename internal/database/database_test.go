package database

import (
	"fmt"
	"testing"

	"concurrency_go_course/internal/compute"
	"concurrency_go_course/internal/config"
	"concurrency_go_course/internal/storage"
	"concurrency_go_course/internal/storage/mock"
	"concurrency_go_course/pkg/logger"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestServiceHandleNeg(t *testing.T) {
	t.Parallel()

	logger.MockLogger()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{
		Engine: &config.EngineConfig{
			Type:             "in_memory",
			PartitionsNumber: 256,
		},
		Network: &config.NetworkConfig{
			Address:        "127.0.0.1:7777",
			MaxConnections: 100,
			MaxMessageSize: "4KB",
			IdleTimeout:    "5m",
		},
		Logging: &config.LoggingConfig{
			Level:  "info",
			Output: "log/output.log",
		},
		Replication: &config.ReplicationConfig{},
	}

	mockEngine := mock.NewMockEngine(ctrl)

	storage, err := storage.New(mockEngine, nil, cfg, "", nil)
	if err != nil {
		t.Errorf("unable to create storage")
	}

	parser := compute.NewRequestParser()
	compute := compute.NewCompute(parser)

	service := NewDatabase(storage, compute)

	tests := map[string]struct {
		in   string
		res  string
		err  error
		exec func()
	}{
		"empty request": {
			in:   "",
			res:  "",
			err:  fmt.Errorf("invalid query length (0)"),
			exec: func() {},
		},
		"GET: no value": {
			in:  "GET unknown",
			res: "",
			exec: func() {
				mockEngine.EXPECT().Get("unknown").Return("", false)
			},
			err: fmt.Errorf("value not found"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			test.exec()
			res, err := service.Handle(test.in)

			assert.Equal(t, err, test.err)
			assert.Equal(t, res, test.res)
		})
	}
}

func TestServiceHandlePos(t *testing.T) {
	t.Parallel()

	logger.MockLogger()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	cfg := &config.Config{
		Engine: &config.EngineConfig{
			Type:             "in_memory",
			PartitionsNumber: 256,
		},
		Network: &config.NetworkConfig{
			Address:        "127.0.0.1:7777",
			MaxConnections: 100,
			MaxMessageSize: "4KB",
			IdleTimeout:    "5m",
		},
		Logging: &config.LoggingConfig{
			Level:  "info",
			Output: "log/output.log",
		},
		Replication: &config.ReplicationConfig{
			ReplicaType: "master",
		},
	}

	mockEngine := mock.NewMockEngine(ctrl)

	storage, err := storage.New(mockEngine, nil, cfg, "master", nil)
	if err != nil {
		t.Errorf("unable to create storage")
	}

	parser := compute.NewRequestParser()
	compute := compute.NewCompute(parser)

	service := NewDatabase(storage, compute)

	tests := map[string]struct {
		in   string
		res  string
		err  error
		exec func()
	}{
		"GET: correct result": {
			in:  "GET key1",
			res: "value1",
			err: nil,
			exec: func() {
				mockEngine.EXPECT().Get("key1").Return("value1", true)
			},
		},
		"SET: correct result": {
			in:  "SET key1 value1",
			res: "OK",
			exec: func() {
				mockEngine.EXPECT().Set("key1", "value1").Return()
			},
			err: nil,
		},
		"DEL: correct result": {
			in:  "DEL key1",
			res: "OK",
			exec: func() {
				mockEngine.EXPECT().Delete("key1").Return()
			},
			err: nil,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			test.exec()
			res, err := service.Handle(test.in)

			assert.Equal(t, err, test.err)
			assert.Equal(t, res, test.res)
		})
	}
}
