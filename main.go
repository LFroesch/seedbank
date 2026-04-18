package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/LFroesch/seedbank/internal/config"
	"github.com/LFroesch/seedbank/internal/generator"
)

func main() {
	cfg := config.Load()

	genName := flag.String("gen", "", "generator name (e.g. person, names, email)")
	schemaPath := flag.String("schema", "", "path to CREATE TABLE SQL file for preview heuristic schema-driven generation")
	fields := flag.String("fields", "", "comma-separated field names (default: all)")
	count := flag.Int("count", cfg.DefaultCount, "number of records")
	format := flag.String("fmt", cfg.DefaultFormat, "output format: json, jsonl, csv, markdown, sql")
	table := flag.String("table", "", "SQL table name (default: derived from generator)")
	seed := flag.Int64("seed", 0, "random seed (0 = random)")
	outPath := flag.String("out", "-", "write output to file path; use - for stdout")
	listGens := flag.Bool("list", false, "list available generators and their fields")
	flag.Parse()

	// List mode: show generators and exit
	if *listGens {
		printGeneratorList()
		os.Exit(0)
	}

	// Pipe mode: if --gen is set, skip TUI
	if *schemaPath != "" {
		runSchema(*schemaPath, *count, *format, *table, *outPath, *seed)
		return
	}

	// Pipe mode: if --gen is set, skip TUI
	if *genName != "" {
		runPipe(*genName, *fields, *count, *format, *table, *outPath, *seed)
		return
	}

	// TUI mode
	m := initialModel()
	p := tea.NewProgram(&m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

func runPipe(genName, fieldStr string, count int, format, table, outPath string, seed int64) {
	gen := generator.Lookup(genName)
	if gen == nil {
		fmt.Fprintf(os.Stderr, "unknown generator: %s\nuse --list to see available generators\n", genName)
		os.Exit(1)
	}

	// Seed
	if seed == 0 {
		seed = time.Now().UnixNano()
	}
	rng := rand.New(rand.NewSource(seed))

	// Generate
	records := gen.Generate(count, rng)

	// Filter fields
	var fieldNames []string
	if fieldStr != "" {
		fieldNames = strings.Split(fieldStr, ",")
		for i := range fieldNames {
			fieldNames[i] = strings.TrimSpace(fieldNames[i])
		}
	} else {
		for _, f := range gen.Fields() {
			fieldNames = append(fieldNames, f.Name)
		}
	}

	// Table name
	if table == "" {
		name := gen.Name()
		if idx := strings.Index(name, "("); idx > 0 {
			name = strings.TrimSpace(name[:idx])
		}
		table = strings.ToLower(strings.ReplaceAll(name, " ", "_"))
	}

	out := formatRecords(records, fieldNames, format, table)

	if err := writeOutput(outPath, out); err != nil {
		fmt.Fprintf(os.Stderr, "write failed: %v\n", err)
		os.Exit(1)
	}
}

func writeOutput(path, content string) error {
	if path == "" || path == "-" {
		fmt.Print(content)
		return nil
	}
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(content), 0644)
}

func printGeneratorList() {
	grouped := generator.Grouped()
	order := []generator.Kind{generator.KindRecord, generator.KindField}
	for _, kind := range order {
		gens := grouped[kind]
		if len(gens) == 0 {
			continue
		}
		fmt.Printf("%s\n", generator.KindLabel(kind))
		for _, g := range gens {
			fmt.Printf("  %s — %s\n", g.Name(), g.Description())
			for _, f := range g.Fields() {
				fmt.Printf("    %-18s %s\n", f.Name, f.Desc)
			}
			fmt.Println()
		}
	}
}
