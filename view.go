package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/LFroesch/seedbank/internal/output"
)

var (
	colorPrimary  = lipgloss.Color("#5AF78E")
	colorAccent   = lipgloss.Color("#57C7FF")
	colorCursorBg = lipgloss.Color("#2A2A40")
	colorWarn     = lipgloss.Color("#FF5555")
	colorDim      = lipgloss.Color("#606060")
	colorText     = lipgloss.Color("#EEEEEE")
	colorMuted    = lipgloss.Color("#B8B8B8")
	colorYellow   = lipgloss.Color("#F3F99D")

	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorPrimary)

	activeTabStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorAccent)

	selectedStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorAccent)

	normalStyle = lipgloss.NewStyle().
			Foreground(colorText)

	dimStyle = lipgloss.NewStyle().
			Foreground(colorDim)

	mutedStyle = lipgloss.NewStyle().
			Foreground(colorMuted)

	accentStyle = lipgloss.NewStyle().
			Foreground(colorAccent)

	checkStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorPrimary)

	statusStyle = lipgloss.NewStyle().
			Foreground(colorYellow)

	warnStyle = lipgloss.NewStyle().
			Foreground(colorWarn)

	panelStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(colorDim).
			Padding(0, 1)

	panelActiveStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(colorAccent).
				Padding(0, 1)

	panelHeaderStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(colorAccent)
)

func (m *model) View() string {
	if m.width == 0 || m.height == 0 {
		return "loading..."
	}

	header := m.renderHeader()
	sep := dimStyle.Render(strings.Repeat("─", m.width))
	body := m.renderBody()
	footer := m.renderFooter()

	parts := []string{header, sep, body}
	if status := m.renderStatusLine(); status != "" {
		parts = append(parts, status)
	}
	parts = append(parts, sep, footer)

	return lipgloss.JoinVertical(lipgloss.Left, parts...)
}

func (m *model) renderBody() string {
	switch m.mode {
	case modeGenerators:
		return m.renderGenerators()
	case modeFields:
		return m.renderFields()
	case modeCount:
		return m.renderCount()
	case modePreview:
		return m.renderPreviewView()
	case modeFormat:
		return m.renderFormatPicker()
	case modeExport:
		return m.renderExportView()
	case modeHelp:
		return m.renderHelp()
	case modeMixSelect:
		return m.renderMixSelect()
	default:
		return ""
	}
}

func (m *model) renderHeader() string {
	title := titleStyle.Render("seedbank")

	tabs := []struct {
		label string
		mode  mode
	}{
		{"Browse", modeGenerators},
		{"Fields", modeFields},
		{"Count", modeCount},
		{"Preview", modePreview},
		{"Format", modeFormat},
		{"Export", modeExport},
	}

	var renderedTabs []string
	for i, tab := range tabs {
		if i > 0 {
			renderedTabs = append(renderedTabs, dimStyle.Render(" │ "))
		}
		label := fmt.Sprintf("%d %s", i+1, tab.label)
		if m.mode == tab.mode {
			renderedTabs = append(renderedTabs, activeTabStyle.Render(label))
		} else {
			renderedTabs = append(renderedTabs, dimStyle.Render(label))
		}
	}
	if m.mode == modeHelp {
		renderedTabs = append(renderedTabs, dimStyle.Render(" │ "), activeTabStyle.Render("Help"))
	}
	if m.mode == modeMixSelect {
		renderedTabs = append(renderedTabs, dimStyle.Render(" │ "), activeTabStyle.Render("Mix"))
	}

	left := title + "  " + strings.Join(renderedTabs, "")

	var rightParts []string
	if m.selectedGen != nil {
		rightParts = append(rightParts, panelHeaderStyle.Render(m.selectedGen.Name()))
		rightParts = append(rightParts, mutedStyle.Render(fmt.Sprintf("%d fields", len(m.selectedGen.Fields()))))
		rightParts = append(rightParts, statusStyle.Render(fmt.Sprintf("%d rows", m.count)))
		rightParts = append(rightParts, accentStyle.Render(output.FormatNames[m.format]))
	} else {
		rightParts = append(rightParts, mutedStyle.Render(fmt.Sprintf("%d generators", len(m.generators))))
	}
	rightParts = append(rightParts, dimStyle.Render(fmt.Sprintf("seed %d", m.seed%10000)))
	right := strings.Join(rightParts, dimStyle.Render(" · "))
	maxRight := m.width / 2
	if maxRight < 12 {
		maxRight = 12
	}
	right = truncateString(right, maxRight)

	gap := m.width - lipgloss.Width(left) - lipgloss.Width(right)
	if gap < 2 {
		left = truncateString(left, m.width-lipgloss.Width(right)-2)
		gap = m.width - lipgloss.Width(left) - lipgloss.Width(right)
		if gap < 1 {
			gap = 1
		}
	}

	return left + strings.Repeat(" ", gap) + right
}

func (m *model) renderStatusLine() string {
	if m.statusMsg == "" {
		return ""
	}
	return statusStyle.Render("  " + m.statusMsg)
}

func (m *model) renderFooter() string {
	leftParts := []string{
		dimStyle.Render(modeString(m.mode)),
		mutedStyle.Render(m.footerContext()),
	}
	left := strings.Join(leftParts, dimStyle.Render("  "))

	rightParts := []string{
		dimStyle.Render(m.footerHints()),
	}
	if m.mode.usesPanelFocus() {
		rightParts = append(rightParts, accentStyle.Render("tab pane"))
	}
	rightParts = append(rightParts, accentStyle.Render("1-6 steps"))
	right := strings.Join(rightParts, dimStyle.Render("  ·  "))

	maxLeft := m.width / 3
	if maxLeft < 18 {
		maxLeft = 18
	}
	left = truncateString(left, maxLeft)

	gap := m.width - lipgloss.Width(left) - lipgloss.Width(right)
	if gap < 2 {
		right = truncateString(right, m.width-lipgloss.Width(left)-2)
		gap = m.width - lipgloss.Width(left) - lipgloss.Width(right)
		if gap < 1 {
			gap = 1
		}
	}

	return left + strings.Repeat(" ", gap) + right
}

func (m *model) renderGenerators() string {
	contentH := m.contentHeight()
	leftW, rightW := m.splitWidths()
	listH := m.generatorListBodyHeight()

	var lines []string
	for i, gen := range m.generators {
		if i < m.genScrollOff || i >= m.genScrollOff+listH {
			continue
		}

		prefix := "  "
		nameStyle := normalStyle
		descStyle := dimStyle
		if i == m.genCursor {
			prefix = accentStyle.Render("▸ ")
			nameStyle = selectedStyle
			descStyle = mutedStyle
		}

		line := prefix + nameStyle.Render(gen.Name())
		desc := truncateString(gen.Description(), leftW-10-lipgloss.Width(gen.Name()))
		if desc != "" {
			line += descStyle.Render("  " + desc)
		}
		lines = append(lines, line)
	}

	if len(lines) == 0 {
		lines = append(lines, dimStyle.Render("  no generators"))
	}

	leftBody := strings.Join(append([]string{
		"",
	}, lines...), "\n")

	rightLines := m.generatorDetailContent()

	rightTitle := "Selected"
	if m.activePane == paneRight {
		rightTitle += " •"
	}

	return m.renderResponsivePanels(
		m.panelTitle("Generators", paneLeft),
		leftBody,
		m.activePane == paneLeft,
		rightTitle,
		m.renderScrolledBody(rightLines, m.genDetailScroll, m.secondaryPanelBodyHeight()),
		m.activePane == paneRight,
		leftW,
		rightW,
		contentH,
	)
}

func (m *model) renderFields() string {
	if m.selectedGen == nil {
		return ""
	}

	contentH := m.contentHeight()
	leftW, rightW := m.splitWidths()

	fields := m.selectedGen.Fields()
	listH := m.primaryListBodyHeight()
	var lines []string
	for i, f := range fields {
		if i < m.fieldScrollOff || i >= m.fieldScrollOff+listH {
			continue
		}
		prefix := "  "
		check := "[ ]"
		nameStyle := normalStyle
		if i == m.fieldCursor {
			prefix = accentStyle.Render("▸ ")
			nameStyle = selectedStyle
		}
		if i < len(m.fieldToggles) && m.fieldToggles[i] {
			check = checkStyle.Render("[x]")
		}
		line := fmt.Sprintf("%s%s %s", prefix, check, nameStyle.Render(f.Name))
		desc := truncateString(f.Desc, leftW-12-lipgloss.Width(f.Name))
		if desc != "" {
			line += dimStyle.Render("  " + desc)
		}
		lines = append(lines, line)
	}

	return m.renderResponsivePanels(
		m.panelTitle(m.selectedGen.Name(), paneLeft),
		strings.Join(lines, "\n"),
		m.activePane == paneLeft,
		m.panelTitle("Summary", paneRight),
		m.renderScrolledBody(m.fieldSummaryContent(), m.fieldSummaryScroll, m.secondaryPanelBodyHeight()),
		m.activePane == paneRight,
		leftW,
		rightW,
		contentH,
	)
}

func (m *model) renderCount() string {
	contentH := m.contentHeight()
	leftW, rightW := m.splitWidths()

	fieldCount := len(m.getSelectedFieldNames())
	leftBody := strings.Join([]string{
		normalStyle.Render("How many rows do you want?"),
		"",
		"  " + m.textInput.View(),
		"",
		dimStyle.Render("Press enter to generate preview output."),
	}, "\n")

	rightBody := strings.Join([]string{
		panelHeaderStyle.Render(m.selectedGen.Name()),
		"  " + mutedStyle.Render(m.selectedGen.Description()),
		"",
		panelHeaderStyle.Render("Current selection"),
		"  " + statusStyle.Render(fmt.Sprintf("%d fields", fieldCount)),
		"  " + accentStyle.Render(output.FormatNames[m.format]),
		"",
		panelHeaderStyle.Render("Navigation"),
		"  " + dimStyle.Render("enter → preview  ·  esc → fields  ·  1-6 → steps"),
	}, "\n")

	return m.renderResponsivePanels("Count", leftBody, true, "Ready To Generate", rightBody, false, leftW, rightW, contentH)
}

func (m *model) renderPreviewView() string {
	contentH := m.contentHeight()
	leftW, rightW := m.splitWidths()

	content := m.preview
	viewLabel := "raw"
	if m.prettyView {
		content = m.formatPretty()
		viewLabel = "pretty"
	}

	lines := strings.Split(content, "\n")
	total := len(lines)
	visibleH := contentH - 4
	if visibleH < 3 {
		visibleH = 3
	}

	maxScroll := total - visibleH
	if maxScroll < 0 {
		maxScroll = 0
	}
	if m.previewScroll > maxScroll {
		m.previewScroll = maxScroll
	}

	end := m.previewScroll + visibleH
	if end > total {
		end = total
	}

	leftBodyLines := []string{
		dimStyle.Render(fmt.Sprintf("%s view · line %d/%d", viewLabel, minInt(m.previewScroll+1, total), total)),
		"",
	}
	leftBodyLines = append(leftBodyLines, strings.Join(lines[m.previewScroll:end], "\n"))
	if total > visibleH {
		leftBodyLines = append(leftBodyLines, "", dimStyle.Render(fmt.Sprintf("%d%% scrolled", percent(m.previewScroll, maxScroll))))
	}

	return m.renderResponsivePanels(
		m.panelTitle("Output", paneLeft),
		strings.Join(leftBodyLines, "\n"),
		m.activePane == paneLeft,
		m.panelTitle("Inspector", paneRight),
		m.renderScrolledBody(m.previewInspectorContent(), m.previewInspectorScroll, m.secondaryPanelBodyHeight()),
		m.activePane == paneRight,
		leftW,
		rightW,
		contentH,
	)
}

func (m *model) renderFormatPicker() string {
	contentH := m.contentHeight()
	leftW, rightW := m.splitWidths()

	var lines []string
	for i, name := range output.FormatNames {
		prefix := "  "
		marker := "  "
		style := normalStyle
		if i == m.formatCursor {
			prefix = accentStyle.Render("▸ ")
			style = selectedStyle
		}
		if output.Format(i) == m.format {
			marker = checkStyle.Render("● ")
		}
		lines = append(lines, prefix+marker+style.Render(name))
	}

	rightBody := strings.Join([]string{
		normalStyle.Render("Switch output format without leaving the generation flow."),
		"",
		panelHeaderStyle.Render("Current"),
		"  " + accentStyle.Render(output.FormatNames[m.format]),
		"",
		panelHeaderStyle.Render("Formats"),
		"  " + dimStyle.Render("JSON / JSONL for APIs"),
		"  " + dimStyle.Render("CSV / Markdown for review"),
		"  " + dimStyle.Render("SQL for direct inserts"),
	}, "\n")

	return m.renderResponsivePanels("Formats", strings.Join(lines, "\n"), true, "Notes", rightBody, false, leftW, rightW, contentH)
}

func (m *model) renderExportView() string {
	contentH := m.contentHeight()
	leftW, rightW := m.splitWidths()

	leftBody := strings.Join([]string{
		normalStyle.Render("Write the current preview to disk."),
		"",
		"  " + m.textInput.View(),
		"",
		dimStyle.Render("Relative paths write into the launch directory unless output_dir is configured."),
	}, "\n")

	rightBody := strings.Join([]string{
		panelHeaderStyle.Render("Export details"),
		"  " + statusStyle.Render(fmt.Sprintf("%d rows", m.count)),
		"  " + accentStyle.Render(output.FormatNames[m.format]),
		"  " + mutedStyle.Render(formatExtension(m.format)),
		"",
		panelHeaderStyle.Render("Keys"),
		"  " + dimStyle.Render("enter → write file"),
		"  " + dimStyle.Render("esc → return to preview"),
	}, "\n")

	return m.renderResponsivePanels("Export", leftBody, true, "Summary", rightBody, false, leftW, rightW, contentH)
}

func (m *model) renderHelp() string {
	lines := m.helpLines()
	body := m.renderScrolledBody(lines, m.helpScroll, m.contentHeight()-2)
	return m.renderPanel("Help", body, m.width, m.contentHeight(), true)
}

func (m *model) helpLines() []string {
	return []string{
		panelHeaderStyle.Render("Seedbank"),
		normalStyle.Render("Fake data generator for terminal-first seeding and fixture work."),
		"",
		panelHeaderStyle.Render("Browse"),
		"  j/k navigate generators",
		"  enter select generator",
		"  m custom mix mode",
		"  tab switch active pane",
		"",
		panelHeaderStyle.Render("Fields"),
		"  space toggle field",
		"  a all/none",
		"  enter continue",
		"  tab switch active pane",
		"",
		panelHeaderStyle.Render("Preview"),
		"  j/k scroll  pgup/pgdn faster  g/G top/bottom",
		"  v raw/pretty  r re-roll",
		"  +/- count  f format",
		"  c copy  e export",
		"  tab switch active pane",
		"",
		panelHeaderStyle.Render("Global"),
		"  1-6 jump between workflow steps",
		"  ? toggle help",
		"  q quit or back",
		"  esc back",
	}
}

func (m *model) renderMixSelect() string {
	contentH := m.contentHeight()
	leftW, rightW := m.splitWidths()

	var lines []string
	for i, gen := range m.generators {
		if i < m.mixScrollOff || i >= m.mixScrollOff+m.primaryListBodyHeight() {
			continue
		}
		prefix := "  "
		check := "[ ]"
		nameStyle := normalStyle
		if i == m.mixCursor {
			prefix = accentStyle.Render("▸ ")
			nameStyle = selectedStyle
		}
		if i < len(m.mixToggles) && m.mixToggles[i] {
			check = checkStyle.Render("[x]")
		}
		lines = append(lines, fmt.Sprintf("%s%s %s", prefix, check, nameStyle.Render(gen.Name())))
	}

	return m.renderResponsivePanels(
		m.panelTitle("Custom Mix", paneLeft),
		strings.Join(lines, "\n"),
		m.activePane == paneLeft,
		m.panelTitle("What This Does", paneRight),
		m.renderScrolledBody(m.mixSummaryContent(), m.mixSummaryScroll, m.secondaryPanelBodyHeight()),
		m.activePane == paneRight,
		leftW,
		rightW,
		contentH,
	)
}

func (m *model) renderPanel(title, body string, width, height int, active bool) string {
	style := panelStyle
	if active {
		style = panelActiveStyle
	}

	frameW := style.GetHorizontalFrameSize()
	frameH := style.GetVerticalFrameSize()
	borderW := style.GetHorizontalBorderSize()
	borderH := style.GetVerticalBorderSize()
	if width < frameW+16 {
		width = frameW + 16
	}
	if height < frameH+3 {
		height = frameH + 3
	}

	contentW := width - frameW
	innerH := height - frameH
	styleW := width - borderW
	styleH := height - borderH
	content := fitPanelContent(title, body, contentW, innerH)

	return style.
		Width(styleW).
		Height(styleH).
		Render(content)
}

func (m *model) splitWidths() (int, int) {
	gap := 1
	usable := m.width - gap
	minPanelWidth := panelStyle.GetHorizontalFrameSize() + 16
	if usable < minPanelWidth*2 {
		usable = minPanelWidth * 2
	}

	left := usable / 2
	right := usable - left
	return left, right
}

func (m *model) shouldStackPanels() bool {
	minSideBySideWidth := (panelStyle.GetHorizontalFrameSize() + 16) * 2
	return m.width < minSideBySideWidth
}

func (m *model) renderResponsivePanels(leftTitle, leftBody string, leftActive bool, rightTitle, rightBody string, rightActive bool, leftW, rightW, contentH int) string {
	if m.shouldStackPanels() {
		stackH := (contentH - 1) / 2
		if stackH < 6 {
			stackH = 6
		}
		bottomH := contentH - stackH
		if bottomH < 6 {
			bottomH = 6
		}

		fullW := m.width
		top := m.renderPanel(leftTitle, leftBody, fullW, stackH, leftActive)
		bottom := m.renderPanel(rightTitle, rightBody, fullW, bottomH, rightActive)
		return lipgloss.JoinVertical(lipgloss.Left, top, bottom)
	}

	leftPanel := m.renderPanel(leftTitle, leftBody, leftW, contentH, leftActive)
	rightPanel := m.renderPanel(rightTitle, rightBody, rightW, contentH, rightActive)
	return lipgloss.JoinHorizontal(lipgloss.Top, leftPanel, " ", rightPanel)
}

func (m *model) contentHeight() int {
	h := m.height - 4
	if m.statusMsg != "" {
		h--
	}
	if m.shouldStackPanels() {
		h--
	}
	if h < 8 {
		h = 8
	}
	return h
}

func (m *model) footerHints() string {
	switch m.mode {
	case modeGenerators:
		return "j/k move  pgup/pgdn jump  enter select  m mix  ? help"
	case modeFields:
		return "j/k move  space toggle  a all/none  enter continue  esc back"
	case modeCount:
		return "type count  enter generate  esc back"
	case modePreview:
		return "j/k scroll  pgup/pgdn jump  v view  r re-roll  +/- count  f format  c copy  e export"
	case modeFormat:
		return "j/k move  enter select  esc back"
	case modeExport:
		return "type filename  enter export  esc back"
	case modeHelp:
		return "j/k scroll  pgup/pgdn jump  ? close  esc close"
	case modeMixSelect:
		return "j/k move  space toggle  pgup/pgdn jump  enter continue  esc back"
	default:
		return ""
	}
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
	default:
		return ""
	}
}

func truncateString(s string, max int) string {
	if max <= 0 {
		return ""
	}
	if lipgloss.Width(s) <= max {
		return s
	}
	runes := []rune(s)
	if max <= 3 {
		return string(runes[:max])
	}
	if len(runes) > max-3 {
		runes = runes[:max-3]
	}
	return string(runes) + "..."
}

func fitPanelContent(title, body string, width, height int) string {
	if width < 1 {
		width = 1
	}
	if height < 1 {
		height = 1
	}

	lines := []string{truncateString(title, width)}
	if height == 1 {
		return panelHeaderStyle.Render(lines[0])
	}

	bodyLines := strings.Split(body, "\n")
	remaining := height - 1
	for i := 0; i < remaining; i++ {
		if i < len(bodyLines) {
			lines = append(lines, truncateString(bodyLines[i], width))
		} else {
			lines = append(lines, "")
		}
	}

	lines[0] = panelHeaderStyle.Render(lines[0])
	return strings.Join(lines, "\n")
}

func (m *model) panelTitle(title string, p pane) string {
	if m.activePane == p {
		return title + " •"
	}
	return title
}

func (m *model) renderScrolledBody(lines []string, scroll, visible int) string {
	if visible < 1 {
		visible = 1
	}
	scroll = clampScroll(scroll, len(lines), visible, false)
	end := scroll + visible
	if end > len(lines) {
		end = len(lines)
	}
	visibleLines := lines[scroll:end]
	if len(visibleLines) == 0 {
		visibleLines = []string{""}
	}
	if len(lines) > visible {
		visibleLines = append(visibleLines, "", dimStyle.Render(fmt.Sprintf("%d%% scrolled", percent(scroll, len(lines)-visible))))
	}
	return strings.Join(visibleLines, "\n")
}

func (m *model) fieldSummaryContent() []string {
	fields := m.selectedGen.Fields()
	selected := m.getSelectedFieldNames()
	currentField := fields[m.fieldCursor]
	lines := []string{
		normalStyle.Render(m.selectedGen.Description()),
		"",
		dimStyle.Render("Field sources are better for Mix. Person and Products already generate coherent records."),
		"",
		panelHeaderStyle.Render("Selection"),
		fmt.Sprintf("  %s of %s fields enabled",
			statusStyle.Render(fmt.Sprintf("%d", len(selected))),
			mutedStyle.Render(fmt.Sprintf("%d", len(fields)))),
		"",
		panelHeaderStyle.Render("Cursor"),
		"  " + accentStyle.Render(currentField.Name),
		"  " + dimStyle.Render(currentField.Desc),
	}
	if len(selected) == 0 {
		lines = append(lines, "", panelHeaderStyle.Render("Enabled"), "  "+warnStyle.Render("none selected, all fields will be used on generate"))
	} else {
		previewCount := len(selected)
		if previewCount > 6 {
			previewCount = 6
		}
		lines = append(lines, "", panelHeaderStyle.Render("Enabled"))
		for _, name := range selected[:previewCount] {
			lines = append(lines, "  "+accentStyle.Render("•")+" "+normalStyle.Render(name))
		}
		if len(selected) > previewCount {
			lines = append(lines, "  "+dimStyle.Render(fmt.Sprintf("+ %d more", len(selected)-previewCount)))
		}
	}
	lines = append(lines, "", panelHeaderStyle.Render("Next"))
	lines = append(lines, "  "+dimStyle.Render("enter → count  ·  a → all/none  ·  tab → pane  ·  esc → back"))
	return lines
}

func (m *model) mixSummaryContent() []string {
	selected := 0
	for _, on := range m.mixToggles {
		if on {
			selected++
		}
	}
	return []string{
		normalStyle.Render("Mix mode combines field sources into one custom record shape."),
		"",
		dimStyle.Render("Use Person or Products when you want coherence out of the box. Use Mix when you want to compose a shape manually."),
		"",
		panelHeaderStyle.Render("Selected"),
		"  " + statusStyle.Render(fmt.Sprintf("%d generators", selected)),
		"",
		panelHeaderStyle.Render("Rule"),
		"  " + dimStyle.Render("Need at least 2 generators before continuing."),
		"",
		panelHeaderStyle.Render("Next"),
		"  " + dimStyle.Render("enter → combined field picker  ·  tab → pane"),
	}
}

func (m *model) footerContext() string {
	if m.selectedGen != nil {
		return fmt.Sprintf("%s · %d rows · %s", m.selectedGen.Name(), m.count, output.FormatNames[m.format])
	}
	return fmt.Sprintf("%d generators", len(m.generators))
}

func percent(pos, max int) int {
	if max <= 0 {
		return 0
	}
	return pos * 100 / max
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
