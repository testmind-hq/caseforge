// internal/webhook/sink.go
// Sink implements event.Sink and dispatches CaseForge events to configured
// webhook endpoints.
package webhook

import (
	"context"
	"fmt"
	"time"

	"github.com/testmind-hq/caseforge/internal/config"
	"github.com/testmind-hq/caseforge/internal/event"
)

// now is a variable so tests can override it for deterministic timestamps.
var now = time.Now

// entry binds a sender to the set of event names it should receive.
type entry struct {
	s      *sender
	events map[EventName]bool
}

// Sink subscribes to the event bus and forwards matching events to webhook
// endpoints. Delivery failures are printed as warnings — they never propagate
// back to the caller or block the gen pipeline.
type Sink struct {
	entries   []entry
	outputDir string // captured from RunCompletePayload context
	totalSent int    // running count of cases seen across all operations
}

// New builds a Sink from the webhook configurations. Entries with an empty
// URL are silently skipped.
func New(cfgs []config.WebhookConfig) *Sink {
	s := &Sink{}
	for _, c := range cfgs {
		if c.URL == "" {
			continue
		}
		evts := make(map[EventName]bool, len(c.Events))
		for _, e := range c.Events {
			evts[EventName(e)] = true
		}
		// Default: subscribe to both events if none specified.
		if len(evts) == 0 {
			evts[EventOnGenerate] = true
			evts[EventOnRunComplete] = true
		}
		s.entries = append(s.entries, entry{
			s:      newSender(c.URL, c.Secret, c.TimeoutSecs, c.MaxRetries),
			events: evts,
		})
	}
	return s
}

// SetOutputDir captures the output directory for the on_run_complete payload.
func (s *Sink) SetOutputDir(dir string) { s.outputDir = dir }

// Emit handles incoming bus events and dispatches webhook POSTs.
func (s *Sink) Emit(e event.Event) {
	switch e.Type {
	case event.EventOperationDone:
		p, ok := e.Payload.(event.OperationDonePayload)
		if !ok {
			return
		}
		s.totalSent += p.CaseCount
		s.dispatch(EventOnGenerate, func() any {
			gp := GeneratePayload{Event: EventOnGenerate}
			gp.Timestamp = now()
			gp.Operation.ID = p.OperationID
			gp.Operation.Method = p.Method
			gp.Operation.Path = p.Path
			gp.CaseCount = p.CaseCount
			return gp
		})

	case event.EventRenderDone:
		s.dispatch(EventOnRunComplete, func() any {
			return RunCompletePayload{
				Event:      EventOnRunComplete,
				Timestamp:  now(),
				TotalCases: s.totalSent,
				OutputDir:  s.outputDir,
			}
		})
	}
}

// dispatch calls payloadFn and POSTs the result to every entry subscribed to evtName.
func (s *Sink) dispatch(evtName EventName, payloadFn func() any) {
	subscribed := false
	for _, e := range s.entries {
		if e.events[evtName] {
			subscribed = true
			break
		}
	}
	if !subscribed {
		return
	}
	payload := payloadFn()
	ctx := context.Background()
	for _, e := range s.entries {
		if !e.events[evtName] {
			continue
		}
		if err := e.s.send(ctx, payload); err != nil {
			fmt.Printf("warning: webhook %s delivery failed: %v\n", evtName, err)
		}
	}
}
