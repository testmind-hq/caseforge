// internal/spec/parser.go
package spec

import (
	"context"
	"fmt"
	"sort"

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
		for name := range doc.Components.SecuritySchemes {
			ps.SecuritySchemes = append(ps.SecuritySchemes, name)
		}
	}

	// Propagate document-level security so per-operation rules can check it.
	for _, req := range doc.Security {
		for schemeName := range req {
			ps.GlobalSecurity = append(ps.GlobalSecurity, schemeName)
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
	// Parse OpenAPI 3.0 Links from each response
	for code, resp := range op.Responses.Map() {
		if resp.Value == nil {
			continue
		}
		linkNames := make([]string, 0, len(resp.Value.Links))
		for name := range resp.Value.Links {
			linkNames = append(linkNames, name)
		}
		sort.Strings(linkNames)
		for _, linkName := range linkNames {
			linkRef := resp.Value.Links[linkName]
			// Skip links that use operationRef (relative JSON pointer); only operationId is supported.
			if linkRef == nil || linkRef.Value == nil || linkRef.Value.OperationID == "" {
				continue
			}
			sl := SpecLink{
				Name:         linkName,
				OperationID:  linkRef.Value.OperationID,
				ResponseCode: code,
				Parameters:   make(map[string]string),
			}
			for paramName, paramExpr := range linkRef.Value.Parameters {
				// Only string runtime expressions (e.g. "$response.body#/id") are captured;
				// non-string literals (integers, booleans) are skipped.
				if s, ok := paramExpr.(string); ok {
					sl.Parameters[paramName] = s
				}
			}
			o.Links = append(o.Links, sl)
		}
	}
	// Sort links for determinism: response-code map iteration order is non-deterministic.
	sort.Slice(o.Links, func(i, j int) bool {
		return o.Links[i].Name < o.Links[j].Name
	})
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
		cmt := &MediaType{}
		if mt.Schema != nil && mt.Schema.Value != nil {
			cmt.Schema = convertSchema(mt.Schema.Value)
		}
		cmt.Example = mt.Example
		if len(mt.Examples) > 0 {
			cmt.Examples = make(map[string]*Example, len(mt.Examples))
			for name, ref := range mt.Examples {
				if ref != nil && ref.Value != nil {
					cmt.Examples[name] = &Example{
						Summary:     ref.Value.Summary,
						Description: ref.Value.Description,
						Value:       ref.Value.Value,
					}
				}
			}
		}
		body.Content[ct] = cmt
	}
	return body
}

func convertResponse(r *openapi3.Response) *Response {
	desc := ""
	if r.Description != nil {
		desc = *r.Description
	}
	resp := &Response{
		Description: desc,
		Content:     make(map[string]*MediaType),
		Headers:     make(map[string]string),
	}
	for ct, mt := range r.Content {
		if mt.Schema != nil && mt.Schema.Value != nil {
			resp.Content[ct] = &MediaType{Schema: convertSchema(mt.Schema.Value)}
		}
	}
	for name, hRef := range r.Headers {
		if hRef.Value != nil && hRef.Value.Schema != nil && hRef.Value.Schema.Value != nil {
			s := hRef.Value.Schema.Value
			typ := ""
			if s.Type != nil && len(*s.Type) > 0 {
				typ = (*s.Type)[0]
			}
			resp.Headers[name] = typ
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
	if s.MinItems > 0 {
		mi := s.MinItems
		cs.MinItems = &mi
	}
	if s.MaxItems != nil {
		cs.MaxItems = s.MaxItems
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
	cs.Example = s.Example
	cs.Pattern = s.Pattern
	return cs
}
