package main

import (
	"fmt"
	"math/rand"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/LFroesch/seedbank/internal/generator"
	"github.com/LFroesch/seedbank/internal/output"
)

type schemaColumn struct {
	Name string
	Type string
}

type schemaParseResult struct {
	Table   string
	Columns []schemaColumn
}

type columnMapping struct {
	ColumnName string
	GenName    string
	FieldName  string
}

func runSchema(schemaPath string, count int, format, table, outPath string, seed int64) {
	content, err := os.ReadFile(schemaPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to read schema file: %v\n", err)
		os.Exit(1)
	}

	parsed, err := parseCreateTableSQL(string(content))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to parse CREATE TABLE: %v\n", err)
		os.Exit(1)
	}

	if table == "" {
		table = parsed.Table
	}

	mappings := mapSchemaColumns(parsed.Columns)

	if seed == 0 {
		seed = time.Now().UnixNano()
	}
	rng := rand.New(rand.NewSource(seed))

	records, fields := generateSchemaRecords(mappings, count, rng)
	out := formatRecords(records, fields, format, table)

	if err := writeOutput(outPath, out); err != nil {
		fmt.Fprintf(os.Stderr, "write failed: %v\n", err)
		os.Exit(1)
	}
}

func parseCreateTableSQL(sql string) (schemaParseResult, error) {
	cleaned := stripSQLComments(sql)
	re := regexp.MustCompile(`(?is)create\s+table\s+(?:if\s+not\s+exists\s+)?([^\s(]+)\s*\(`)
	m := re.FindStringSubmatchIndex(cleaned)
	if m == nil {
		return schemaParseResult{}, fmt.Errorf("no CREATE TABLE statement found")
	}

	rawName := cleaned[m[2]:m[3]]
	tableName := normalizeIdent(rawName)
	openParenIdx := m[1] - 1

	closeParenIdx, err := findMatchingParen(cleaned, openParenIdx)
	if err != nil {
		return schemaParseResult{}, err
	}

	definitionBlock := cleaned[openParenIdx+1 : closeParenIdx]
	parts := splitTopLevelByComma(definitionBlock)

	var cols []schemaColumn
	for _, part := range parts {
		p := strings.TrimSpace(part)
		if p == "" || isConstraintDefinition(p) {
			continue
		}
		col, ok := parseColumnDefinition(p)
		if ok {
			cols = append(cols, col)
		}
	}

	if len(cols) == 0 {
		return schemaParseResult{}, fmt.Errorf("no columns found in CREATE TABLE")
	}

	return schemaParseResult{
		Table:   tableName,
		Columns: cols,
	}, nil
}

func stripSQLComments(s string) string {
	lineComment := regexp.MustCompile(`(?m)--.*$`)
	blockComment := regexp.MustCompile(`(?s)/\*.*?\*/`)
	out := lineComment.ReplaceAllString(s, "")
	return blockComment.ReplaceAllString(out, "")
}

func findMatchingParen(s string, openIdx int) (int, error) {
	depth := 0
	for i := openIdx; i < len(s); i++ {
		switch s[i] {
		case '(':
			depth++
		case ')':
			depth--
			if depth == 0 {
				return i, nil
			}
		}
	}
	return -1, fmt.Errorf("unbalanced parentheses in CREATE TABLE")
}

func splitTopLevelByComma(s string) []string {
	var out []string
	start := 0
	depth := 0
	for i := 0; i < len(s); i++ {
		switch s[i] {
		case '(':
			depth++
		case ')':
			if depth > 0 {
				depth--
			}
		case ',':
			if depth == 0 {
				out = append(out, s[start:i])
				start = i + 1
			}
		}
	}
	if start < len(s) {
		out = append(out, s[start:])
	}
	return out
}

func isConstraintDefinition(def string) bool {
	s := strings.ToLower(strings.TrimSpace(def))
	prefixes := []string{
		"primary key",
		"foreign key",
		"unique",
		"constraint",
		"check",
		"index",
		"key ",
	}
	for _, p := range prefixes {
		if strings.HasPrefix(s, p) {
			return true
		}
	}
	return false
}

func parseColumnDefinition(def string) (schemaColumn, bool) {
	tokens := strings.Fields(def)
	if len(tokens) < 2 {
		return schemaColumn{}, false
	}

	colName := normalizeIdent(tokens[0])
	rest := strings.ToLower(strings.TrimSpace(def[len(tokens[0]):]))
	typeName := extractTypeName(rest)
	if typeName == "" {
		return schemaColumn{}, false
	}

	return schemaColumn{
		Name: colName,
		Type: typeName,
	}, true
}

func extractTypeName(rest string) string {
	r := strings.TrimSpace(rest)
	if r == "" {
		return ""
	}

	stopWords := []string{
		" not null", " null", " default ", " primary key", " unique",
		" references ", " check ", " constraint ", " collate ", " generated ",
	}
	cut := len(r)
	for _, w := range stopWords {
		if idx := strings.Index(r, w); idx >= 0 && idx < cut {
			cut = idx
		}
	}
	r = strings.TrimSpace(r[:cut])
	// Keep type base while preserving common multi-word type names.
	multiWord := []string{
		"double precision",
		"timestamp with time zone",
		"timestamp without time zone",
		"time with time zone",
		"time without time zone",
		"character varying",
	}
	for _, t := range multiWord {
		if strings.HasPrefix(r, t) {
			return t
		}
	}
	if idx := strings.Index(r, "("); idx >= 0 {
		r = r[:idx]
	}
	fields := strings.Fields(r)
	if len(fields) == 0 {
		return ""
	}
	return fields[0]
}

func normalizeIdent(raw string) string {
	out := strings.TrimSpace(raw)
	if out == "" {
		return out
	}
	if strings.Contains(out, ".") {
		parts := strings.Split(out, ".")
		out = parts[len(parts)-1]
	}
	out = strings.Trim(out, "`\"[]")
	return out
}

func mapSchemaColumns(cols []schemaColumn) []columnMapping {
	mappings := make([]columnMapping, 0, len(cols))
	for _, col := range cols {
		genName, fieldName := pickGeneratorField(col)
		mappings = append(mappings, columnMapping{
			ColumnName: col.Name,
			GenName:    genName,
			FieldName:  fieldName,
		})
	}
	return mappings
}

func pickGeneratorField(col schemaColumn) (string, string) {
	name := strings.ToLower(col.Name)
	typ := strings.ToLower(col.Type)

	containsAny := func(s string, keys ...string) bool {
		for _, k := range keys {
			if strings.Contains(s, k) {
				return true
			}
		}
		return false
	}

	// Name-driven hints first.
	switch {
	case name == "first_name":
		return "names", "first_name"
	case name == "last_name":
		return "names", "last_name"
	case name == "full_name" || name == "name":
		return "names", "full_name"
	case containsAny(name, "email"):
		return "emails", "email"
	case containsAny(name, "phone", "mobile", "tel"):
		return "phone", "phone"
	case containsAny(name, "street"):
		return "addresses", "street"
	case containsAny(name, "city"):
		return "addresses", "city"
	case name == "state":
		return "addresses", "state"
	case containsAny(name, "zip", "postal"):
		return "addresses", "zip"
	case containsAny(name, "country"):
		return "addresses", "country"
	case containsAny(name, "company", "employer"):
		return "companies", "company"
	case containsAny(name, "department", "dept"):
		return "companies", "department"
	case containsAny(name, "title", "role"):
		return "companies", "job_title"
	case containsAny(name, "website", "site", "domain", "hostname"):
		return "network", "hostname"
	case containsAny(name, "url", "uri"):
		return "network", "url"
	case containsAny(name, "user_agent", "useragent"):
		return "network", "user_agent"
	case containsAny(name, "mac"):
		return "network", "mac_address"
	case containsAny(name, "ipv6"):
		return "network", "ipv6"
	case containsAny(name, "ip", "ipv4"):
		return "network", "ipv4"
	case containsAny(name, "address"):
		return "addresses", "full_address"
	case containsAny(name, "currency_code"):
		return "numbers", "currency_code"
	case containsAny(name, "price", "amount", "cost", "total", "subtotal", "balance", "salary"):
		return "numbers", "currency_amount"
	case containsAny(name, "percent", "ratio"):
		return "numbers", "percentage"
	case containsAny(name, "active", "enabled", "is_", "has_", "deleted", "verified"):
		return "numbers", "boolean"
	case name == "id":
		if containsAny(typ, "uuid") {
			return "identifiers", "uuid"
		}
		if containsAny(typ, "bigint") {
			return "numbers", "bigint"
		}
		return "numbers", "integer"
	case strings.HasSuffix(name, "_id"):
		if containsAny(typ, "uuid") {
			return "identifiers", "uuid"
		}
		if containsAny(typ, "bigint") {
			return "numbers", "bigint"
		}
		return "numbers", "integer"
	case containsAny(name, "uuid"):
		return "identifiers", "uuid"
	case containsAny(name, "dob", "birth"):
		return "dates", "date"
	case containsAny(name, "created_at", "updated_at", "deleted_at", "timestamp"):
		return "dates", "datetime"
	case containsAny(name, "date"):
		return "dates", "date"
	case containsAny(name, "time"):
		return "dates", "time"
	}

	// Type-driven fallback.
	switch {
	case containsAny(typ, "uuid"):
		return "identifiers", "uuid"
	case containsAny(typ, "bool"):
		return "numbers", "boolean"
	case containsAny(typ, "bigint"):
		return "numbers", "bigint"
	case containsAny(typ, "int", "serial"):
		return "numbers", "integer"
	case containsAny(typ, "decimal", "numeric"):
		return "numbers", "decimal"
	case containsAny(typ, "float", "double", "real"):
		return "numbers", "float"
	case containsAny(typ, "timestamp", "datetime"):
		return "dates", "datetime"
	case typ == "date":
		return "dates", "date"
	case containsAny(typ, "time"):
		return "dates", "time"
	case containsAny(typ, "json", "text", "char", "varchar", "string"):
		if containsAny(name, "name") {
			return "names", "full_name"
		}
		return "lorem", "sentence"
	default:
		return "lorem", "word"
	}
}

func generateSchemaRecords(mappings []columnMapping, count int, rng *rand.Rand) ([]map[string]any, []string) {
	fields := make([]string, 0, len(mappings))
	for _, m := range mappings {
		fields = append(fields, m.ColumnName)
	}

	sourceCache := map[string][]map[string]any{}
	for _, m := range mappings {
		if _, ok := sourceCache[m.GenName]; ok {
			continue
		}
		gen := generator.Lookup(m.GenName)
		if gen == nil {
			continue
		}
		sourceCache[m.GenName] = gen.Generate(count, rng)
	}

	records := make([]map[string]any, count)
	for i := 0; i < count; i++ {
		rec := make(map[string]any, len(mappings))
		for _, m := range mappings {
			if rows, ok := sourceCache[m.GenName]; ok && i < len(rows) {
				rec[m.ColumnName] = rows[i][m.FieldName]
			}
		}
		records[i] = rec
	}
	return records, fields
}

func formatRecords(records []map[string]any, fieldNames []string, format, table string) string {
	switch strings.ToLower(format) {
	case "json":
		return output.FormatJSON(records, fieldNames)
	case "jsonl":
		return output.FormatJSONLines(records, fieldNames)
	case "csv":
		return output.FormatCSV(records, fieldNames)
	case "markdown", "md":
		return output.FormatMarkdown(records, fieldNames)
	case "sql":
		return output.FormatSQL(records, fieldNames, table)
	default:
		fmt.Fprintf(os.Stderr, "unknown format: %s\nvalid: json, jsonl, csv, markdown, sql\n", format)
		os.Exit(1)
		return ""
	}
}
