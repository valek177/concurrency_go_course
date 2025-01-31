package config

import (
	"testing"

	"concurrency_go_course/pkg/logger"

	"github.com/stretchr/testify/assert"
)

func TestNewServerConfig(t *testing.T) {
	t.Parallel()

	logger.MockLogger()

	defaultCfg := DefaultConfig()

	tests := []struct {
		name          string
		cfgPath       string
		result        *Config
		expectedError error
	}{
		{
			name:          "unable to read config, apply default config",
			cfgPath:       "",
			result:        defaultCfg,
			expectedError: nil,
		},
		{
			name:          "unable to read config, apply default config",
			cfgPath:       "some_incorrect_path",
			result:        defaultCfg,
			expectedError: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg, err := NewConfig(tt.cfgPath)
			assert.Nil(t, err)
			assert.Equal(t, tt.result, cfg)
		})
	}
}
