package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/LFroesch/seedbank/internal/output"
)

// Style palette — matching scout's purple accent theme
var (
	purple    = lipgloss.Color("99")
	white     = lipgloss.Color("255")
	dim       = lipgloss.Color("243")
	green     = lipgloss.Color("76")
	orange    = lipgloss.Color("214")
	bg        = lipgloss.Color("235")
	darkBg    = lipgloss.Color("233")

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(purple).
			Background(darkBg).
			Padding(0, 1)

	headerStyle = lipgloss.NewStyle().
			Foreground(white).
			Background(darkBg)

	statusStyle = lipgloss.NewStyle().
			Foreground(white).
			Background(bg).
			Padding(0, 1)

	selectedStyle = lipgloss.NewStyle().
			Foreground(purple).
			Bold(true)

	normalStyle = lipgloss.NewStyle().
			Foreground(white)

	dimStyle = lipgloss.NewStyle().
			Foreground(dim)

	accentStyle = lipgloss.NewStyle().
			Foreground(orange)

	checkStyle = lipgloss.NewStyle().
			Foreground(green).
			Bold(true)

	borderStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(purple)

	previewBorderStyle = lipgloss.NewStyle().
				Border(lipgloss.NormalBorder()).
				BorderForeground(dim)
)

func (m *model) View() string {
	if m.width == 0 || m.height == 0 {
		return "loading..."
	}

	if m.width < minTerminalWidth || m.height < minTerminalHeight {
		return lipgloss.NewStyle().
			Foreground(orange).Bold(true).Padding(1).
			Render(fmt.Sprintf("terminal too small: %dx%d\nminimum: %dx%d",
				m.width, m.height, minTerminalWidth, minTerminalHeight))
	}

	header := m.renderHeader()
	var body string

	switch m.mode {
	case modeGenerators:
		body = m.renderGenerators()
	case modeFields:
		body = m.renderFields()
	case modeCount:
		body = m.renderCount()
	case modePreview:
		body = m.renderPreviewView()
	case modeFormat:
		body = m.renderFormatPicker()
	case modeExport:
		body = m.renderExportView()
	case modeHelp:
		body = m.renderHelp()
	case modeMixSelect:
		body = m.renderMixSelect()
	}

	status := m.renderStatus()

	return header + "\n" + body + "\n" + status
}

func (m *model) renderHeader() string {
	w := m.getSafeWidth()

	title := titleStyle.Render(" seedbank ")
	modeLabel := dimStyle.Render(modeString(m.mode))

	gap := w - lipgloss.Width(title) - lipgloss.Width(modeLabel) - 2
	if gap < 1 {
		gap = 1
	}

	return headerStyle.Width(w).Render(
		title + strings.Repeat(" ", gap) + modeLabel,
	)
}

func (m *model) renderStatus() string {
	w := m.getSafeWidth()

	left := ""
	if m.statusMsg != "" {
		left = accentStyle.Render(m.statusMsg)
	}

	right := dimStyle.Render(fmt.Sprintf("seed:%d  count:%d  fmt:%s",
		m.seed%10000, m.count, output.FormatNames[m.format]))

	gap := w - lipgloss.Width(left) - lipgloss.Width(right) - 2
	if gap < 1 {
		gap = 1
	}

	return statusStyle.Width(w).Render(
		left + strings.Repeat(" ", gap) + right,
	)
}

func (m *model) renderGenerators() string {
	w := m.getSafeWidth()
	h := m.getContentHeight()

	var sb strings.Builder
	sb.WriteString(dimStyle.Render("  Select a data generator:") + "\n\n")

	for i, gen := range m.generators {
		if i < m.genScrollOff || i >= m.genScrollOff+h-2 {
			continue
		}

		cursor := "  "
		style := normalStyle
		if i == m.genCursor {
			cursor = "> "
			style = selectedStyle
		}

		name := style.Render(gen.Name())
		desc := dimStyle.Render(" - " + gen.Description())

		line := cursor + name + desc
		// Truncate if too wide
		if lipgloss.Width(line) > w-2 {
			line = line[:w-4] + ".."
		}
		sb.WriteString(line + "\n")
	}

	hint := dimStyle.Render("\n  j/k navigate  enter select  m mix mode  ? help  q quit")
	sb.WriteString(hint)

	return sb.String()
}

func (m *model) renderFields() string {
	if m.selectedGen == nil {
		return ""
	}

	fields := m.selectedGen.Fields()
	var sb strings.Builder

	sb.WriteString(selectedStyle.Render("  "+m.selectedGen.Name()) +
		dimStyle.Render(" - select fields:") + "\n\n")

	for i, f := range fields {
		cursor := "  "
		if i == m.fieldCursor {
			cursor = "> "
		}

		check := "[ ]"
		if i < len(m.fieldToggles) && m.fieldToggles[i] {
			check = checkStyle.Render("[x]")
		}

		name := normalStyle.Render(f.Name)
		desc := dimStyle.Render(" - " + f.Desc)

		if i == m.fieldCursor {
			name = selectedStyle.Render(f.Name)
		}

		sb.WriteString(fmt.Sprintf("%s%s %s%s\n", cursor, check, name, desc))
	}

	hint := dimStyle.Render("\n  j/k navigate  space toggle  a all/none  enter continue  esc back")
	sb.WriteString(hint)

	return sb.String()
}

func (m *model) renderCount() string {
	var sb strings.Builder

	sb.WriteString(selectedStyle.Render("  "+m.selectedGen.Name()) +
		dimStyle.Render(" - how many records?") + "\n\n")

	sb.WriteString("  " + m.textInput.View() + "\n")

	hint := dimStyle.Render("\n  enter generate  esc back")
	sb.WriteString(hint)

	return sb.String()
}

func (m *model) renderPreviewView() string {
	w := m.getSafeWidth()
	h := m.getContentHeight()

	// Pick content based on view mode
	content := m.preview
	if m.prettyView {
		content = m.formatPretty()
	}

	lines := strings.Split(content, "\n")
	total := len(lines)

	// Clamp scroll
	maxScroll := total - h + 2
	if maxScroll < 0 {
		maxScroll = 0
	}
	if m.previewScroll > maxScroll {
		m.previewScroll = maxScroll
	}

	var sb strings.Builder

	// Info bar
	viewLabel := "raw"
	if m.prettyView {
		viewLabel = "pretty"
	}
	info := selectedStyle.Render("  "+m.selectedGen.Name()) +
		dimStyle.Render(fmt.Sprintf("  %d records  %s  [%s]", m.count, output.FormatNames[m.format], viewLabel))
	sb.WriteString(info + "\n")

	// Preview content
	end := m.previewScroll + h - 2
	if end > total {
		end = total
	}
	visible := lines[m.previewScroll:end]

	previewW := w - 4
	if previewW < 20 {
		previewW = 20
	}

	box := previewBorderStyle.Width(previewW).Render(strings.Join(visible, "\n"))
	sb.WriteString(box + "\n")

	// Scroll indicator
	if total > h-2 {
		pct := 0
		if maxScroll > 0 {
			pct = m.previewScroll * 100 / maxScroll
		}
		sb.WriteString(dimStyle.Render(fmt.Sprintf("  %d/%d lines (%d%%)", m.previewScroll+1, total, pct)))
	}

	hint := dimStyle.Render("  j/k scroll  v view  r re-roll  +/- count  f format  c copy  e export  esc back")
	sb.WriteString("\n" + hint)

	return sb.String()
}

func (m *model) renderFormatPicker() string {
	var sb strings.Builder

	sb.WriteString(dimStyle.Render("  Select output format:") + "\n\n")

	for i, name := range output.FormatNames {
		cursor := "  "
		style := normalStyle
		if i == m.formatCursor {
			cursor = "> "
			style = selectedStyle
		}

		marker := "  "
		if output.Format(i) == m.format {
			marker = checkStyle.Render("* ")
		}

		sb.WriteString(cursor + marker + style.Render(name) + "\n")
	}

	hint := dimStyle.Render("\n  j/k navigate  enter select  esc back")
	sb.WriteString(hint)

	return sb.String()
}

func (m *model) renderExportView() string {
	var sb strings.Builder

	sb.WriteString(dimStyle.Render("  Export to file:") + "\n\n")
	sb.WriteString("  " + m.textInput.View() + "\n")
	sb.WriteString(dimStyle.Render(fmt.Sprintf("\n  format: %s  records: %d", output.FormatNames[m.format], m.count)) + "\n")

	hint := dimStyle.Render("\n  enter export  esc cancel")
	sb.WriteString(hint)

	return sb.String()
}

func (m *model) renderHelp() string {
	help := `
  SEEDBANK - Fake Data Generator

  GENERATOR LIST
    j/k       Navigate generators
    enter     Select generator
    m         Mix mode (combine generators)
    ?         Toggle this help
    q         Quit

  FIELD SELECTION
    j/k       Navigate fields
    space     Toggle field on/off
    a         Toggle all fields
    enter     Continue to count
    esc       Back

  COUNT INPUT
    type      Enter number of records
    enter     Generate data
    esc       Back

  PREVIEW
    j/k       Scroll preview
    g/G       Top/bottom
    v         Toggle raw/pretty view
    r         Re-roll (new random seed)
    +/-       Double/halve record count
    f         Change output format
    c         Copy to clipboard
    e         Export to file
    esc       Back

  FORMAT SELECTION
    j/k       Navigate formats
    enter     Select format
    esc       Back

  EXPORT
    type      Enter filename
    enter     Write file
    esc       Cancel

  FORMATS: JSON, JSONL, CSV, Markdown, SQL
`
	lines := strings.Split(help, "\n")
	h := m.getContentHeight()

	maxScroll := len(lines) - h
	if maxScroll < 0 {
		maxScroll = 0
	}
	if m.helpScroll > maxScroll {
		m.helpScroll = maxScroll
	}

	end := m.helpScroll + h
	if end > len(lines) {
		end = len(lines)
	}

	visible := lines[m.helpScroll:end]
	return dimStyle.Render(strings.Join(visible, "\n"))
}

func (m *model) renderMixSelect() string {
	var sb strings.Builder

	sb.WriteString(selectedStyle.Render("  Custom Mix") +
		dimStyle.Render(" — select generators to combine:") + "\n\n")

	selected := 0
	for _, on := range m.mixToggles {
		if on {
			selected++
		}
	}

	for i, gen := range m.generators {
		cursor := "  "
		if i == m.mixCursor {
			cursor = "> "
		}

		check := "[ ]"
		if i < len(m.mixToggles) && m.mixToggles[i] {
			check = checkStyle.Render("[x]")
		}

		name := normalStyle.Render(gen.Name())
		if i == m.mixCursor {
			name = selectedStyle.Render(gen.Name())
		}

		sb.WriteString(fmt.Sprintf("%s%s %s\n", cursor, check, name))
	}

	sb.WriteString(dimStyle.Render(fmt.Sprintf("\n  %d selected", selected)))
	hint := dimStyle.Render("\n  j/k navigate  space toggle  enter continue (2+ required)  esc back")
	sb.WriteString(hint)

	return sb.String()
}

func modeString(m mode) string {
	switch m {
	case modeGenerators:
		return "generators"
	case modeFields:
		return "fields"
	case modeCount:
		return "count"
	case modePreview:
		return "preview"
	case modeFormat:
		return "format"
	case modeExport:
		return "export"
	case modeHelp:
		return "help"
	case modeMixSelect:
		return "mix"
	}
	return ""
}
