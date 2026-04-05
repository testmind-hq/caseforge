// internal/webhook/webhook_test.go
package webhook

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testmind-hq/caseforge/internal/config"
	"github.com/testmind-hq/caseforge/internal/event"
)

// ─────────────────────────────────────────────────────────
// sender tests
// ─────────────────────────────────────────────────────────

func TestSender_PostsJSONBody(t *testing.T) {
	var got []byte
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got, _ = io.ReadAll(r.Body)
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	s := newSender(srv.URL, "", 5, 1)
	err := s.send(t.Context(), map[string]string{"hello": "world"})
	require.NoError(t, err)
	assert.Contains(t, string(got), "hello")
}

func TestSender_SetsContentTypeJSON(t *testing.T) {
	var ct string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ct = r.Header.Get("Content-Type")
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	s := newSender(srv.URL, "", 5, 1)
	require.NoError(t, s.send(t.Context(), map[string]string{}))
	assert.Equal(t, "application/json", ct)
}

func TestSender_SignsPayloadWhenSecretSet(t *testing.T) {
	secret := "topsecret"
	var gotSig string
	var gotBody []byte
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotSig = r.Header.Get(signatureHeader)
		json.NewDecoder(r.Body).Decode(&gotBody)
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	payload := map[string]string{"key": "value"}
	s := newSender(srv.URL, secret, 5, 1)
	require.NoError(t, s.send(t.Context(), payload))

	body, _ := json.Marshal(payload)
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(body)
	expected := "sha256=" + hex.EncodeToString(mac.Sum(nil))
	assert.Equal(t, expected, gotSig)
}

func TestSender_NoSignatureHeaderWhenNoSecret(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Empty(t, r.Header.Get(signatureHeader))
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	s := newSender(srv.URL, "", 5, 1)
	require.NoError(t, s.send(t.Context(), map[string]string{}))
}

func TestSender_RetriesOn5xxThenSucceeds(t *testing.T) {
	var calls int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		n := atomic.AddInt32(&calls, 1)
		if n < 3 {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	s := newSender(srv.URL, "", 5, 3)
	s.backoff = func(int) time.Duration { return 0 } // zero-delay for test speed
	err := s.send(t.Context(), map[string]string{})
	require.NoError(t, err)
	assert.Equal(t, int32(3), atomic.LoadInt32(&calls))
}

func TestSender_NoRetryOn4xx(t *testing.T) {
	var calls int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&calls, 1)
		w.WriteHeader(http.StatusBadRequest)
	}))
	defer srv.Close()

	s := newSender(srv.URL, "", 5, 3)
	err := s.send(t.Context(), map[string]string{})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "400")
	assert.Equal(t, int32(1), atomic.LoadInt32(&calls), "must not retry on 4xx")
}

func TestSender_ReturnsErrorAfterAllRetriesExhausted(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadGateway)
	}))
	defer srv.Close()

	s := newSender(srv.URL, "", 5, 2)
	s.backoff = func(int) time.Duration { return 0 }
	err := s.send(t.Context(), map[string]string{})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "all 3 attempts failed")
}

// ─────────────────────────────────────────────────────────
// Sink tests
// ─────────────────────────────────────────────────────────

func TestSink_SkipsEntryWithEmptyURL(t *testing.T) {
	s := New([]config.WebhookConfig{{URL: "", Events: []string{"on_generate"}}})
	assert.Empty(t, s.entries)
}

func TestSink_DefaultsToAllEventsWhenNoneSpecified(t *testing.T) {
	s := New([]config.WebhookConfig{{URL: "http://example.com"}})
	require.Len(t, s.entries, 1)
	assert.True(t, s.entries[0].events[EventOnGenerate])
	assert.True(t, s.entries[0].events[EventOnRunComplete])
}

func TestSink_OnGeneratePostsCorrectPayload(t *testing.T) {
	var got GeneratePayload
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&got)
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	// Fix timestamp for determinism.
	fixed := time.Date(2026, 4, 5, 12, 0, 0, 0, time.UTC)
	now = func() time.Time { return fixed }
	t.Cleanup(func() { now = time.Now })

	s := New([]config.WebhookConfig{{URL: srv.URL, Events: []string{"on_generate"}}})
	s.Emit(event.Event{
		Type: event.EventOperationDone,
		Payload: event.OperationDonePayload{
			OperationID: "createUser",
			Method:      "POST",
			Path:        "/users",
			CaseCount:   13,
		},
	})

	assert.Equal(t, EventOnGenerate, got.Event)
	assert.Equal(t, fixed, got.Timestamp)
	assert.Equal(t, "createUser", got.Operation.ID)
	assert.Equal(t, "POST", got.Operation.Method)
	assert.Equal(t, "/users", got.Operation.Path)
	assert.Equal(t, 13, got.CaseCount)
}

func TestSink_OnRunCompletePostsCorrectPayload(t *testing.T) {
	var got RunCompletePayload
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewDecoder(r.Body).Decode(&got)
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	fixed := time.Date(2026, 4, 5, 12, 1, 0, 0, time.UTC)
	now = func() time.Time { return fixed }
	t.Cleanup(func() { now = time.Now })

	s := New([]config.WebhookConfig{{URL: srv.URL, Events: []string{"on_run_complete"}}})
	s.SetOutputDir("./cases")

	// Simulate two operations completing first.
	s.Emit(event.Event{Type: event.EventOperationDone, Payload: event.OperationDonePayload{CaseCount: 10}})
	s.Emit(event.Event{Type: event.EventOperationDone, Payload: event.OperationDonePayload{CaseCount: 6}})
	s.Emit(event.Event{Type: event.EventRenderDone})

	assert.Equal(t, EventOnRunComplete, got.Event)
	assert.Equal(t, fixed, got.Timestamp)
	assert.Equal(t, 16, got.TotalCases)
	assert.Equal(t, "./cases", got.OutputDir)
}

func TestSink_EventFilteringRespected(t *testing.T) {
	var generateCalls, completeCalls int32
	genSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&generateCalls, 1)
		w.WriteHeader(http.StatusOK)
	}))
	defer genSrv.Close()
	completeSrv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&completeCalls, 1)
		w.WriteHeader(http.StatusOK)
	}))
	defer completeSrv.Close()

	s := New([]config.WebhookConfig{
		{URL: genSrv.URL, Events: []string{"on_generate"}},
		{URL: completeSrv.URL, Events: []string{"on_run_complete"}},
	})

	s.Emit(event.Event{Type: event.EventOperationDone, Payload: event.OperationDonePayload{CaseCount: 5}})
	s.Emit(event.Event{Type: event.EventRenderDone})

	assert.Equal(t, int32(1), atomic.LoadInt32(&generateCalls))
	assert.Equal(t, int32(1), atomic.LoadInt32(&completeCalls))
}

func TestSink_DeliveryFailureNeverPanics(t *testing.T) {
	// Point at a URL that immediately refuses connections.
	s := New([]config.WebhookConfig{{URL: "http://127.0.0.1:1", MaxRetries: 1, Events: []string{"on_generate"}}})
	// Speed up: zero-delay backoff on the underlying sender.
	for i := range s.entries {
		s.entries[i].s.backoff = func(int) time.Duration { return 0 }
	}
	// Must not panic or return error — failures are warnings only.
	assert.NotPanics(t, func() {
		s.Emit(event.Event{
			Type:    event.EventOperationDone,
			Payload: event.OperationDonePayload{CaseCount: 1},
		})
	})
}

func TestSink_MultipleEndpointsAllReceiveEvent(t *testing.T) {
	var calls int32
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt32(&calls, 1)
		w.WriteHeader(http.StatusOK)
	})
	srv1 := httptest.NewServer(handler)
	defer srv1.Close()
	srv2 := httptest.NewServer(handler)
	defer srv2.Close()

	s := New([]config.WebhookConfig{
		{URL: srv1.URL, Events: []string{"on_run_complete"}},
		{URL: srv2.URL, Events: []string{"on_run_complete"}},
	})
	s.Emit(event.Event{Type: event.EventRenderDone})

	assert.Equal(t, int32(2), atomic.LoadInt32(&calls))
}
