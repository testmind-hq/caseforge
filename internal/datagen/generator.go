// internal/datagen/generator.go
package datagen

import (
	"regexp"
	"strings"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/testmind-hq/caseforge/internal/spec"
)

type BoundaryKind int

const (
	BoundaryMin        BoundaryKind = iota
	BoundaryMinMinusOne             // min - 1 (invalid)
	BoundaryMax
	BoundaryMaxPlusOne // max + 1 (invalid)
)

type Generator struct {
	// llm reserved for Phase 2 semantic data generation
	llm any
	// Pool holds observed real API values; takes priority over synthetic generation.
	Pool *DataPool
}

func NewGenerator(llm any) *Generator {
	return &Generator{llm: llm}
}

// Generate produces a valid fake value for the given schema and field name.
// Three-tier fallback: format → field name → random by type.
func (g *Generator) Generate(s *spec.Schema, fieldName string) any {
	if s == nil {
		return gofakeit.Word()
	}

	// Tier 0: enum — always pick from enum
	if len(s.Enum) > 0 {
		return s.Enum[gofakeit.Number(0, len(s.Enum)-1)]
	}

	// Tier 0.5: pool — prefer real observed values for non-enum fields
	if g != nil && g.Pool != nil {
		if val, ok := g.Pool.ValueFor(fieldName); ok {
			return val
		}
	}

	// Tier 1: format-aware
	if s.Format != "" {
		if val, ok := generateByFormat(s.Format); ok {
			return val
		}
	}

	// Tier 1.5: pattern-aware string generation
	if s.Type == "string" && s.Pattern != "" {
		if val, ok := generateByPattern(s.Pattern); ok {
			return val
		}
	}

	// Tier 2: field name semantic (description provides disambiguation context)
	if fieldName != "" {
		if val, ok := generateByFieldName(fieldName, s.Description); ok {
			return val
		}
	}

	// Tier 3: fallback by type
	return generateByType(s)
}

func generateByType(s *spec.Schema) any {
	switch s.Type {
	case "string":
		return gofakeit.Word()
	case "integer":
		min, max := int64(0), int64(1000)
		if s.Minimum != nil {
			min = int64(*s.Minimum)
		}
		if s.Maximum != nil {
			max = int64(*s.Maximum)
		}
		return int64(gofakeit.Number(int(min), int(max)))
	case "number":
		return gofakeit.Float64Range(0, 1000)
	case "boolean":
		return gofakeit.Bool()
	case "array":
		// Return a single-element array
		if s.Items != nil {
			return []any{generateByType(s.Items)}
		}
		return []any{}
	case "object":
		result := map[string]any{}
		for name, prop := range s.Properties {
			result[name] = generateByType(prop)
		}
		return result
	default:
		return gofakeit.Word()
	}
}

// generateByPattern attempts to produce a string that matches the given regex pattern.
// It uses a best-effort approach: walk the pattern tokens and produce a concrete candidate,
// then verify it matches. Falls back to gofakeit.Word() if the pattern is invalid or
// the candidate doesn't match.
func generateByPattern(pattern string) (string, bool) {
	compiled, err := regexp.Compile(pattern)
	if err != nil {
		// Pattern is invalid; caller falls through to lower tiers.
		return "", false
	}

	candidate := buildPatternCandidate(pattern)

	// Verify the candidate matches
	if compiled.MatchString(candidate) {
		return candidate, true
	}

	// Candidate does not match (e.g. non-capturing groups (?:...), look-ahead
	// (?=...), or other constructs not handled by buildPatternCandidate);
	// caller falls through to lower tiers.
	return "", false
}

// buildPatternCandidate walks a regex pattern string and produces a rough concrete candidate.
// It handles common OpenAPI patterns (anchors, \d/\w/\s escapes, [a-z] classes,
// quantifiers, alternation, grouping) without requiring external packages.
// Non-capturing groups (?:...) and zero-width assertions are not handled and
// will cause the verification step in generateByPattern to fail, triggering fallback.
func buildPatternCandidate(pattern string) string {
	var out strings.Builder
	i := 0
	n := len(pattern)

	// readToken reads the next "atom" (a character, escape, char class, or group)
	// and returns (the representative string, next index).
	var readToken func(idx int) (string, int)
	readToken = func(idx int) (string, int) {
		if idx >= n {
			return "", idx
		}
		ch := pattern[idx]
		switch ch {
		case '^', '$':
			// anchors — skip
			return "", idx + 1
		case '.':
			return "a", idx + 1
		case '\\':
			if idx+1 >= n {
				return "", idx + 1
			}
			next := pattern[idx+1]
			switch next {
			case 'd':
				return "1", idx + 2
			case 'D':
				return "a", idx + 2
			case 'w':
				return "a", idx + 2
			case 'W':
				return "1", idx + 2
			case 's':
				return " ", idx + 2
			case 'S':
				return "a", idx + 2
			default:
				return string(next), idx + 2
			}
		case '[':
			// character class
			end := strings.Index(pattern[idx:], "]")
			if end < 0 {
				return "a", idx + 1
			}
			class := pattern[idx+1 : idx+end]
			newIdx := idx + end + 1
			if len(class) == 0 {
				return "a", newIdx
			}
			negated := false
			if class[0] == '^' {
				negated = true
				class = class[1:]
			}
			if negated || len(class) == 0 {
				return "a", newIdx
			}
			// Pick first printable char in class (skip range markers)
			if len(class) >= 3 && class[1] == '-' {
				// range like a-z: return start of range
				return string(class[0]), newIdx
			}
			return string(class[0]), newIdx
		case '(':
			// group: find matching ), generate inner content
			depth := 0
			j := idx
			for j < n {
				if pattern[j] == '(' {
					depth++
				} else if pattern[j] == ')' {
					depth--
					if depth == 0 {
						break
					}
				}
				j++
			}
			inner := pattern[idx+1 : j]
			// Handle alternation in group: take first alternative
			if pipe := strings.Index(inner, "|"); pipe >= 0 {
				inner = inner[:pipe]
			}
			return buildPatternCandidate(inner), j + 1
		case '|':
			// alternation at top level: stop (take what we have so far)
			return "", n // signal end
		default:
			return string(ch), idx + 1
		}
	}

	for i < n {
		token, next := readToken(i)
		if next == n && token == "" && pattern[i] == '|' {
			// top-level alternation: stop
			break
		}

		// Check for quantifier after token
		q := ""
		if next < n {
			switch pattern[next] {
			case '+':
				q = "+"
				next++
			case '*':
				q = "*"
				next++
			case '?':
				q = "?"
				next++
			case '{':
				// {n} or {n,m}
				end := strings.Index(pattern[next:], "}")
				if end >= 0 {
					q = pattern[next : next+end+1]
					next = next + end + 1
				}
			}
		}

		// Determine repeat count from quantifier
		repeat := 1
		if q != "" {
			switch q {
			case "+", "*":
				repeat = 1
			case "?":
				repeat = 1
			default:
				// {n} or {n,m}
				inner := q[1 : len(q)-1]
				if comma := strings.Index(inner, ","); comma >= 0 {
					// {n,m} — use n (minimum), but at least 1
					nStr := inner[:comma]
					repeat = parseIntOr(nStr, 1)
					if repeat == 0 {
						repeat = 1
					}
				} else {
					repeat = parseIntOr(inner, 1)
					if repeat == 0 {
						repeat = 1
					}
				}
			}
		}

		for r := 0; r < repeat; r++ {
			out.WriteString(token)
		}
		i = next
	}

	return out.String()
}

// parseIntOr parses a decimal string, returning def on error or empty input.
func parseIntOr(s string, def int) int {
	if len(s) == 0 {
		return def
	}
	val := 0
	for _, c := range s {
		if c < '0' || c > '9' {
			return def
		}
		val = val*10 + int(c-'0')
	}
	return val
}

// GenerateBoundary produces a boundary value for numeric/string schemas.
func (g *Generator) GenerateBoundary(s *spec.Schema, kind BoundaryKind) any {
	if s == nil {
		return nil
	}
	switch s.Type {
	case "integer":
		switch kind {
		case BoundaryMin:
			if s.Minimum != nil {
				return int64(*s.Minimum)
			}
			return int64(0)
		case BoundaryMinMinusOne:
			if s.Minimum != nil {
				return int64(*s.Minimum) - 1
			}
			return int64(-1)
		case BoundaryMax:
			if s.Maximum != nil {
				return int64(*s.Maximum)
			}
			return int64(1000)
		case BoundaryMaxPlusOne:
			if s.Maximum != nil {
				return int64(*s.Maximum) + 1
			}
			return int64(1001)
		}
	case "string":
		switch kind {
		case BoundaryMin:
			if s.MinLength != nil {
				return gofakeit.LetterN(uint(*s.MinLength))
			}
			return ""
		case BoundaryMinMinusOne:
			if s.MinLength != nil && *s.MinLength > 0 {
				return gofakeit.LetterN(uint(*s.MinLength - 1))
			}
			return ""
		case BoundaryMax:
			if s.MaxLength != nil {
				return gofakeit.LetterN(uint(*s.MaxLength))
			}
			return gofakeit.LetterN(255)
		case BoundaryMaxPlusOne:
			if s.MaxLength != nil {
				return gofakeit.LetterN(uint(*s.MaxLength + 1))
			}
			return gofakeit.LetterN(256)
		}
	}
	return nil
}
