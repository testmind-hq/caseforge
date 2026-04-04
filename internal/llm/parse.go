// internal/llm/parse.go
package llm

import "strings"

// ExtractJSON extracts a JSON value (object or array) from LLM response text.
// It handles three common LLM output patterns:
//  1. Bare JSON:           [{"a":1}]
//  2. Markdown-fenced:     ```json\n[{"a":1}]\n```
//  3. Text preamble:       "Here are the results:\n[{"a":1}]"
//
// Extraction uses bracket-count tracking, so nested structures are handled correctly.
// Returns the input unchanged if no JSON structure is found.
func ExtractJSON(text string) string {
	text = strings.TrimSpace(text)
	if text == "" {
		return text
	}

	// Step 1: strip markdown fences (```json ... ``` or ``` ... ```)
	if strings.HasPrefix(text, "```") {
		if idx := strings.Index(text, "\n"); idx != -1 {
			text = text[idx+1:]
			if idx2 := strings.LastIndex(text, "```"); idx2 != -1 {
				text = strings.TrimSpace(text[:idx2])
			}
		}
	}

	// Step 2: find the first JSON token character
	start := -1
	var open, close rune
	for i, r := range text {
		if r == '[' || r == '{' {
			start = i
			open = r
			if r == '[' {
				close = ']'
			} else {
				close = '}'
			}
			break
		}
	}
	if start == -1 {
		return text // no JSON structure — return as-is, json.Unmarshal will error
	}

	// Step 3: bracket-count scan to find matching close
	depth := 0
	for i, r := range text[start:] {
		switch r {
		case open:
			depth++
		case close:
			depth--
			if depth == 0 {
				return text[start : start+i+1]
			}
		}
	}
	// Unclosed bracket — return from start to end (let caller handle the error)
	return text[start:]
}
