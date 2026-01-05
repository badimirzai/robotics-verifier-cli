package cmd

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func runInitCommand(t *testing.T, args ...string) (string, error) {
	t.Helper()
	cmd := newInitCmd()
	var buf bytes.Buffer
	cmd.SetOut(&buf)
	cmd.SetErr(&buf)
	cmd.SetArgs(args)
	err := cmd.Execute()
	return buf.String(), err
}

func TestInitListTemplates(t *testing.T) {
	output, err := runInitCommand(t, "--list")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	for _, name := range []string{"4wd-problem", "4wd-clean"} {
		if !strings.Contains(output, name) {
			t.Fatalf("expected template %q in output, got %q", name, output)
		}
	}
}

func TestInitWritesTemplate(t *testing.T) {
	tmpDir := t.TempDir()
	outPath := filepath.Join(tmpDir, "robot.yaml")

	if _, err := runInitCommand(t, "--template", "4wd-problem", "--out", outPath); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("read output: %v", err)
	}
	if len(data) == 0 {
		t.Fatalf("expected non-empty output file")
	}
}

func TestInitExistingFileWithoutForceFails(t *testing.T) {
	tmpDir := t.TempDir()
	outPath := filepath.Join(tmpDir, "robot.yaml")

	if err := os.WriteFile(outPath, []byte("old"), 0o644); err != nil {
		t.Fatalf("write temp file: %v", err)
	}

	_, err := runInitCommand(t, "--template", "4wd-problem", "--out", outPath)
	if err == nil {
		t.Fatalf("expected error when output exists without --force")
	}
	if !strings.Contains(err.Error(), "--force") {
		t.Fatalf("expected --force hint in error, got %v", err)
	}

	data, readErr := os.ReadFile(outPath)
	if readErr != nil {
		t.Fatalf("read output: %v", readErr)
	}
	if string(data) != "old" {
		t.Fatalf("expected existing file to remain untouched")
	}
}

func TestInitForceOverwrites(t *testing.T) {
	tmpDir := t.TempDir()
	outPath := filepath.Join(tmpDir, "robot.yaml")

	if err := os.WriteFile(outPath, []byte("old"), 0o644); err != nil {
		t.Fatalf("write temp file: %v", err)
	}

	if _, err := runInitCommand(t, "--template", "4wd-problem", "--out", outPath, "--force"); err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	data, err := os.ReadFile(outPath)
	if err != nil {
		t.Fatalf("read output: %v", err)
	}
	if string(data) == "old" {
		t.Fatalf("expected file contents to be overwritten")
	}
	if len(data) == 0 {
		t.Fatalf("expected non-empty output file")
	}
}
