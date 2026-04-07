package output

import (
	"strings"
	"testing"
)

var testRecords = []map[string]any{
	{"name": "Alice", "age": 30, "active": true},
	{"name": "Bob", "age": 25, "active": false},
}

var testFields = []string{"name", "age", "active"}

func TestFormatJSON(t *testing.T) {
	out := FormatJSON(testRecords, testFields)
	if !strings.Contains(out, `"name": "Alice"`) {
		t.Error("JSON output missing Alice")
	}
	if !strings.HasPrefix(out, "[") || !strings.HasSuffix(strings.TrimSpace(out), "]") {
		t.Error("JSON output not wrapped in array brackets")
	}
}

func TestFormatJSONLines(t *testing.T) {
	out := FormatJSONLines(testRecords, testFields)
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != 2 {
		t.Errorf("expected 2 JSONL lines, got %d", len(lines))
	}
}

func TestFormatCSV(t *testing.T) {
	out := FormatCSV(testRecords, testFields)
	lines := strings.Split(strings.TrimSpace(out), "\n")
	if len(lines) != 3 { // header + 2 data rows
		t.Errorf("expected 3 CSV lines (header + 2), got %d", len(lines))
	}
	if !strings.Contains(lines[0], "name") {
		t.Error("CSV header missing field names")
	}
}

func TestFormatMarkdown(t *testing.T) {
	out := FormatMarkdown(testRecords, testFields)
	if !strings.Contains(out, "|") {
		t.Error("markdown output missing table pipes")
	}
	if !strings.Contains(out, "---") {
		t.Error("markdown output missing separator")
	}
}

func TestFormatSQL(t *testing.T) {
	out := FormatSQL(testRecords, testFields, "users")
	if !strings.Contains(out, "INSERT INTO users") {
		t.Error("SQL output missing INSERT INTO users")
	}
	if !strings.Contains(out, "'Alice'") {
		t.Error("SQL output missing quoted string value")
	}
	if strings.Count(out, "INSERT INTO") != 2 {
		t.Errorf("expected 2 INSERT INTO statements, got %d", strings.Count(out, "INSERT INTO"))
	}
}

func TestFormatPretty(t *testing.T) {
	out := FormatPretty(testRecords, testFields)
	if !strings.Contains(out, "name") {
		t.Error("pretty output missing field labels")
	}
	if !strings.Contains(out, "Alice") {
		t.Error("pretty output missing values")
	}
}

func TestFormatEmptyRecords(t *testing.T) {
	for _, fn := range []struct {
		name string
		call func() string
	}{
		{"JSON", func() string { return FormatJSON(nil, testFields) }},
		{"JSONL", func() string { return FormatJSONLines(nil, testFields) }},
		{"CSV", func() string { return FormatCSV(nil, testFields) }},
		{"SQL", func() string { return FormatSQL(nil, testFields, "t") }},
	} {
		t.Run(fn.name, func(t *testing.T) {
			out := fn.call()
			// Should not panic, should return something reasonable
			_ = out
		})
	}
}

func TestFieldFiltering(t *testing.T) {
	out := FormatJSON(testRecords, []string{"name"})
	if strings.Contains(out, `"age"`) {
		t.Error("filtered JSON should not contain age field")
	}
	if !strings.Contains(out, `"name"`) {
		t.Error("filtered JSON should contain name field")
	}
}
