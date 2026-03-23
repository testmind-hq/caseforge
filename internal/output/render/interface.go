// internal/output/render/interface.go
package render

import "github.com/testmind-hq/caseforge/internal/output/schema"

type Renderer interface {
	Format() string // "hurl"|"markdown"|"csv"|"postman"|"k6"
	Render(cases []schema.TestCase, outDir string) error
}
