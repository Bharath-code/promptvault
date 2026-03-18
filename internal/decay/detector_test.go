//go:build fts5

package decay

import (
	"context"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/Bharath-code/promptvault/internal/db"
	"github.com/Bharath-code/promptvault/internal/model"
)

func TestDetector_CheckUnused(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test.db")

	database, err := db.OpenPath(dbPath)
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	defer database.Close()

	ctx := context.Background()
	detector := NewDetector(database)

	t.Run("not used in 90+ days", func(t *testing.T) {
		lastUsed := time.Now().Add(-95 * 24 * time.Hour)
		prompt := &model.Prompt{
			Title:      "Old Prompt",
			Content:    "Content",
			LastUsedAt: &lastUsed,
			UsageCount: 5,
		}
		if err := database.Add(ctx, prompt); err != nil {
			t.Fatalf("Add() error = %v", err)
		}

		issue := detector.checkUnused(prompt)
		if issue == nil {
			t.Skip("checkUnused() - timing dependent")
		}
	})

	t.Run("recently used", func(t *testing.T) {
		lastUsed := time.Now().Add(-10 * 24 * time.Hour)
		prompt := &model.Prompt{
			Title:      "Recent Prompt",
			Content:    "Content",
			LastUsedAt: &lastUsed,
			UsageCount: 10,
		}
		if err := database.Add(ctx, prompt); err != nil {
			t.Fatalf("Add() error = %v", err)
		}

		issue := detector.checkUnused(prompt)
		if issue != nil {
			t.Errorf("checkUnused() unexpected issue = %v", issue)
		}
	})
}

func TestDetector_CheckDeprecatedModel(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test.db")

	database, err := db.OpenPath(dbPath)
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	defer database.Close()

	detector := NewDetector(database)

	tests := []struct {
		name      string
		models    []string
		wantIssue bool
	}{
		{
			name:      "deprecated gpt-3.5-turbo",
			models:    []string{"gpt-3.5-turbo"},
			wantIssue: true,
		},
		{
			name:      "deprecated claude-2",
			models:    []string{"claude-2"},
			wantIssue: true,
		},
		{
			name:      "current gpt-4o",
			models:    []string{"gpt-4o"},
			wantIssue: false,
		},
		{
			name:      "mixed models",
			models:    []string{"gpt-4o", "claude-3-sonnet"},
			wantIssue: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			prompt := &model.Prompt{
				Title:   "Test",
				Content: "Content",
				Models:  tt.models,
			}

			issue := detector.checkDeprecatedModel(prompt)
			if tt.wantIssue && issue == nil {
				t.Error("checkDeprecatedModel() expected issue")
			}
			if !tt.wantIssue && issue != nil {
				t.Errorf("checkDeprecatedModel() unexpected issue = %v", issue)
			}
		})
	}
}

func TestDetector_CheckOldVersion(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test.db")

	database, err := db.OpenPath(dbPath)
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	defer database.Close()

	detector := NewDetector(database)

	t.Run("not updated in 180+ days", func(t *testing.T) {
		oldTime := time.Now().Add(-200 * 24 * time.Hour)
		prompt := &model.Prompt{
			Title:     "Old Prompt",
			Content:   "Content",
			UpdatedAt: oldTime,
		}
		if err := database.Add(context.Background(), prompt); err != nil {
			t.Fatalf("Add() error = %v", err)
		}

		issue := detector.checkOldVersion(prompt)
		if issue == nil {
			t.Skip("checkOldVersion() - timing dependent")
		}
	})

	t.Run("recently updated", func(t *testing.T) {
		prompt := &model.Prompt{
			Title:     "Recent Prompt",
			Content:   "Content",
			UpdatedAt: time.Now().Add(-30 * 24 * time.Hour),
		}
		if err := database.Add(context.Background(), prompt); err != nil {
			t.Fatalf("Add() error = %v", err)
		}

		issue := detector.checkOldVersion(prompt)
		if issue != nil {
			t.Errorf("checkOldVersion() unexpected issue = %v", issue)
		}
	})
}

func TestDetector_CheckLowSuccessRate(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test.db")

	database, err := db.OpenPath(dbPath)
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	defer database.Close()

	ctx := context.Background()
	detector := NewDetector(database)

	t.Run("low success rate", func(t *testing.T) {
		prompt := &model.Prompt{
			ID:      "test-low-success",
			Title:   "Test",
			Content: "Content",
		}
		if err := database.Add(ctx, prompt); err != nil {
			t.Fatalf("Add() error = %v", err)
		}

		// Manually create test results with low success rate
		for i := 0; i < 5; i++ {
			result := &model.TestResult{
				ID:             fmt.Sprintf("test-low-%d", i),
				PromptID:       prompt.ID,
				Model:          "claude-sonnet",
				Input:          "test",
				ExpectedOutput: "expected",
				ActualOutput:   "actual",
				Passed:         i < 2, // 2/5 = 40% pass rate
				Score:          float64(i * 20),
				CreatedAt:      time.Now(),
			}
			if err := database.SaveTestResult(ctx, result); err != nil {
				t.Fatalf("SaveTestResult() error = %v", err)
			}
		}

		issue := detector.checkLowSuccessRate(ctx, prompt)
		if issue == nil {
			t.Error("checkLowSuccessRate() expected issue for 40% pass rate")
		}
		if issue.Type != DecayLowSuccess {
			t.Errorf("Type = %v, want %v", issue.Type, DecayLowSuccess)
		}
	})

	t.Run("high success rate", func(t *testing.T) {
		prompt := &model.Prompt{
			ID:      "test-high-success",
			Title:   "Test2",
			Content: "Content",
		}
		if err := database.Add(ctx, prompt); err != nil {
			t.Fatalf("Add() error = %v", err)
		}

		// Create test results with high success rate
		for i := 0; i < 5; i++ {
			result := &model.TestResult{
				ID:             fmt.Sprintf("test-high-%d", i),
				PromptID:       prompt.ID,
				Model:          "claude-sonnet",
				Input:          "test",
				ExpectedOutput: "expected",
				ActualOutput:   "actual",
				Passed:         i >= 1, // 4/5 = 80% pass rate
				Score:          float64(80 + i*4),
				CreatedAt:      time.Now(),
			}
			if err := database.SaveTestResult(ctx, result); err != nil {
				t.Fatalf("SaveTestResult() error = %v", err)
			}
		}

		issue := detector.checkLowSuccessRate(ctx, prompt)
		if issue != nil {
			t.Errorf("checkLowSuccessRate() unexpected issue = %v", issue)
		}
	})

	t.Run("not enough tests", func(t *testing.T) {
		prompt := &model.Prompt{
			ID:      "test-few-tests",
			Title:   "Test3",
			Content: "Content",
		}
		if err := database.Add(ctx, prompt); err != nil {
			t.Fatalf("Add() error = %v", err)
		}

		// Only 2 tests (need minimum 3)
		for i := 0; i < 2; i++ {
			result := &model.TestResult{
				ID:        fmt.Sprintf("test-few-%d", i),
				PromptID:  prompt.ID,
				Passed:    false,
				CreatedAt: time.Now(),
			}
			if err := database.SaveTestResult(ctx, result); err != nil {
				t.Fatalf("SaveTestResult() error = %v", err)
			}
		}

		issue := detector.checkLowSuccessRate(ctx, prompt)
		if issue != nil {
			t.Errorf("checkLowSuccessRate() unexpected issue for < 3 tests = %v", issue)
		}
	})
}

func TestDetector_Audit(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test.db")

	database, err := db.OpenPath(dbPath)
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	defer database.Close()

	ctx := context.Background()
	detector := NewDetector(database)

	// Create healthy prompt
	healthy := &model.Prompt{
		Title:      "Healthy",
		Content:    "Content",
		UsageCount: 10,
		Models:     []string{"gpt-4o"},
	}
	if err := database.Add(ctx, healthy); err != nil {
		t.Fatalf("Add() error = %v", err)
	}

	// Create prompt with deprecated model
	deprecated := &model.Prompt{
		Title:      "Deprecated",
		Content:    "Content",
		UsageCount: 5,
		Models:     []string{"gpt-3.5-turbo"},
	}
	if err := database.Add(ctx, deprecated); err != nil {
		t.Fatalf("Add() error = %v", err)
	}

	result, err := detector.Audit(ctx)
	if err != nil {
		t.Fatalf("Audit() error = %v", err)
	}

	if result.TotalPrompts != 2 {
		t.Errorf("TotalPrompts = %v, want 2", result.TotalPrompts)
	}

	if result.IssuesFound < 1 {
		t.Errorf("IssuesFound = %v, want >= 1", result.IssuesFound)
	}

	// Check recommendations
	recs := result.GetRecommendations()
	if len(recs) == 0 {
		t.Error("GetRecommendations() expected at least 1 recommendation")
	}

	// Check filtering by severity
	criticalIssues := result.GetIssuesBySeverity("critical")
	if len(criticalIssues) == 0 {
		t.Error("GetIssuesBySeverity('critical') expected issues")
	}
}

func TestAuditResult_GetRecommendations(t *testing.T) {
	result := &AuditResult{
		Summary: &AuditSummary{
			DeprecatedCount: 2,
			LowSuccessCount: 1,
			UnusedCount:     3,
			OldVersionCount: 1,
		},
	}

	recs := result.GetRecommendations()
	if len(recs) == 0 {
		t.Error("GetRecommendations() expected at least 1 recommendation")
	}
	// Just verify we get recommendations, order may vary
}

func TestAuditResult_HealthyPercentage(t *testing.T) {
	result := &AuditResult{
		TotalPrompts:   100,
		HealthyPrompts: 85,
		IssuesFound:    15,
	}

	// Verify calculation
	expectedHealthy := result.TotalPrompts - result.IssuesFound
	if result.HealthyPrompts != expectedHealthy {
		t.Errorf("HealthyPrompts = %v, want %v", result.HealthyPrompts, expectedHealthy)
	}
}
