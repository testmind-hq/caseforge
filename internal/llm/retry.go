// internal/llm/retry.go
package llm

import (
	"context"
	"strings"
	"time"
)

// retryDelays defines the wait before attempt i+1 for normal transient errors.
var retryDelays = []time.Duration{0, 500 * time.Millisecond}

// rateLimitDelays defines exponential backoff waits for 429 rate-limit responses.
var rateLimitDelays = []time.Duration{
	5 * time.Second,
	15 * time.Second,
	30 * time.Second,
	60 * time.Second,
}

// isRateLimitErr reports whether err looks like a 429 / rate-limit response.
// Checked by string because each SDK wraps HTTP errors differently.
func isRateLimitErr(err error) bool {
	if err == nil {
		return false
	}
	msg := strings.ToLower(err.Error())
	return strings.Contains(msg, "status 429") ||
		strings.Contains(msg, "status code 429") ||
		strings.Contains(msg, "rate limit") ||
		strings.Contains(msg, "too many requests") ||
		strings.Contains(msg, "ratelimit")
}

// Retry calls fn up to maxAttempts times, returning the first successful result.
// Between attempts it waits according to retryDelays for normal errors, or
// rateLimitDelays for 429 / rate-limit responses.
// Context cancellation aborts the wait. Returns (nil, nil) when maxAttempts <= 0.
func Retry(ctx context.Context, maxAttempts int, fn func() (*CompletionResponse, error)) (*CompletionResponse, error) {
	var lastErr error
	rlIdx := 0 // index into rateLimitDelays for successive 429s
	for i := 0; i < maxAttempts; i++ {
		if i > 0 {
			var delay time.Duration
			if isRateLimitErr(lastErr) {
				delay = rateLimitDelays[min(rlIdx, len(rateLimitDelays)-1)]
				rlIdx++
			} else {
				rlIdx = 0
				idx := i - 1
				if idx < len(retryDelays) {
					delay = retryDelays[idx]
				} else {
					delay = retryDelays[len(retryDelays)-1]
				}
			}
			if delay > 0 {
				select {
				case <-ctx.Done():
					return nil, ctx.Err()
				case <-time.After(delay):
				}
			} else {
				select {
				case <-ctx.Done():
					return nil, ctx.Err()
				default:
				}
			}
		}
		resp, err := fn()
		if err == nil {
			return resp, nil
		}
		lastErr = err
	}
	return nil, lastErr
}

