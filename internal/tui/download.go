package tui

import (
	"context"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jhonroun/cmrd/pkg/cmrd"
)

type progressMsg struct {
	Event cmrd.ProgressEvent
}

type streamClosedMsg struct{}

type model struct {
	bar      progress.Model
	updates  <-chan cmrd.ProgressEvent
	phase    string
	message  string
	percent  float64
	total    int
	showHelp bool
	finished bool
	err      error
}

func newModel(updates <-chan cmrd.ProgressEvent) model {
	bar := progress.New(progress.WithDefaultGradient())
	return model{
		bar:     bar,
		updates: updates,
		phase:   "init",
		message: "starting",
	}
}

func waitForUpdate(updates <-chan cmrd.ProgressEvent) tea.Cmd {
	return func() tea.Msg {
		event, ok := <-updates
		if !ok {
			return streamClosedMsg{}
		}
		return progressMsg{Event: event}
	}
}

func (m model) Init() tea.Cmd {
	return waitForUpdate(m.updates)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch typed := msg.(type) {
	case tea.KeyMsg:
		switch typed.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "h", "?":
			m.showHelp = !m.showHelp
			return m, nil
		}
	case streamClosedMsg:
		m.finished = true
		return m, tea.Quit
	case progressMsg:
		event := typed.Event
		m.phase = event.Phase
		m.message = event.Message
		m.total = event.TotalFiles
		if event.Percent > 0 {
			m.percent = event.Percent
		}
		if event.Err != nil {
			m.err = event.Err
			m.finished = true
			return m, tea.Quit
		}
		if event.Done {
			m.finished = true
			return m, tea.Quit
		}
		barCmd := m.bar.SetPercent(m.percent / 100.0)
		return m, tea.Batch(waitForUpdate(m.updates), barCmd)
	}
	return m, nil
}

func (m model) View() string {
	title := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("45")).Render("CMRD TUI Downloader")
	hint := lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render("Keys: h/? help, q/Ctrl+C quit")
	phase := fmt.Sprintf("Phase: %s", strings.ToUpper(m.phase))
	total := fmt.Sprintf("Resolved files: %d", m.total)
	progressValue := fmt.Sprintf("Progress: %.1f%%", m.percent)
	status := fmt.Sprintf("Status: %s", m.message)

	if m.err != nil {
		status = fmt.Sprintf("Status: error: %v", m.err)
	}
	if m.finished && m.err == nil {
		status = "Status: completed"
	}

	lines := []string{
		title,
		"",
		phase,
		total,
		progressValue,
		m.bar.ViewAs(m.percent / 100.0),
		status,
		"",
		hint,
	}
	if m.showHelp {
		help := lipgloss.NewStyle().Foreground(lipgloss.Color("212")).Render("Help: TUI updates live progress from resolver/aria2 events; press q to exit.")
		lines = append(lines, help)
	}
	return strings.Join(lines, "\n")
}

// RunDownload starts download and renders progress in Bubble Tea UI.
func RunDownload(ctx context.Context, client *cmrd.Client, links []string) error {
	updates := make(chan cmrd.ProgressEvent, 64)
	errCh := make(chan error, 1)

	go func() {
		defer close(updates)
		err := client.Download(ctx, links, func(event cmrd.ProgressEvent) {
			select {
			case updates <- event:
			case <-ctx.Done():
			}
		})
		errCh <- err
	}()

	program := tea.NewProgram(newModel(updates))
	if _, err := program.Run(); err != nil {
		return err
	}

	return <-errCh
}
