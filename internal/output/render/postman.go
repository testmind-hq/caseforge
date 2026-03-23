// internal/output/render/postman.go
package render

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
	"github.com/testmind-hq/caseforge/internal/output/schema"
)

// PostmanRenderer renders test cases as a Postman Collection v2.1 JSON file.
type PostmanRenderer struct{}

func NewPostmanRenderer() *PostmanRenderer { return &PostmanRenderer{} }
func (r *PostmanRenderer) Format() string  { return "postman" }

func (r *PostmanRenderer) Render(cases []schema.TestCase, outDir string) error {
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return fmt.Errorf("creating output dir: %w", err)
	}

	collection := map[string]any{
		"info": map[string]any{
			"name":        "CaseForge Generated",
			"_postman_id": uuid.New().String(),
			"schema":      "https://schema.getpostman.com/json/collection/v2.1.0/collection.json",
		},
		"variable": []any{
			map[string]any{"key": "base_url", "value": "http://localhost", "type": "string"},
		},
		"item": buildPostmanItems(cases),
	}

	data, err := json.MarshalIndent(collection, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling collection: %w", err)
	}
	return os.WriteFile(filepath.Join(outDir, "collection.json"), data, 0644)
}

func buildPostmanItems(cases []schema.TestCase) []any {
	var items []any
	for _, tc := range cases {
		if tc.Kind == "chain" {
			items = append(items, buildPostmanFolder(tc))
		} else {
			if len(tc.Steps) > 0 {
				items = append(items, buildPostmanRequest(tc.Title, tc.Steps[0]))
			}
		}
	}
	return items
}

func buildPostmanFolder(tc schema.TestCase) map[string]any {
	var subItems []any
	for _, step := range tc.Steps {
		subItems = append(subItems, buildPostmanRequest(step.Title, step))
	}
	return map[string]any{
		"name": tc.Title,
		"item": subItems,
	}
}

func buildPostmanRequest(name string, step schema.Step) map[string]any {
	req := map[string]any{
		"method": step.Method,
		"url":    buildPostmanURL(step.Path),
		"header": buildPostmanHeaders(step.Headers),
	}
	if step.Body != nil {
		data, err := json.Marshal(step.Body)
		rawBody := string(data)
		if err != nil {
			rawBody = fmt.Sprintf("// ERROR: body serialization failed: %v", err)
		}
		req["body"] = map[string]any{
			"mode": "raw",
			"raw":  rawBody,
			"options": map[string]any{
				"raw": map[string]any{"language": "json"},
			},
		}
	}

	item := map[string]any{
		"name":    name,
		"request": req,
	}

	script := buildTestScript(step)
	if script != "" {
		item["event"] = []any{
			map[string]any{
				"listen": "test",
				"script": map[string]any{
					"type": "text/javascript",
					"exec": strings.Split(script, "\n"),
				},
			},
		}
	}
	return item
}

func buildPostmanURL(path string) map[string]any {
	rawURL := "{{base_url}}" + path
	segments := strings.Split(strings.Trim(path, "/"), "/")
	pathParts := make([]any, 0, len(segments))
	for _, seg := range segments {
		if strings.HasPrefix(seg, "{") && strings.HasSuffix(seg, "}") {
			// OpenAPI {param} → Postman :param
			pathParts = append(pathParts, ":"+seg[1:len(seg)-1])
		} else {
			pathParts = append(pathParts, seg)
		}
	}
	return map[string]any{
		"raw":  rawURL,
		"host": []any{"{{base_url}}"},
		"path": pathParts,
	}
}

func buildPostmanHeaders(headers map[string]string) []any {
	var result []any
	for k, v := range headers {
		result = append(result, map[string]any{"key": k, "value": v})
	}
	return result
}

func buildTestScript(step schema.Step) string {
	var lines []string

	for _, a := range step.Assertions {
		switch {
		case a.Target == "status_code":
			code := 200
			if f, ok := a.Expected.(float64); ok {
				code = int(f)
			} else if i, ok := a.Expected.(int); ok {
				code = i
			}
			lines = append(lines,
				fmt.Sprintf("pm.test(\"status is %d\", function () {", code),
				fmt.Sprintf("    pm.response.to.have.status(%d);", code),
				"});",
			)
		case strings.HasPrefix(a.Target, "jsonpath "):
			expr := strings.TrimPrefix(a.Target, "jsonpath ")
			fieldPath := strings.TrimPrefix(expr, "$.")
			valJS := postmanJSValue(a.Expected)
			switch a.Operator {
			case "eq":
				lines = append(lines,
					fmt.Sprintf("pm.test(%q, function () {", a.Target+" eq "+fmt.Sprint(a.Expected)),
					"    var jsonData = pm.response.json();",
					fmt.Sprintf("    pm.expect(jsonData.%s).to.eql(%s);", fieldPath, valJS),
					"});",
				)
			case "ne":
				if a.Expected == nil {
					lines = append(lines,
						fmt.Sprintf("pm.test(%q, function () {", a.Target+" not exists"),
						"    var jsonData = pm.response.json();",
						fmt.Sprintf("    pm.expect(jsonData.%s).to.not.exist;", fieldPath),
						"});",
					)
				} else {
					lines = append(lines,
						fmt.Sprintf("pm.test(%q, function () {", a.Target+" ne "+fmt.Sprint(a.Expected)),
						"    var jsonData = pm.response.json();",
						fmt.Sprintf("    pm.expect(jsonData.%s).to.not.eql(%s);", fieldPath, valJS),
						"});",
					)
				}
			case "contains":
				lines = append(lines,
					fmt.Sprintf("pm.test(%q, function () {", a.Target+" contains "+fmt.Sprint(a.Expected)),
					"    var jsonData = pm.response.json();",
					fmt.Sprintf("    pm.expect(String(jsonData.%s)).to.include(%s);", fieldPath, valJS),
					"});",
				)
			}
		case strings.HasPrefix(a.Target, "header "):
			headerName := strings.TrimPrefix(a.Target, "header ")
			valJS := postmanJSValue(a.Expected)
			switch a.Operator {
			case "eq":
				lines = append(lines,
					fmt.Sprintf("pm.test(%q, function () {", a.Target+" eq "+fmt.Sprint(a.Expected)),
					fmt.Sprintf("    pm.expect(pm.response.headers.get(%q)).to.eql(%s);", headerName, valJS),
					"});",
				)
			case "ne":
				lines = append(lines,
					fmt.Sprintf("pm.test(%q, function () {", a.Target+" ne "+fmt.Sprint(a.Expected)),
					fmt.Sprintf("    pm.expect(pm.response.headers.get(%q)).to.not.eql(%s);", headerName, valJS),
					"});",
				)
			}
		}
	}

	// Add capture scripts
	if len(step.Captures) > 0 {
		lines = append(lines, "var jsonData = pm.response.json();")
		for _, cap := range step.Captures {
			if strings.HasPrefix(cap.From, "jsonpath $.") {
				fieldPath := strings.TrimPrefix(cap.From, "jsonpath $.")
				lines = append(lines,
					fmt.Sprintf("pm.environment.set(%q, jsonData.%s);", cap.Name, fieldPath),
				)
			}
		}
	}

	return strings.Join(lines, "\n")
}

func postmanJSValue(v any) string {
	if v == nil {
		return "null"
	}
	switch val := v.(type) {
	case string:
		return fmt.Sprintf("%q", val)
	case bool:
		if val {
			return "true"
		}
		return "false"
	default:
		return fmt.Sprintf("%v", v)
	}
}
