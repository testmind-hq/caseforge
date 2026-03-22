// internal/methodology/owasp_spec.go
package methodology

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/testmind-hq/caseforge/internal/output/schema"
	"github.com/testmind-hq/caseforge/internal/security"
	"github.com/testmind-hq/caseforge/internal/spec"
)

// SecuritySpecTechnique generates cross-operation OWASP cases.
// Covers: API5 (function-level auth), API8 (CORS), API9 (asset management).
type SecuritySpecTechnique struct{}

func NewSecuritySpecTechnique() *SecuritySpecTechnique { return &SecuritySpecTechnique{} }
func (t *SecuritySpecTechnique) Name() string          { return "owasp_api_top10_spec" }

func (t *SecuritySpecTechnique) Generate(s *spec.ParsedSpec) ([]schema.TestCase, error) {
	var cases []schema.TestCase
	cases = append(cases, buildAPI5Cases(s)...)
	cases = append(cases, buildAPI8Cases(s)...)
	cases = append(cases, buildAPI9Cases(s)...)
	return cases, nil
}

// API5: Function-Level Authorization
func buildAPI5Cases(s *spec.ParsedSpec) []schema.TestCase {
	var lowPrivPaths, highPrivPaths []*spec.Operation
	for _, op := range s.Operations {
		p := op.Path
		if strings.Contains(p, "/me") || strings.Contains(p, "/profile") {
			lowPrivPaths = append(lowPrivPaths, op)
		} else if strings.Contains(p, "admin") || op.Method == "DELETE" {
			highPrivPaths = append(highPrivPaths, op)
		}
	}
	if len(lowPrivPaths) == 0 || len(highPrivPaths) == 0 {
		return nil
	}

	var cases []schema.TestCase
	for _, op := range highPrivPaths {
		step := schema.Step{
			ID: "step-1", Title: "access privileged endpoint with regular user token",
			Type: "test", Method: op.Method, Path: op.Path,
			Headers: map[string]string{"Authorization": "Bearer {{user_token}}"},
			Assertions: []schema.Assertion{
				{Target: "status_code", Operator: "eq", Expected: 403},
			},
		}
		id := fmt.Sprintf("TC-%s", uuid.New().String()[:8])
		cases = append(cases, schema.TestCase{
			Schema: schema.SchemaBaseURL, Version: "1", ID: id,
			Title:    fmt.Sprintf("[OWASP-API5] %s %s — 功能级授权缺失", op.Method, op.Path),
			Kind:     "single", Priority: "P0",
			Tags:     []string{"security", "owasp", "api5-function-level-auth"},
			Source:   schema.CaseSource{Technique: "owasp_api_top10_spec", SpecPath: fmt.Sprintf("%s %s", op.Method, op.Path), Rationale: "用普通用户 token 访问高权限接口应返回 403"},
			Steps:       []schema.Step{step},
			GeneratedAt: time.Now(),
		})
	}
	return cases
}

// API8: CORS Misconfiguration (one OPTIONS case per unique path)
func buildAPI8Cases(s *spec.ParsedSpec) []schema.TestCase {
	seenPaths := map[string]bool{}
	var cases []schema.TestCase
	for _, op := range s.Operations {
		if seenPaths[op.Path] {
			continue
		}
		seenPaths[op.Path] = true
		step := schema.Step{
			ID: "step-1", Title: "OPTIONS preflight request",
			Type: "test", Method: "OPTIONS", Path: op.Path,
			Headers: map[string]string{"Origin": "https://evil.example.com"},
			Assertions: []schema.Assertion{
				{Target: "header Access-Control-Allow-Origin", Operator: "ne", Expected: "*"},
			},
		}
		id := fmt.Sprintf("TC-%s", uuid.New().String()[:8])
		cases = append(cases, schema.TestCase{
			Schema: schema.SchemaBaseURL, Version: "1", ID: id,
			Title:    fmt.Sprintf("[OWASP-API8] OPTIONS %s — CORS 安全配置", op.Path),
			Kind:     "single", Priority: "P0",
			Tags:     []string{"security", "owasp", "api8-cors"},
			Source:   schema.CaseSource{Technique: "owasp_api_top10_spec", SpecPath: fmt.Sprintf("OPTIONS %s", op.Path), Rationale: "CORS 响应头 Access-Control-Allow-Origin 不应为 *"},
			Steps:       []schema.Step{step},
			GeneratedAt: time.Now(),
		})
	}
	return cases
}

// API9: Improper Asset Management (old versioned paths should return 404)
func buildAPI9Cases(s *spec.ParsedSpec) []schema.TestCase {
	v1Paths, _ := security.FindVersionedPaths(s.Operations)
	if len(v1Paths) == 0 {
		return nil
	}

	var cases []schema.TestCase
	seen := map[string]bool{}
	for _, op := range s.Operations {
		if !seen[op.Path] && containsPath(v1Paths, op.Path) {
			seen[op.Path] = true
			step := schema.Step{
				ID: "step-1", Title: "access deprecated v1 endpoint",
				Type: "test", Method: op.Method, Path: op.Path,
				Headers: map[string]string{},
				Assertions: []schema.Assertion{
					{Target: "status_code", Operator: "eq", Expected: 404},
				},
			}
			id := fmt.Sprintf("TC-%s", uuid.New().String()[:8])
			cases = append(cases, schema.TestCase{
				Schema: schema.SchemaBaseURL, Version: "1", ID: id,
				Title:    fmt.Sprintf("[OWASP-API9] %s %s — 旧版本 API 应已下线", op.Method, op.Path),
				Kind:     "single", Priority: "P0",
				Tags:     []string{"security", "owasp", "api9-asset-management"},
				Source:   schema.CaseSource{Technique: "owasp_api_top10_spec", SpecPath: fmt.Sprintf("%s %s", op.Method, op.Path), Rationale: "Spec 同时含有 v1/v2 路径，旧版本路径应已下线返回 404"},
				Steps:       []schema.Step{step},
				GeneratedAt: time.Now(),
			})
		}
	}
	return cases
}

func containsPath(paths []string, p string) bool {
	for _, v := range paths {
		if v == p {
			return true
		}
	}
	return false
}
