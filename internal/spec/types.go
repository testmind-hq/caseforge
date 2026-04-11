// internal/spec/types.go
package spec

// Parameter, RequestBody, Response, Schema are thin wrappers
// over kin-openapi types. Fields added as needed during implementation.
// See parser.go for conversion from kin-openapi to these types.

type Parameter struct {
	Name     string
	In       string // "query"|"path"|"header"|"cookie"
	Required bool
	Schema   *Schema
}

type RequestBody struct {
	Required bool
	Content  map[string]*MediaType // key: "application/json" etc.
}

type MediaType struct {
	Schema   *Schema
	Example  any                 // mediaType-level example (single value)
	Examples map[string]*Example // mediaType-level named examples
}

// Example represents a named OpenAPI example object.
type Example struct {
	Summary     string
	Description string
	Value       any // nil when ExternalValue is used (not supported here)
}

type Response struct {
	Description string
	Content     map[string]*MediaType
	Headers     map[string]string // key=header name, value=schema type ("string", "integer", etc.)
}

// SpecLink represents an OpenAPI 3.0 link declared on a response.
// It expresses that a parameter of another operation can be populated
// from a value in this response — a machine-readable producer→consumer edge.
type SpecLink struct {
	Name         string            // link object name (e.g. "GetUserById")
	OperationID  string            // target operationId (e.g. "getUser")
	Parameters   map[string]string // paramName → "$response.body#/fieldPath" or literal
	ResponseCode string            // response code this link appears on (e.g. "201")
}

type Schema struct {
	Type        string
	Format      string
	Description string
	Properties  map[string]*Schema
	Items       *Schema // for array type
	Enum        []any
	Minimum     *float64
	Maximum     *float64
	MinLength   *int64
	MaxLength   *int64
	MinItems    *uint64
	MaxItems    *uint64
	Required    []string
	Nullable    bool
	ReadOnly    bool
	WriteOnly   bool
	Ref         string // original $ref path if applicable
	Example     any
	Pattern     string // regex pattern constraint for string fields (OpenAPI `pattern`)
}
