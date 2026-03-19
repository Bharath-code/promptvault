package cmd

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/Bharath-code/promptvault/internal/ai"
	"github.com/Bharath-code/promptvault/internal/model"
	"github.com/spf13/cobra"
)

// create command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a new prompt with AI assistance",
	Long: `Create a new prompt with AI-assisted authoring.

This command provides an interactive experience with:
- Variable detection ({{variable}} syntax)
- Tag and stack recommendations
- Anti-pattern detection
- Quality scoring
- Improvement suggestions

Examples:
  # Interactive AI-assisted creation
  promptvault create

  # Quick creation without AI
  promptvault add "My prompt" --content "..." --stack frontend/react
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		aiFlag, _ := cmd.Flags().GetBool("ai")

		var prompt *model.Prompt
		var err error

		if aiFlag {
			// AI-assisted creation
			assistant := ai.NewAssistant()
			prompt, err = assistant.InteractiveCreate(ctx)
		} else {
			// Standard interactive creation
			prompt, err = interactiveCreate()
		}

		if err != nil {
			return err
		}

		// Save the prompt
		if err := database.Add(ctx, prompt); err != nil {
			printError("Failed to save prompt: %v", err)
			return err
		}

		printSuccess("Created prompt: %s", prompt.Title)
		logInfo("ID: %s", prompt.ID)
		if prompt.Stack != "" {
			logInfo("Stack: %s", prompt.Stack)
		}
		if len(prompt.Tags) > 0 {
			logInfo("Tags: %s", strings.Join(prompt.Tags, ", "))
		}

		return nil
	},
}

func init() {
	createCmd.Flags().Bool("ai", false, "Use AI-assisted creation")
}

// interactiveCreate runs standard interactive prompt creation
func interactiveCreate() (*model.Prompt, error) {
	if !quiet {
		fmt.Println()
		fmt.Println("📝 Create New Prompt")
		fmt.Println(strings.Repeat("─", 60))
		fmt.Println()
	}

	// Get title
	if !quiet {
		fmt.Print("📝 Title: ")
	}
	var title string
	fmt.Scanln(&title)

	if strings.TrimSpace(title) == "" {
		return nil, fmt.Errorf("title is required")
	}

	// Get content
	if !quiet {
		fmt.Println()
		fmt.Println("📄 Content (type 'DONE' on a new line to finish):")
		fmt.Println(strings.Repeat("─", 60))
	}

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

	// Get stack
	if !quiet {
		fmt.Println()
		fmt.Print("📚 Tech stack (e.g., frontend/react/hooks): ")
	}
	var stack string
	fmt.Scanln(&stack)

	// Get tags
	if !quiet {
		fmt.Print("🏷️  Tags (comma-separated): ")
	}
	var tagsStr string
	fmt.Scanln(&tagsStr)

	// Get models
	if !quiet {
		fmt.Print("🤖 Models (comma-separated): ")
	}
	var modelsStr string
	fmt.Scanln(&modelsStr)

	// Parse
	tags := parseCommaList(tagsStr)
	models := parseCommaList(modelsStr)

	return &model.Prompt{
		Title:   title,
		Content: content,
		Stack:   stack,
		Tags:    tags,
		Models:  models,
	}, nil
}

// Simple scanner for reading multi-line input
type scanner struct{}

func newScanner() *scanner {
	return &scanner{}
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
