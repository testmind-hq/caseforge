// internal/event/bus.go
package event

import "sync"

// Bus is a thread-safe fan-out event sink. It implements Sink itself so
// it can be passed wherever a Sink is accepted.
type Bus struct {
	mu    sync.RWMutex
	sinks []Sink
}

// NewBus returns an empty Bus with no subscribers.
func NewBus() *Bus { return &Bus{} }

// Subscribe registers a sink to receive all future events.
func (b *Bus) Subscribe(s Sink) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.sinks = append(b.sinks, s)
}

// Emit broadcasts the event to all registered sinks.
func (b *Bus) Emit(e Event) {
	b.mu.RLock()
	sinks := make([]Sink, len(b.sinks))
	copy(sinks, b.sinks)
	b.mu.RUnlock()

	for _, s := range sinks {
		s.Emit(e)
	}
}
