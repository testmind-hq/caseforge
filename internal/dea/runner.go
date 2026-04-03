// internal/dea/runner.go
package dea

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// probeClient is a dedicated HTTP client with a per-probe timeout.
// Using a dedicated client avoids mutating http.DefaultClient globals.
var probeClient = &http.Client{Timeout: 10 * time.Second}

// RunProbe executes a single HTTP probe against targetURL and returns the evidence.
// targetURL is the API base URL (e.g. "http://localhost:8080"); probe.Path is appended to it.
func RunProbe(ctx context.Context, targetURL string, probe Probe) (*Evidence, error) {
	var bodyReader io.Reader
	if probe.Body != nil {
		data, err := json.Marshal(probe.Body)
		if err != nil {
			return nil, fmt.Errorf("marshal probe body: %w", err)
		}
		bodyReader = bytes.NewReader(data)
	}

	url := targetURL + probe.Path
	req, err := http.NewRequestWithContext(ctx, probe.Method, url, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	for k, v := range probe.Headers {
		req.Header.Set(k, v)
	}
	q := req.URL.Query()
	for k, v := range probe.QueryParams {
		q.Set(k, v)
	}
	req.URL.RawQuery = q.Encode()

	start := time.Now()
	resp, err := probeClient.Do(req)
	duration := time.Since(start)
	if err != nil {
		return nil, fmt.Errorf("execute probe %s %s: %w", probe.Method, probe.Path, err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response body: %w", err)
	}

	headers := make(map[string]string)
	for k, vs := range resp.Header {
		if len(vs) > 0 {
			headers[k] = vs[0]
		}
	}

	return &Evidence{
		ActualStatus:  resp.StatusCode,
		ActualBody:    string(respBody),
		ActualHeaders: headers,
		DurationMs:    duration.Milliseconds(),
	}, nil
}
