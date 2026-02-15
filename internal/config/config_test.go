package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestConfigDefaults(t *testing.T) {
	tests := []struct {
		name     string
		field    string
		expected interface{}
	}{
		{
			name:     "default output dir",
			field:    "OutputDir",
			expected: ".",
		},
		{
			name:     "default concurrency",
			field:    "Concurrency",
			expected: 3,
		},
		{
			name:     "default timeout",
			field:    "Timeout",
			expected: 30,
		},
		{
			name:     "default retry",
			field:    "Retry",
			expected: 3,
		},
		{
			name:     "default user agent",
			field:    "UserAgent",
			expected: "GoGetIt/0.1.0",
		},
	}

	cfg := Get()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			switch tt.field {
			case "OutputDir":
				if cfg.OutputDir != tt.expected.(string) {
					t.Errorf("OutputDir = %q, expected %q", cfg.OutputDir, tt.expected)
				}
			case "Concurrency":
				if cfg.Concurrency != tt.expected.(int) {
					t.Errorf("Concurrency = %d, expected %d", cfg.Concurrency, tt.expected)
				}
			case "Timeout":
				if cfg.Timeout != tt.expected.(int) {
					t.Errorf("Timeout = %d, expected %d", cfg.Timeout, tt.expected)
				}
			case "Retry":
				if cfg.Retry != tt.expected.(int) {
					t.Errorf("Retry = %d, expected %d", cfg.Retry, tt.expected)
				}
			case "UserAgent":
				if cfg.UserAgent != tt.expected.(string) {
					t.Errorf("UserAgent = %q, expected %q", cfg.UserAgent, tt.expected)
				}
			}
		})
	}
}

func TestConfigNegativeConcurrency(t *testing.T) {
	tests := []struct {
		name           string
		concurrency    int
		expectedResult int
	}{
		{
			name:           "zero concurrency",
			concurrency:    0,
			expectedResult: 3,
		},
		{
			name:           "negative concurrency",
			concurrency:    -5,
			expectedResult: -5,
		},
		{
			name:           "positive concurrency",
			concurrency:    5,
			expectedResult: 5,
		},
		{
			name:           "large concurrency",
			concurrency:    100,
			expectedResult: 100,
		},
		{
			name:           "one concurrency",
			concurrency:    1,
			expectedResult: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			configFile := filepath.Join(tmpDir, ".gogetit.yaml")

			content := "concurrency: " + itoa(tt.concurrency)
			if err := os.WriteFile(configFile, []byte(content), 0644); err != nil {
				t.Fatalf("failed to write config file: %v", err)
			}

			originalConfig := currentConfig
			currentConfig = &Config{}
			defer func() { currentConfig = originalConfig }()

			SetConfigFile(configFile)
			Load()

			if currentConfig.Concurrency != tt.expectedResult {
				t.Errorf("concurrency = %d, expected %d after Load()",
					currentConfig.Concurrency, tt.expectedResult)
			}
		})
	}
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	neg := false
	if n < 0 {
		neg = true
		n = -n
	}
	var digits []byte
	for n > 0 {
		digits = append([]byte{byte('0' + n%10)}, digits...)
		n /= 10
	}
	if neg {
		digits = append([]byte{'-'}, digits...)
	}
	return string(digits)
}
