package db

import (
	"context"
	"path/filepath"
	"testing"
	"time"

	"github.com/Bharath-code/promptvault/internal/model"
)

func TestDB_Versioning(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test.db")

	database, err := OpenPath(dbPath)
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	defer database.Close()

	ctx := context.Background()

	// Create a prompt
	prompt := &model.Prompt{
		Title:   "Test Prompt",
		Content: "Original content",
		Stack:   "test/stack",
	}
	if err := database.Add(ctx, prompt); err != nil {
		t.Fatalf("failed to add prompt: %v", err)
	}

	t.Run("CreateVersion", func(t *testing.T) {
		// Create first version
		if err := database.CreateVersion(ctx, prompt, "Initial version", "testuser"); err != nil {
			t.Fatalf("CreateVersion() error = %v", err)
		}

		// Verify version was created
		versions, err := database.GetPromptHistory(ctx, prompt.ID)
		if err != nil {
			t.Fatalf("GetPromptHistory() error = %v", err)
		}

		if len(versions) != 1 {
			t.Errorf("GetPromptHistory() got %d versions, want 1", len(versions))
		}

		if versions[0].Version != 1 {
			t.Errorf("Version = %v, want 1", versions[0].Version)
		}

		if versions[0].Title != "Test Prompt" {
			t.Errorf("Title = %v, want 'Test Prompt'", versions[0].Title)
		}
	})

	t.Run("AutoVersionOnUpdate", func(t *testing.T) {
		// Update prompt (should auto-create version)
		prompt.Content = "Updated content"
		if err := database.Update(ctx, prompt, "Updated content", "testuser"); err != nil {
			t.Fatalf("Update() error = %v", err)
		}

		// Verify version was auto-created
		versions, err := database.GetPromptHistory(ctx, prompt.ID)
		if err != nil {
			t.Fatalf("GetPromptHistory() error = %v", err)
		}

		if len(versions) != 2 {
			t.Errorf("GetPromptHistory() got %d versions, want 2", len(versions))
		}

		if versions[0].Version != 2 {
			t.Errorf("Latest Version = %v, want 2", versions[0].Version)
		}
	})

	t.Run("GetCurrentVersion", func(t *testing.T) {
		version, err := database.GetCurrentVersion(ctx, prompt.ID)
		if err != nil {
			t.Fatalf("GetCurrentVersion() error = %v", err)
		}

		if version != 2 {
			t.Errorf("GetCurrentVersion() = %v, want 2", version)
		}
	})

	t.Run("GetPromptVersion", func(t *testing.T) {
		v1, err := database.GetPromptVersion(ctx, prompt.ID, 1)
		if err != nil {
			t.Fatalf("GetPromptVersion(1) error = %v", err)
		}

		if v1.Content != "Original content" {
			t.Errorf("V1 Content = %v, want 'Original content'", v1.Content)
		}

		v2, err := database.GetPromptVersion(ctx, prompt.ID, 2)
		if err != nil {
			t.Fatalf("GetPromptVersion(2) error = %v", err)
		}

		if v2.Content != "Updated content" {
			t.Errorf("V2 Content = %v, want 'Updated content'", v2.Content)
		}
	})

	t.Run("DeletePromptVersions", func(t *testing.T) {
		// Delete versions
		if err := database.DeletePromptVersions(ctx, prompt.ID); err != nil {
			t.Fatalf("DeletePromptVersions() error = %v", err)
		}

		// Verify versions are deleted
		versions, err := database.GetPromptHistory(ctx, prompt.ID)
		if err != nil {
			t.Fatalf("GetPromptHistory() error = %v", err)
		}

		if len(versions) != 0 {
			t.Errorf("GetPromptHistory() got %d versions, want 0", len(versions))
		}
	})
}

func TestDB_VersionCascadeDelete(t *testing.T) {
	t.Skip("Cascade delete requires SQLite foreign keys enabled - skipping for now")
	// This test would require: PRAGMA foreign_keys = ON;
	// Which needs to be set on every connection
}

func TestDB_VersionTimestamps(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test.db")

	database, err := OpenPath(dbPath)
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	defer database.Close()

	ctx := context.Background()

	prompt := &model.Prompt{
		Title:   "Test",
		Content: "Content",
	}
	if err := database.Add(ctx, prompt); err != nil {
		t.Fatalf("failed to add prompt: %v", err)
	}

	before := time.Now()
	if err := database.CreateVersion(ctx, prompt, "v1", "user"); err != nil {
		t.Fatalf("CreateVersion() error = %v", err)
	}
	after := time.Now()

	versions, _ := database.GetPromptHistory(ctx, prompt.ID)
	if len(versions) != 1 {
		t.Fatalf("Expected 1 version")
	}

	createdAt := versions[0].CreatedAt
	if createdAt.Before(before) || createdAt.After(after) {
		t.Errorf("CreatedAt %v not between %v and %v", createdAt, before, after)
	}
}
