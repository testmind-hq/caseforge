package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
	"github.com/testmind-hq/caseforge/internal/event"
)

func TestProgressModelInitialState(t *testing.T) {
	m := NewProgressModel(5)
	assert.Equal(t, 0, m.done)
	assert.Equal(t, 5, m.total)
	assert.False(t, m.finished)
}

func TestProgressModelUpdate_CaseGenerated(t *testing.T) {
	m := NewProgressModel(3)
	updated, _ := m.Update(EventMsg{event.Event{Type: event.EventOperationDone}})
	pm := updated.(ProgressModel)
	assert.Equal(t, 1, pm.done)
}

func TestProgressModelUpdate_RenderDone(t *testing.T) {
	m := NewProgressModel(3)
	updated, cmd := m.Update(EventMsg{event.Event{Type: event.EventRenderDone}})
	pm := updated.(ProgressModel)
	assert.True(t, pm.finished)
	assert.NotNil(t, cmd) // should return tea.Quit
}

func TestProgressModelView(t *testing.T) {
	m := NewProgressModel(5)
	view := m.View()
	assert.Contains(t, view, "0/5")
}

func TestProgressModelViewFinished(t *testing.T) {
	m := NewProgressModel(3)
	m.done = 3
	m.finished = true
	view := m.View()
	assert.Contains(t, view, "Done")
}

// Ensure ProgressModel satisfies tea.Model interface at compile time.
var _ tea.Model = ProgressModel{}
