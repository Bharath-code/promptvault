package export

import (
	"strings"
	"testing"

	"github.com/Bharath-code/promptvault/internal/model"
)

func TestExportFormats(t *testing.T) {
	prompts := []*model.Prompt{
		{
			Title:   "Test Prompt 1",
			Content: "Content 1",
			Stack:   "frontend/react",
			Tags:    []string{"ui", "react"},
		},
		{
			Title:   "Test Prompt 2",
			Content: "Content 2",
			Stack:   "backend/go",
		},
	}

	exporter := New(prompts)

	tests := []struct {
		format   Format
		contains string
	}{
		{FormatSkillMD, "Test Prompt 1"},
		{FormatAgentsMD, "Test Prompt 2"},
		{FormatClaudeMD, "CLAUDE.md"},
		{FormatCursorRules, ".cursorrules"},
		{FormatWindsurf, "Windsurf"},
		{FormatJSON, `"title": "Test Prompt 1"`},
		{FormatPlainText, "=== Test Prompt 1 ==="},
		{FormatMarkdown, "PromptVault Export"},
	}

	for _, tt := range tests {
		t.Run(string(tt.format), func(t *testing.T) {
			res, err := exporter.Export(tt.format)
			if err != nil {
				t.Fatalf("unexpected error exporting %s: %v", tt.format, err)
			}
			if !strings.Contains(res, tt.contains) {
				t.Errorf("expected exported %s to contain %q\nGot: %s", tt.format, tt.contains, res)
			}
		})
	}
}
