package config

import (
	"encoding/json"
	"errors"
	"os"
)

const (
	DefaultListenAddress = "127.0.0.1:9000"
	DefaultDataDirectory = "./data"
)

// Config contains the runtime configuration for a r/uqny node.
type Config struct {
	NodeID        string `json:"node_id"`
	ListenAddress string `json:"listen_address"`
	DataDirectory string `json:"data_directory"`
}

// Default returns a configuration with safe local defaults.
func Default() Config {
	return Config{
		NodeID:        "",
		ListenAddress: DefaultListenAddress,
		DataDirectory: DefaultDataDirectory,
	}
}

// Validate checks whether the configuration is valid.
func (c Config) Validate() error {
	if c.ListenAddress == "" {
		return errors.New("listen address cannot be empty")
	}

	if c.DataDirectory == "" {
		return errors.New("data directory cannot be empty")
	}

	return nil
}

// Load reads a configuration from a JSON file.
func Load(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}

	var cfg Config

	if err := json.Unmarshal(data, &cfg); err != nil {
		return Config{}, err
	}

	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

// Save writes the configuration to a JSON file.
func Save(path string, cfg Config) error {
	if err := cfg.Validate(); err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	data = append(data, '\n')

	return os.WriteFile(path, data, 0644)
}
