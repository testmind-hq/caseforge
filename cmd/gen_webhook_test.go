// cmd/gen_webhook_test.go
// Integration test: verifies that caseforge gen fires webhook events when
// webhooks are configured.
package cmd

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"sync/atomic"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testmind-hq/caseforge/internal/webhook"
)

// setWebhookConfig injects a webhook config into viper for the duration of
// the test. It resets viper on cleanup so other tests are unaffected.
func setWebhookConfig(t *testing.T, url string, events []string) {
	t.Helper()
	viper.Set("webhooks", []map[string]any{
		{"url": url, "events": events, "max_retries": 1},
	})
	t.Cleanup(func() { viper.Reset() })
}

// TestGenWebhook_OnGenerateFiredPerOperation verifies that gen calls the
// on_generate webhook once for each completed operation.
func TestGenWebhook_OnGenerateFiredPerOperation(t *testing.T) {
	t.Cleanup(resetGenGlobals(t))

	var calls int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var p webhook.GeneratePayload
		b, _ := io.ReadAll(r.Body)
		if json.Unmarshal(b, &p) == nil && p.Event == webhook.EventOnGenerate {
			atomic.AddInt32(&calls, 1)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	setWebhookConfig(t, srv.URL, []string{"on_generate"})

	outDir := runMiniGen(t)
	_ = outDir

	// mini.yaml has 3 operations → expect 3 on_generate calls.
	assert.Equal(t, int32(3), atomic.LoadInt32(&calls),
		"on_generate must fire once per completed operation")
}

// TestGenWebhook_OnRunCompleteFiredOnce verifies that gen calls the
// on_run_complete webhook exactly once after rendering finishes.
func TestGenWebhook_OnRunCompleteFiredOnce(t *testing.T) {
	t.Cleanup(resetGenGlobals(t))

	var calls int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var p webhook.RunCompletePayload
		b, _ := io.ReadAll(r.Body)
		if json.Unmarshal(b, &p) == nil && p.Event == webhook.EventOnRunComplete {
			atomic.AddInt32(&calls, 1)
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	setWebhookConfig(t, srv.URL, []string{"on_run_complete"})

	outDir := runMiniGen(t)
	_ = outDir

	assert.Equal(t, int32(1), atomic.LoadInt32(&calls),
		"on_run_complete must fire exactly once per run")
}

// TestGenWebhook_PayloadShapeOnGenerate verifies the on_generate payload
// contains the expected fields (event name, operation, case_count).
func TestGenWebhook_PayloadShapeOnGenerate(t *testing.T) {
	t.Cleanup(resetGenGlobals(t))

	var mu sync.Mutex
	var payloads []webhook.GeneratePayload
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var p webhook.GeneratePayload
		b, _ := io.ReadAll(r.Body)
		if json.Unmarshal(b, &p) == nil && p.Event == webhook.EventOnGenerate {
			mu.Lock()
			payloads = append(payloads, p)
			mu.Unlock()
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	setWebhookConfig(t, srv.URL, []string{"on_generate"})
	outDir := runMiniGen(t)
	_ = outDir

	mu.Lock()
	defer mu.Unlock()
	require.NotEmpty(t, payloads)
	for _, p := range payloads {
		assert.Equal(t, webhook.EventOnGenerate, p.Event)
		assert.False(t, p.Timestamp.IsZero(), "timestamp must be set")
		assert.NotEmpty(t, p.Operation.Method)
		assert.NotEmpty(t, p.Operation.Path)
		assert.Greater(t, p.CaseCount, 0)
	}
}

// TestGenWebhook_NoWebhookWhenConfigEmpty verifies that gen runs without
// error when no webhooks are configured (the common case).
func TestGenWebhook_NoWebhookWhenConfigEmpty(t *testing.T) {
	t.Cleanup(resetGenGlobals(t))
	outDir := runMiniGen(t)
	cases := readCases(t, outDir)
	require.NotEmpty(t, cases)
}

// TestGenWebhook_ConcurrencySafeOnGenerateCount verifies that totalSent is
// accumulated correctly under --concurrency 3 (race detector must pass).
func TestGenWebhook_ConcurrencySafeOnGenerateCount(t *testing.T) {
	t.Cleanup(resetGenGlobals(t))

	var totalReceived int32
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var p webhook.GeneratePayload
		b, _ := io.ReadAll(r.Body)
		if json.Unmarshal(b, &p) == nil && p.Event == webhook.EventOnGenerate {
			atomic.AddInt32(&totalReceived, int32(p.CaseCount))
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	setWebhookConfig(t, srv.URL, []string{"on_generate", "on_run_complete"})
	genConcurrency = 3
	outDir := runMiniGen(t)
	cases := readCases(t, outDir)
	require.NotEmpty(t, cases)

	// Verify on_generate events fired (race detector enforces no data race on totalSent).
	assert.Greater(t, atomic.LoadInt32(&totalReceived), int32(0),
		"on_generate must fire at least once under concurrent execution")
}
