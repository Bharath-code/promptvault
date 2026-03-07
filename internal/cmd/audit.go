package cmd

import (
	"context"
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/spf13/cobra"
	"github.com/Bharath-code/promptvault/internal/decay"
)

// audit command
var auditCmd = &cobra.Command{
	Use:   "audit",
	Short: "Audit prompts for decay and quality issues",
	Long: `Audit your prompt vault for decay and quality issues.

This command checks for:
- Prompts not used in 90+ days
- Deprecated model usage
- Low test success rates (< 50%)
- Prompts not updated in 180+ days

Examples:
  # Full audit
  promptvault audit

  # Show only critical issues
  promptvault audit --severity critical

  # JSON output for scripting
  promptvault audit --json
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		jsonOutput, _ := cmd.Flags().GetBool("json")
		severity, _ := cmd.Flags().GetString("severity")

		// Run audit
		detector := decay.NewDetector(database)
		result, err := detector.Audit(ctx)
		if err != nil {
			printError("Audit failed: %v", err)
			return err
		}

		// Output results
		if jsonOutput {
			return outputJSON(result)
		}

		displayAudit(result, severity)
		return nil
	},
}

func init() {
	auditCmd.Flags().Bool("json", false, "Output as JSON")
	auditCmd.Flags().String("severity", "", "Filter by severity (critical, warning, info)")
}

// displayAudit shows audit results in terminal
func displayAudit(result *decay.AuditResult, severityFilter string) {
	fmt.Println()
	fmt.Printf("🔍 PromptVault Audit Report\n")
	fmt.Println(strings.Repeat("─", 70))
	fmt.Printf("Generated: %s\n", result.GeneratedAt.Format("2006-01-02 15:04:05"))
	fmt.Println()

	// Summary
	fmt.Printf("📊 Summary:\n")
	fmt.Printf("   Total Prompts:   %d\n", result.TotalPrompts)
	fmt.Printf("   Healthy:         %d (%.1f%%)\n", result.HealthyPrompts, 
		float64(result.HealthyPrompts)/float64(result.TotalPrompts)*100)
	fmt.Printf("   Issues Found:    %d\n", result.IssuesFound)
	fmt.Println()

	// Recommendations
	fmt.Printf("💡 Recommendations:\n")
	for _, rec := range result.GetRecommendations() {
		fmt.Printf("   %s\n", rec)
	}
	fmt.Println()

	// Filter by severity if requested
	issues := result.Issues
	if severityFilter != "" {
		issues = result.GetIssuesBySeverity(severityFilter)
	}

	if len(issues) == 0 {
		if severityFilter != "" {
			fmt.Printf("No %s severity issues found\n", severityFilter)
		}
		return
	}

	// Sort by severity
	sortIssuesBySeverity(issues)

	// Display issues
	fmt.Printf("📋 Issues Found (%d):\n", len(issues))
	fmt.Println(strings.Repeat("─", 70))

	for i, issue := range issues {
		if i >= 10 {
			fmt.Printf("   ... and %d more issues\n", len(issues)-10)
			break
		}

		icon := getSeverityIcon(issue.Severity)
		fmt.Printf("\n%s [%s] %s\n", icon, issue.Type, issue.Prompt.Title)
		fmt.Printf("   %s\n", issue.Description)
		fmt.Printf("   💡 %s\n", issue.Suggestion)
		
		if details, ok := issue.Details["days_since_used"]; ok {
			fmt.Printf("   📅 Last used: %d days ago\n", details)
		}
		if details, ok := issue.Details["deprecated_model"]; ok {
			fmt.Printf("   🤖 Deprecated model: %s\n", details)
		}
		if details, ok := issue.Details["pass_rate"]; ok {
			fmt.Printf("   📊 Pass rate: %.1f%%\n", details)
		}
	}

	fmt.Println()
	fmt.Printf("Showing %d of %d issues\n", min(10, len(issues)), len(issues))
	fmt.Println()
	fmt.Printf("Run 'promptvault audit --severity critical' to see only critical issues\n")
}

// Helper functions
func getSeverityIcon(severity string) string {
	switch severity {
	case "critical":
		return "🔴"
	case "warning":
		return "🟡"
	case "info":
		return "🟢"
	default:
		return "⚪"
	}
}

func sortIssuesBySeverity(issues []*decay.DecayIssue) {
	severityOrder := map[string]int{"critical": 0, "warning": 1, "info": 2}
	sort.Slice(issues, func(i, j int) bool {
		return severityOrder[issues[i].Severity] < severityOrder[issues[j].Severity]
	})
}

func outputJSON(result *decay.AuditResult) error {
	// Simplified JSON output
	fmt.Printf(`{
  "total_prompts": %d,
  "healthy_prompts": %d,
  "issues_found": %d,
  "generated_at": "%s",
  "recommendations": [
`, result.TotalPrompts, result.HealthyPrompts, result.IssuesFound, 
		result.GeneratedAt.Format(time.RFC3339))

	recs := result.GetRecommendations()
	for i, rec := range recs {
		comma := ","
		if i == len(recs)-1 {
			comma = ""
		}
		fmt.Printf("    %q%s\n", rec, comma)
	}

	fmt.Printf("  ]\n}")
	return nil
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
