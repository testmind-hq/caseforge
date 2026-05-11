// internal/mutation/render.go
package mutation

import (
	"fmt"
	"html"
	"sort"
	"strings"
)

// RenderMarkdown returns a Markdown mutation report string.
func RenderMarkdown(run MutationRun) string {
	pct := 0
	if run.TotalRuns > 0 {
		pct = int(run.MutationScore * 100)
	}
	var b strings.Builder
	fmt.Fprintf(&b, "# Mutation Report — %s\n\n", run.GeneratedAt)
	fmt.Fprintf(&b, "Mutation Score: **%d/%d killed (%d%%)** · Survivors: **%d**\n\n",
		run.Killed, run.TotalRuns, pct, run.Survivors)

	b.WriteString(renderOperationTable(run))

	if len(run.Clusters) > 0 {
		b.WriteString("## Survivor Summary (by risk)\n\n")
		b.WriteString("| Case | Title | Risk | Survived Operators |\n")
		b.WriteString("|------|-------|------|--------------------|\n")
		for _, c := range run.Clusters {
			fmt.Fprintf(&b, "| %s | %s | %.2f | %s |\n",
				c.CaseID, c.Title, c.RiskScore, strings.Join(c.Operators, ", "))
		}
		b.WriteString("\n")
	}

	if len(run.Feedback) > 0 {
		b.WriteString("## Suggested Assertions\n\n")
		for _, item := range run.Feedback {
			fmt.Fprintf(&b, "### %s — %s\n\n", item.CaseID, item.Title)
			for _, a := range item.SuggestedAssertions {
				if a.Expected != nil {
					fmt.Fprintf(&b, "- `%s` `%s` `%v`\n", a.Target, a.Operator, a.Expected)
				} else {
					fmt.Fprintf(&b, "- `%s` `%s`\n", a.Target, a.Operator)
				}
			}
			b.WriteString("\n")
		}
	}

	return b.String()
}

// renderOperationTable returns the per-operation Markdown table, or "" if no scores.
func renderOperationTable(run MutationRun) string {
	if len(run.OperationScores) == 0 {
		return ""
	}
	var b strings.Builder
	fmt.Fprintf(&b, "## Per-Operation Mutation Score\n\n")
	fmt.Fprintf(&b, "| Operation | Score | Killed | Survivors |\n")
	fmt.Fprintf(&b, "|-----------|-------|--------|-----------|\n")
	for _, op := range run.OperationScores {
		pct := int(op.MutationScore * 100)
		badge := ""
		switch {
		case op.MutationScore == 1.0:
			badge = " ✓"
		case op.MutationScore < 0.7:
			badge = " ⚠"
		}
		fmt.Fprintf(&b, "| %s | %d%%%s | %d/%d | %d |\n",
			op.Operation, pct, badge, op.Killed, op.TotalRuns, op.Survivors)
	}
	b.WriteString("\n")
	return b.String()
}

// RenderHTML returns a self-contained HTML mutation report (no external resources).
func RenderHTML(run MutationRun) string {
	// build operator×case result grid
	type cell struct{ survived, exists bool }
	grid := map[string]map[string]cell{}
	caseTitle := map[string]string{}
	for _, r := range run.Results {
		if grid[r.Operator] == nil {
			grid[r.Operator] = map[string]cell{}
		}
		grid[r.Operator][r.CaseID] = cell{survived: r.Survived, exists: true}
		caseTitle[r.CaseID] = r.Title
	}
	caseIDs := make([]string, 0, len(caseTitle))
	for id := range caseTitle {
		caseIDs = append(caseIDs, id)
	}
	sort.Strings(caseIDs)

	pct := 0
	if run.TotalRuns > 0 {
		pct = int(run.MutationScore * 100)
	}

	var b strings.Builder
	b.WriteString("<!DOCTYPE html><html><head><meta charset=\"utf-8\">\n")
	b.WriteString("<title>Mutation Report</title>\n")
	b.WriteString(`<style>
body{font-family:sans-serif;margin:2rem;color:#111}
h1{margin-bottom:.5rem}
.stats{display:flex;gap:2rem;margin:1rem 0}
.stat-label{color:#666;font-size:.8rem}
table{border-collapse:collapse;display:block;overflow-x:auto;font-size:.8rem}
th{background:#f3f4f6;padding:6px 8px;text-align:left;border:1px solid #e5e7eb;white-space:nowrap}
td{padding:4px 8px;border:1px solid #e5e7eb;white-space:nowrap}
.killed{background:#dcfce7;color:#166534;text-align:center}
.survived{background:#fee2e2;color:#991b1b;text-align:center}
.na{background:#f9fafb;color:#9ca3af;text-align:center}
h2{margin-top:2rem}
</style></head><body>
`)
	b.WriteString("<h1>Mutation Report</h1>\n")
	b.WriteString("<div class=\"stats\">")
	fmt.Fprintf(&b, "<div><div class=\"stat-label\">Score</div><strong>%d%%</strong></div>", pct)
	fmt.Fprintf(&b, "<div><div class=\"stat-label\">Killed/Total</div><strong>%d/%d</strong></div>", run.Killed, run.TotalRuns)
	fmt.Fprintf(&b, "<div><div class=\"stat-label\">Survivors</div><strong>%d</strong></div>", run.Survivors)
	fmt.Fprintf(&b, "<div><div class=\"stat-label\">Generated</div><strong>%s</strong></div>", html.EscapeString(run.GeneratedAt))
	b.WriteString("</div>\n")

	b.WriteString("<h2>Operator × Case Heatmap</h2>\n<table>\n<tr><th>Operator</th>")
	for _, id := range caseIDs {
		fmt.Fprintf(&b, "<th title=\"%s\">%s</th>", html.EscapeString(caseTitle[id]), html.EscapeString(id))
	}
	b.WriteString("</tr>\n")
	for _, op := range run.Operators {
		fmt.Fprintf(&b, "<tr><td>%s</td>", html.EscapeString(op))
		for _, id := range caseIDs {
			c := grid[op][id]
			if !c.exists {
				b.WriteString("<td class=\"na\">—</td>")
			} else if c.survived {
				b.WriteString("<td class=\"survived\">✗</td>")
			} else {
				b.WriteString("<td class=\"killed\">✓</td>")
			}
		}
		b.WriteString("</tr>\n")
	}
	b.WriteString("</table>\n")

	if len(run.Clusters) > 0 {
		b.WriteString("<h2>Survivors by Risk</h2>\n<table>\n<tr><th>Case</th><th>Title</th><th>Risk</th><th>Operators</th></tr>\n")
		for _, c := range run.Clusters {
			fmt.Fprintf(&b, "<tr><td>%s</td><td>%s</td><td>%.0f%%</td><td>%s</td></tr>\n",
				html.EscapeString(c.CaseID),
				html.EscapeString(c.Title),
				c.RiskScore*100,
				html.EscapeString(strings.Join(c.Operators, ", ")))
		}
		b.WriteString("</table>\n")
	}

	b.WriteString("</body></html>\n")
	return b.String()
}
