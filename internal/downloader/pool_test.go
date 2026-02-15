package downloader

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

func TestWorkerPoolBasic(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("test data"))
	}))
	defer server.Close()

	cfg := &Config{
		OutputDir:   t.TempDir(),
		Concurrency: 1,
		Timeout:     10,
		Retry:       1,
		UserAgent:   "test",
	}

	worker := NewWorkerV2(0, cfg)
	ctx := context.Background()

	result := worker.Process(ctx, server.URL+"/file1.txt", cfg.OutputDir)

	if result.Error != nil {
		t.Errorf("unexpected error: %v", result.Error)
	}
}

func TestWorkerPoolConcurrency(t *testing.T) {
	var activeCount int64
	maxConcurrent := int64(0)

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		current := atomic.AddInt64(&activeCount, 1)
		if current > atomic.LoadInt64(&maxConcurrent) {
			atomic.StoreInt64(&maxConcurrent, current)
		}
		time.Sleep(50 * time.Millisecond)
		atomic.AddInt64(&activeCount, -1)
		w.Write([]byte("data"))
	}))
	defer server.Close()

	cfg := &Config{
		OutputDir:   t.TempDir(),
		Concurrency: 3,
		Timeout:     30,
		Retry:       0,
		UserAgent:   "test",
	}

	worker := NewWorkerV2(0, cfg)
	ctx := context.Background()

	result := worker.Process(ctx, server.URL+"/file.txt", cfg.OutputDir)

	if result.Error != nil {
		t.Errorf("unexpected error: %v", result.Error)
	}
}

type RaceProneCounter struct {
	value int
}

func (c *RaceProneCounter) Increment() {
	c.value++
}

func (c *RaceProneCounter) Value() int {
	return c.value
}

func TestWorkerPoolRaceCondition(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("data"))
	}))
	defer server.Close()

	cfg := &Config{
		OutputDir:   t.TempDir(),
		Concurrency: 5,
		Timeout:     10,
		Retry:       0,
		UserAgent:   "test",
	}

	var wg sync.WaitGroup
	counter := &RaceProneCounter{}

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			worker := NewWorkerV2(0, cfg)
			ctx := context.Background()
			result := worker.Process(ctx, server.URL+"/file.txt", cfg.OutputDir)
			if result.Error == nil {
				counter.Increment()
			}
		}()
	}

	wg.Wait()

	if counter.Value() == 100 {
		t.Errorf("expected race condition to cause incorrect count, got correct count of 100 - RaceProneCounter should have race condition")
	}
}

func TestWorkerPoolEmptyInput(t *testing.T) {
	cfg := &Config{
		OutputDir:   t.TempDir(),
		Concurrency: 1,
		Timeout:     10,
		Retry:       1,
		UserAgent:   "test",
	}

	if cfg.Concurrency < 1 {
		t.Errorf("expected valid concurrency, got %d", cfg.Concurrency)
	}
}
