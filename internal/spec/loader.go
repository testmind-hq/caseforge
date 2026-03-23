// internal/spec/loader.go
package spec

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

type SpecLoader interface {
	Load(src string) (*ParsedSpec, error)
}

type ParsedSpec struct {
	Title      string
	Version    string
	Operations []*Operation
	Schemas    map[string]*Schema
}

type Operation struct {
	OperationID  string
	Method       string
	Path         string
	Summary      string
	Description  string
	Parameters   []*Parameter
	RequestBody  *RequestBody
	Responses    map[string]*Response
	Tags         []string
	Security     []string            // names of security schemes declared on this operation
	SemanticInfo *SemanticAnnotation // filled by LLM pre-processing
}

type SemanticAnnotation struct {
	ResourceType    string
	ActionType      string // "create"|"read"|"update"|"delete"|"action"
	HasStateMachine bool
	StateField      string
	UniqueFields    []string
	ImplicitRules   []string
}

// NewLoader returns a composite loader that dispatches by source prefix.
func NewLoader() *CompositeLoader {
	return &CompositeLoader{}
}

type CompositeLoader struct{}

func (c *CompositeLoader) Load(src string) (*ParsedSpec, error) {
	return c.loaderFor(src).Load(src)
}

func (c *CompositeLoader) loaderFor(src string) SpecLoader {
	if strings.HasPrefix(src, "http://") || strings.HasPrefix(src, "https://") {
		return NewURLLoader()
	}
	return NewLocalLoader()
}

// LocalLoader reads spec from the filesystem.
type LocalLoader struct{}

func NewLocalLoader() *LocalLoader { return &LocalLoader{} }

func (l *LocalLoader) Load(src string) (*ParsedSpec, error) {
	data, err := os.ReadFile(src)
	if err != nil {
		return nil, fmt.Errorf("loading spec from %s: %w", src, err)
	}
	return parseRawSpec(data, src)
}

// URLLoader fetches spec over HTTP/HTTPS.
type URLLoader struct{ client *http.Client }

func NewURLLoader() *URLLoader {
	return &URLLoader{client: &http.Client{}}
}

func (u *URLLoader) Load(src string) (*ParsedSpec, error) {
	resp, err := u.client.Get(src)
	if err != nil {
		return nil, fmt.Errorf("fetching spec from %s: %w", src, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("fetching spec: HTTP %d from %s", resp.StatusCode, src)
	}
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return parseRawSpec(data, src)
}
