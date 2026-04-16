package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/LFroesch/seedbank/internal/generator"
)

func TestWriteOutputCreatesParentDirectories(t *testing.T) {
	tmp := t.TempDir()
	target := filepath.Join(tmp, "nested", "seed", "out.json")

	if err := writeOutput(target, "hello"); err != nil {
		t.Fatalf("writeOutput() error = %v", err)
	}

	data, err := os.ReadFile(target)
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	if string(data) != "hello" {
		t.Fatalf("file contents = %q, want %q", string(data), "hello")
	}
}

func TestLookupGeneratorAliases(t *testing.T) {
	tests := map[string]string{
		"person":     "Person",
		"product":    "Products",
		"email":      "Emails",
		"phone":      "Phone Numbers",
		"image":      "Photos",
		"text":       "Lorem Ipsum",
		"uuid":       "Identifiers",
		"identifier": "Identifiers",
		"net":        "Network",
		"internet":   "Network",
	}

	for input, want := range tests {
		got := generator.Lookup(input)
		if got == nil {
			t.Fatalf("Lookup(%q) returned nil", input)
		}
		if got.Name() != want {
			t.Fatalf("Lookup(%q) = %q, want %q", input, got.Name(), want)
		}
	}
}
