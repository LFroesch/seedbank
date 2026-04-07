package output

// filterFields returns records with only the specified fields.
// If fields is empty, returns records as-is.
func filterFields(records []map[string]any, fields []string) []map[string]any {
	if len(fields) == 0 {
		return records
	}

	result := make([]map[string]any, len(records))
	for i, rec := range records {
		filtered := make(map[string]any, len(fields))
		for _, f := range fields {
			if v, ok := rec[f]; ok {
				filtered[f] = v
			}
		}
		result[i] = filtered
	}
	return result
}

// Format is the supported output format type.
type Format int

const (
	JSON Format = iota
	JSONLines
	CSV
	Markdown
	SQL
)

// FormatNames maps format to display name.
var FormatNames = []string{"JSON", "JSONL", "CSV", "Markdown", "SQL"}
