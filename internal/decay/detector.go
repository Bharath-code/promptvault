package decay

import (
	"context"
	"fmt"
	"time"

	"github.com/Bharath-code/promptvault/internal/db"
	"github.com/Bharath-code/promptvault/internal/model"
)

// DecayType represents the type of decay detected
type DecayType string

const (
	DecayUnused      DecayType = "unused"       // Not used in 90+ days
	DecayLowSuccess  DecayType = "low_success"  // Test success rate < 50%
	DecayDeprecated  DecayType = "deprecated"   // Uses deprecated model
	DecayOldVersion  DecayType = "old_version"  // Not updated in 180+ days
	DecayLowQuality  DecayType = "low_quality"  // Quality score < 50
)

// DecayIssue represents a detected decay issue
type DecayIssue struct {
	Prompt      *model.Prompt
	Type        DecayType
	Severity    string // "critical", "warning", "info"
	Description string
	Suggestion  string
	Details     map[string]interface{}
	DetectedAt  time.Time
}

// AuditResult contains the full audit results
type AuditResult struct {
	TotalPrompts   int
	HealthyPrompts int
	IssuesFound    int
	Issues         []*DecayIssue
	Summary        *AuditSummary
	GeneratedAt    time.Time
}

// AuditSummary provides aggregate statistics
type AuditSummary struct {
	UnusedCount       int
	LowSuccessCount   int
	DeprecatedCount   int
	OldVersionCount   int
	LowQualityCount   int
	AverageQuality    float64
	AverageLastUsed   time.Time
	MostUsedPrompt    *model.Prompt
	LeastUsedPrompt   *model.Prompt
}

// Detector performs decay detection on prompts
type Detector struct {
	db *db.DB
}

// NewDetector creates a new decay detector
func NewDetector(database *db.DB) *Detector {
	return &Detector{
		db: database,
	}
}

// Audit runs a full decay audit on all prompts
func (d *Detector) Audit(ctx context.Context) (*AuditResult, error) {
	prompts, err := d.db.List(ctx, "")
	if err != nil {
		return nil, fmt.Errorf("listing prompts: %w", err)
	}

	result := &AuditResult{
		TotalPrompts: len(prompts),
		Issues:       make([]*DecayIssue, 0),
		GeneratedAt:  time.Now().UTC(),
	}

	summary := &AuditSummary{}

	for _, prompt := range prompts {
		// Check for decay issues
		issues := d.detectIssues(ctx, prompt)
		result.Issues = append(result.Issues, issues...)
	}

	// Calculate summary
	result.IssuesFound = len(result.Issues)
	result.HealthyPrompts = result.TotalPrompts - result.IssuesFound

	// Count by type
	for _, issue := range result.Issues {
		switch issue.Type {
		case DecayUnused:
			summary.UnusedCount++
		case DecayLowSuccess:
			summary.LowSuccessCount++
		case DecayDeprecated:
			summary.DeprecatedCount++
		case DecayOldVersion:
			summary.OldVersionCount++
		case DecayLowQuality:
			summary.LowQualityCount++
		}
	}

	result.Summary = summary
	return result, nil
}

// detectIssues checks a single prompt for decay issues
func (d *Detector) detectIssues(ctx context.Context, prompt *model.Prompt) []*DecayIssue {
	var issues []*DecayIssue

	// Check if unused
	if issue := d.checkUnused(prompt); issue != nil {
		issues = append(issues, issue)
	}

	// Check for deprecated models
	if issue := d.checkDeprecatedModel(prompt); issue != nil {
		issues = append(issues, issue)
	}

	// Check if old version
	if issue := d.checkOldVersion(prompt); issue != nil {
		issues = append(issues, issue)
	}

	// Check test success rate
	if issue := d.checkLowSuccessRate(ctx, prompt); issue != nil {
		issues = append(issues, issue)
	}

	return issues
}

// checkUnused detects prompts not used in 90+ days
func (d *Detector) checkUnused(prompt *model.Prompt) *DecayIssue {
	if prompt.UsageCount == 0 {
		daysSinceCreated := time.Since(prompt.CreatedAt).Hours() / 24
		if daysSinceCreated > 30 {
			return &DecayIssue{
				Prompt:   prompt,
				Type:     DecayUnused,
				Severity: "warning",
				Description: fmt.Sprintf("Never used in %.0f days", daysSinceCreated),
				Suggestion: "Consider testing this prompt or removing if obsolete",
				Details: map[string]interface{}{
					"days_since_created": int(daysSinceCreated),
					"usage_count":        prompt.UsageCount,
				},
				DetectedAt: time.Now().UTC(),
			}
		}
	}

	if prompt.LastUsedAt != nil {
		daysSinceUsed := time.Since(*prompt.LastUsedAt).Hours() / 24
		if daysSinceUsed > 90 {
			return &DecayIssue{
				Prompt:   prompt,
				Type:     DecayUnused,
				Severity: "warning",
				Description: fmt.Sprintf("Not used in %.0f days", daysSinceUsed),
				Suggestion: "Review if this prompt is still relevant",
				Details: map[string]interface{}{
					"days_since_used": int(daysSinceUsed),
					"last_used":       prompt.LastUsedAt,
				},
				DetectedAt: time.Now().UTC(),
			}
		}
	}

	return nil
}

// checkDeprecatedModel detects use of deprecated models
func (d *Detector) checkDeprecatedModel(prompt *model.Prompt) *DecayIssue {
	deprecatedModels := map[string]string{
		"gpt-3.5-turbo":      "Use gpt-4o or gpt-4-turbo instead",
		"gpt-4-turbo":        "Use gpt-4o instead",
		"claude-2":           "Use claude-3-sonnet or claude-3-opus instead",
		"claude-instant":     "Use claude-3-haiku instead",
		"text-davinci-003":   "Use gpt-4o instead",
		"text-davinci-002":   "Use gpt-4o instead",
		"code-davinci-002":   "Use gpt-4o instead",
	}

	for _, model := range prompt.Models {
		if suggestion, ok := deprecatedModels[model]; ok {
			return &DecayIssue{
				Prompt:   prompt,
				Type:     DecayDeprecated,
				Severity: "critical",
				Description: fmt.Sprintf("Uses deprecated model: %s", model),
				Suggestion: suggestion,
				Details: map[string]interface{}{
					"deprecated_model": model,
					"all_models":       prompt.Models,
				},
				DetectedAt: time.Now().UTC(),
			}
		}
	}

	return nil
}

// checkOldVersion detects prompts not updated in 180+ days
func (d *Detector) checkOldVersion(prompt *model.Prompt) *DecayIssue {
	daysSinceUpdated := time.Since(prompt.UpdatedAt).Hours() / 24
	if daysSinceUpdated > 180 {
		return &DecayIssue{
			Prompt:   prompt,
			Type:     DecayOldVersion,
			Severity: "info",
			Description: fmt.Sprintf("Not updated in %.0f days", daysSinceUpdated),
			Suggestion: "Review and update if needed",
			Details: map[string]interface{}{
				"days_since_updated": int(daysSinceUpdated),
				"last_updated":       prompt.UpdatedAt,
			},
			DetectedAt: time.Now().UTC(),
		}
	}

	return nil
}

// checkLowSuccessRate detects prompts with low test success rate
func (d *Detector) checkLowSuccessRate(ctx context.Context, prompt *model.Prompt) *DecayIssue {
	suite, err := d.db.GetPromptTestSuite(ctx, prompt.ID)
	if err != nil || suite == nil {
		return nil
	}

	if len(suite.Tests) < 3 {
		// Not enough tests to determine pattern
		return nil
	}

	if suite.PassRate < 50 {
		return &DecayIssue{
			Prompt:   prompt,
			Type:     DecayLowSuccess,
			Severity: "critical",
			Description: fmt.Sprintf("Low test success rate: %.1f%%", suite.PassRate),
			Suggestion: "Review test failures and update prompt",
			Details: map[string]interface{}{
				"pass_rate":   suite.PassRate,
				"avg_score":   suite.AvgScore,
				"total_tests": len(suite.Tests),
			},
			DetectedAt: time.Now().UTC(),
		}
	}

	return nil
}

// GetRecommendations returns prioritized recommendations based on audit
func (r *AuditResult) GetRecommendations() []string {
	var recommendations []string

	if r.Summary.DeprecatedCount > 0 {
		recommendations = append(recommendations,
			fmt.Sprintf("🔴 Critical: Update %d prompts using deprecated models", r.Summary.DeprecatedCount))
	}

	if r.Summary.LowSuccessCount > 0 {
		recommendations = append(recommendations,
			fmt.Sprintf("🔴 Critical: Fix %d prompts with low test success rates", r.Summary.LowSuccessCount))
	}

	if r.Summary.UnusedCount > 0 {
		recommendations = append(recommendations,
			fmt.Sprintf("🟡 Warning: Review %d unused prompts", r.Summary.UnusedCount))
	}

	if r.Summary.OldVersionCount > 0 {
		recommendations = append(recommendations,
			fmt.Sprintf("🟢 Info: Update %d outdated prompts", r.Summary.OldVersionCount))
	}

	if len(recommendations) == 0 {
		recommendations = append(recommendations, "✅ All prompts are healthy!")
	}

	return recommendations
}

// GetIssuesBySeverity filters issues by severity
func (r *AuditResult) GetIssuesBySeverity(severity string) []*DecayIssue {
	var filtered []*DecayIssue
	for _, issue := range r.Issues {
		if issue.Severity == severity {
			filtered = append(filtered, issue)
		}
	}
	return filtered
}
