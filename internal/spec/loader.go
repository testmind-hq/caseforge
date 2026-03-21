// internal/spec/loader.go
package spec

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
