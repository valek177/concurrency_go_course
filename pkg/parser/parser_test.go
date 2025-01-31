package parser

import (
	"testing"

	"concurrency_go_course/pkg/logger"

	"github.com/stretchr/testify/assert"
)

func TestParseSize(t *testing.T) {
	t.Parallel()

	logger.MockLogger()

	tests := []struct {
		name       string
		msgSizeStr string
		result     int
	}{
		{
			name:       "3 bytes is 3",
			msgSizeStr: "3B",
			result:     3,
		},
		{
			name:       "5KB is 5 * 1024",
			msgSizeStr: "5KB",
			result:     5 * 1024,
		},
		{
			name:       "10KB is 10 * 1024",
			msgSizeStr: "10KB",
			result:     10 * 1024,
		},
		{
			name:       "10kb is 10 * 1024",
			msgSizeStr: "10kb",
			result:     10 * 1024,
		},
		{
			name:       "10MB is 10 * 1024 * 1024",
			msgSizeStr: "10MB",
			result:     10 * 1024 * 1024,
		},
		{
			name:       "incorrect value",
			msgSizeStr: "1nn",
			result:     0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			res, _ := ParseSize(tt.msgSizeStr)
			assert.Equal(t, tt.result, res)
		})
	}
}
