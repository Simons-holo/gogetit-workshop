package progress

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	progress   progress.Model
	items      []itemState
	total      int
	finished   int
	width      int
	updateChan chan interface{}
}

type itemState struct {
	url        string
	total      int64
	downloaded int64
	status     string
	percentage float64
}

func newModel(total int) model {
	p := progress.New(progress.WithDefaultGradient())
	return model{
		progress:   p,
		items:      make([]itemState, 0),
		total:      total,
		updateChan: make(chan interface{}, 5),
	}
}

func (m model) Init() tea.Cmd {
	return tea.Batch(
		waitForUpdate(m.updateChan),
	)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "q" || msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		return m, nil

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.progress.Width = msg.Width - 20
		return m, nil

	case ProgressMsg:
		found := false
		for i, item := range m.items {
			if item.url == msg.URL {
				m.items[i].downloaded += msg.Delta
				m.items[i].total = msg.Total
				if msg.Total > 0 {
					m.items[i].percentage = float64(m.items[i].downloaded) / float64(msg.Total)
				}
				found = true
				break
			}
		}
		if !found {
			newItem := itemState{
				url:        msg.URL,
				total:      msg.Total,
				downloaded: msg.Delta,
				status:     "downloading",
			}
			if msg.Total > 0 {
				newItem.percentage = float64(msg.Delta) / float64(msg.Total)
			}
			m.items = append(m.items, newItem)
		}
		return m, nil

	case TotalMsg:
		for i, item := range m.items {
			if item.url == msg.URL {
				m.items[i].total = msg.Total
				if msg.Total > 0 {
					m.items[i].percentage = float64(m.items[i].downloaded) / float64(msg.Total)
				}
				break
			}
		}
		return m, nil

	case DoneMsg:
		for i, item := range m.items {
			if item.url == msg.URL {
				m.items[i].status = "finished"
				m.items[i].percentage = 1.0
				break
			}
		}
		m.finished++
		if m.finished >= m.total {
			return m, tea.Quit
		}
		return m, nil

	default:
		return m, nil
	}
}

func (m model) View() string {
	var b strings.Builder

	b.WriteString(styleHeader.Render("  GoGetIt Downloader"))
	b.WriteString("\n\n")

	for _, item := range m.items {
		statusIcon := " "
		if item.status == "finished" {
			statusIcon = "+"
		} else if strings.HasPrefix(item.status, "error") {
			statusIcon = "x"
		}

		shortURL := item.url
		if len(shortURL) > 50 {
			shortURL = shortURL[:47] + "..."
		}

		b.WriteString(fmt.Sprintf("  %s %s\n", statusIcon, styleURL.Render(shortURL)))

		if item.total > 0 {
			perc := item.percentage
			if perc > 1.0 {
				perc = 1.0
			}
			m.progress.SetPercent(perc)
			bar := m.progress.View()
			downloadedMB := float64(item.downloaded) / 1024 / 1024
			totalMB := float64(item.total) / 1024 / 1024
			b.WriteString(fmt.Sprintf("  %s %.1f/%.1f MB\n", bar, downloadedMB, totalMB))
		} else {
			b.WriteString("  Downloading...\n")
		}
		b.WriteString("\n")
	}

	b.WriteString(fmt.Sprintf("\n  Progress: %d/%d completed\n", m.finished, m.total))

	return b.String()
}

func waitForUpdate(ch chan interface{}) tea.Cmd {
	return func() tea.Msg {
		return <-ch
	}
}
