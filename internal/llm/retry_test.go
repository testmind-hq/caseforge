// internal/llm/retry_test.go
package llm

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRetry_SucceedsFirstAttempt(t *testing.T) {
	calls := 0
	resp, err := Retry(context.Background(), 3, func() (*CompletionResponse, error) {
		calls++
		return &CompletionResponse{Text: "ok"}, nil
	})
	require.NoError(t, err)
	assert.Equal(t, "ok", resp.Text)
	assert.Equal(t, 1, calls)
}

func TestRetry_SucceedsOnSecondAttempt(t *testing.T) {
	calls := 0
	resp, err := Retry(context.Background(), 3, func() (*CompletionResponse, error) {
		calls++
		if calls < 2 {
			return nil, errors.New("transient error")
		}
		return &CompletionResponse{Text: "ok"}, nil
	})
	require.NoError(t, err)
	assert.Equal(t, "ok", resp.Text)
	assert.Equal(t, 2, calls)
}

func TestRetry_AllAttemptsFail_ReturnsLastError(t *testing.T) {
	calls := 0
	sentinel := errors.New("always fails")
	_, err := Retry(context.Background(), 3, func() (*CompletionResponse, error) {
		calls++
		return nil, sentinel
	})
	require.Error(t, err)
	assert.Equal(t, sentinel, err)
	assert.Equal(t, 3, calls)
}

func TestRetry_ContextCancelledDuringWait(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	calls := 0
	_, err := Retry(ctx, 3, func() (*CompletionResponse, error) {
		calls++
		if calls == 1 {
			cancel() // cancel during the wait before attempt 2
		}
		return nil, errors.New("fail")
	})
	require.Error(t, err)
	assert.True(t, errors.Is(err, context.Canceled), "expected context.Canceled, got %v", err)
	assert.Equal(t, 1, calls)
}

func TestRetry_ContextCancelledDuringBlockingWait(t *testing.T) {
	// Covers the delay>0 (500ms) branch: cancel after attempt 2 triggers the
	// blocking select in the wait before attempt 3.
	ctx, cancel := context.WithCancel(context.Background())
	calls := 0
	_, err := Retry(ctx, 3, func() (*CompletionResponse, error) {
		calls++
		if calls == 2 {
			cancel() // cancel before the 500ms wait before attempt 3
		}
		return nil, errors.New("fail")
	})
	require.Error(t, err)
	assert.True(t, errors.Is(err, context.Canceled), "expected context.Canceled, got %v", err)
	assert.Equal(t, 2, calls)
}

func TestRetry_ZeroMaxAttempts_ReturnsNil(t *testing.T) {
	called := false
	resp, err := Retry(context.Background(), 0, func() (*CompletionResponse, error) {
		called = true
		return &CompletionResponse{Text: "x"}, nil
	})
	assert.Nil(t, resp)
	assert.Nil(t, err)
	assert.False(t, called)
}

func TestRetry_RateLimitError_UsesRateLimitDelayPath(t *testing.T) {
	// Swap rateLimitDelays for zero-duration values so the test doesn't sleep.
	origRL := rateLimitDelays
	rateLimitDelays = []time.Duration{0, 0, 0, 0}
	defer func() { rateLimitDelays = origRL }()

	calls := 0
	resp, err := Retry(context.Background(), 4, func() (*CompletionResponse, error) {
		calls++
		if calls < 3 {
			return nil, errors.New("status 429 rate limit exceeded")
		}
		return &CompletionResponse{Text: "ok"}, nil
	})
	require.NoError(t, err)
	assert.Equal(t, "ok", resp.Text)
	assert.Equal(t, 3, calls, "must retry after 429 errors and succeed on 3rd attempt")
}

func TestRetry_RateLimitThenNormalError_ResetsRLIndex(t *testing.T) {
	// Verify rlIdx resets when a non-429 error follows a 429 error.
	origRL := rateLimitDelays
	origR := retryDelays
	rateLimitDelays = []time.Duration{0, 0}
	retryDelays = []time.Duration{0, 0}
	defer func() { rateLimitDelays = origRL; retryDelays = origR }()

	callN := 0
	_, err := Retry(context.Background(), 3, func() (*CompletionResponse, error) {
		callN++
		if callN == 1 {
			return nil, errors.New("status 429 rate limit exceeded")
		}
		return nil, errors.New("internal server error") // non-429
	})
	require.Error(t, err)
	assert.Equal(t, 3, callN)
}

func TestIsRateLimitErr(t *testing.T) {
	assert.True(t, isRateLimitErr(errors.New("request failed with status 429")))
	assert.True(t, isRateLimitErr(errors.New("status code 429 too many requests")))
	assert.True(t, isRateLimitErr(errors.New("Rate limit exceeded")))
	assert.True(t, isRateLimitErr(errors.New("too many requests")))
	assert.True(t, isRateLimitErr(errors.New("ratelimit reached")))
	// Must NOT match bare numbers that aren't HTTP status codes
	assert.False(t, isRateLimitErr(errors.New("error code 4290")))
	assert.False(t, isRateLimitErr(errors.New("row 429 of 1000 processed")))
	assert.False(t, isRateLimitErr(errors.New("connection refused")))
	assert.False(t, isRateLimitErr(errors.New("timeout")))
	assert.False(t, isRateLimitErr(nil))
}
