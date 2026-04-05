package tui

import (
	"fmt"
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
	updated, _ := m.Update(EventMsg{event.Event{
		Type:    event.EventOperationDone,
		Payload: event.OperationDonePayload{Method: "GET", Path: "/pets", CaseCount: 4},
	}})
	pm := updated.(ProgressModel)
	assert.Equal(t, 1, pm.done)
	assert.Len(t, pm.rows, 1)
	assert.Equal(t, "GET /pets", pm.rows[0].label)
	assert.Equal(t, 4, pm.rows[0].caseCount)
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

func TestProgressModel_ViewShowsCompletedOperations(t *testing.T) {
	m := NewProgressModel(3)
	m.rows = []opRow{
		{label: "GET /pets", caseCount: 4},
		{label: "POST /pets", caseCount: 6},
	}
	m.done = 2
	view := m.View()
	assert.Contains(t, view, "GET /pets")
	assert.Contains(t, view, "POST /pets")
	assert.Contains(t, view, "4 cases")
	assert.Contains(t, view, "6 cases")
}

func TestProgressModel_ViewScrollsToLast12Rows(t *testing.T) {
	m := NewProgressModel(20)
	for i := 0; i < 15; i++ {
		m.rows = append(m.rows, opRow{label: fmt.Sprintf("GET /op%d", i), caseCount: 1})
	}
	m.done = 15
	view := m.View()
	// First 3 rows should be scrolled away (only last 12 visible)
	assert.NotContains(t, view, "GET /op0")
	assert.NotContains(t, view, "GET /op2")
	// Last 12 rows should be visible
	assert.Contains(t, view, "GET /op14")
	assert.Contains(t, view, "GET /op3")
}

func TestProgressModel_WindowSizeMsg(t *testing.T) {
	// WindowSizeMsg must be handled without panic and not change done/total.
	m := NewProgressModel(3)
	updated, _ := m.Update(tea.WindowSizeMsg{Width: 120, Height: 40})
	pm := updated.(ProgressModel)
	assert.Equal(t, 0, pm.done)
	assert.Equal(t, 3, pm.total)
}

func TestProgressModel_OperationDone_NoPayload(t *testing.T) {
	// Payload-less EventOperationDone must still increment done without panic.
	m := NewProgressModel(3)
	updated, _ := m.Update(EventMsg{event.Event{Type: event.EventOperationDone}})
	pm := updated.(ProgressModel)
	assert.Equal(t, 1, pm.done)
	assert.Empty(t, pm.rows) // no row appended without payload
}

// Ensure ProgressModel satisfies tea.Model interface at compile time.
var _ tea.Model = ProgressModel{}
