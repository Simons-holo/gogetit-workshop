package utils

import (
	"context"
	"net/http"
	"time"
)

type RetryConfig struct {
	MaxAttempts int
	Delay       time.Duration
}

var DefaultRetryConfig = RetryConfig{
	MaxAttempts: 3,
	Delay:       time.Second,
}

func DoWithRetry(ctx context.Context, client *http.Client, req *http.Request, cfg RetryConfig) (*http.Response, error) {
	var lastErr error

	for attempt := 0; attempt < cfg.MaxAttempts; attempt++ {
		if attempt > 0 {
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(cfg.Delay):
			}
		}

		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}

		if resp.StatusCode >= 500 {
			resp.Body.Close()
			lastErr = &HTTPError{StatusCode: resp.StatusCode, Message: "server error"}
			continue
		}

		return resp, nil
	}

	return nil, lastErr
}

type HTTPError struct {
	StatusCode int
	Message    string
}

func (e *HTTPError) Error() string {
	return e.Message
}
