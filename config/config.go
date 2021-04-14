package config

import (
	"github.com/BurntSushi/toml"
	"github.com/feitianlove/golib/common/logger"
	golibVar "github.com/feitianlove/golib/config"
	"path/filepath"
)

type Config struct {
	Kafka      *KafkaConfig
	TailF      *TailFConfig
	CtrlLog    *logger.LogConf
	ConsoleLog *logger.LogConf
}
type KafkaConfig struct {
	Address string
	Topic   string
}

type TailFConfig struct {
	Dir string
}

func DefaultConfig() *Config {
	return &Config{
		Kafka: &KafkaConfig{
			Address: "127.0.0.1:9290",
			Topic:   "logtransfersAgent_one",
		},
		CtrlLog: &logger.LogConf{
			LogLevel:      "info",
			LogPath:       filepath.Join(golibVar.LaunchDir, "../log/ctrl.log"),
			LogReserveDay: 1,
			ReportCaller:  false,
		},
		ConsoleLog: &logger.LogConf{
			LogLevel:      "info",
			LogPath:       filepath.Join(golibVar.LaunchDir, "../log/console.log"),
			LogReserveDay: 1,
			ReportCaller:  false,
		},
		TailF: &TailFConfig{
			Dir: filepath.Join(golibVar.LaunchDir, "../log/access.log"),
		},
	}
}

func NewConfig(filePath string) (*Config, error) {
	cfg := DefaultConfig()
	if _, err := toml.DecodeFile(filePath, cfg); err != nil {
		return nil, err
	}
	return cfg, nil
}
