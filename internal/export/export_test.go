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

func TestExportBulk(t *testing.T) {
	prompts := []*model.Prompt{
		{Title: "React Hooks", Content: "Use useState wisely", Stack: "frontend/react"},
		{Title: "Go API", Content: "Return errors properly", Stack: "backend/go"},
	}

	exporter := New(prompts)
	files, err := exporter.ExportBulk()
	if err != nil {
		t.Fatalf("ExportBulk failed: %v", err)
	}

	if len(files) != 2 {
		t.Errorf("expected 2 files, got %d", len(files))
	}

	// Check filenames are sanitized
	for _, f := range files {
		if strings.ContainsAny(f.Filename, `/\:*?"<>|`) {
			t.Errorf("filename %q contains invalid characters", f.Filename)
		}
		if !strings.HasSuffix(f.Filename, ".md") {
			t.Errorf("filename %q should end with .md", f.Filename)
		}
	}

	// Check content structure
	for _, f := range files {
		if !strings.Contains(f.Content, "# ") {
			t.Errorf("expected content to start with header for %s", f.Filename)
		}
	}
}

func TestImportFromJSON(t *testing.T) {
	importer := NewImporter()

	tests := []struct {
		name     string
		input    string
		wantLen  int
		wantErrs int
	}{
		{
			name:    "valid single object",
			input:   `{"title": "Test", "content": "Hello world"}`,
			wantLen: 1,
		},
		{
			name:    "valid array",
			input:   `[{"title": "A", "content": "Content A"}, {"title": "B", "content": "Content B"}]`,
			wantLen: 2,
		},
		{
			name:     "invalid JSON",
			input:    `{invalid}`,
			wantLen:  0,
			wantErrs: 1,
		},
		{
			name:    "empty array",
			input:   `[]`,
			wantLen: 0,
		},
		{
			name:    "missing title",
			input:   `{"content": "No title"}`,
			wantLen: 0,
		},
		{
			name:    "missing content",
			input:   `{"title": "No content"}`,
			wantLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := importer.ImportFromJSON(tt.input)
			if len(result.Prompts) != tt.wantLen {
				t.Errorf("expected %d prompts, got %d", tt.wantLen, len(result.Prompts))
			}
			if len(result.Errors) != tt.wantErrs {
				t.Errorf("expected %d errors, got %d", tt.wantErrs, len(result.Errors))
			}
		})
	}
}

func TestImportFromJSONWithMetadata(t *testing.T) {
	importer := NewImporter()

	input := `{
		"title": "React Component",
		"content": "Create a component",
		"stack": "frontend/react",
		"tags": ["react", "components"],
		"models": ["claude", "gpt-4"]
	}`

	result := importer.ImportFromJSON(input)

	if len(result.Prompts) != 1 {
		t.Fatalf("expected 1 prompt, got %d", len(result.Prompts))
	}

	p := result.Prompts[0]
	if p.Title != "React Component" {
		t.Errorf("expected title 'React Component', got %q", p.Title)
	}
	if p.Stack != "frontend/react" {
		t.Errorf("expected stack 'frontend/react', got %q", p.Stack)
	}
	if len(p.Tags) != 2 {
		t.Errorf("expected 2 tags, got %d", len(p.Tags))
	}
	if len(p.Models) != 2 {
		t.Errorf("expected 2 models, got %d", len(p.Models))
	}
}

func TestImportFromMarkdown(t *testing.T) {
	importer := NewImporter()

	tests := []struct {
		name    string
		input   string
		wantLen int
	}{
		{
			name: "single prompt",
			input: `## My Prompt

Some content here.
`,
			wantLen: 1,
		},
		{
			name: "multiple prompts",
			input: `## First Prompt

Content for first.

## Second Prompt

Content for second.
`,
			wantLen: 2,
		},
		{
			name:    "empty content",
			input:   `## Only Title`,
			wantLen: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := importer.ImportFromMarkdown(tt.input)
			if len(result.Prompts) != tt.wantLen {
				t.Errorf("expected %d prompts, got %d", tt.wantLen, len(result.Prompts))
			}
		})
	}
}

func TestImportFromMarkdownWithStack(t *testing.T) {
	importer := NewImporter()

	input := "## React Hooks\n\n**Stack:** `frontend/react`\n\nThis is the hook content.\n"

	result := importer.ImportFromMarkdown(input)

	if len(result.Prompts) == 0 {
		t.Fatal("expected at least 1 prompt")
	}

	p := result.Prompts[0]
	if p.Stack != "frontend/react" {
		t.Errorf("expected stack 'frontend/react', got %q", p.Stack)
	}
}

func TestSanitizeFilename(t *testing.T) {
	tests := []struct {
		input string
		check func(string) bool // custom check function
	}{
		{"Simple Title", func(s string) bool { return s == "Simple_Title" }},
		{"Title with/slash", func(s string) bool { return s == "Title_with-slash" }},
		{"Title with:colon", func(s string) bool { return s == "Title_with-colon" }},
		{"Title with*asterisk", func(s string) bool { return s == "Title_withasterisk" }},
		{"Title with?question", func(s string) bool { return s == "Title_withquestion" }},
		{"Title with\"quotes\"", func(s string) bool { return s == "Title_withquotes" }},
		{"Title with<angle>", func(s string) bool { return s == "Title_withangle" }},
		{"Title with|pipe", func(s string) bool { return s == "Title_with-pipe" }},
		{"  Leading and trailing  ", func(s string) bool { return s == "Leading_and_trailing" }},
		{"Very Long Title That Exceeds One Hundred Characters And Should Be Truncated Because We Have A Limit Of 100 Characters In Filenames To Keep Things Manageable",
			func(s string) bool { return len(s) <= 100 && s != "" }},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			got := sanitizeFilename(tt.input)
			if !tt.check(got) {
				t.Errorf("sanitizeFilename(%q) = %q, failed check", tt.input, got)
			}
		})
	}
}
