package progress

import (
	"fmt"
	"sync"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type Progress struct {
	mu       sync.RWMutex
	items    map[string]*ProgressItem
	total    int
	finished int
	program  *tea.Program
	done     chan struct{}
}

type ProgressItem struct {
	URL        string
	Total      int64
	Downloaded int64
	Percentage int
	Status     string
}

type ProgressMsg struct {
	URL   string
	Delta int64
	Total int64
}

type TotalMsg struct {
	URL   string
	Total int64
}

type DoneMsg struct {
	URL string
}

func New(total int) *Progress {
	return &Progress{
		items: make(map[string]*ProgressItem),
		total: total,
		done:  make(chan struct{}),
	}
}

func (p *Progress) Start() {
	model := newModel(p.total)
	p.program = tea.NewProgram(model)

	go func() {
		if _, err := p.program.Run(); err != nil {
			fmt.Printf("Progress error: %v\n", err)
		}
	}()

	<-p.done
}

func (p *Progress) Stop() {
	close(p.done)
	if p.program != nil {
		p.program.Quit()
	}
}

func (p *Progress) Update(url string, delta int64) {
	p.mu.Lock()
	defer p.mu.Unlock()

	item, exists := p.items[url]
	if !exists {
		item = &ProgressItem{URL: url, Status: "downloading"}
		p.items[url] = item
	}

	item.Downloaded += delta
	if item.Total > 0 {
		item.Percentage = int((float64(item.Downloaded) / float64(item.Total)) * 100)
		if item.Percentage > 100 {
			item.Percentage = 100
		}
	}

	if p.program != nil {
		p.program.Send(ProgressMsg{URL: url, Delta: delta, Total: item.Total})
	}
}

func (p *Progress) UpdateTotal(url string, total int64) {
	p.mu.Lock()
	defer p.mu.Unlock()

	item, exists := p.items[url]
	if !exists {
		item = &ProgressItem{URL: url, Status: "downloading"}
		p.items[url] = item
	}

	item.Total = total
	if total > 0 {
		item.Percentage = int((float64(item.Downloaded) / float64(total)) * 100)
	}

	if p.program != nil {
		p.program.Send(TotalMsg{URL: url, Total: total})
	}
}

func (p *Progress) Finish(url string) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if item, exists := p.items[url]; exists {
		item.Status = "finished"
		item.Percentage = 100
	}
	p.finished++

	if p.program != nil {
		p.program.Send(DoneMsg{URL: url})
	}
}

func (p *Progress) Error(url string, err error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	if item, exists := p.items[url]; exists {
		item.Status = fmt.Sprintf("error: %v", err)
	}
	p.finished++

	if p.program != nil {
		p.program.Send(DoneMsg{URL: url})
	}
}

func (p *Progress) GetStats() (int, int) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.finished, p.total
}

var (
	styleURL      = lipgloss.NewStyle().Foreground(lipgloss.Color("36")).Bold(true)
	styleProgress = lipgloss.NewStyle().Foreground(lipgloss.Color("35"))
	styleFinished = lipgloss.NewStyle().Foreground(lipgloss.Color("82"))
	styleError    = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
	styleHeader   = lipgloss.NewStyle().Foreground(lipgloss.Color("226")).Bold(true)
	styleBarBg    = lipgloss.NewStyle().Background(lipgloss.Color("236"))
	styleBarFill  = lipgloss.NewStyle().Background(lipgloss.Color("82"))
)
