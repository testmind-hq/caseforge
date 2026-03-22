// internal/tui/model.go
// Package tui provides a Bubble Tea progress display for caseforge gen/run.
package tui

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/testmind-hq/caseforge/internal/event"
)

// EventMsg wraps a CaseForge event for the Bubble Tea message loop.
type EventMsg struct{ event.Event }

// ProgressModel is the Bubble Tea model that shows generation progress.
type ProgressModel struct {
	total    int
	done     int
	finished bool
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
		case event.EventCaseGenerated, event.EventOperationDone:
			m.done++
		case event.EventRenderDone:
			m.finished = true
			return m, tea.Quit
		}
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	}
	return m, nil
}

// View renders the current state.
func (m ProgressModel) View() string {
	if m.finished {
		return fmt.Sprintf("Done — %d/%d\n", m.done, m.total)
	}
	return fmt.Sprintf("Generating cases... %d/%d\n", m.done, m.total)
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
