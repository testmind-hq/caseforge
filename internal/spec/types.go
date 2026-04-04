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
	Ref         string // original $ref path if applicable
	Example     any
}
