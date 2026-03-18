package tui

import (
	"regexp"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	highlightColor = lipgloss.Color("#A78BFA") // Purple - keywords
	stringColor    = lipgloss.Color("#34D399") // Green - strings
	commentColor   = lipgloss.Color("#64748B") // Gray - comments
	numberColor    = lipgloss.Color("#F59E0B") // Amber - numbers
	typeColor      = lipgloss.Color("#06B6D4") // Cyan - types
	functionColor  = lipgloss.Color("#3B82F6") // Blue - functions
	operatorColor  = lipgloss.Color("#EF4444") // Red - operators
	variableColor  = lipgloss.Color("#F97316") // Orange - variables
)

type SyntaxRule struct {
	Pattern *regexp.Regexp
	Color   lipgloss.Color
	Style   lipgloss.Style
}

var syntaxRules []SyntaxRule

func init() {
	syntaxRules = []SyntaxRule{
		// Comments (must be first to catch line comments before other patterns)
		{
			Pattern: regexp.MustCompile(`(?m)^(\s*)(//.*|#.*|--.*|/\*.*?\*/)`),
			Color:   commentColor,
		},
		// Strings (double and single quotes)
		{
			Pattern: regexp.MustCompile(`"(?:[^"\\]|\\.)*"|'(?:[^'\\]|\\.)*'|` + "`" + `(?:[^` + "`" + `\\]|\\.)*` + "`"),
			Color:   stringColor,
		},
		// Keywords
		{
			Pattern: regexp.MustCompile(`\b(func|const|var|type|struct|interface|package|import|return|if|else|for|range|switch|case|default|break|continue|defer|go|select|chan|map|make|new|nil|true|false|this|self|class|def|async|await|try|catch|finally|throw|const|let|var|function|export|default|from|as|of|in|on|with|using|namespace|public|private|protected|static|readonly|enum|implements|extends|abstract|virtual|override|void|null|undefined|let|yield|match|where|when|impl|trait|mod|pub|crate|use|mut|ref|move|loop|while|given|and|or|not|is|then|end|do|begin|rescue|ensure|raise|elsif|unless|until|module|alias|yield|and|or|not|in|end|def|class|module|if|unless|case|when|while|until|for|begin|do|retry|rescue|ensure|raise|return|break|next|redo|unless|super|nil|self|true|false|and|or|not|is|a|an)\b`),
			Color:   highlightColor,
		},
		// Types
		{
			Pattern: regexp.MustCompile(`\b(string|int|int8|int16|int32|int64|uint|uint8|uint16|uint32|uint64|float32|float64|bool|byte|rune|error|any|void|number|boolean|string|symbol|array|object|dict|list|tuple|set|map|slice|chan|struct|interface|enum|union|optional|nullable|int64|int32|uint64|uint32|float64|float32|complex64|complex128|any|unknown|never|object|Array|String|Number|Boolean|Object|Function|Symbol|Map|Set|WeakMap|WeakSet|Promise|Record|Partial|Required|Readonly|Pick|Omit|Exclude|Extract|NonNullable|ReturnType|Parameters)\b`),
			Color:   typeColor,
		},
		// Numbers
		{
			Pattern: regexp.MustCompile(`\b(\d+\.?\d*|\.\d+)\b`),
			Color:   numberColor,
		},
		// Function calls
		{
			Pattern: regexp.MustCompile(`\b([a-zA-Z_][a-zA-Z0-9_]*)\s*\(`),
			Color:   functionColor,
		},
		// Variables in {{}} template syntax
		{
			Pattern: regexp.MustCompile(`\{\{([^}]+)\}\}`),
			Color:   variableColor,
		},
	}
}

// HighlightCode applies syntax highlighting to code content
func HighlightCode(code string) string {
	result := code

	for _, rule := range syntaxRules {
		result = rule.Pattern.ReplaceAllStringFunc(result, func(match string) string {
			return rule.Style.Foreground(rule.Color).Render(match)
		})
	}

	return result
}

// HighlightMarkdownContent applies highlighting to markdown code blocks
func HighlightMarkdownContent(content string) string {
	lines := strings.Split(content, "\n")
	var result []string

	inCodeBlock := false

	for _, line := range lines {
		// Code block start/end
		if strings.HasPrefix(line, "```") {
			if inCodeBlock {
				// End code block
				result = append(result, lipgloss.NewStyle().
					Foreground(lipgloss.Color("#475569")).
					Render(line))
				inCodeBlock = false
			} else {
				// Start code block
				result = append(result, lipgloss.NewStyle().
					Foreground(lipgloss.Color("#475569")).
					Render(line))
				inCodeBlock = true
			}
			continue
		}

		if inCodeBlock {
			// Highlight code within block
			result = append(result, HighlightCode(line))
		} else {
			// Regular markdown content
			result = append(result, highlightMarkdownLine(line))
		}
	}

	return strings.Join(result, "\n")
}

func highlightMarkdownLine(line string) string {
	result := line

	// Bold text **text**
	result = regexp.MustCompile(`\*\*(.+?)\*\*`).
		ReplaceAllStringFunc(result, func(match string) string {
			return lipgloss.NewStyle().Bold(true).Render(match)
		})

	// Italic text *text*
	result = regexp.MustCompile(`(?<!\*)\*(?!\*)(.+?)(?<!\*)\*(?!\*)`).
		ReplaceAllStringFunc(result, func(match string) string {
			return lipgloss.NewStyle().Italic(true).Render(match)
		})

	// Inline code `code`
	result = regexp.MustCompile("`([^`]+)`").
		ReplaceAllStringFunc(result, func(match string) string {
			return lipgloss.NewStyle().
				Background(lipgloss.Color("#1E293B")).
				Foreground(lipgloss.Color("#A78BFA")).
				Render(match)
		})

	// Headers # ## ###
	if strings.HasPrefix(line, "#") {
		level := 0
		for _, ch := range line {
			if ch == '#' {
				level++
			} else {
				break
			}
		}
		colors := []string{"#7C3AED", "#8B5CF6", "#A78BFA", "#C4B5FD"}
		colorIdx := level - 1
		if colorIdx >= len(colors) {
			colorIdx = len(colors) - 1
		}
		return lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color(colors[colorIdx])).
			Render(line)
	}

	// List items - and numbered lists
	if matched, _ := regexp.MatchString(`^(\s*)([-*+]|\d+\.)\s`, line); matched {
		return lipgloss.NewStyle().
			Foreground(lipgloss.Color("#06B6D4")).
			Render(line)
	}

	// Links [text](url)
	result = regexp.MustCompile(`\[([^\]]+)\]\([^\)]+\)`).
		ReplaceAllStringFunc(result, func(match string) string {
			return lipgloss.NewStyle().
				Underline(true).
				Foreground(lipgloss.Color("#3B82F6")).
				Render(match)
		})

	return result
}

// HighlightPromptContent highlights content for the preview pane
func HighlightPromptContent(content string, maxLines int) string {
	lines := strings.Split(content, "\n")

	// Truncate if too many lines
	if len(lines) > maxLines {
		lines = lines[:maxLines]
		content = strings.Join(lines, "\n") + "\n..."
	} else {
		content = strings.Join(lines, "\n")
	}

	// Check if content has code-like patterns
	hasCode := false
	for _, line := range lines {
		if strings.Contains(line, "```") || strings.Contains(line, "func ") ||
			strings.Contains(line, "const ") || strings.Contains(line, "class ") ||
			strings.Contains(line, "def ") || strings.Contains(line, "function ") ||
			strings.Contains(line, "import ") || strings.Contains(line, "export ") {
			hasCode = true
			break
		}
	}

	if hasCode {
		return HighlightMarkdownContent(content)
	}

	// Just do markdown highlighting for regular content
	return highlightMarkdownLine(content)
}
