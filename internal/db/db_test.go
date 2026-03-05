//go:build fts5

package db

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/Bharath-code/promptvault/internal/model"
)

func TestDB(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test.db")

	database, err := OpenPath(dbPath)
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	defer database.Close()

	ctx := context.Background()

	// 1. Test Add
	p := &model.Prompt{
		Title:   "Test",
		Content: "Hello World",
		Stack:   "test/stack",
	}
	if err := database.Add(ctx, p); err != nil {
		t.Fatalf("failed to add prompt: %v", err)
	}

	// 2. Test Get
	fetched, err := database.Get(ctx, p.ID)
	if err != nil {
		t.Fatalf("failed to get prompt: %v", err)
	}
	if fetched.Title != p.Title {
		t.Errorf("expected title %q, got %q", p.Title, fetched.Title)
	}

	// 3. Test Usage Increment
	if err := database.IncrementUsage(ctx, p.ID); err != nil {
		t.Fatalf("failed to increment usage: %v", err)
	}

	fetched, err = database.Get(ctx, p.ID)
	if err != nil {
		t.Fatalf("failed to get prompt after usage increment: %v", err)
	}
	if fetched.UsageCount != 1 {
		t.Errorf("expected usage 1, got %d", fetched.UsageCount)
	}

	// 4. Test Search
	results, err := database.Search(ctx, "World")
	if err != nil {
		t.Fatalf("search failed: %v", err)
	}
	if len(results) != 1 {
		t.Errorf("expected 1 search result, got %d", len(results))
	}

	// 5. Test Delete
	if err := database.Delete(ctx, p.ID); err != nil {
		t.Fatalf("delete failed: %v", err)
	}

	if _, err := database.Get(ctx, p.ID); err == nil {
		t.Errorf("expected error getting deleted prompt")
	}
}
