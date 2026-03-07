package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/Bharath-code/promptvault/internal/model"
)

// Assistant provides AI-assisted prompt authoring
type Assistant struct {
	claudeAPIKey string
	openaiAPIKey string
	httpClient   *httpClient
}

// NewAssistant creates a new AI assistant
func NewAssistant() *Assistant {
	return &Assistant{
		claudeAPIKey: os.Getenv("ANTHROPIC_API_KEY"),
		openaiAPIKey: os.Getenv("OPENAI_API_KEY"),
		httpClient:   &httpClient{timeout: 60 * time.Second},
	}
}

// PromptAnalysis is the result of analyzing a prompt
type PromptAnalysis struct {
	Variables       []string `json:"variables"`
	SuggestedTags   []string `json:"suggested_tags"`
	SuggestedStack  string   `json:"suggested_stack"`
	AntiPatterns    []string `json:"anti_patterns"`
	Improvements    []string `json:"improvements"`
	QualityScore    int      `json:"quality_score"` // 0-100
	EstimatedTokens int      `json:"estimated_tokens"`
}

// AnalyzePrompt analyzes a prompt and provides recommendations
func (a *Assistant) AnalyzePrompt(ctx context.Context, title, content string) (*PromptAnalysis, error) {
	analysis := &PromptAnalysis{
		Variables:     detectVariables(content),
		AntiPatterns:  detectAntiPatterns(content),
		Improvements:  suggestImprovements(content),
		QualityScore:  calculateQualityScore(content),
	}

	// Get AI-powered recommendations
	if a.claudeAPIKey != "" {
		aiAnalysis, err := a.getAIAnalysis(ctx, title, content)
		if err == nil {
			analysis.SuggestedTags = aiAnalysis.SuggestedTags
			analysis.SuggestedStack = aiAnalysis.SuggestedStack
			analysis.Improvements = append(analysis.Improvements, aiAnalysis.Improvements...)
		}
	}

	// Fallback to rule-based if no API key
	if analysis.SuggestedTags == nil {
		analysis.SuggestedTags = suggestTags(content)
	}
	if analysis.SuggestedStack == "" {
		analysis.SuggestedStack = suggestStack(content)
	}

	analysis.EstimatedTokens = estimateTokens(content)

	return analysis, nil
}

// detectVariables finds {{variable}} patterns in content
func detectVariables(content string) []string {
	re := regexp.MustCompile(`\{\{\s*([^}]+)\s*\}\}`)
	matches := re.FindAllStringSubmatch(content, -1)

	variables := make([]string, 0, len(matches))
	seen := make(map[string]bool)

	for _, match := range matches {
		varName := strings.TrimSpace(match[1])
		if !seen[varName] {
			variables = append(variables, varName)
			seen[varName] = true
		}
	}

	return variables
}

// detectAntiPatterns finds common prompt anti-patterns
func detectAntiPatterns(content string) []string {
	var patterns []string

	// Too vague
	if len(strings.Fields(content)) < 20 {
		patterns = append(patterns, "⚠️ Prompt is too short - add more context and specifics")
	}

// No examples
	if !strings.Contains(strings.ToLower(content), "example") &&
		!strings.Contains(strings.ToLower(content), "for instance") {
		patterns = append(patterns, "💡 Consider adding examples for clarity")
	}

	// No output format
	if !strings.Contains(strings.ToLower(content), "output") &&
		!strings.Contains(strings.ToLower(content), "format") &&
		!strings.Contains(strings.ToLower(content), "return") {
		patterns = append(patterns, "📝 Specify the expected output format")
	}

	// All caps (shouting)
	if strings.ToUpper(content) == content && len(content) > 50 {
		patterns = append(patterns, "🔊 Avoid ALL CAPS - use normal casing")
	}

	// Too many nested instructions
	nestedCount := strings.Count(content, "if") + strings.Count(content, "then") + strings.Count(content, "else")
	if nestedCount > 5 {
		patterns = append(patterns, "🌲 Too many nested conditions - simplify the logic")
	}

	// No constraints
	if !strings.Contains(strings.ToLower(content), "don't") &&
		!strings.Contains(strings.ToLower(content), "avoid") &&
		!strings.Contains(strings.ToLower(content), "never") {
		patterns = append(patterns, "🚫 Add constraints to prevent unwanted behavior")
	}

	return patterns
}

// suggestImprovements provides rule-based improvements
func suggestImprovements(content string) []string {
	var improvements []string

	// Check for step-by-step structure
	if !strings.Contains(content, "1.") && !strings.Contains(content, "2.") &&
		!strings.Contains(content, "- ") && !strings.Contains(content, "* ") {
		improvements = append(improvements, "Consider breaking into numbered steps for clarity")
	}

	// Check for role definition
	if !strings.Contains(strings.ToLower(content), "you are") &&
		!strings.Contains(strings.ToLower(content), "act as") &&
		!strings.Contains(strings.ToLower(content), "your role") {
		improvements = append(improvements, "Add a role definition (e.g., 'You are an expert React developer')")
	}

	// Check for context
	if len(strings.Fields(content)) < 100 && !strings.Contains(strings.ToLower(content), "context") {
		improvements = append(improvements, "Provide more context about the use case")
	}

	return improvements
}

// calculateQualityScore computes a simple quality score
func calculateQualityScore(content string) int {
	score := 50 // Base score

	// Length bonus (good prompts are detailed)
	wordCount := len(strings.Fields(content))
	if wordCount > 50 {
		score += 10
	}
	if wordCount > 100 {
		score += 10
	}
	if wordCount > 200 {
		score += 5
	}

	// Structure bonus
	if strings.Contains(content, "1.") || strings.Contains(content, "- ") {
		score += 10
	}

	// Examples bonus
	if strings.Contains(strings.ToLower(content), "example") {
		score += 10
	}

	// Variables bonus (shows reusability)
	if strings.Contains(content, "{{") {
		score += 5
	}

	// Constraints bonus
	if strings.Contains(strings.ToLower(content), "don't") ||
		strings.Contains(strings.ToLower(content), "avoid") {
		score += 5
	}

	// Cap at 100
	if score > 100 {
		score = 100
	}

	return score
}

// suggestTags suggests tags based on content
func suggestTags(content string) []string {
	contentLower := strings.ToLower(content)

	// Common patterns
	tagKeywords := map[string][]string{
		"react":        {"react", "frontend", "hooks"},
		"typescript":   {"typescript", "types", "javascript"},
		"python":       {"python", "backend", "scripting"},
		"docker":       {"docker", "containers", "devops"},
		"test":         {"testing", "unit-test", "quality"},
		"debug":        {"debugging", "troubleshooting", "fix"},
		"optimize":     {"performance", "optimization", "speed"},
		"security":     {"security", "authentication", "authorization"},
		"api":          {"api", "rest", "graphql"},
		"database":     {"database", "sql", "query"},
	}

	for keyword, tagList := range tagKeywords {
		if strings.Contains(contentLower, keyword) {
			return tagList[:min(3, len(tagList))]
		}
	}

	return []string{"general"}
}

// suggestStack suggests a tech stack based on content
func suggestStack(content string) string {
	contentLower := strings.ToLower(content)

	stackPatterns := map[string]string{
		"react":         "frontend/react",
		"vue":           "frontend/vue",
		"angular":       "frontend/angular",
		"typescript":    "frontend/typescript",
		"nextjs":        "frontend/react/nextjs",
		"node":          "backend/node",
		"express":       "backend/node/express",
		"python":        "backend/python",
		"django":        "backend/python/django",
		"fastapi":       "backend/python/fastapi",
		"go":            "backend/go",
		"docker":        "devops/docker",
		"kubernetes":    "devops/kubernetes",
		"terraform":     "devops/terraform",
		"aws":           "devops/aws",
		"postgresql":    "database/postgresql",
		"mongodb":       "database/mongodb",
		"prisma":        "database/prisma",
	}

	for keyword, stack := range stackPatterns {
		if strings.Contains(contentLower, keyword) {
			return stack
		}
	}

	return ""
}

// estimateTokens estimates token count
func estimateTokens(content string) int {
	// Rough estimate: 1 token ≈ 4 characters
	return len(content) / 4
}

// getAIAnalysis gets AI-powered analysis from Claude
func (a *Assistant) getAIAnalysis(ctx context.Context, title, content string) (*PromptAnalysis, error) {
	prompt := fmt.Sprintf(`Analyze this AI prompt and provide recommendations:

Title: %s

Content:
%s

Provide your analysis in JSON format:
{
  "suggested_tags": ["tag1", "tag2"],
  "suggested_stack": "frontend/react",
  "improvements": ["suggestion 1", "suggestion 2"]
}

Be concise and practical.`, title, content)

	response, err := a.callClaude(ctx, prompt)
	if err != nil {
		return nil, err
	}

	var analysis PromptAnalysis
	if err := json.Unmarshal([]byte(response), &analysis); err != nil {
		return nil, fmt.Errorf("parsing AI response: %w", err)
	}

	return &analysis, nil
}

// InteractiveCreate runs an interactive prompt creation session
func (a *Assistant) InteractiveCreate(ctx context.Context) (*model.Prompt, error) {
	fmt.Println()
	fmt.Println("🤖 AI-Assisted Prompt Creation")
	fmt.Println(strings.Repeat("─", 60))
	fmt.Println()

	// Get title
	fmt.Print("📝 Prompt Title: ")
	var title string
	fmt.Scanln(&title)

	if title == "" {
		return nil, fmt.Errorf("title is required")
	}

	// Get content
	fmt.Println()
	fmt.Println("📄 Enter prompt content (type 'DONE' on a new line to finish):")
	fmt.Println(strings.Repeat("─", 60))

	var contentBuilder strings.Builder
	scanner := newScanner()
	for {
		line, err := scanner.ReadString('\n')
		if err != nil {
			return nil, err
		}

		line = strings.TrimSpace(line)
		if strings.ToUpper(line) == "DONE" {
			break
		}

		contentBuilder.WriteString(line + "\n")
	}

	content := strings.TrimSpace(contentBuilder.String())
	if content == "" {
		return nil, fmt.Errorf("content is required")
	}

	// Analyze with AI
	fmt.Println()
	fmt.Println("🔍 Analyzing prompt...")

	analysis, err := a.AnalyzePrompt(ctx, title, content)
	if err != nil {
		fmt.Printf("⚠️  AI analysis failed: %v\n", err)
		fmt.Println("Continuing with basic analysis...")
	}

	// Show analysis
	if analysis != nil {
		fmt.Println()
		fmt.Println("📊 Analysis Results:")
		fmt.Println(strings.Repeat("─", 60))

		if len(analysis.Variables) > 0 {
			fmt.Printf("🏷️  Variables detected: %s\n", strings.Join(analysis.Variables, ", "))
		}

		if analysis.SuggestedStack != "" {
			fmt.Printf("📚 Suggested stack: %s\n", analysis.SuggestedStack)
		}

		if len(analysis.SuggestedTags) > 0 {
			fmt.Printf("🏷️  Suggested tags: %s\n", strings.Join(analysis.SuggestedTags, ", "))
		}

		fmt.Printf("⭐ Quality score: %d/100\n", analysis.QualityScore)

		if len(analysis.AntiPatterns) > 0 {
			fmt.Println()
			fmt.Println("⚠️  Anti-patterns found:")
			for _, ap := range analysis.AntiPatterns {
				fmt.Printf("   %s\n", ap)
			}
		}

		if len(analysis.Improvements) > 0 {
			fmt.Println()
			fmt.Println("💡 Suggested improvements:")
			for _, imp := range analysis.Improvements {
				fmt.Printf("   • %s\n", imp)
			}
		}

		fmt.Println()
		fmt.Print("Continue with these suggestions? [Y/n]: ")
		var confirm string
		fmt.Scanln(&confirm)

		if strings.ToLower(confirm) == "n" {
			fmt.Println("Continuing with original content...")
		} else {
			// Apply improvements (simplified - in real version would use AI to rewrite)
			if analysis.SuggestedStack != "" {
				fmt.Printf("Using suggested stack: %s\n", analysis.SuggestedStack)
			}
		}
	}

	// Get stack
	fmt.Println()
	fmt.Print("📚 Tech stack (e.g., frontend/react/hooks): ")
	var stack string
	fmt.Scanln(&stack)

	// Get tags
	fmt.Print("🏷️  Tags (comma-separated): ")
	var tagsStr string
	fmt.Scanln(&tagsStr)

	// Get models
	fmt.Print("🤖 Models (comma-separated, e.g., claude-sonnet,gpt-4o): ")
	var modelsStr string
	fmt.Scanln(&modelsStr)

	// Parse tags and models
	tags := parseCommaList(tagsStr)
	models := parseCommaList(modelsStr)

	// Use suggestions if empty
	if analysis != nil {
		if stack == "" && analysis.SuggestedStack != "" {
			stack = analysis.SuggestedStack
		}
		if len(tags) == 0 && len(analysis.SuggestedTags) > 0 {
			tags = analysis.SuggestedTags
		}
	}

	return &model.Prompt{
		Title:   title,
		Content: content,
		Stack:   stack,
		Tags:    tags,
		Models:  models,
	}, nil
}

// Helper functions
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func parseCommaList(s string) []string {
	if s == "" {
		return nil
	}

	parts := strings.Split(s, ",")
	result := make([]string, 0, len(parts))

	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p != "" {
			result = append(result, p)
		}
	}

	return result
}

// Simple scanner for reading input
type scanner struct {
	buf []byte
}

func newScanner() *scanner {
	return &scanner{buf: make([]byte, 0, 1024)}
}

func (s *scanner) ReadString(delim byte) (string, error) {
	var line []byte
	for {
		b := make([]byte, 1)
		_, err := os.Stdin.Read(b)
		if err != nil {
			return "", err
		}
		line = append(line, b[0])
		if b[0] == delim {
			break
		}
	}
	return string(line), nil
}

// HTTP client for API calls
type httpClient struct {
	timeout time.Duration
}

func (a *Assistant) callClaude(ctx context.Context, prompt string) (string, error) {
	// Simplified - in production would use proper HTTP client
	return "", fmt.Errorf("AI analysis not yet implemented - using rule-based analysis")
}
