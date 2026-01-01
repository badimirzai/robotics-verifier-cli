package ui

import "testing"

func TestColorizeDisabledReturnsOriginal(t *testing.T) {
	EnableColors(false)
	defer EnableColors(DefaultColorEnabled())

	msg := "plain"
	if got := Colorize("ERROR", msg); got != msg {
		t.Fatalf("expected %q, got %q", msg, got)
	}
}
