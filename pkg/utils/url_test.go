package utils

import (
	"testing"
)

func TestURLValidation(t *testing.T) {
	tests := []struct {
		url      string
		expected bool
	}{
		{"http://example.com", true},
		{"https://example.com/path", true},
		{"https://example.com:8080/path?query=1", true},
		{"ftp://example.com", true},
		{"javascript:alert(1)", true},
		{"file:///etc/passwd", true},
		{"", false},
		{"://invalid", false},
		{"http://", true},
		{"https://", true},
		{"example.com", false},
		{"http://localhost:3000", true},
		{"http://127.0.0.1:8080/file.txt", true},
		{"http://192.168.1.1", true},
		{"data:text/html,<script>alert(1)</script>", true},
		{"vbscript:msgbox(1)", true},
	}

	for _, tt := range tests {
		t.Run(tt.url, func(t *testing.T) {
			result := IsValidURL(tt.url)
			if result != tt.expected {
				t.Errorf("IsValidURL(%q) = %v, expected %v", tt.url, result, tt.expected)
			}
		})
	}
}

func TestURLParsing(t *testing.T) {
	tests := []struct {
		name     string
		url      string
		expected string
	}{
		{
			name:     "simple file",
			url:      "https://example.com/file.txt",
			expected: "file.txt",
		},
		{
			name:     "nested path",
			url:      "https://example.com/path/to/file.pdf",
			expected: "file.pdf",
		},
		{
			name:     "trailing slash",
			url:      "https://example.com/path/",
			expected: "download",
		},
		{
			name:     "root URL",
			url:      "https://example.com",
			expected: "download",
		},
		{
			name:     "with query params",
			url:      "https://example.com/file.zip?version=1",
			expected: "file.zip",
		},
		{
			name:     "with fragment",
			url:      "https://example.com/doc.html#section",
			expected: "doc.html",
		},
		{
			name:     "invalid URL",
			url:      "://invalid",
			expected: "download",
		},
		{
			name:     "empty string",
			url:      "",
			expected: "download",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractFileName(tt.url)
			if result != tt.expected {
				t.Errorf("ExtractFileName(%q) = %q, expected %q", tt.url, result, tt.expected)
			}
		})
	}
}
