// internal/llm/retry_test.go
package llm

import (
	"context"
	"errors"
	"testing"

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
