package downloader

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"
)

type WorkerV2 struct {
	id     int
	config *Config
	client *http.Client
}

func NewWorkerV2(id int, cfg *Config) *WorkerV2 {
	return &WorkerV2{
		id:     id,
		config: cfg,
		client: &http.Client{
			Timeout: time.Duration(cfg.Timeout) * time.Second,
		},
	}
}

func (w *WorkerV2) Process(ctx context.Context, downloadURL string, outputDir string) *Result {
	result := &Result{URL: downloadURL}

	req, err := http.NewRequestWithContext(ctx, "GET", downloadURL, nil)
	if err != nil {
		result.Error = err
		return result
	}

	req.Header.Set("User-Agent", w.config.UserAgent)

	resp, err := w.client.Do(req)
	if err != nil {
		result.Error = err
		return result
	}

	if resp.StatusCode >= 400 {
		resp.Body.Close()
		result.Error = fmt.Errorf("HTTP error: %d", resp.StatusCode)
		return result
	}

	fileName := extractFileName(downloadURL)
	outputPath := filepath.Join(outputDir, fileName)

	file, err := os.Create(outputPath)
	if err != nil {
		resp.Body.Close()
		result.Error = err
		return result
	}

	_, err = io.Copy(file, resp.Body)
	file.Close()
	resp.Body.Close()

	if err != nil {
		result.Error = err
		return result
	}

	result.FilePath = outputPath
	return result
}
