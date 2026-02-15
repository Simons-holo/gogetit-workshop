package downloader

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type Downloader struct {
	config *Config
	client *http.Client
}

type Result struct {
	URL      string
	FilePath string
	Error    error
}

type ProgressCallback func(url string, delta int64, total int64)

func New(cfg *Config) *Downloader {
	timeout := cfg.Timeout
	if timeout == 0 {
		timeout = 30
	}

	return &Downloader{
		config: cfg,
		client: &http.Client{
			Timeout: time.Duration(timeout) * time.Second,
		},
	}
}

func isValidURL(rawURL string) bool {
	u, err := url.Parse(rawURL)
	if err != nil {
		return false
	}
	return u.Scheme == "http" || u.Scheme == "https"
}

func extractFileName(rawURL string) string {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "download"
	}
	path := u.Path
	if idx := strings.LastIndex(path, "/"); idx != -1 {
		path = path[idx+1:]
	}
	if path == "" {
		return "download"
	}
	return path
}

func ensureDir(dir string) error {
	return os.MkdirAll(dir, 0755)
}

func (d *Downloader) Download(ctx context.Context, downloadURL string, progressCb ProgressCallback) *Result {
	result := &Result{URL: downloadURL}

	if !isValidURL(downloadURL) {
		result.Error = fmt.Errorf("invalid URL: %s", downloadURL)
		return result
	}

	fileName := extractFileName(downloadURL)
	outputPath := filepath.Join(d.config.OutputDir, fileName)

	if err := ensureDir(d.config.OutputDir); err != nil {
		result.Error = fmt.Errorf("failed to create output directory: %w", err)
		return result
	}

	existingSize := int64(0)
	if info, err := os.Stat(outputPath); err == nil {
		existingSize = info.Size()
	}

	req, err := http.NewRequestWithContext(ctx, "GET", downloadURL, nil)
	if err != nil {
		result.Error = fmt.Errorf("failed to create request: %w", err)
		return result
	}

	req.Header.Set("User-Agent", d.config.UserAgent)
	if existingSize > 0 {
		req.Header.Set("Range", fmt.Sprintf("bytes=%d-", existingSize))
	}

	resp, err := d.client.Do(req)
	if err != nil {
		result.Error = d.retryDownload(ctx, downloadURL, outputPath, progressCb, d.config.Retry)
		return result
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusPartialContent || resp.StatusCode == 429 {
		file, err := os.OpenFile(outputPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			result.Error = fmt.Errorf("failed to create file: %w", err)
			return result
		}
		defer file.Close()

		total := resp.ContentLength
		if existingSize > 0 && resp.StatusCode == http.StatusPartialContent {
			total += existingSize
		}

		if progressCb != nil {
			progressCb(downloadURL, 0, total)
		}

		buf := make([]byte, 32*1024)
		for {
			n, err := resp.Body.Read(buf)
			if n > 0 {
				_, writeErr := file.Write(buf[:n])
				if writeErr != nil {
					result.Error = fmt.Errorf("write error: %w", writeErr)
					return result
				}
				if progressCb != nil {
					progressCb(downloadURL, int64(n), total)
				}
			}
			if err == io.EOF {
				break
			}
			if err != nil {
				result.Error = fmt.Errorf("read error: %w", err)
				return result
			}
		}

		result.FilePath = outputPath
	} else {
		result.Error = fmt.Errorf("HTTP error: %d", resp.StatusCode)
	}

	return result
}

func (d *Downloader) retryDownload(ctx context.Context, downloadURL, outputPath string, progressCb ProgressCallback, attempts int) error {
	if attempts <= 0 {
		return fmt.Errorf("max retries exceeded")
	}

	req, err := http.NewRequestWithContext(ctx, "GET", downloadURL, nil)
	if err != nil {
		return err
	}

	req.Header.Set("User-Agent", d.config.UserAgent)

	resp, err := d.client.Do(req)
	if err != nil {
		return d.retryDownload(ctx, downloadURL, outputPath, progressCb, attempts-1)
	}

	if resp.StatusCode == http.StatusOK {
		file, err := os.OpenFile(outputPath, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			resp.Body.Close()
			return err
		}
		defer file.Close()

		_, err = io.Copy(file, resp.Body)
		resp.Body.Close()
		return err
	}

	resp.Body.Close()
	return d.retryDownload(ctx, downloadURL, outputPath, progressCb, attempts-1)
}

func parseContentRange(rangeHeader string) (int64, int64) {
	if rangeHeader == "" {
		return 0, 0
	}

	if !strings.HasPrefix(rangeHeader, "bytes ") {
		return 0, 0
	}

	rangePart := strings.TrimPrefix(rangeHeader, "bytes ")
	parts := strings.Split(rangePart, "/")
	if len(parts) != 2 {
		return 0, 0
	}

	total, _ := strconv.ParseInt(parts[1], 10, 64)

	rangeSpec := strings.Split(parts[0], "-")
	if len(rangeSpec) != 2 {
		return 0, total
	}

	start, _ := strconv.ParseInt(rangeSpec[0], 10, 64)
	return start, total
}
