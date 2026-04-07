package output

import (
	"fmt"
	"strings"
)

// FormatSQL formats records as SQL INSERT statements.
func FormatSQL(records []map[string]any, fields []string, tableName string) string {
	if len(records) == 0 {
		return ""
	}

	cols := fields
	if len(cols) == 0 {
		cols = sortedKeys(records[0])
	}

	var sb strings.Builder

	// CREATE TABLE hint
	sb.WriteString(fmt.Sprintf("-- Table: %s\n", tableName))
	sb.WriteString(fmt.Sprintf("-- INSERT %d rows\n\n", len(records)))

	for _, rec := range records {
		vals := make([]string, len(cols))
		for j, col := range cols {
			vals[j] = sqlValue(rec[col])
		}
		sb.WriteString(fmt.Sprintf("INSERT INTO %s (%s) VALUES (%s);\n",
			tableName,
			strings.Join(cols, ", "),
			strings.Join(vals, ", "),
		))
	}

	return sb.String()
}

func sqlValue(v any) string {
	switch val := v.(type) {
	case int:
		return fmt.Sprintf("%d", val)
	case int64:
		return fmt.Sprintf("%d", val)
	case float64:
		return fmt.Sprintf("%g", val)
	case bool:
		if val {
			return "TRUE"
		}
		return "FALSE"
	default:
		s := fmt.Sprintf("%v", val)
		s = strings.ReplaceAll(s, "'", "''")
		return "'" + s + "'"
	}
}
