package main

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
)

func TestRenderPanelUsesOuterDimensions(t *testing.T) {
	m := initialModel()

	panel := m.renderPanel("Title", "line 1\nline 2\nline 3", 40, 12, true)
	if got := lipgloss.Width(panel); got != 40 {
		t.Fatalf("panel width = %d, want 40", got)
	}
	if got := lipgloss.Height(panel); got != 12 {
		t.Fatalf("panel height = %d, want 12", got)
	}
}

func TestViewKeepsHeaderVisibleAtWideWidth(t *testing.T) {
	m := initialModel()
	m.width = 100
	m.height = 28

	view := m.View()
	if !strings.Contains(view, "seedbank") {
		t.Fatal("view does not contain header title")
	}
	if !strings.Contains(view, "1 Browse") {
		t.Fatal("view does not contain numbered workflow tabs")
	}
	if got := lipgloss.Height(view); got != m.height {
		t.Fatalf("view height = %d, want %d", got, m.height)
	}
}

func TestViewKeepsHeaderVisibleInStackedLayout(t *testing.T) {
	m := initialModel()
	m.width = 48
	m.height = 24

	view := m.View()
	if !strings.Contains(view, "seedbank") {
		t.Fatal("view does not contain header title")
	}
	if got := lipgloss.Height(view); got != m.height {
		t.Fatalf("stacked view height = %d, want %d", got, m.height)
	}
}

func TestSplitWidthsStayBalanced(t *testing.T) {
	m := initialModel()
	m.width = 101

	left, right := m.splitWidths()
	if diff := left - right; diff < -1 || diff > 1 {
		t.Fatalf("split widths too uneven: left=%d right=%d", left, right)
	}
	if left+right+1 != m.width {
		t.Fatalf("split widths + gap = %d, want %d", left+right+1, m.width)
	}
}

func TestFooterShowsStepHints(t *testing.T) {
	m := initialModel()
	m.width = 160

	footer := m.renderFooter()
	if !strings.Contains(footer, "1-6 steps") {
		t.Fatal("footer missing step hints")
	}
}
