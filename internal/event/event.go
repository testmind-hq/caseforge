// internal/event/event.go
// Package event defines event types for CaseForge's event-driven architecture.
package event

// EventType identifies what happened.
type EventType string

const (
	EventSpecLoaded          EventType = "spec.loaded"
	EventOperationAnnotating EventType = "operation.annotating"
	EventOperationDone       EventType = "operation.done"
	EventCaseGenerated       EventType = "case.generated"
	EventRenderDone          EventType = "render.done"
	EventError               EventType = "error"
)

// Event carries data about something that happened.
type Event struct {
	Type    EventType
	Payload any
}

// OperationDonePayload is the structured payload for EventOperationDone and
// EventOperationAnnotating. CaseCount is zero for annotation events.
type OperationDonePayload struct {
	OperationID string
	Method      string
	Path        string
	CaseCount   int
}

// Sink receives events. Implement this to observe CaseForge progress.
type Sink interface {
	Emit(e Event)
}

// NoopSink discards all events. Used as default until Phase 2 TUI is wired in.
type NoopSink struct{}

func (s *NoopSink) Emit(_ Event) {}

// SinkFunc is a function adapter that implements Sink.
type SinkFunc func(Event)

func (f SinkFunc) Emit(e Event) { f(e) }
