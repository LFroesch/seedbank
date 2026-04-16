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

	"github.com/LFroesch/seedbank/internal/generator"
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
		if m.handleStepJump(msg.String()) {
			return m, nil
		}
		if msg.String() == "tab" && !m.textInputFocused() && m.mode.usesPanelFocus() {
			m.togglePane()
			return m, nil
		}

		// Global keys
		switch msg.String() {
		case "q":
			if m.mode == modeGenerators {
				return m, tea.Quit
			}
			if m.textInputFocused() {
				// Don't quit while typing count
			} else {
				// Go back instead of quit
				return m, m.goBack()
			}
		case "?":
			if m.mode != modeExport && m.mode != modeCount {
				if m.mode == modeHelp {
					m.mode = m.prevMode
				} else {
					m.prevMode = m.mode
					m.mode = modeHelp
				}
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
		m.resetPaneFocus()
		m.mode = modeGenerators
	case modeCount:
		m.textInput.Blur()
		m.resetPaneFocus()
		m.mode = modeFields
	case modePreview:
		m.resetPaneFocus()
		m.mode = modeCount
	case modeFormat:
		m.resetPaneFocus()
		m.mode = modePreview
	case modeExport:
		m.textInput.Blur()
		m.resetPaneFocus()
		m.mode = modePreview
	case modeHelp:
		m.mode = m.prevMode
	case modeMixSelect:
		m.resetPaneFocus()
		m.mode = modeGenerators
	}
	return nil
}

func (m *model) updateGenerators(msg tea.KeyMsg) tea.Cmd {
	detailLines := 0
	if len(m.generators) > 0 {
		detailLines = len(m.generatorDetailContent())
	}
	detailH := m.secondaryPanelBodyHeight()

	switch msg.String() {
	case "up", "k":
		if m.activePane == paneRight {
			m.genDetailScroll--
		} else if m.genCursor > 0 {
			m.genCursor--
		}
	case "down", "j":
		if m.activePane == paneRight {
			m.genDetailScroll++
		} else if m.genCursor < len(m.generators)-1 {
			m.genCursor++
		}
	case "pgup":
		m.genDetailScroll -= detailH
	case "pgdown":
		m.genDetailScroll += detailH
	case "enter":
		m.selectedGen = m.generators[m.genCursor]
		fields := m.selectedGen.Fields()
		m.fieldToggles = make([]bool, len(fields))
		for i := range m.fieldToggles {
			m.fieldToggles[i] = true // All on by default
		}
		m.fieldCursor = 0
		m.fieldScrollOff = 0
		m.fieldSummaryScroll = 0
		m.selectedFields = nil
		m.resetPaneFocus()
		m.mode = modeFields
	case "m":
		// Enter mix mode
		m.mixToggles = make([]bool, len(m.generators))
		m.mixCursor = 0
		m.mixScrollOff = 0
		m.mixSummaryScroll = 0
		m.resetPaneFocus()
		m.mode = modeMixSelect
	case "home", "g":
		if m.activePane == paneRight {
			m.genDetailScroll = 0
		} else {
			m.genCursor = 0
		}
	case "end", "G":
		if m.activePane == paneRight {
			m.genDetailScroll = clampScroll(m.genDetailScroll, detailLines, detailH, true)
		} else {
			m.genCursor = len(m.generators) - 1
		}
	}
	if m.genCursor < 0 {
		m.genCursor = 0
	}
	if m.genCursor >= len(m.generators) {
		m.genCursor = len(m.generators) - 1
	}
	// Keep scroll in bounds
	m.ensureGenScroll()
	m.genDetailScroll = clampScroll(m.genDetailScroll, detailLines, detailH, false)
	return nil
}

func (m *model) updateFields(msg tea.KeyMsg) tea.Cmd {
	fields := m.selectedGen.Fields()
	summaryLines := len(m.fieldSummaryContent())
	summaryH := m.secondaryPanelBodyHeight()
	switch msg.String() {
	case "up", "k":
		if m.activePane == paneRight {
			m.fieldSummaryScroll--
		} else if m.fieldCursor > 0 {
			m.fieldCursor--
		}
	case "down", "j":
		if m.activePane == paneRight {
			m.fieldSummaryScroll++
		} else if m.fieldCursor < len(fields)-1 {
			m.fieldCursor++
		}
	case "pgup":
		if m.activePane == paneRight {
			m.fieldSummaryScroll -= summaryH
		} else {
			m.fieldCursor -= m.primaryListBodyHeight()
		}
	case "pgdown":
		if m.activePane == paneRight {
			m.fieldSummaryScroll += summaryH
		} else {
			m.fieldCursor += m.primaryListBodyHeight()
		}
	case " ", "space":
		if m.activePane == paneLeft && m.fieldCursor < len(m.fieldToggles) {
			m.fieldToggles[m.fieldCursor] = !m.fieldToggles[m.fieldCursor]
		}
	case "a":
		if m.activePane == paneLeft {
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
		}
	case "enter":
		m.resetPaneFocus()
		m.mode = modeCount
		m.textInput.SetValue(strconv.Itoa(m.count))
		m.textInput.Placeholder = "number of records..."
		m.textInput.Focus()
	case "esc":
		m.resetPaneFocus()
		m.mode = modeGenerators
	}
	if m.fieldCursor < 0 {
		m.fieldCursor = 0
	}
	if m.fieldCursor >= len(fields) {
		m.fieldCursor = len(fields) - 1
	}
	m.ensureFieldScroll()
	m.fieldSummaryScroll = clampScroll(m.fieldSummaryScroll, summaryLines, summaryH, false)
	m.ensureFieldSummaryScroll()
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
		m.previewInspectorScroll = 0
		m.resetPaneFocus()
		m.mode = modePreview
		return nil
	case "esc":
		m.textInput.Blur()
		m.resetPaneFocus()
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
	viewH := m.previewBodyHeight()
	inspectorLines := m.previewInspectorLines()
	inspectorH := m.secondaryPanelBodyHeight()

	switch msg.String() {
	case "up", "k":
		if m.activePane == paneRight {
			m.scrollPreviewInspector(-1, inspectorLines, inspectorH)
		} else if m.previewScroll > 0 {
			m.previewScroll--
		}
	case "down", "j":
		if m.activePane == paneRight {
			m.scrollPreviewInspector(1, inspectorLines, inspectorH)
		} else if m.previewScroll < lines-viewH {
			m.previewScroll++
		}
	case "pgup":
		if m.activePane == paneRight {
			m.scrollPreviewInspector(-inspectorH, inspectorLines, inspectorH)
		} else {
			m.previewScroll -= viewH
		}
	case "pgdown":
		if m.activePane == paneRight {
			m.scrollPreviewInspector(inspectorH, inspectorLines, inspectorH)
		} else {
			m.previewScroll += viewH
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
		m.previewInspectorScroll = 0
		m.statusMsg = "re-rolled with new seed"
		m.statusExp = time.Now().Add(2 * time.Second)
	case "f":
		m.formatCursor = int(m.format)
		m.resetPaneFocus()
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
		m.previewInspectorScroll = 0
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
			m.previewInspectorScroll = 0
			m.statusMsg = fmt.Sprintf("count: %d", m.count)
			m.statusExp = time.Now().Add(2 * time.Second)
		}
	case "esc":
		m.resetPaneFocus()
		m.mode = modeCount
		m.textInput.SetValue(strconv.Itoa(m.count))
		m.textInput.Placeholder = "number of records..."
		m.textInput.Focus()
	case "home", "g":
		if m.activePane == paneRight {
			m.previewInspectorScroll = 0
		} else {
			m.previewScroll = 0
		}
	case "end", "G":
		if m.activePane == paneRight {
			m.previewInspectorScroll = clampScroll(m.previewInspectorScroll, inspectorLines, inspectorH, true)
		} else {
			max := lines - viewH
			if max < 0 {
				max = 0
			}
			m.previewScroll = max
		}
	}
	m.previewScroll = clampScroll(m.previewScroll, lines, viewH, false)
	m.previewInspectorScroll = clampScroll(m.previewInspectorScroll, inspectorLines, inspectorH, false)
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
		m.previewInspectorScroll = 0
		m.resetPaneFocus()
		m.mode = modePreview
	case "esc":
		m.resetPaneFocus()
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

		outPath := m.resolveOutputPath(filename)
		if err := os.MkdirAll(filepath.Dir(outPath), 0755); err != nil {
			m.statusMsg = "export failed: " + err.Error()
			m.statusExp = time.Now().Add(4 * time.Second)
			m.mode = modePreview
			return nil
		}

		err := os.WriteFile(outPath, []byte(m.preview), 0644)
		if err != nil {
			m.statusMsg = "export failed: " + err.Error()
		} else {
			m.statusMsg = fmt.Sprintf("exported %d records to %s", m.count, outPath)
		}
		m.statusExp = time.Now().Add(4 * time.Second)
		m.resetPaneFocus()
		m.mode = modePreview
		return nil
	case "esc":
		m.textInput.Blur()
		m.resetPaneFocus()
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
		m.mode = m.prevMode
		m.helpScroll = 0
	case "up", "k":
		if m.helpScroll > 0 {
			m.helpScroll--
		}
	case "down", "j":
		m.helpScroll++
	case "pgup":
		m.helpScroll -= m.secondaryPanelBodyHeight()
	case "pgdown":
		m.helpScroll += m.secondaryPanelBodyHeight()
	case "home", "g":
		m.helpScroll = 0
	case "end", "G":
		m.helpScroll = 1 << 30
	}
	m.helpScroll = clampScroll(m.helpScroll, len(m.helpLines()), m.contentHeight()-2, false)
	return nil
}

func (m *model) updateMixSelect(msg tea.KeyMsg) tea.Cmd {
	summaryLines := len(m.mixSummaryContent())
	summaryH := m.secondaryPanelBodyHeight()
	switch msg.String() {
	case "up", "k":
		if m.activePane == paneRight {
			m.mixSummaryScroll--
		} else if m.mixCursor > 0 {
			m.mixCursor--
		}
	case "down", "j":
		if m.activePane == paneRight {
			m.mixSummaryScroll++
		} else if m.mixCursor < len(m.generators)-1 {
			m.mixCursor++
		}
	case "pgup":
		if m.activePane == paneRight {
			m.mixSummaryScroll -= summaryH
		} else {
			m.mixCursor -= m.primaryListBodyHeight()
		}
	case "pgdown":
		if m.activePane == paneRight {
			m.mixSummaryScroll += summaryH
		} else {
			m.mixCursor += m.primaryListBodyHeight()
		}
	case " ", "space":
		if m.activePane == paneLeft && m.mixCursor < len(m.mixToggles) {
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
		m.fieldScrollOff = 0
		m.fieldSummaryScroll = 0
		m.resetPaneFocus()
		m.mode = modeFields
	case "esc":
		m.resetPaneFocus()
		m.mode = modeGenerators
	}
	if m.mixCursor < 0 {
		m.mixCursor = 0
	}
	if m.mixCursor >= len(m.generators) {
		m.mixCursor = len(m.generators) - 1
	}
	m.ensureMixScroll()
	m.mixSummaryScroll = clampScroll(m.mixSummaryScroll, summaryLines, summaryH, false)
	m.ensureMixSummaryScroll()
	return nil
}

func (m *model) ensureGenScroll() {
	viewH := m.generatorListBodyHeight()
	if m.genCursor < m.genScrollOff {
		m.genScrollOff = m.genCursor
	}
	if m.genCursor >= m.genScrollOff+viewH {
		m.genScrollOff = m.genCursor - viewH + 1
	}
}

func (m *model) ensureFieldScroll() {
	viewH := m.primaryListBodyHeight()
	if m.fieldCursor < m.fieldScrollOff {
		m.fieldScrollOff = m.fieldCursor
	}
	if m.fieldCursor >= m.fieldScrollOff+viewH {
		m.fieldScrollOff = m.fieldCursor - viewH + 1
	}
}

func (m *model) ensureMixScroll() {
	viewH := m.primaryListBodyHeight()
	if m.mixCursor < m.mixScrollOff {
		m.mixScrollOff = m.mixCursor
	}
	if m.mixCursor >= m.mixScrollOff+viewH {
		m.mixScrollOff = m.mixCursor - viewH + 1
	}
}

func (m mode) usesPanelFocus() bool {
	switch m {
	case modeGenerators, modeFields, modePreview, modeMixSelect:
		return true
	default:
		return false
	}
}

func (m *model) handleStepJump(key string) bool {
	if m.textInputFocused() {
		return false
	}

	target := mode(-1)
	switch key {
	case "1":
		target = modeGenerators
	case "2":
		target = modeFields
	case "3":
		target = modeCount
	case "4":
		target = modePreview
	case "5":
		target = modeFormat
	case "6":
		target = modeExport
	default:
		return false
	}

	m.resetPaneFocus()
	if !m.prepareStep(target) {
		return false
	}
	return true
}

func (m *model) resolveOutputPath(filename string) string {
	if filepath.IsAbs(filename) {
		return filename
	}
	base := "."
	if m.config != nil && strings.TrimSpace(m.config.OutputDir) != "" {
		base = m.config.OutputDir
	}
	return filepath.Join(base, filename)
}

func (m *model) previewBodyHeight() int {
	h := m.contentHeight() - 4
	if h < 3 {
		h = 3
	}
	return h
}

func (m *model) generatorListBodyHeight() int {
	h := m.primaryListBodyHeight() - 2
	if h < 3 {
		h = 3
	}
	return h
}

func (m *model) primaryListBodyHeight() int {
	h := m.contentHeight() - 4
	if h < 4 {
		h = 4
	}
	return h
}

func (m *model) secondaryPanelBodyHeight() int {
	h := m.contentHeight() - 3
	if h < 3 {
		h = 3
	}
	return h
}

func (m *model) previewInspectorLines() int {
	return len(m.previewInspectorContent())
}

func (m *model) previewInspectorContent() []string {
	selectedFields := m.getSelectedFieldNames()
	lines := []string{
		panelHeaderStyle.Render("Preview"),
		"  " + normalStyle.Render(m.selectedGen.Name()),
		"  " + statusStyle.Render(fmt.Sprintf("%d rows", m.count)),
		"  " + accentStyle.Render(output.FormatNames[m.format]),
	}
	viewLabel := "raw view"
	if m.prettyView {
		viewLabel = "pretty view"
	}
	lines = append(lines, "  "+mutedStyle.Render(viewLabel), "")
	lines = append(lines, panelHeaderStyle.Render("Fields"))
	for _, field := range selectedFields {
		lines = append(lines, "  "+accentStyle.Render("•")+" "+field)
	}
	lines = append(lines, "", panelHeaderStyle.Render("Actions"))
	lines = append(lines,
		"  "+dimStyle.Render("tab switch pane"),
		"  "+dimStyle.Render("j/k or pgup/pgdn scroll"),
		"  "+dimStyle.Render("v toggle view"),
		"  "+dimStyle.Render("r re-roll seed"),
		"  "+dimStyle.Render("+/- adjust count"),
		"  "+dimStyle.Render("f change format"),
		"  "+dimStyle.Render("c copy"),
		"  "+dimStyle.Render("e export"),
	)
	return lines
}

func (m *model) scrollPreviewInspector(delta, total, visible int) {
	m.previewInspectorScroll += delta
	m.previewInspectorScroll = clampScroll(m.previewInspectorScroll, total, visible, false)
}

func clampScroll(current, total, visible int, toEnd bool) int {
	max := total - visible
	if max < 0 {
		max = 0
	}
	if toEnd {
		return max
	}
	if current < 0 {
		return 0
	}
	if current > max {
		return max
	}
	return current
}

func (m *model) ensureFieldSummaryScroll() {
	lines := len(m.fieldSummaryContent())
	m.fieldSummaryScroll = clampScroll(m.fieldSummaryScroll, lines, m.secondaryPanelBodyHeight(), false)
}

func (m *model) ensureMixSummaryScroll() {
	lines := len(m.mixSummaryContent())
	m.mixSummaryScroll = clampScroll(m.mixSummaryScroll, lines, m.secondaryPanelBodyHeight(), false)
}

func (m *model) prepareStep(target mode) bool {
	switch target {
	case modeGenerators:
		m.textInput.Blur()
		m.mode = modeGenerators
		return true
	case modeFields:
		m.ensureSelectedGenerator()
		m.textInput.Blur()
		m.mode = modeFields
		return true
	case modeCount:
		m.ensureSelectedGenerator()
		m.textInput.SetValue(strconv.Itoa(m.count))
		m.textInput.Placeholder = "number of records..."
		m.textInput.Focus()
		m.mode = modeCount
		return true
	case modePreview:
		m.ensurePreviewReady()
		m.textInput.Blur()
		m.mode = modePreview
		return true
	case modeFormat:
		m.ensurePreviewReady()
		m.textInput.Blur()
		m.formatCursor = int(m.format)
		m.mode = modeFormat
		return true
	case modeExport:
		m.ensurePreviewReady()
		m.textInput.Placeholder = "output filename..."
		m.textInput.SetValue("seedbank_output" + formatExtension(m.format))
		m.textInput.Focus()
		m.mode = modeExport
		return true
	default:
		return false
	}
}

func (m *model) ensureSelectedGenerator() {
	if m.selectedGen != nil {
		return
	}
	if len(m.generators) == 0 {
		return
	}
	m.selectedGen = m.generators[m.genCursor]
	fields := m.selectedGen.Fields()
	m.fieldToggles = make([]bool, len(fields))
	for i := range m.fieldToggles {
		m.fieldToggles[i] = true
	}
	m.fieldCursor = 0
	m.fieldScrollOff = 0
	m.fieldSummaryScroll = 0
}

func (m *model) ensurePreviewReady() {
	m.ensureSelectedGenerator()
	if m.selectedGen == nil {
		return
	}
	if len(m.records) == 0 {
		m.generate()
	}
	m.previewScroll = 0
	m.previewInspectorScroll = 0
}

func (m *model) generatorDetailContent() []string {
	if len(m.generators) == 0 {
		return []string{dimStyle.Render("no generators")}
	}
	gen := m.generators[m.genCursor]
	kindLabel := "Field source"
	usageLine := "Use Mix when you want to compose this with other field generators."
	if gen.Kind() == generator.KindRecord {
		kindLabel = "Coherent record builder"
		usageLine = "Use this when you want one ready-to-export record with matching fields."
	}
	lines := []string{
		normalStyle.Render(gen.Description()),
		"",
		panelHeaderStyle.Render("Type"),
		"  " + accentStyle.Render(kindLabel),
		"  " + dimStyle.Render(usageLine),
		"",
		panelHeaderStyle.Render("Fields"),
	}
	for _, field := range gen.Fields() {
		lines = append(lines, "  "+selectedStyle.Render(field.Name)+dimStyle.Render(" · ")+dimStyle.Render(field.Desc))
	}
	return lines
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
