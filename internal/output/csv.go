package output

import (
	"fmt"
	"sort"
	"strings"
)

// FormatCSV formats records as CSV with a header row.
func FormatCSV(records []map[string]any, fields []string) string {
	if len(records) == 0 {
		return ""
	}

	cols := fields
	if len(cols) == 0 {
		cols = sortedKeys(records[0])
	}

	var sb strings.Builder

	// Header
	sb.WriteString(strings.Join(cols, ","))
	sb.WriteByte('\n')

	// Rows
	for _, rec := range records {
		vals := make([]string, len(cols))
		for j, col := range cols {
			vals[j] = csvEscape(fmt.Sprintf("%v", rec[col]))
		}
		sb.WriteString(strings.Join(vals, ","))
		sb.WriteByte('\n')
	}

	return sb.String()
}

func csvEscape(s string) string {
	if strings.ContainsAny(s, ",\"\n") {
		return "\"" + strings.ReplaceAll(s, "\"", "\"\"") + "\""
	}
	return s
}

func sortedKeys(m map[string]any) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
