package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Config struct {
	OutputDir   string `yaml:"output_dir"`
	Concurrency int    `yaml:"concurrency"`
	Timeout     int    `yaml:"timeout"`
	Retry       int    `yaml:"retry"`
	UserAgent   string `yaml:"user_agent"`
}

var defaultConfig = &Config{
	OutputDir:   ".",
	Concurrency: 3,
	Timeout:     30,
	Retry:       3,
	UserAgent:   "GoGetIt/0.1.0",
}

var currentConfig = defaultConfig

func Get() *Config {
	return currentConfig
}

func SetConfigFile(path string) {
	currentConfig = &Config{}
	loadFromFile(path)
}

func Load() {
	home, err := os.UserHomeDir()
	if err != nil {
		return
	}

	configPath := filepath.Join(home, ".gogetit.yaml")
	if _, err := os.Stat(configPath); err == nil {
		loadFromFile(configPath)
	}

	if currentConfig.Concurrency <= 0 {
		currentConfig.Concurrency = 3
	}
	if currentConfig.Timeout <= 0 {
		currentConfig.Timeout = 30
	}
	if currentConfig.Retry <= 0 {
		currentConfig.Retry = 3
	}
}

func loadFromFile(path string) {
	data, err := os.ReadFile(path)
	if err != nil {
		return
	}

	yaml.Unmarshal(data, currentConfig)
}
