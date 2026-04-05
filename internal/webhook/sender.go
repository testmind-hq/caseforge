// internal/webhook/sender.go
// HTTP sender with exponential-backoff retry and optional HMAC-SHA256 signing.
package webhook

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	defaultTimeoutSecs = 10
	defaultMaxRetries  = 3
	signatureHeader    = "X-CaseForge-Signature-256"
)

// sender posts JSON payloads to a single URL with retry and optional signing.
type sender struct {
	url        string
	secret     string
	maxRetries int
	client     *http.Client
	// backoff returns the wait duration before attempt i (1-indexed).
	// Overridable in tests for fast execution.
	backoff func(attempt int) time.Duration
}

func newSender(url, secret string, timeoutSecs, maxRetries int) *sender {
	if timeoutSecs <= 0 {
		timeoutSecs = defaultTimeoutSecs
	}
	if maxRetries <= 0 {
		maxRetries = defaultMaxRetries
	}
	return &sender{
		url:        url,
		secret:     secret,
		maxRetries: maxRetries,
		client:     &http.Client{Timeout: time.Duration(timeoutSecs) * time.Second},
		backoff:    func(attempt int) time.Duration { return time.Duration(1<<(attempt-1)) * time.Second },
	}
}

// send marshals payload to JSON and posts it, retrying on transient errors.
// It returns an error on failure; callers are responsible for treating it as non-fatal.
func (s *sender) send(ctx context.Context, payload any) error {
	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("webhook: marshal payload: %w", err)
	}

	var lastErr error
	for attempt := 0; attempt <= s.maxRetries; attempt++ {
		if attempt > 0 {
			wait := s.backoff(attempt)
			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(wait):
			}
		}

		req, err := http.NewRequestWithContext(ctx, http.MethodPost, s.url, bytes.NewReader(body))
		if err != nil {
			return fmt.Errorf("webhook: build request: %w", err)
		}
		req.Header.Set("Content-Type", "application/json")
		if s.secret != "" {
			req.Header.Set(signatureHeader, "sha256="+sign(body, s.secret))
		}

		resp, err := s.client.Do(req)
		if err != nil {
			lastErr = fmt.Errorf("attempt %d: %w", attempt+1, err)
			continue // network error → retry
		}
		// Drain body before closing so the connection can be reused.
		io.Copy(io.Discard, resp.Body) //nolint:errcheck
		resp.Body.Close()

		if resp.StatusCode >= 200 && resp.StatusCode < 300 {
			return nil // success
		}
		if resp.StatusCode >= 400 && resp.StatusCode < 500 {
			// Client error — don't retry
			return fmt.Errorf("webhook: server returned %d (client error, not retrying)", resp.StatusCode)
		}
		// 5xx → retry
		lastErr = fmt.Errorf("attempt %d: server returned %d", attempt+1, resp.StatusCode)
	}
	return fmt.Errorf("webhook: all %d attempts failed; last error: %w", s.maxRetries+1, lastErr)
}

// sign returns a hex-encoded HMAC-SHA256 of body using key.
func sign(body []byte, key string) string {
	mac := hmac.New(sha256.New, []byte(key))
	mac.Write(body)
	return hex.EncodeToString(mac.Sum(nil))
}
