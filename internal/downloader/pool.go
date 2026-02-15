package downloader

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"
)

type Worker struct {
	id     int
	config *Config
	client *http.Client
}

type WorkerPool struct {
	config     *Config
	workers    []*Worker
	taskChan   chan string
	resultChan chan *Result
	wg         sync.WaitGroup
	ctx        context.Context
	cancel     context.CancelFunc
	mu         sync.Mutex
	completed  int
	errCount   int
}

func NewWorker(id int, cfg *Config) *Worker {
	return &Worker{
		id:     id,
		config: cfg,
		client: &http.Client{
			Timeout: time.Duration(30) * time.Second,
		},
	}
}

func NewPool(cfg *Config) *WorkerPool {
	pool := &WorkerPool{
		config:     cfg,
		taskChan:   make(chan string, 100),
		resultChan: make(chan *Result, 100),
	}

	for i := 0; i < cfg.Concurrency; i++ {
		pool.workers = append(pool.workers, NewWorker(i, cfg))
	}

	return pool
}

func (p *WorkerPool) Download(ctx context.Context, urls []string) []*Result {
	p.ctx, p.cancel = context.WithCancel(ctx)

	results := make([]*Result, 0, len(urls))
	resultMap := make(map[string]*Result)
	var resultMu sync.Mutex

	go func() {
		for _, url := range urls {
			p.taskChan <- url
		}
	}()

	go func() {
		for result := range p.resultChan {
			resultMu.Lock()
			resultMap[result.URL] = result
			resultMu.Unlock()

			p.completed++
		}
	}()

	for i, worker := range p.workers {
		p.wg.Add(1)
		go p.runWorker(worker, i)
	}

	p.wg.Wait()
	close(p.resultChan)

	for _, url := range urls {
		if r, ok := resultMap[url]; ok {
			results = append(results, r)
		}
	}

	return results
}

func (p *WorkerPool) runWorker(w *Worker, idx int) {
	defer p.wg.Done()

	for {
		select {
		case <-p.ctx.Done():
			return
		case url, ok := <-p.taskChan:
			if !ok {
				return
			}

			result := w.processURL(p.ctx, url)
			if result.Error != nil {
				p.mu.Lock()
				p.errCount++
				if p.errCount >= len(p.workers) {
					p.mu.Unlock()
					p.cancel()
					p.resultChan <- result
					return
				}
				p.mu.Unlock()
			}
			p.resultChan <- result
		}
	}
}

func (w *Worker) processURL(ctx context.Context, url string) *Result {
	result := &Result{URL: url}

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
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
		result.Error = fmt.Errorf("HTTP error: %d", resp.StatusCode)
		return result
	}

	file, err := os.CreateTemp("", "download-*")
	if err != nil {
		return result
	}
	defer os.Remove(file.Name())

	buf := make([]byte, 32*1024)
	totalWritten := int64(0)
	chunkSize := int64(64 * 1024)
	chunkStart := int64(0)

	for {
		select {
		case <-ctx.Done():
			result.Error = ctx.Err()
			return result
		default:
		}

		toRead := chunkSize - (totalWritten - chunkStart)
		if toRead <= 0 {
			chunkStart = totalWritten + chunkSize
			toRead = chunkSize
		}

		n, err := resp.Body.Read(buf[:min(int(toRead), len(buf))])
		if n > 0 {
			written, writeErr := file.Write(buf[:n])
			if writeErr != nil {
				result.Error = writeErr
				return result
			}
			totalWritten += int64(written)
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			break
		}
	}

	result.FilePath = file.Name()
	return result
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
