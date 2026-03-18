//go:build fts5

package db

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/Bharath-code/promptvault/internal/model"
)

func TestSearchHistory(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test.db")

	database, err := OpenPath(dbPath)
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	defer database.Close()

	ctx := context.Background()

	// Add a prompt first (needed for foreign key if any)
	p := &model.Prompt{
		Title:   "Test Prompt",
		Content: "Test content for search",
	}
	if err := database.Add(ctx, p); err != nil {
		t.Fatalf("failed to add prompt: %v", err)
	}

	// Test: Clear history first
	if err := database.ClearSearchHistory(ctx); err != nil {
		t.Fatalf("failed to clear history: %v", err)
	}

	// Test: Empty history
	history, err := database.GetSearchHistory(ctx, 10)
	if err != nil {
		t.Fatalf("failed to get history: %v", err)
	}
	if len(history) != 0 {
		t.Errorf("expected empty history, got %d items", len(history))
	}

	// Test: Add search queries
	queries := []string{"react hooks", "api design", "authentication"}
	for _, q := range queries {
		if err := database.AddSearchHistory(ctx, q); err != nil {
			t.Fatalf("failed to add search history: %v", err)
		}
	}

	// Test: Get history
	history, err = database.GetSearchHistory(ctx, 10)
	if err != nil {
		t.Fatalf("failed to get history: %v", err)
	}
	if len(history) != 3 {
		t.Errorf("expected 3 history items, got %d", len(history))
	}

	// Test: Duplicate query increments count
	if err := database.AddSearchHistory(ctx, "react hooks"); err != nil {
		t.Fatalf("failed to add duplicate search: %v", err)
	}
	history, err = database.GetSearchHistory(ctx, 10)
	if err != nil {
		t.Fatalf("failed to get history after duplicate: %v", err)
	}
	if len(history) != 3 {
		t.Errorf("expected 3 items (duplicate should increment), got %d", len(history))
	}

	// Test: Limit history
	history, err = database.GetSearchHistory(ctx, 2)
	if err != nil {
		t.Fatalf("failed to get limited history: %v", err)
	}
	if len(history) != 2 {
		t.Errorf("expected 2 items with limit, got %d", len(history))
	}

	// Test: Delete single item
	if err := database.DeleteSearchHistoryItem(ctx, "api design"); err != nil {
		t.Fatalf("failed to delete history item: %v", err)
	}
	history, err = database.GetSearchHistory(ctx, 10)
	if err != nil {
		t.Fatalf("failed to get history after delete: %v", err)
	}
	if len(history) != 2 {
		t.Errorf("expected 2 items after delete, got %d", len(history))
	}

	// Test: Clear all history
	if err := database.ClearSearchHistory(ctx); err != nil {
		t.Fatalf("failed to clear history: %v", err)
	}
	history, err = database.GetSearchHistory(ctx, 10)
	if err != nil {
		t.Fatalf("failed to get history after clear: %v", err)
	}
	if len(history) != 0 {
		t.Errorf("expected 0 items after clear, got %d", len(history))
	}

	// Test: Empty query is ignored
	if err := database.AddSearchHistory(ctx, ""); err != nil {
		t.Fatalf("failed to handle empty query: %v", err)
	}
	history, err = database.GetSearchHistory(ctx, 10)
	if err != nil {
		t.Fatalf("failed to get history: %v", err)
	}
	if len(history) != 0 {
		t.Errorf("expected empty history for empty query, got %d", len(history))
	}

	// Test: Whitespace-only query is ignored
	if err := database.AddSearchHistory(ctx, "   "); err != nil {
		t.Fatalf("failed to handle whitespace query: %v", err)
	}
	history, err = database.GetSearchHistory(ctx, 10)
	if err != nil {
		t.Fatalf("failed to get history: %v", err)
	}
	if len(history) != 0 {
		t.Errorf("expected empty history for whitespace query, got %d", len(history))
	}
}
