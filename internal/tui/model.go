// internal/tui/model.go
// Package tui provides a Bubble Tea progress display for caseforge gen/run.
package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/testmind-hq/caseforge/internal/event"
)

// Styles for the progress list.
var (
	styleDone    = lipgloss.NewStyle().Foreground(lipgloss.Color("2")) // green
	styleSpinner = lipgloss.NewStyle().Foreground(lipgloss.Color("3")) // yellow
	styleCount   = lipgloss.NewStyle().Foreground(lipgloss.Color("8")) // dim gray
	styleSummary = lipgloss.NewStyle().Bold(true)
)

const maxVisibleRows = 12 // maximum rows shown in the scrolling list

// opRow holds the display state for a single completed operation.
type opRow struct {
	label     string // e.g. "POST /pets"
	caseCount int
}

// EventMsg wraps a CaseForge event for the Bubble Tea message loop.
type EventMsg struct{ event.Event }

// PrintMsg prints a line above the TUI view via tea.Println.
// Always send with prog.Send (not prog.Printf) — Send has a ctx.Done guard
// that prevents deadlocks after the program exits.
type PrintMsg struct{ Text string }

// ProgressModel is the Bubble Tea model that shows generation progress.
type ProgressModel struct {
	total      int
	annotated  int // operations annotated by LLM so far
	done       int // operations fully generated so far
	finished   bool
	rows       []opRow
}

// NewProgressModel creates a model that expects `total` operations to complete.
func NewProgressModel(total int) ProgressModel {
	return ProgressModel{total: total}
}

// Init satisfies tea.Model — no initial command needed.
func (m ProgressModel) Init() tea.Cmd { return nil }

// Update handles incoming messages.
func (m ProgressModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case EventMsg:
		switch msg.Type {
		case event.EventOperationAnnotating:
			m.annotated++
		case event.EventOperationDone:
			m.done++
			if p, ok := msg.Payload.(event.OperationDonePayload); ok {
				label := strings.TrimSpace(p.Method + " " + p.Path)
				if label == "" {
					label = p.OperationID
				}
				m.rows = append(m.rows, opRow{label: label, caseCount: p.CaseCount})
			}
		case event.EventRenderDone:
			m.finished = true
			return m, tea.Quit
		}
	case PrintMsg:
		return m, tea.Println(msg.Text)
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	}
	return m, nil
}

// View renders the current state.
func (m ProgressModel) View() string {
	var b strings.Builder

	// Scrolling operation list: show the last maxVisibleRows completed rows.
	visible := m.rows
	if len(visible) > maxVisibleRows {
		hidden := len(visible) - maxVisibleRows
		b.WriteString(styleCount.Render(fmt.Sprintf("  … %d more above\n", hidden)))
		visible = visible[len(visible)-maxVisibleRows:]
	}
	for _, r := range visible {
		line := styleDone.Render("✓ "+r.label) +
			styleCount.Render(fmt.Sprintf("  (%d cases)", r.caseCount))
		b.WriteString(line + "\n")
	}

	// Progress / summary line.
	spinner := styleSpinner.Render("⠋")
	switch {
	case m.finished:
		b.WriteString(styleSummary.Render(fmt.Sprintf("Done — %d/%d operations", m.done, m.total)) + "\n")
	case m.annotated > 0 && m.done == 0:
		b.WriteString(fmt.Sprintf("%s Annotating with AI... %d/%d\n", spinner, m.annotated, m.total))
	default:
		b.WriteString(fmt.Sprintf("%s Generating cases... %d/%d\n", spinner, m.done, m.total))
	}

	return b.String()
}

// TUISink bridges the event.Sink interface to a running tea.Program.
type TUISink struct {
	prog *tea.Program
}

func NewTUISink(prog *tea.Program) *TUISink {
	return &TUISink{prog: prog}
}

// Emit sends the event to the Bubble Tea program's message loop.
func (s *TUISink) Emit(e event.Event) {
	s.prog.Send(EventMsg{e})
}
