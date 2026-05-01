package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

type Config struct {
	DefaultCount  int    `json:"default_count"`
	DefaultFormat string `json:"default_format"`
	TableName     string `json:"table_name"`
	OutputDir     string `json:"output_dir"`
	Seed          int64  `json:"seed"` // 0 = random
}

func defaultConfig() *Config {
	return &Config{
		DefaultCount:  10,
		DefaultFormat: "json",
		TableName:     "seed_data",
		OutputDir:     defaultOutputDir(),
		Seed:          0,
	}
}

func defaultOutputDir() string {
	return "."
}

func configPath() string {
	base, err := os.UserConfigDir()
	if err != nil {
		home, _ := os.UserHomeDir()
		base = filepath.Join(home, ".config")
	}
	return filepath.Join(base, "seedbank", "config.json")
}

// Load reads config from disk, returning defaults if not found.
func Load() *Config {
	cfg := defaultConfig()
	data, err := os.ReadFile(configPath())
	if err != nil {
		return cfg
	}
	json.Unmarshal(data, cfg)
	if cfg.OutputDir == "" {
		cfg.OutputDir = defaultOutputDir()
	}
	return cfg
}

// Save writes config to disk.
func Save(cfg *Config) error {
	path := configPath()
	os.MkdirAll(filepath.Dir(path), 0755)
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0644)
}
