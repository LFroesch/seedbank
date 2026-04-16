package main

import (
	"path/filepath"
	"testing"
)

func TestHandleStepJumpIgnoresFocusedTextInput(t *testing.T) {
	m := initialModel()
	m.mode = modeCount
	m.textInput.Focus()

	if jumped := m.handleStepJump("1"); jumped {
		t.Fatal("expected step jump to be ignored while typing")
	}
	if m.mode != modeCount {
		t.Fatalf("mode = %v, want count", m.mode)
	}
}

func TestHandleStepJumpBuildsPreviewState(t *testing.T) {
	m := initialModel()

	if jumped := m.handleStepJump("4"); !jumped {
		t.Fatal("expected preview jump to succeed")
	}
	if m.mode != modePreview {
		t.Fatalf("mode = %v, want preview", m.mode)
	}
	if m.selectedGen == nil {
		t.Fatal("expected step jump to select a generator")
	}
	if len(m.records) == 0 {
		t.Fatal("expected preview jump to generate records")
	}
}

func TestHandleStepJumpMovesToPreviewWhenReady(t *testing.T) {
	m := initialModel()
	m.selectedGen = m.generators[0]
	m.fieldToggles = make([]bool, len(m.selectedGen.Fields()))
	for i := range m.fieldToggles {
		m.fieldToggles[i] = true
	}
	m.count = 2
	m.generate()

	if jumped := m.handleStepJump("4"); !jumped {
		t.Fatal("expected preview jump to succeed")
	}
	if m.mode != modePreview {
		t.Fatalf("mode = %v, want preview", m.mode)
	}
}

func TestHandleStepJumpCanOpenExport(t *testing.T) {
	m := initialModel()

	if jumped := m.handleStepJump("6"); !jumped {
		t.Fatal("expected export jump to succeed")
	}
	if m.mode != modeExport {
		t.Fatalf("mode = %v, want export", m.mode)
	}
	if !m.textInput.Focused() {
		t.Fatal("expected export input to be focused")
	}
}

func TestResolveOutputPathUsesConfiguredBase(t *testing.T) {
	m := initialModel()
	m.config.OutputDir = filepath.Join("fixtures", "seedbank")

	got := m.resolveOutputPath("users/output.json")
	want := filepath.Join("fixtures", "seedbank", "users", "output.json")
	if got != want {
		t.Fatalf("resolveOutputPath() = %q, want %q", got, want)
	}
}

func TestResolveOutputPathPreservesAbsolutePaths(t *testing.T) {
	m := initialModel()
	abs := filepath.Join(string(filepath.Separator), "tmp", "seedbank.json")

	if got := m.resolveOutputPath(abs); got != abs {
		t.Fatalf("resolveOutputPath() = %q, want %q", got, abs)
	}
}
