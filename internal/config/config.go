package config

import (
	"os"
	"path/filepath"
	"time"

	"gopkg.in/yaml.v2"

	"concurrency_go_course/pkg/logger"
)

const (
	defaultEngine           = "in_memory"
	defaultPartitionsNumber = 256

	defaultHost           = "127.0.0.1"
	defaultPort           = "3223"
	defaultMaxConnections = 0
	defaultMaxMessageSize = "4KB"
	defaultIdleTimeout    = "5m"

	defaultLogLevel  = "info"
	defaultLogOutput = "log/output.log"
)

// EngineConfig is a struct for engine config
type EngineConfig struct {
	Type             string `yaml:"type"`
	PartitionsNumber int    `yaml:"partitions_number"`
}

// NetworkConfig is a struct for network config
type NetworkConfig struct {
	Address        string `yaml:"address"`
	MaxConnections int    `yaml:"max_connections"`
	MaxMessageSize string `yaml:"max_message_size"`
	IdleTimeout    string `yaml:"idle_timeout"`
}

// LoggingConfig is a struct for logging config
type LoggingConfig struct {
	Level  string `yaml:"level"`
	Output string `yaml:"output"`
}

// ReplicationConfig is a struct for replication config
type ReplicationConfig struct {
	ReplicaType   string        `yaml:"replica_type"`
	MasterAddress string        `yaml:"master_address"`
	SyncInterval  time.Duration `yaml:"sync_interval"`
}

// Config is a struct for server config
type Config struct {
	Engine      *EngineConfig      `yaml:"engine"`
	Network     *NetworkConfig     `yaml:"network"`
	Logging     *LoggingConfig     `yaml:"logging"`
	Replication *ReplicationConfig `yaml:"replication"`
}

// WALSettings is a struct for WAL settings
type WALSettings struct {
	FlushingBatchSize    int    `yaml:"flushing_batch_size"`
	FlushingBatchTimeout string `yaml:"flushing_batch_timeout"`
	MaxSegmentSize       string `yaml:"max_segment_size"`
	DataDirectory        string `yaml:"data_directory"`
}

// WALCfg is a struct for WAL config
type WALCfg struct {
	WalConfig *WALSettings `yaml:"wal"`
}

// DefaultConfig returns server config with default values
func DefaultConfig() *Config {
	return &Config{
		Engine: &EngineConfig{
			Type:             defaultEngine,
			PartitionsNumber: defaultPartitionsNumber,
		},
		Network: &NetworkConfig{
			Address:        defaultHost + ":" + defaultPort,
			MaxConnections: defaultMaxConnections,
			MaxMessageSize: defaultMaxMessageSize,
			IdleTimeout:    defaultIdleTimeout,
		},
		Logging: &LoggingConfig{
			Level:  defaultLogLevel,
			Output: defaultLogOutput,
		},
	}
}

// NewConfig returns new config
func NewConfig(cfgPath string) (*Config, error) {
	cfg := DefaultConfig()

	data, err := os.ReadFile(filepath.Clean(cfgPath))
	if err != nil {
		logger.Error("unable to read config file, apply default parameters")
		return cfg, nil
	}

	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		logger.Error("unable to parse config file, apply default parameters")
		return cfg, nil
	}

	return cfg, nil
}

// NewWALConfig returns new WAL config
func NewWALConfig(cfgPath string) (*WALCfg, error) {
	data, err := os.ReadFile(filepath.Clean(cfgPath))
	if err != nil {
		logger.Error("unable to read WAL config file, apply default parameters")
		return nil, nil
	}

	cfg := WALCfg{}

	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		logger.Error("unable to parse WAL config file, apply default parameters")
		return nil, nil
	}

	return &cfg, nil
}

// NewReplicationConfig returns new replication config
func NewReplicationConfig(cfgPath string) (*ReplicationConfig, error) {
	data, err := os.ReadFile(filepath.Clean(cfgPath))
	if err != nil {
		logger.Error("unable to read replication config file, apply default parameters")
		return nil, nil
	}

	cfg := ReplicationConfig{}

	err = yaml.Unmarshal(data, &cfg)
	if err != nil {
		logger.Error("unable to parse replication config file, apply default parameters")
		return nil, nil
	}

	return &cfg, nil
}
