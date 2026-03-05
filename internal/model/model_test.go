package model

import (
	"testing"
)

func TestSeedPrompts(t *testing.T) {
	seeds := SeedPrompts()
	if len(seeds) == 0 {
		t.Fatal("expected seed prompts, got 0")
	}

	for _, p := range seeds {
		if p.Title == "" {
			t.Errorf("seed prompt has empty title")
		}
		if p.Content == "" {
			t.Errorf("seed prompt %q has empty content", p.Title)
		}
		if p.Stack == "" {
			t.Errorf("seed prompt %q has empty stack", p.Title)
		}
	}
}
