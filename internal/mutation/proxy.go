// internal/mutation/proxy.go
package mutation

import (
	"bytes"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
)

// Registry returns all 12 built-in mutation operators in a fixed order.
func Registry() []Operator {
	return []Operator{
		NewFieldDropOperator(),
		NewFieldTypeSwapOperator(),
		NewArrayToNullOperator(),
		NewNullToArrayOperator(),
		NewStatusSwap2xxOperator(),
		NewErrorInflationOperator(),
		NewPaginationOffByOneOperator(),
		NewEmptyResultInjectionOperator(),
		NewContentTypeSwapOperator(),
		NewHeaderDropOperator(),
		NewDateFormatSwapOperator(),
		NewNumericPrecisionLossOperator(),
	}
}

// MutationProxy wraps httputil.ReverseProxy and intercepts responses to apply mutations.
type MutationProxy struct {
	server   *http.Server
	listener net.Listener
	mu       sync.RWMutex
	active   Operator // nil = passthrough
}

// NewProxy creates a MutationProxy that forwards to target (e.g. "http://api:8080").
// Call Close() when done.
func NewProxy(target string) (*MutationProxy, error) {
	targetURL, err := url.Parse(target)
	if err != nil {
		return nil, fmt.Errorf("invalid target URL: %w", err)
	}

	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return nil, fmt.Errorf("starting proxy listener: %w", err)
	}

	mp := &MutationProxy{listener: ln}

	rp := httputil.NewSingleHostReverseProxy(targetURL)
	rp.ModifyResponse = mp.modifyResponse

	mp.server = &http.Server{Handler: rp}
	go mp.server.Serve(ln) //nolint:errcheck
	return mp, nil
}

// Addr returns the proxy's local listen address (e.g. "127.0.0.1:34521").
func (mp *MutationProxy) Addr() string {
	return mp.listener.Addr().String()
}

// SetActive sets the operator applied to every response until changed.
// Pass nil for passthrough mode.
func (mp *MutationProxy) SetActive(op Operator) {
	mp.mu.Lock()
	defer mp.mu.Unlock()
	mp.active = op
}

// Close shuts down the proxy server.
func (mp *MutationProxy) Close() error {
	return mp.server.Close()
}

func (mp *MutationProxy) modifyResponse(resp *http.Response) error {
	mp.mu.RLock()
	op := mp.active
	mp.mu.RUnlock()

	if op == nil {
		return nil
	}

	body, err := io.ReadAll(resp.Body)
	_ = resp.Body.Close()
	if err != nil {
		return err
	}

	mutated, err := op.Apply(resp, body)
	if err != nil {
		mutated = body // operator failed: pass original body through
	}

	resp.Body = io.NopCloser(bytes.NewReader(mutated))
	resp.ContentLength = int64(len(mutated))
	resp.Header.Set("Content-Length", fmt.Sprintf("%d", len(mutated)))
	return nil
}
