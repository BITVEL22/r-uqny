package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefault(t *testing.T) {
	cfg := Default()

	if cfg.ListenAddress != DefaultListenAddress {
		t.Fatalf(
			"expected listen address %q, got %q",
			DefaultListenAddress,
			cfg.ListenAddress,
		)
	}

	if cfg.DataDirectory != DefaultDataDirectory {
		t.Fatalf(
			"expected data directory %q, got %q",
			DefaultDataDirectory,
			cfg.DataDirectory,
		)
	}
}

func TestDefaultIsValid(t *testing.T) {
	cfg := Default()

	if err := cfg.Validate(); err != nil {
		t.Fatalf("default config should be valid: %v", err)
	}
}

func TestEmptyListenAddress(t *testing.T) {
	cfg := Default()
	cfg.ListenAddress = ""

	if err := cfg.Validate(); err == nil {
		t.Fatal("expected validation error for empty listen address")
	}
}

func TestEmptyDataDirectory(t *testing.T) {
	cfg := Default()
	cfg.DataDirectory = ""

	if err := cfg.Validate(); err == nil {
		t.Fatal("expected validation error for empty data directory")
	}
}

func TestSaveAndLoad(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")

	expected := Config{
		NodeID:        "test-node-01",
		ListenAddress: "127.0.0.1:9100",
		DataDirectory: "./test-data",
	}

	if err := Save(path, expected); err != nil {
		t.Fatalf("failed to save config: %v", err)
	}

	actual, err := Load(path)
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}

	if actual != expected {
		t.Fatalf("loaded config differs from saved config:\nexpected: %+v\ngot: %+v", expected, actual)
	}
}

func TestLoadInvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "invalid.json")

	if err := os.WriteFile(path, []byte("{invalid"), 0644); err != nil {
		t.Fatalf("failed to create invalid config: %v", err)
	}

	if _, err := Load(path); err == nil {
		t.Fatal("expected error when loading invalid JSON")
	}
}
