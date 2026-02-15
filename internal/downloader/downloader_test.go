package downloader

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestDownloadSingleFile(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Length", "12")
		w.Write([]byte("test content"))
	}))
	defer server.Close()

	cfg := &Config{
		OutputDir:   t.TempDir(),
		Concurrency: 1,
		Timeout:     10,
		Retry:       1,
		UserAgent:   "test",
	}

	d := New(cfg)

	result := d.Download(context.Background(), server.URL, nil)

	if result.Error != nil {
		t.Errorf("expected no error, got: %v", result.Error)
	}
	if result.FilePath == "" {
		t.Error("expected file path to be set")
	}
	if result.URL != server.URL {
		t.Errorf("expected URL %s, got %s", server.URL, result.URL)
	}
}

func TestDownloadMultipleFiles(t *testing.T) {
	callCount := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		callCount++
		w.Header().Set("Content-Length", "12")
		w.Write([]byte("test content"))
	}))
	defer server.Close()

	urls := []string{
		server.URL + "/file1.txt",
		server.URL + "/file2.txt",
		server.URL + "/file3.txt",
	}

	cfg := &Config{
		OutputDir:   t.TempDir(),
		Concurrency: 1,
		Timeout:     10,
		Retry:       1,
		UserAgent:   "test",
	}

	d := New(cfg)

	var results []*Result
	for _, url := range urls {
		result := d.Download(context.Background(), url, nil)
		results = append(results, result)
	}

	if len(results) != len(urls) {
		t.Errorf("expected %d results, got %d", len(urls), len(results))
	}
	if callCount != len(urls) {
		t.Errorf("expected %d calls, got %d", len(urls), callCount)
	}

	for i, result := range results {
		if result.Error != nil {
			t.Errorf("result %d: unexpected error: %v", i, result.Error)
		}
	}
}

func TestDownloadInvalidURL(t *testing.T) {
	tests := []struct {
		name string
		url  string
	}{
		{name: "empty URL", url: ""},
		{name: "invalid scheme", url: "ftp://example.com/file"},
		{name: "malformed URL", url: "://invalid"},
	}

	cfg := &Config{
		OutputDir:   t.TempDir(),
		Concurrency: 1,
		Timeout:     10,
		Retry:       1,
		UserAgent:   "test",
	}

	d := New(cfg)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := d.Download(context.Background(), tt.url, nil)
			if result.Error == nil {
				t.Errorf("expected error for URL %q", tt.url)
			}
		})
	}
}

func TestDownloadTimeout(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		w.Write([]byte("delayed response"))
	}))
	defer server.Close()

	cfg := &Config{
		OutputDir:   t.TempDir(),
		Concurrency: 1,
		Timeout:     5,
		Retry:       0,
		UserAgent:   "test",
	}

	d := New(cfg)

	result := d.Download(context.Background(), server.URL, nil)

	if result.Error != nil {
		t.Errorf("expected download to succeed with 5s timeout for 200ms request, got error: %v", result.Error)
	}
	if result.FilePath == "" {
		t.Error("expected file path to be set")
	}
}

func TestDownloadRetry(t *testing.T) {
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts < 3 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write([]byte("success after retry"))
	}))
	defer server.Close()

	cfg := &Config{
		OutputDir:   t.TempDir(),
		Concurrency: 1,
		Timeout:     10,
		Retry:       3,
		UserAgent:   "test",
	}

	d := New(cfg)

	result := d.Download(context.Background(), server.URL, nil)

	if result.Error != nil {
		t.Errorf("expected success after retries, got error: %v", result.Error)
	}
	if attempts < 3 {
		t.Errorf("expected at least 3 attempts, got %d", attempts)
	}
}

func TestDownloadHTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	cfg := &Config{
		OutputDir:   t.TempDir(),
		Concurrency: 1,
		Timeout:     10,
		Retry:       0,
		UserAgent:   "test",
	}

	d := New(cfg)

	result := d.Download(context.Background(), server.URL, nil)

	if result.Error == nil {
		t.Error("expected HTTP error")
	}
	expected := fmt.Sprintf("HTTP error: %d", http.StatusNotFound)
	if result.Error.Error() != expected {
		t.Errorf("expected error %q, got %q", expected, result.Error.Error())
	}
}
