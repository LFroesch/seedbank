package main

import (
	"math/rand"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textinput"

	"github.com/LFroesch/seedbank/internal/config"
	"github.com/LFroesch/seedbank/internal/generator"
	"github.com/LFroesch/seedbank/internal/output"
)

// Terminal dimension constants
const (
	minTerminalWidth  = 60
	minTerminalHeight = 20
	uiOverhead        = 7 // header + status + borders + padding
)

type mode int

const (
	modeGenerators mode = iota // Browse generator list
	modeFields                 // Select fields for chosen generator
	modeCount                  // Enter record count
	modePreview                // Preview generated data
	modeFormat                 // Choose output format
	modeExport                 // Export confirmation / file path
	modeHelp                   // Help screen
	modeMixSelect              // Multi-select generators for mix mode
)

type model struct {
	mode mode

	// Terminal
	width  int
	height int

	// Generator selection
	generators    []generator.Generator
	genCursor     int
	genScrollOff  int
	selectedGen   generator.Generator
	selectedFields []string // Which fields are toggled on
	fieldCursor   int
	fieldToggles  []bool // Parallel to selectedGen.Fields()

	// Generation
	count     int
	records   []map[string]any
	rng       *rand.Rand
	seed      int64

	// Output
	format        output.Format
	formatCursor  int
	preview       string
	previewLines  []string
	previewScroll int
	prettyView    bool // true = human-readable cards, false = raw copyable output

	// Export
	textInput  textinput.Model
	statusMsg  string
	statusExp  time.Time

	// Config
	config *config.Config

	// Help
	helpScroll int

	// Mix mode
	mixToggles []bool                // parallel to generators list
	mixCursor  int
	mixGens    []generator.Generator // selected generators for mix
}

func initialModel() model {
	cfg := config.Load()

	ti := textinput.New()
	ti.Placeholder = "output filename..."
	ti.CharLimit = 256
	ti.Width = 40

	seed := cfg.Seed
	if seed == 0 {
		seed = time.Now().UnixNano()
	}

	m := model{
		mode:         modeGenerators,
		generators:   generator.Registry,
		genCursor:    0,
		count:        cfg.DefaultCount,
		rng:          rand.New(rand.NewSource(seed)),
		seed:         seed,
		format:       output.JSON,
		formatCursor: 0,
		textInput:    ti,
		config:       cfg,
	}

	return m
}

func (m *model) getSafeWidth() int {
	if m.width < minTerminalWidth {
		return minTerminalWidth
	}
	return m.width
}

func (m *model) getSafeHeight() int {
	if m.height < minTerminalHeight {
		return minTerminalHeight
	}
	return m.height
}

func (m *model) getContentHeight() int {
	h := m.getSafeHeight() - uiOverhead
	if h < 3 {
		h = 3
	}
	return h
}

// buildMixGen creates a MixGenerator from selected generators and sets up field toggles.
func (m *model) buildMixGen() {
	var selected []generator.Generator
	for i, on := range m.mixToggles {
		if on {
			selected = append(selected, m.generators[i])
		}
	}
	m.mixGens = selected
	fields := generator.BuildMixFields(selected)
	m.selectedGen = &generator.MixGenerator{
		Gens:    selected,
		Fields_: fields,
	}
	m.fieldToggles = make([]bool, len(fields))
	for i := range m.fieldToggles {
		m.fieldToggles[i] = true
	}
	m.fieldCursor = 0
}

// tableName returns a SQL-friendly table name from the current generator.
func (m *model) tableName() string {
	if m.selectedGen == nil {
		return m.config.TableName
	}
	name := m.selectedGen.Name()
	// Strip parenthetical like "(Linked)"
	if idx := strings.Index(name, "("); idx > 0 {
		name = strings.TrimSpace(name[:idx])
	}
	name = strings.ToLower(name)
	name = strings.ReplaceAll(name, " ", "_")
	return name
}

// generate runs the selected generator and builds both preview strings.
func (m *model) generate() {
	if m.selectedGen == nil {
		return
	}

	// Collect selected field names
	fields := m.getSelectedFieldNames()

	// Generate records
	m.records = m.selectedGen.Generate(m.count, m.rng)

	// Format preview
	m.formatPreview(fields)
}

func (m *model) getSelectedFieldNames() []string {
	if m.selectedGen == nil {
		return nil
	}
	allFields := m.selectedGen.Fields()
	var names []string
	for i, f := range allFields {
		if i < len(m.fieldToggles) && m.fieldToggles[i] {
			names = append(names, f.Name)
		}
	}
	// If none selected, use all
	if len(names) == 0 {
		for _, f := range allFields {
			names = append(names, f.Name)
		}
	}
	return names
}

func (m *model) formatPreview(fields []string) {
	switch m.format {
	case output.JSON:
		m.preview = output.FormatJSON(m.records, fields)
	case output.JSONLines:
		m.preview = output.FormatJSONLines(m.records, fields)
	case output.CSV:
		m.preview = output.FormatCSV(m.records, fields)
	case output.Markdown:
		m.preview = output.FormatMarkdown(m.records, fields)
	case output.SQL:
		m.preview = output.FormatSQL(m.records, fields, m.tableName())
	}
}

// formatPretty builds a human-readable card view of the records.
func (m *model) formatPretty() string {
	fields := m.getSelectedFieldNames()
	return output.FormatPretty(m.records, fields)
}

// reseed creates a new RNG with a fresh seed for re-rolling data.
func (m *model) reseed() {
	m.seed = time.Now().UnixNano()
	m.rng = rand.New(rand.NewSource(m.seed))
}
