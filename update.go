package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/LFroesch/seedbank/internal/output"
)

func (m *model) Init() tea.Cmd {
	return tea.SetWindowTitle("seedbank")
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// Clear expired status
	if m.statusMsg != "" && time.Now().After(m.statusExp) {
		m.statusMsg = ""
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
		// Global keys
		switch msg.String() {
		case "q":
			if m.mode == modeGenerators {
				return m, tea.Quit
			}
			if m.mode == modeExport && m.textInput.Focused() {
				// Don't quit while typing
			} else if m.mode == modeCount && m.textInput.Focused() {
				// Don't quit while typing count
			} else {
				// Go back instead of quit
				return m, m.goBack()
			}
		case "?":
			if m.mode != modeExport && m.mode != modeCount {
				m.mode = modeHelp
				m.helpScroll = 0
				return m, nil
			}
		}

		switch m.mode {
		case modeGenerators:
			return m, m.updateGenerators(msg)
		case modeFields:
			return m, m.updateFields(msg)
		case modeCount:
			return m, m.updateCount(msg)
		case modePreview:
			return m, m.updatePreview(msg)
		case modeFormat:
			return m, m.updateFormat(msg)
		case modeExport:
			return m, m.updateExport(msg)
		case modeHelp:
			return m, m.updateHelp(msg)
		case modeMixSelect:
			return m, m.updateMixSelect(msg)
		}
	}

	return m, nil
}

func (m *model) goBack() tea.Cmd {
	switch m.mode {
	case modeFields:
		m.mode = modeGenerators
	case modeCount:
		m.textInput.Blur()
		m.mode = modeFields
	case modePreview:
		m.mode = modeCount
	case modeFormat:
		m.mode = modePreview
	case modeExport:
		m.textInput.Blur()
		m.mode = modeFormat
	case modeHelp:
		m.mode = modeGenerators
	case modeMixSelect:
		m.mode = modeGenerators
	}
	return nil
}

func (m *model) updateGenerators(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "up", "k":
		if m.genCursor > 0 {
			m.genCursor--
		}
	case "down", "j":
		if m.genCursor < len(m.generators)-1 {
			m.genCursor++
		}
	case "enter":
		m.selectedGen = m.generators[m.genCursor]
		fields := m.selectedGen.Fields()
		m.fieldToggles = make([]bool, len(fields))
		for i := range m.fieldToggles {
			m.fieldToggles[i] = true // All on by default
		}
		m.fieldCursor = 0
		m.selectedFields = nil
		m.mode = modeFields
	case "m":
		// Enter mix mode
		m.mixToggles = make([]bool, len(m.generators))
		m.mixCursor = 0
		m.mode = modeMixSelect
	case "home", "g":
		m.genCursor = 0
	case "end", "G":
		m.genCursor = len(m.generators) - 1
	}
	// Keep scroll in bounds
	m.ensureGenScroll()
	return nil
}

func (m *model) updateFields(msg tea.KeyMsg) tea.Cmd {
	fields := m.selectedGen.Fields()
	switch msg.String() {
	case "up", "k":
		if m.fieldCursor > 0 {
			m.fieldCursor--
		}
	case "down", "j":
		if m.fieldCursor < len(fields)-1 {
			m.fieldCursor++
		}
	case " ", "space":
		if m.fieldCursor < len(m.fieldToggles) {
			m.fieldToggles[m.fieldCursor] = !m.fieldToggles[m.fieldCursor]
		}
	case "a":
		allOn := true
		for _, t := range m.fieldToggles {
			if !t {
				allOn = false
				break
			}
		}
		for i := range m.fieldToggles {
			m.fieldToggles[i] = !allOn
		}
	case "enter":
		m.mode = modeCount
		m.textInput.SetValue(strconv.Itoa(m.count))
		m.textInput.Placeholder = "number of records..."
		m.textInput.Focus()
	case "esc":
		m.mode = modeGenerators
	}
	return nil
}

func (m *model) updateCount(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "enter":
		val := strings.TrimSpace(m.textInput.Value())
		if n, err := strconv.Atoi(val); err == nil && n > 0 {
			m.count = n
		}
		m.textInput.Blur()
		m.generate()
		m.previewScroll = 0
		m.mode = modePreview
		return nil
	case "esc":
		m.textInput.Blur()
		m.mode = modeFields
		return nil
	}

	// Pass to text input
	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	return cmd
}

func (m *model) updatePreview(msg tea.KeyMsg) tea.Cmd {
	lines := strings.Count(m.preview, "\n") + 1
	viewH := m.getContentHeight()

	switch msg.String() {
	case "up", "k":
		if m.previewScroll > 0 {
			m.previewScroll--
		}
	case "down", "j":
		if m.previewScroll < lines-viewH {
			m.previewScroll++
		}
	case "v":
		m.prettyView = !m.prettyView
		m.previewScroll = 0
		if m.prettyView {
			m.statusMsg = "pretty view"
		} else {
			m.statusMsg = "raw view"
		}
		m.statusExp = time.Now().Add(2 * time.Second)
	case "r":
		// Re-roll with new seed
		m.reseed()
		m.generate()
		m.previewScroll = 0
		m.statusMsg = "re-rolled with new seed"
		m.statusExp = time.Now().Add(2 * time.Second)
	case "f":
		m.formatCursor = int(m.format)
		m.mode = modeFormat
	case "c":
		// Copy preview to clipboard
		content := m.preview
		if m.prettyView {
			content = m.formatPretty()
		}
		if err := clipboard.WriteAll(content); err != nil {
			m.statusMsg = "copy failed: " + err.Error()
		} else {
			m.statusMsg = "copied to clipboard"
		}
		m.statusExp = time.Now().Add(2 * time.Second)
	case "e":
		m.textInput.Placeholder = "output filename..."
		ext := formatExtension(m.format)
		m.textInput.SetValue("seedbank_output" + ext)
		m.textInput.Focus()
		m.mode = modeExport
	case "+":
		m.count = m.count * 2
		m.generate()
		m.previewScroll = 0
		m.statusMsg = fmt.Sprintf("count: %d", m.count)
		m.statusExp = time.Now().Add(2 * time.Second)
	case "-":
		if m.count > 1 {
			m.count = m.count / 2
			if m.count < 1 {
				m.count = 1
			}
			m.generate()
			m.previewScroll = 0
			m.statusMsg = fmt.Sprintf("count: %d", m.count)
			m.statusExp = time.Now().Add(2 * time.Second)
		}
	case "esc":
		m.mode = modeCount
		m.textInput.SetValue(strconv.Itoa(m.count))
		m.textInput.Placeholder = "number of records..."
		m.textInput.Focus()
	case "home", "g":
		m.previewScroll = 0
	case "end", "G":
		max := lines - viewH
		if max < 0 {
			max = 0
		}
		m.previewScroll = max
	}
	return nil
}

func (m *model) updateFormat(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "up", "k":
		if m.formatCursor > 0 {
			m.formatCursor--
		}
	case "down", "j":
		if m.formatCursor < len(output.FormatNames)-1 {
			m.formatCursor++
		}
	case "enter":
		m.format = output.Format(m.formatCursor)
		fields := m.getSelectedFieldNames()
		m.formatPreview(fields)
		m.previewScroll = 0
		m.mode = modePreview
	case "esc":
		m.mode = modePreview
	}
	return nil
}

func (m *model) updateExport(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "enter":
		filename := strings.TrimSpace(m.textInput.Value())
		if filename == "" {
			filename = "seedbank_output" + formatExtension(m.format)
		}
		m.textInput.Blur()

		// Build output dir path
		outPath := filename
		if !filepath.IsAbs(filename) {
			outPath = filepath.Join(m.config.OutputDir, filename)
		}

		err := os.WriteFile(outPath, []byte(m.preview), 0644)
		if err != nil {
			m.statusMsg = "export failed: " + err.Error()
		} else {
			m.statusMsg = fmt.Sprintf("exported %d records to %s", m.count, outPath)
		}
		m.statusExp = time.Now().Add(4 * time.Second)
		m.mode = modePreview
		return nil
	case "esc":
		m.textInput.Blur()
		m.mode = modePreview
		return nil
	}

	var cmd tea.Cmd
	m.textInput, cmd = m.textInput.Update(msg)
	return cmd
}

func (m *model) updateHelp(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "esc", "q", "?":
		m.mode = modeGenerators
	case "up", "k":
		if m.helpScroll > 0 {
			m.helpScroll--
		}
	case "down", "j":
		m.helpScroll++
	}
	return nil
}

func (m *model) updateMixSelect(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "up", "k":
		if m.mixCursor > 0 {
			m.mixCursor--
		}
	case "down", "j":
		if m.mixCursor < len(m.generators)-1 {
			m.mixCursor++
		}
	case " ", "space":
		if m.mixCursor < len(m.mixToggles) {
			m.mixToggles[m.mixCursor] = !m.mixToggles[m.mixCursor]
		}
	case "enter":
		// Need at least 2 generators selected for a mix
		count := 0
		for _, on := range m.mixToggles {
			if on {
				count++
			}
		}
		if count < 2 {
			m.statusMsg = "select at least 2 generators to mix"
			m.statusExp = time.Now().Add(2 * time.Second)
			return nil
		}
		m.buildMixGen()
		m.mode = modeFields
	case "esc":
		m.mode = modeGenerators
	}
	return nil
}

func (m *model) ensureGenScroll() {
	viewH := m.getContentHeight()
	if m.genCursor < m.genScrollOff {
		m.genScrollOff = m.genCursor
	}
	if m.genCursor >= m.genScrollOff+viewH {
		m.genScrollOff = m.genCursor - viewH + 1
	}
}

func formatExtension(f output.Format) string {
	switch f {
	case output.JSON:
		return ".json"
	case output.JSONLines:
		return ".jsonl"
	case output.CSV:
		return ".csv"
	case output.Markdown:
		return ".md"
	case output.SQL:
		return ".sql"
	}
	return ".txt"
}
