package config

import (
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v2"

	"concurrency_go_course/pkg/logger"
)

const (
	configPathName = "CONFIG_PATH"

	defaultEngine = "in_memory"

	defaultHost           = "127.0.0.1"
	defaultPort           = "3223"
	defaultMaxConnections = 100
	defaultMaxMessageSize = "4KB"
	defaultIdleTimeout    = "5m"

	defaultLogLevel  = "info"
	defaultLogOutput = "log/output.log"
)

var blockSizes = map[string]int{
	"B":  1,
	"KB": 1024,
	"MB": 1024 * 1024,
	"GB": 1024 * 1024 * 1024,
}

type EngineConfig struct {
	Type string `yaml:"type"`
}

type NetworkConfig struct {
	Address        string `yaml:"address"`
	MaxConnections int    `yaml:"max_connections"`
	MaxMessageSize string `yaml:"max_message_size"`
	IdleTimeout    string `yaml:"idle_timeout"`
}

type LoggingConfig struct {
	Level  string `yaml:"level"`
	Output string `yaml:"output"`
}

type ServerConfig struct {
	Engine  *EngineConfig  `yaml:"engine"`
	Network *NetworkConfig `yaml:"network"`
	Logging *LoggingConfig `yaml:"logging"`
}

func DefaultServerConfig() *ServerConfig {
	return &ServerConfig{
		Engine: &EngineConfig{
			Type: defaultEngine,
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

// NewServerConfig returns new server config
func NewServerConfig() (*ServerConfig, error) {
	cfg := DefaultServerConfig()

	err := godotenv.Load("local.env")
	if err != nil {
		logger.Error("unable to load local.env, apply default parameters")
		return cfg, nil
	}

	cfgPath := os.Getenv(configPathName)
	if len(cfgPath) == 0 {
		logger.Error("config path not found, apply default parameters")
		return cfg, nil
	}

	data, err := os.ReadFile(cfgPath)
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

func ParseMaxMessageSize(msgSizeStr string) int {
	msgSizeStr = strings.TrimSpace(strings.ToUpper(msgSizeStr))

	re := regexp.MustCompile(`([0-9]+)(\w+)`)
	res := re.FindAllStringSubmatch(msgSizeStr, -1)

	for k, v := range blockSizes {
		if !strings.HasSuffix(msgSizeStr, k) {
			continue
		}

		if res[0][2] == k {
			size, err := strconv.Atoi(res[0][1])
			if err != nil {
				logger.Error("unable to parse max message size")
				continue
			}
			return size * v
		}
	}

	logger.Debug("unable to convert max message size")

	return 0
}
