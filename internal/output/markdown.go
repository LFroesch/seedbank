package output

import (
	"fmt"
	"strings"
)

// FormatMarkdown formats records as a markdown table.
func FormatMarkdown(records []map[string]any, fields []string) string {
	if len(records) == 0 {
		return ""
	}

	cols := fields
	if len(cols) == 0 {
		cols = sortedKeys(records[0])
	}

	var sb strings.Builder

	// Header
	sb.WriteString("| " + strings.Join(cols, " | ") + " |")
	sb.WriteByte('\n')

	// Separator
	seps := make([]string, len(cols))
	for i := range seps {
		seps[i] = "---"
	}
	sb.WriteString("| " + strings.Join(seps, " | ") + " |")
	sb.WriteByte('\n')

	// Rows
	for _, rec := range records {
		vals := make([]string, len(cols))
		for j, col := range cols {
			v := fmt.Sprintf("%v", rec[col])
			// Escape pipes in values
			v = strings.ReplaceAll(v, "|", "\\|")
			vals[j] = v
		}
		sb.WriteString("| " + strings.Join(vals, " | ") + " |")
		sb.WriteByte('\n')
	}

	return sb.String()
}
