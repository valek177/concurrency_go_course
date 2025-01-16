package compute

import (
	"fmt"
	"testing"

	mstorage "concurrency_go_course/internal/storage/mock"
	"concurrency_go_course/pkg/logger"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestHandleCompute(t *testing.T) {
	t.Parallel()

	logger.MockLogger()

	ctrl := gomock.NewController(t)
	defer t.Cleanup(ctrl.Finish)

	mockStorage := mstorage.NewMockStorage(ctrl)

	parser := NewRequestParser()
	handler := NewCompute(mockStorage, parser)

	negTests := map[string]struct {
		in   string
		res  string
		exec func()
		err  error
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

	for name, test := range negTests {
		t.Run(name, func(t *testing.T) {
			test.exec()
			_, err := handler.Handle(test.in)
			assert.Equal(t, err, test.err)
		})
	}

	posTests := map[string]struct {
		in   string
		res  string
		exec func()
		err  error
	}{
		"GET: existing value": {
			in:  "GET key1",
			res: "value1",
			exec: func() {
				mockStorage.EXPECT().Get("key1").Return("value1", true)
			},
			err: nil,
		},
		"SET: new value": {
			in:  "SET key2 value2",
			res: "",
			exec: func() {
				mockStorage.EXPECT().Set("key2", "value2").Return()
			},
			err: nil,
		},
		"DEL: existing value": {
			in:  "DEL key1",
			res: "",
			exec: func() {
				mockStorage.EXPECT().Delete("key1").Return()
			},
			err: nil,
		},
	}

	for name, test := range posTests {
		t.Run(name, func(t *testing.T) {
			test.exec()
			res, _ := handler.Handle(test.in)
			assert.Equal(t, res, test.res)
		})
	}
}
