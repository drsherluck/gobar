package main

import (
	"os"

	"github.com/pelletier/go-toml/v2"
)

type NetworkConfig struct {
	Interface string
}

type Config struct {
	Network NetworkConfig `toml:"omitempty"`
	Modules []string
}

func NewDefaultConfig() *Config {
	return &Config{
		Network: NetworkConfig{"enp4s0"},
		Modules: []string{"network", "volume", "cputemp", "memory", "wheater", "battery", "time"},
	}
}

func NewCustomConfig(path string) (*Config, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	err = toml.Unmarshal(file, &cfg)
	if err != nil {
		return nil, err
	}

	if len(cfg.Network.Interface) == 0 {
		cfg.Network.Interface = "enp4s0"
	}

	return &cfg, nil
}
