// internal/llm/retry.go
package llm

import (
	"context"
	"time"
)

// retryDelays defines the wait before attempt i+1 (index = attempt number, 0-based).
// Attempt 0 runs immediately. Before attempt 1: 0 ms. Before attempt 2: 500 ms.
var retryDelays = []time.Duration{0, 500 * time.Millisecond}

// Retry calls fn up to maxAttempts times, returning the first successful result.
// Between attempts it waits according to retryDelays; context cancellation aborts the wait.
// Returns (nil, nil) when maxAttempts <= 0.
func Retry(ctx context.Context, maxAttempts int, fn func() (*CompletionResponse, error)) (*CompletionResponse, error) {
	var lastErr error
	for i := 0; i < maxAttempts; i++ {
		if i > 0 {
			delay := retryDelays[len(retryDelays)-1]
			if i-1 < len(retryDelays) {
				delay = retryDelays[i-1]
			}
			if delay > 0 {
				select {
				case <-ctx.Done():
					return nil, ctx.Err()
				case <-time.After(delay):
				}
			} else {
				// Zero-delay retry: non-blocking check so a cancelled context
				// is detected without racing against time.After(0).
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
