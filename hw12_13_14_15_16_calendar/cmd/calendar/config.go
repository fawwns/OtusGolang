package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server  ServerConf    `yaml:"server"`
	Logger  LoggerConf    `yaml:"logger"`
	Storage StorageConfig `yaml:"storage"`
}

type ServerConf struct {
	Host string `yaml:"host"`
	Port int    `yaml:"port"`
}

type LoggerConf struct {
	Level string `yaml:"level"`
}
type StorageConfig struct {
	Type string `yaml:"type"` // "sql" или "inmemory".
	DSN  string `yaml:"dsn"`  // Используется только если Type="sql".
	Code string `yaml:"code"` // Опционально, если требуется.
}
type SQLConfig struct {
	DSN string `yaml:"dsn"`
}

type InMemoryConf struct {
	Enabled bool `yaml:"enabled"`
}

func LoadConfig(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("cannot open config file: %w", err)
	}
	defer f.Close()

	var cfg Config
	if err := yaml.NewDecoder(f).Decode(&cfg); err != nil {
		return nil, fmt.Errorf("cannot decode config: %w", err)
	}

	return &cfg, nil
}
