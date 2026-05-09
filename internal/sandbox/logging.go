// internal/sandbox/logging.go
package sandbox

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"os"
)

// multiHandler fans log records out to multiple slog.Handler implementations.
type multiHandler struct {
	handlers []slog.Handler
}

func (m *multiHandler) Enabled(ctx context.Context, l slog.Level) bool {
	for _, h := range m.handlers {
		if h.Enabled(ctx, l) {
			return true
		}
	}
	return false
}

func (m *multiHandler) Handle(ctx context.Context, r slog.Record) error {
	var errs []error
	for _, h := range m.handlers {
		if h.Enabled(ctx, r.Level) {
			if err := h.Handle(ctx, r.Clone()); err != nil {
				errs = append(errs, err)
			}
		}
	}
	return errors.Join(errs...)
}

func (m *multiHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	hs := make([]slog.Handler, len(m.handlers))
	for i, h := range m.handlers {
		hs[i] = h.WithAttrs(attrs)
	}
	return &multiHandler{handlers: hs}
}

func (m *multiHandler) WithGroup(name string) slog.Handler {
	hs := make([]slog.Handler, len(m.handlers))
	for i, h := range m.handlers {
		hs[i] = h.WithGroup(name)
	}
	return &multiHandler{handlers: hs}
}

// buildLogger constructs a slog.Logger from Options.
// Returns the logger and a cleanup function (closes the log file if opened).
func buildLogger(opts Options) (*slog.Logger, func()) {
	level := parseLevel(opts.logLevelOrDefault())

	var handlers []slog.Handler
	var closers []io.Closer

	if opts.logLevelOrDefault() != "silent" {
		handlers = append(handlers, slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{Level: level}))
	}

	if opts.LogFile != "" {
		f, err := os.OpenFile(opts.LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err == nil {
			handlers = append(handlers, slog.NewJSONHandler(f, &slog.HandlerOptions{Level: level}))
			closers = append(closers, f)
		}
	}

	var h slog.Handler
	switch len(handlers) {
	case 0:
		h = slog.NewTextHandler(io.Discard, &slog.HandlerOptions{Level: slog.LevelError + 100})
	case 1:
		h = handlers[0]
	default:
		h = &multiHandler{handlers: handlers}
	}

	cleanup := func() {
		for _, c := range closers {
			_ = c.Close()
		}
	}
	return slog.New(h), cleanup
}

func parseLevel(level string) slog.Level {
	switch level {
	case "warn":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	case "silent":
		return slog.LevelError + 100 // effectively disabled
	default: // "info"
		return slog.LevelInfo
	}
}
