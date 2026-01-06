package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestBuildPartsStore_UsesCLIDirs(t *testing.T) {
	tmp := t.TempDir()
	cwd := filepath.Join(tmp, "project")
	partsDir := filepath.Join(tmp, "custom_parts")
	if err := os.MkdirAll(filepath.Join(partsDir, "motors"), 0o755); err != nil {
		t.Fatalf("mkdir custom parts: %v", err)
	}
	if err := os.MkdirAll(cwd, 0o755); err != nil {
		t.Fatalf("mkdir cwd: %v", err)
	}

	partYAML := `part_id: motors/custom_motor
type: motor
name: CLI Motor
motor:
  voltage_min_v: 6
  voltage_max_v: 12
  nominal_current_a: 0.4
  stall_current_a: 1.0
`
	if err := os.WriteFile(filepath.Join(partsDir, "motors", "custom_motor.yaml"), []byte(partYAML), 0o644); err != nil {
		t.Fatalf("write custom part: %v", err)
	}

	oldWD, err := os.Getwd()
	if err != nil {
		t.Fatalf("get wd: %v", err)
	}
	if err := os.Chdir(cwd); err != nil {
		t.Fatalf("chdir: %v", err)
	}
	t.Cleanup(func() {
		_ = os.Chdir(oldWD)
	})

	store, err := buildPartsStore([]string{partsDir}, "")
	if err != nil {
		t.Fatalf("buildPartsStore: %v", err)
	}
	part, err := store.LoadMotor("motors/custom_motor")
	if err != nil {
		t.Fatalf("LoadMotor: %v", err)
	}
	if part.Name != "CLI Motor" {
		t.Fatalf("expected CLI part to load, got %q", part.Name)
	}
}

func TestPartsSearchDirs_Order(t *testing.T) {
	tmp := t.TempDir()
	cwd := filepath.Join(tmp, "project")
	cliDir := filepath.Join(tmp, "cli_parts")
	envDir := filepath.Join(tmp, "env_parts")

	got := partsSearchDirs(cwd, []string{cliDir}, envDir)
	if len(got) < 4 {
		t.Fatalf("expected at least 4 search dirs, got %d", len(got))
	}
	if got[0] != filepath.Join(cwd, "rv_parts") {
		t.Fatalf("expected rv_parts first, got %q", got[0])
	}
	if got[1] != filepath.Join(cwd, "parts") {
		t.Fatalf("expected built-in parts second, got %q", got[1])
	}
	if got[2] != filepath.Clean(cliDir) {
		t.Fatalf("expected cli dir third, got %q", got[2])
	}
	if got[3] != filepath.Clean(envDir) {
		t.Fatalf("expected env dir fourth, got %q", got[3])
	}
}
