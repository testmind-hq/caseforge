// internal/output/writer/writer.go
package writer

import "github.com/testmind-hq/caseforge/internal/output/schema"

type SchemaWriter interface {
	Write(cases []schema.TestCase, outDir string) error
	Read(indexPath string) ([]schema.TestCase, error)
}
