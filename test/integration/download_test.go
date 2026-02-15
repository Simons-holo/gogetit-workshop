package main

import (
	"context"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestDownloadSingleFile(t *testing.T) {
	content := []byte("test file content for download")
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write(content)
	}))
	defer server.Close()

	tmpDir, err := os.MkdirTemp("", "download-test")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)

	client := &http.Client{Timeout: 30 * time.Second}
	ctx := context.Background()

	req, err := http.NewRequestWithContext(ctx, "GET", server.URL, nil)
	if err != nil {
		t.Fatal(err)
	}

	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Errorf("expected status 200, got %d", resp.StatusCode)
	}

	file, err := os.CreateTemp(tmpDir, "downloaded-*")
	if err != nil {
		t.Fatal(err)
	}

	buf := make([]byte, len(content))
	n, _ := resp.Body.Read(buf)
	file.Write(buf[:n])
	file.Close()

	data, err := os.ReadFile(file.Name())
	if err != nil {
		t.Fatal(err)
	}

	expected := 30
	if len(data) != expected {
		t.Errorf("expected %d bytes, got %d", expected, len(data))
	}
}

func TestDownloadMultipleFiles(t *testing.T) {
	files := map[string][]byte{
		"/file1.txt": []byte("file 1 content"),
		"/file2.txt": []byte("file 2 content"),
		"/file3.txt": []byte("file 3 content"),
	}

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if content, ok := files[r.URL.Path]; ok {
			w.WriteHeader(http.StatusOK)
			w.Write(content)
			return
		}
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	for path, expected := range files {
		resp, err := http.Get(server.URL + path)
		if err != nil {
			t.Errorf("failed to fetch %s: %v", path, err)
			continue
		}

		if resp.StatusCode != http.StatusOK {
			t.Errorf("expected 200 for %s, got %d", path, resp.StatusCode)
			resp.Body.Close()
			continue
		}

		buf := make([]byte, len(expected)+10)
		n, _ := resp.Body.Read(buf)
		resp.Body.Close()

		if n != len(expected) {
			t.Errorf("expected %d bytes for %s, got %d", len(expected), path, n)
		}
	}
}

func TestDownloadWithRetry(t *testing.T) {
	attempts := 0
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempts++
		if attempts < 3 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("success after retry"))
	}))
	defer server.Close()

	var lastErr error
	var resp *http.Response

	for i := 0; i < 5; i++ {
		resp, lastErr = http.Get(server.URL)
		if lastErr == nil && resp.StatusCode == http.StatusOK {
			break
		}
		if resp != nil {
			resp.Body.Close()
		}
		time.Sleep(10 * time.Millisecond)
	}

	if lastErr != nil {
		t.Errorf("download failed after retries: %v", lastErr)
	}

	if attempts != 3 {
		t.Errorf("expected 3 attempts, got %d", attempts)
	}
}

func TestDownloadCancellation(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, "GET", server.URL, nil)
	if err != nil {
		t.Fatal(err)
	}

	_, err = http.DefaultClient.Do(req)
	if err == nil {
		t.Error("expected timeout error, got nil")
	}
}
