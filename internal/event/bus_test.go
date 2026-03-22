package event

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

type recordSink struct {
	mu     sync.Mutex
	events []Event
}

func (s *recordSink) Emit(e Event) {
	s.mu.Lock()
	s.events = append(s.events, e)
	s.mu.Unlock()
}

func TestBusEmitsToSubscribers(t *testing.T) {
	bus := NewBus()
	sink1 := &recordSink{}
	sink2 := &recordSink{}
	bus.Subscribe(sink1)
	bus.Subscribe(sink2)

	bus.Emit(Event{Type: EventCaseGenerated, Payload: "tc-1"})
	bus.Emit(Event{Type: EventRenderDone})

	assert.Len(t, sink1.events, 2)
	assert.Len(t, sink2.events, 2)
	assert.Equal(t, EventCaseGenerated, sink1.events[0].Type)
}

func TestBusNoSubscribersNoPanic(t *testing.T) {
	bus := NewBus()
	assert.NotPanics(t, func() {
		bus.Emit(Event{Type: EventSpecLoaded})
	})
}

func TestBusConcurrentEmit(t *testing.T) {
	bus := NewBus()
	sink := &recordSink{}
	bus.Subscribe(sink)

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			bus.Emit(Event{Type: EventCaseGenerated})
		}()
	}
	wg.Wait()
	assert.Len(t, sink.events, 100)
}
