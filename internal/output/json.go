package output

import (
	"encoding/json"
	"strings"
)

// FormatJSON formats records as pretty-printed JSON.
func FormatJSON(records []map[string]any, fields []string) string {
	filtered := filterFields(records, fields)
	b, err := json.MarshalIndent(filtered, "", "  ")
	if err != nil {
		return "error: " + err.Error()
	}
	return string(b)
}

// FormatJSONLines formats records as newline-delimited JSON.
func FormatJSONLines(records []map[string]any, fields []string) string {
	filtered := filterFields(records, fields)
	var sb strings.Builder
	for _, rec := range filtered {
		b, _ := json.Marshal(rec)
		sb.Write(b)
		sb.WriteByte('\n')
	}
	return sb.String()
}
