package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/LFroesch/seedbank/internal/generator"
	"github.com/LFroesch/seedbank/internal/output"
)

func main() {
	genName := flag.String("gen", "", "generator name (e.g. person, names, email)")
	fields := flag.String("fields", "", "comma-separated field names (default: all)")
	count := flag.Int("count", 10, "number of records")
	format := flag.String("fmt", "json", "output format: json, jsonl, csv, markdown, sql")
	table := flag.String("table", "", "SQL table name (default: derived from generator)")
	seed := flag.Int64("seed", 0, "random seed (0 = random)")
	listGens := flag.Bool("list", false, "list available generators and their fields")
	flag.Parse()

	// List mode: show generators and exit
	if *listGens {
		for _, g := range generator.Registry {
			fmt.Printf("%s — %s\n", g.Name(), g.Description())
			for _, f := range g.Fields() {
				fmt.Printf("  %-20s %s\n", f.Name, f.Desc)
			}
			fmt.Println()
		}
		os.Exit(0)
	}

	// Pipe mode: if --gen is set, skip TUI
	if *genName != "" {
		runPipe(*genName, *fields, *count, *format, *table, *seed)
		return
	}

	// TUI mode
	m := initialModel()
	p := tea.NewProgram(&m, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}

func runPipe(genName, fieldStr string, count int, format, table string, seed int64) {
	// Find generator by name (case-insensitive, partial match)
	var gen generator.Generator
	genLower := strings.ToLower(genName)
	for _, g := range generator.Registry {
		name := strings.ToLower(g.Name())
		// Exact match or contains
		if name == genLower || strings.Contains(name, genLower) {
			gen = g
			break
		}
	}
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

	// Format and output
	var out string
	switch strings.ToLower(format) {
	case "json":
		out = output.FormatJSON(records, fieldNames)
	case "jsonl":
		out = output.FormatJSONLines(records, fieldNames)
	case "csv":
		out = output.FormatCSV(records, fieldNames)
	case "markdown", "md":
		out = output.FormatMarkdown(records, fieldNames)
	case "sql":
		out = output.FormatSQL(records, fieldNames, table)
	default:
		fmt.Fprintf(os.Stderr, "unknown format: %s\nvalid: json, jsonl, csv, markdown, sql\n", format)
		os.Exit(1)
	}

	fmt.Print(out)
}
