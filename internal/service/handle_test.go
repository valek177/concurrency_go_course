package service

import (
	"fmt"
	"testing"

	"concurrency_go_course/internal/compute"
	mstorage "concurrency_go_course/internal/storage/mock"
	"concurrency_go_course/pkg/logger"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestServiceHandleNeg(t *testing.T) {
	t.Parallel()

	logger.MockLogger()

	ctrl := gomock.NewController(t)
	defer t.Cleanup(ctrl.Finish)

	mockStorage := mstorage.NewMockStorage(ctrl)

	parser := compute.NewRequestParser()
	compute := compute.NewCompute(parser)

	service := NewService(mockStorage, compute)

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
				mockStorage.EXPECT().Get("unknown").Return("", false)
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
	defer t.Cleanup(ctrl.Finish)

	mockStorage := mstorage.NewMockStorage(ctrl)

	parser := compute.NewRequestParser()
	compute := compute.NewCompute(parser)

	service := NewService(mockStorage, compute)

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
				mockStorage.EXPECT().Get("key1").Return("value1", true)
			},
		},
		"SET: correct result": {
			in:  "SET key1 value1",
			res: "",
			exec: func() {
				mockStorage.EXPECT().Set("key1", "value1").Return()
			},
			err: nil,
		},
		"DEL: correct result": {
			in:  "DEL key1",
			res: "",
			exec: func() {
				mockStorage.EXPECT().Delete("key1").Return()
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
