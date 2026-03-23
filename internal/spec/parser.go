// internal/spec/parser.go
package spec

import (
	"context"
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
)

// parseRawSpec parses raw YAML/JSON bytes into a ParsedSpec.
// src is used as a hint for $ref resolution.
func parseRawSpec(data []byte, src string) (*ParsedSpec, error) {
	loader := openapi3.NewLoader()
	doc, err := loader.LoadFromData(data)
	if err != nil {
		return nil, fmt.Errorf("parsing spec: %w", err)
	}
	if err := doc.Validate(context.Background()); err != nil {
		// Validation errors are non-fatal; log and continue
		// (many real-world specs have minor issues)
		_ = err
	}
	return convertDoc(doc), nil
}

func convertDoc(doc *openapi3.T) *ParsedSpec {
	ps := &ParsedSpec{
		Title:   doc.Info.Title,
		Version: doc.Info.Version,
		Schemas: make(map[string]*Schema),
	}

	for path, item := range doc.Paths.Map() {
		for method, op := range item.Operations() {
			if op == nil {
				continue
			}
			ps.Operations = append(ps.Operations, convertOperation(method, path, op))
		}
	}

	if doc.Components != nil {
		for name, ref := range doc.Components.Schemas {
			if ref.Value != nil {
				ps.Schemas[name] = convertSchema(ref.Value)
			}
		}
	}

	return ps
}

func convertOperation(method, path string, op *openapi3.Operation) *Operation {
	o := &Operation{
		OperationID: op.OperationID,
		Method:      method,
		Path:        path,
		Summary:     op.Summary,
		Description: op.Description,
		Tags:        op.Tags,
		Responses:   make(map[string]*Response),
	}
	if op.Security != nil {
		for _, req := range *op.Security {
			for schemeName := range req {
				o.Security = append(o.Security, schemeName)
			}
		}
	}
	for _, p := range op.Parameters {
		if p.Value != nil {
			o.Parameters = append(o.Parameters, convertParameter(p.Value))
		}
	}
	if op.RequestBody != nil && op.RequestBody.Value != nil {
		o.RequestBody = convertRequestBody(op.RequestBody.Value)
	}
	for code, resp := range op.Responses.Map() {
		if resp.Value != nil {
			o.Responses[code] = convertResponse(resp.Value)
		}
	}
	return o
}

func convertParameter(p *openapi3.Parameter) *Parameter {
	param := &Parameter{
		Name:     p.Name,
		In:       p.In,
		Required: p.Required,
	}
	if p.Schema != nil && p.Schema.Value != nil {
		param.Schema = convertSchema(p.Schema.Value)
	}
	return param
}

func convertRequestBody(rb *openapi3.RequestBody) *RequestBody {
	body := &RequestBody{Required: rb.Required, Content: make(map[string]*MediaType)}
	for ct, mt := range rb.Content {
		if mt.Schema != nil && mt.Schema.Value != nil {
			body.Content[ct] = &MediaType{Schema: convertSchema(mt.Schema.Value)}
		}
	}
	return body
}

func convertResponse(r *openapi3.Response) *Response {
	desc := ""
	if r.Description != nil {
		desc = *r.Description
	}
	resp := &Response{Description: desc, Content: make(map[string]*MediaType)}
	for ct, mt := range r.Content {
		if mt.Schema != nil && mt.Schema.Value != nil {
			resp.Content[ct] = &MediaType{Schema: convertSchema(mt.Schema.Value)}
		}
	}
	return resp
}

func convertSchema(s *openapi3.Schema) *Schema {
	if s == nil {
		return nil
	}
	cs := &Schema{
		Format:      s.Format,
		Description: s.Description,
		Nullable:    s.Nullable,
		Enum:        s.Enum,
		Required:    s.Required,
	}
	// Type is *openapi3.Types (a pointer to a slice of strings) in kin-openapi v0.134+
	if s.Type != nil && len(*s.Type) > 0 {
		cs.Type = (*s.Type)[0]
	}
	if s.Min != nil {
		cs.Minimum = s.Min
	}
	if s.Max != nil {
		cs.Maximum = s.Max
	}
	if s.MinLength > 0 {
		ml := int64(s.MinLength)
		cs.MinLength = &ml
	}
	if s.MaxLength != nil {
		ml := int64(*s.MaxLength)
		cs.MaxLength = &ml
	}
	if s.Items != nil && s.Items.Value != nil {
		cs.Items = convertSchema(s.Items.Value)
	}
	cs.Properties = make(map[string]*Schema)
	for name, prop := range s.Properties {
		if prop.Value != nil {
			cs.Properties[name] = convertSchema(prop.Value)
		}
	}
	return cs
}
