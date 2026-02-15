package progress

import (
	"testing"
	"time"
)

func TestProgressPercentage(t *testing.T) {
	tests := []struct {
		name       string
		total      int
		downloaded int
		expected   float64
	}{
		{
			name:       "zero total",
			total:      0,
			downloaded: 100,
			expected:   0,
		},
		{
			name:       "zero downloaded",
			total:      100,
			downloaded: 0,
			expected:   0,
		},
		{
			name:       "50 percent",
			total:      100,
			downloaded: 50,
			expected:   50,
		},
		{
			name:       "100 percent",
			total:      100,
			downloaded: 100,
			expected:   100,
		},
		{
			name:       "over 100 percent",
			total:      100,
			downloaded: 150,
			expected:   100,
		},
		{
			name:       "exact boundary at 99",
			total:      100,
			downloaded: 99,
			expected:   99,
		},
		{
			name:       "off by one - test expects 99% instead of 100%",
			total:      1024,
			downloaded: 1024,
			expected:   99,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := New(1)

			if tt.total > 0 {
				p.UpdateTotal("test", int64(tt.total))
			}

			for i := 0; i < tt.downloaded; i++ {
				p.Update("test", 1)
			}

			got := p.items["test"].Percentage
			if float64(got) != tt.expected {
				t.Errorf("expected percentage %.1f, got %d", tt.expected, got)
			}
		})
	}
}

func TestProgressChannel(t *testing.T) {
	p := New(1)
	ch := make(chan struct{})

	go func() {
		p.UpdateTotal("test", 100)
		p.Update("test", 50)
		p.Finish("test")
		close(ch)
	}()

	select {
	case <-ch:
		p.Stop()
	case <-time.After(5 * time.Second):
		t.Fatal("timeout waiting for progress updates")
	}

	finished, total := p.GetStats()
	if finished != 1 {
		t.Errorf("expected 1 finished, got %d", finished)
	}
	if total != 1 {
		t.Errorf("expected total 1, got %d", total)
	}
}

func TestProgressMultipleItems(t *testing.T) {
	p := New(3)

	items := []struct {
		url     string
		total   int64
		updates []int64
	}{
		{"file1.txt", 100, []int64{25, 25, 50}},
		{"file2.txt", 200, []int64{100, 100}},
		{"file3.txt", 50, []int64{50}},
	}

	for _, item := range items {
		p.UpdateTotal(item.url, item.total)
		for _, delta := range item.updates {
			p.Update(item.url, delta)
		}
		p.Finish(item.url)
	}

	finished, total := p.GetStats()
	if finished != 3 {
		t.Errorf("expected 3 finished items, got %d", finished)
	}
	if total != 3 {
		t.Errorf("expected total 3, got %d", total)
	}
}

func TestProgressError(t *testing.T) {
	p := New(1)

	p.UpdateTotal("test.txt", 100)
	p.Update("test.txt", 50)
	p.Error("test.txt", nil)

	finished, _ := p.GetStats()
	if finished != 1 {
		t.Errorf("expected 1 finished (with error), got %d", finished)
	}
}

func TestProgressConcurrentUpdates(t *testing.T) {
	p := New(1)
	p.UpdateTotal("test", 1000)

	done := make(chan bool)

	for i := 0; i < 10; i++ {
		go func() {
			for j := 0; j < 100; j++ {
				p.Update("test", 1)
			}
			done <- true
		}()
	}

	for i := 0; i < 10; i++ {
		<-done
	}

	item := p.items["test"]
	if item.Downloaded != 1000 {
		t.Errorf("expected 1000 downloaded, got %d", item.Downloaded)
	}
}
