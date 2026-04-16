package main

import (
	"math/rand"
	"testing"
)

func TestParseCreateTableSQL(t *testing.T) {
	sql := `
CREATE TABLE IF NOT EXISTS public.users (
  id UUID PRIMARY KEY,
  first_name VARCHAR(100) NOT NULL,
  email TEXT UNIQUE,
  created_at TIMESTAMP NOT NULL,
  is_active BOOLEAN DEFAULT true,
  amount DECIMAL(10,2),
  ip_address INET,
  CONSTRAINT fk_users_org FOREIGN KEY (id) REFERENCES orgs(id)
);`

	got, err := parseCreateTableSQL(sql)
	if err != nil {
		t.Fatalf("parseCreateTableSQL() error = %v", err)
	}
	if got.Table != "users" {
		t.Fatalf("table = %q, want %q", got.Table, "users")
	}
	if len(got.Columns) != 7 {
		t.Fatalf("column count = %d, want 7", len(got.Columns))
	}
}

func TestMapSchemaColumnsHeuristics(t *testing.T) {
	cols := []schemaColumn{
		{Name: "id", Type: "uuid"},
		{Name: "customer_id", Type: "bigint"},
		{Name: "email", Type: "text"},
		{Name: "first_name", Type: "varchar"},
		{Name: "created_at", Type: "timestamp"},
		{Name: "price", Type: "decimal"},
		{Name: "is_active", Type: "boolean"},
		{Name: "ip_address", Type: "varchar"},
		{Name: "user_agent", Type: "text"},
	}
	mapped := mapSchemaColumns(cols)

	want := map[string]struct {
		gen   string
		field string
	}{
		"id":          {gen: "identifiers", field: "uuid"},
		"customer_id": {gen: "numbers", field: "bigint"},
		"email":       {gen: "emails", field: "email"},
		"first_name":  {gen: "names", field: "first_name"},
		"created_at":  {gen: "dates", field: "datetime"},
		"price":       {gen: "numbers", field: "currency_amount"},
		"is_active":   {gen: "numbers", field: "boolean"},
		"ip_address":  {gen: "network", field: "ipv4"},
		"user_agent":  {gen: "network", field: "user_agent"},
	}

	for _, m := range mapped {
		w, ok := want[m.ColumnName]
		if !ok {
			t.Fatalf("unexpected mapped column %q", m.ColumnName)
		}
		if m.GenName != w.gen || m.FieldName != w.field {
			t.Fatalf("column %q mapped to %s.%s, want %s.%s", m.ColumnName, m.GenName, m.FieldName, w.gen, w.field)
		}
	}
}

func TestGenerateSchemaRecordsPreservesColumnNames(t *testing.T) {
	mappings := []columnMapping{
		{ColumnName: "user_id", GenName: "identifiers", FieldName: "uuid"},
		{ColumnName: "email_address", GenName: "emails", FieldName: "email"},
		{ColumnName: "signup_ip", GenName: "network", FieldName: "ipv4"},
		{ColumnName: "created_at", GenName: "dates", FieldName: "datetime"},
	}

	records, fields := generateSchemaRecords(mappings, 3, rand.New(rand.NewSource(42)))

	if len(fields) != 4 {
		t.Fatalf("fields len = %d, want 4", len(fields))
	}
	if len(records) != 3 {
		t.Fatalf("records len = %d, want 3", len(records))
	}

	for i, rec := range records {
		for _, f := range []string{"user_id", "email_address", "signup_ip", "created_at"} {
			if rec[f] == nil {
				t.Fatalf("record %d missing field %q", i, f)
			}
		}
	}
}
