// internal/output/render/markdown.go
package render

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/testmind-hq/caseforge/internal/output/schema"
)

type MarkdownRenderer struct{}

func NewMarkdownRenderer() *MarkdownRenderer { return &MarkdownRenderer{} }
func (r *MarkdownRenderer) Format() string   { return "markdown" }

func (r *MarkdownRenderer) Render(cases []schema.TestCase, outDir string) error {
	if err := os.MkdirAll(outDir, 0755); err != nil {
		return err
	}
	var b strings.Builder
	b.WriteString("# CaseForge Test Case Report\n\n")
	b.WriteString(fmt.Sprintf("**Total cases:** %d\n\n---\n\n", len(cases)))

	for _, tc := range cases {
		b.WriteString(fmt.Sprintf("## %s — %s\n\n", tc.ID, tc.Title))
		b.WriteString(fmt.Sprintf("| Field | Value |\n|---|---|\n"))
		b.WriteString(fmt.Sprintf("| Priority | %s |\n", tc.Priority))
		b.WriteString(fmt.Sprintf("| Technique | %s |\n", tc.Source.Technique))
		b.WriteString(fmt.Sprintf("| Spec Path | `%s` |\n", tc.Source.SpecPath))
		b.WriteString(fmt.Sprintf("| Rationale | %s |\n\n", tc.Source.Rationale))

		for _, step := range tc.Steps {
			b.WriteString(fmt.Sprintf("### %s `%s`\n\n", step.Method, step.Path))
			b.WriteString("**Assertions:**\n\n")
			for _, a := range step.Assertions {
				b.WriteString(fmt.Sprintf("- `%s` %s `%v`\n", a.Target, a.Operator, a.Expected))
			}
			b.WriteString("\n")
		}
		b.WriteString("---\n\n")
	}

	return os.WriteFile(filepath.Join(outDir, "report.md"), []byte(b.String()), 0644)
}
