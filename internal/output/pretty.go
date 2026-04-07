package output

import (
	"fmt"
	"strings"
)

// FormatPretty renders records as human-readable cards, one per record.
// Each card shows field labels aligned with their values, separated by dividers.
func FormatPretty(records []map[string]any, fields []string) string {
	if len(records) == 0 {
		return "(no records)"
	}

	cols := fields
	if len(cols) == 0 && len(records) > 0 {
		cols = sortedKeys(records[0])
	}

	// Find max label width for alignment
	maxLabel := 0
	for _, c := range cols {
		if len(c) > maxLabel {
			maxLabel = len(c)
		}
	}

	var sb strings.Builder
	divider := strings.Repeat("─", maxLabel+20)

	for i, rec := range records {
		sb.WriteString(fmt.Sprintf("  #%d\n", i+1))
		for _, col := range cols {
			label := col + strings.Repeat(" ", maxLabel-len(col))
			val := fmt.Sprintf("%v", rec[col])
			sb.WriteString(fmt.Sprintf("  %s  %s\n", label, val))
		}
		if i < len(records)-1 {
			sb.WriteString("  " + divider + "\n")
		}
	}

	return sb.String()
}
