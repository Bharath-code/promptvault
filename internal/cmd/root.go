package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/atotto/clipboard"
	"github.com/spf13/cobra"
	"github.com/Bharath-code/promptvault/internal/db"
	"github.com/Bharath-code/promptvault/internal/export"
	"github.com/Bharath-code/promptvault/internal/mcp"
	"github.com/Bharath-code/promptvault/internal/model"
	gistsync "github.com/Bharath-code/promptvault/internal/sync"
	"github.com/Bharath-code/promptvault/internal/tui"
)

var database *db.DB

// Constants for clean code
const (
	stdinReadBuffer = 64 * 1024 // 64KB buffer for stdin reads
)

// DX: Color codes for richer output
var (
	colorSuccess = "\033[38;5;2m"    // Green
	colorError   = "\033[38;5;1m"    // Red
	colorWarning = "\033[38;5;3m"    // Yellow
	colorInfo    = "\033[38;5;6m"    // Cyan
	colorPrimary = "\033[38;5;129m"  // Purple
	colorMuted   = "\033[38;5;245m"  // Gray
	colorReset   = "\033[0m"
	
	// Icons
	iconSuccess = "✓"
	iconError   = "✗"
	iconWarning = "⚠"
	iconInfo    = "ℹ"
	iconSparkle = "⚡"
)

// DX: Verbose/Debug mode flags
var (
	verbose bool
	debug   bool
)

// DX: Logging helpers for verbose mode
func logInfo(format string, args ...interface{}) {
	if verbose {
		fmt.Fprintf(os.Stderr, colorInfo+iconInfo+"  "+format+colorReset+"\n", args...)
	}
}

func logDebug(format string, args ...interface{}) {
	if debug {
		timestamp := time.Now().Format("15:04:05.000")
		fmt.Fprintf(os.Stderr, "\033[38;5;245m🔍 [%s] %s\033[0m\n", timestamp, fmt.Sprintf(format, args...))
	}
}

func printSuccess(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, colorSuccess+iconSuccess+" "+format+colorReset+"\n", args...)
}

func printError(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, colorError+iconError+" "+format+colorReset+"\n", args...)
}

func printWarning(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, colorWarning+iconWarning+" "+format+colorReset+"\n", args...)
}

func printInfo(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, colorInfo+iconInfo+"  "+format+colorReset+"\n", args...)
}

// DX: Smart error suggestions
func suggestFix(err error) string {
	if err == nil {
		return ""
	}
	
	errStr := err.Error()
	var suggestions []string
	
	switch {
	case strings.Contains(errStr, "title is required"):
		suggestions = append(suggestions,
			"Pass title as argument: promptvault add \"My prompt title\"",
			"Or use --title flag: promptvault add --title \"My prompt title\"",
		)
	case strings.Contains(errStr, "content is required"):
		suggestions = append(suggestions,
			"Pipe content via stdin: cat prompt.txt | promptvault add \"Title\"",
			"Or use --content flag: promptvault add \"Title\" --content \"Your prompt here\"",
		)
	case strings.Contains(errStr, "GitHub token"):
		suggestions = append(suggestions,
			"Create a token at: https://github.com/settings/tokens",
			"Required scope: gist (write)",
			"Then run: promptvault sync push --token <your_token>",
		)
	case strings.Contains(errStr, "no Gist ID"):
		suggestions = append(suggestions,
			"Run 'promptvault sync push --token <token>' first to create a backup",
			"Or set PROMPTVAULT_GIST_ID environment variable",
		)
	case strings.Contains(errStr, "prompt not found"):
		suggestions = append(suggestions,
			"Run 'promptvault list' to see all prompts",
			"Or 'promptvault search <query>' to find by content",
		)
	case strings.Contains(errStr, "no prompts to export"):
		suggestions = append(suggestions,
			"Run 'promptvault init' to add starter prompts",
			"Or 'promptvault add' to create your first prompt",
		)
	}
	
	if len(suggestions) == 0 {
		return ""
	}
	
	var sb strings.Builder
	sb.WriteString("\n")
	sb.WriteString(colorInfo + "💡 Tips:" + colorReset + "\n")
	for i, s := range suggestions {
		if i > 2 { // Show max 3 suggestions
			break
		}
		sb.WriteString(fmt.Sprintf("   %s %s\n", colorMuted+"•"+colorReset, s))
	}
	
	return sb.String()
}

// wrapError adds suggestions to error messages
func wrapError(err error, message string) error {
	if err == nil {
		return nil
	}
	
	suggestion := suggestFix(err)
	if suggestion != "" {
		fmt.Fprint(os.Stderr, suggestion)
	}
	
	return fmt.Errorf("%s: %w", message, err)
}

// DX: Smart defaults detection
func detectStackFromPath() string {
	// Check for package.json (Node.js/React)
	if _, err := os.Stat("package.json"); err == nil {
		if _, err := os.Stat("src"); err == nil {
			// Check for React
			if _, err := os.Stat("src/components"); err == nil {
				return "frontend/react"
			}
			if _, err := os.Stat("src/pages"); err == nil {
				return "frontend/nextjs"
			}
			return "backend/node"
		}
		return "backend/node"
	}
	
	// Check for go.mod (Go)
	if _, err := os.Stat("go.mod"); err == nil {
		return "backend/go"
	}
	
	// Check for Python
	if _, err := os.Stat("requirements.txt"); err == nil {
		return "backend/python"
	}
	if _, err := os.Stat("setup.py"); err == nil {
		return "backend/python"
	}
	if _, err := os.Stat("pyproject.toml"); err == nil {
		return "backend/python"
	}
	
	// Check for Terraform
	if _, err := os.Stat("main.tf"); err == nil {
		return "devops/terraform"
	}
	
	// Check for Docker
	if _, err := os.Stat("Dockerfile"); err == nil {
		return "devops/docker"
	}
	
	// Check for Kubernetes
	if _, err := os.Stat("k8s"); err == nil {
		return "devops/kubernetes"
	}
	
	return ""
}

// getGitBranch returns the current Git branch name
func getGitBranch() string {
	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	output, err := cmd.Output()
	if err != nil {
		logDebug("Not a git repository or git not available")
		return ""
	}
	return strings.TrimSpace(string(output))
}

// getCurrentDirectory returns the current working directory name
func getCurrentDirectory() string {
	dir, err := os.Getwd()
	if err != nil {
		return ""
	}
	return filepath.Base(dir)
}

// centerText centers text within a given width
func centerText(text string, width int) string {
	padding := (width - len(text)) / 2
	if padding < 0 {
		padding = 0
	}
	return strings.Repeat(" ", padding) + text + strings.Repeat(" ", width-padding-len(text))
}

// spinner frames for loading indicator
var spinnerFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

// showLoading displays a loading spinner with a message
// Returns a done function that should be called when the operation is complete
func showLoading(msg string) func() {
	done := make(chan struct{})
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		i := 0
		for {
			select {
			case <-done:
				// Clear the loading line completely
				fmt.Printf("\r\033[2K\r")
				return
			default:
				fmt.Printf("\r%s %s", spinnerFrames[i%len(spinnerFrames)], msg)
				i++
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()

	return func() {
		close(done)
		wg.Wait()
	}
}

// Root command
var rootCmd = &cobra.Command{
	Use:   "promptvault",
	Short: "⚡ PromptVault — The universal prompt OS for developers",
	Long: `
⚡ PromptVault — Manage AI prompts by tech stack, right from your terminal.

Store, search, version, and deploy prompts across every AI tool.
Run without arguments to open the interactive TUI.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		logDebug("Opening database at %s", database.Path())
		return tui.Run(database)
	},
}

func init() {
	// DX: Add verbose and debug flags
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")
	rootCmd.PersistentFlags().BoolVarP(&debug, "debug", "d", false, "Enable debug output")
}

// add command
var addCmd = &cobra.Command{
	Use:   "add [title]",
	Short: "Add a new prompt",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		logDebug("Executing add command")
		
		title, _ := cmd.Flags().GetString("title")
		content, _ := cmd.Flags().GetString("content")
		stack, _ := cmd.Flags().GetString("stack")
		tagsStr, _ := cmd.Flags().GetString("tags")
		modelsStr, _ := cmd.Flags().GetString("models")
		verified, _ := cmd.Flags().GetBool("verified")
		preview, _ := cmd.Flags().GetBool("preview")

		if len(args) > 0 && title == "" {
			title = args[0]
		}

		if title == "" {
			printError("Title is required")
			fmt.Fprint(os.Stderr, suggestFix(fmt.Errorf("title is required")))
			os.Exit(1)
		}
		if content == "" {
			// Try reading from stdin
			stat, _ := os.Stdin.Stat()
			if (stat.Mode() & os.ModeCharDevice) == 0 {
				buf := make([]byte, stdinReadBuffer)
				n, _ := os.Stdin.Read(buf)
				content = string(buf[:n])
				logDebug("Read %d bytes from stdin", n)
			}
		}
		if content == "" {
			printError("Content is required")
			fmt.Fprint(os.Stderr, suggestFix(fmt.Errorf("content is required")))
			os.Exit(1)
		}

		// DX: Auto-detect stack from current directory if not specified
		if stack == "" {
			detectedStack := detectStackFromPath()
			if detectedStack != "" {
				stack = detectedStack
				logDebug("Auto-detected stack: %s", stack)
			}
		}

		// DX: Auto-add Git branch as tag
		var tags []string
		for _, t := range strings.Split(tagsStr, ",") {
			if t = strings.TrimSpace(t); t != "" {
				tags = append(tags, t)
			}
		}
		
		// Add git branch as tag if in a repo
		branch := getGitBranch()
		if branch != "" && branch != "HEAD" {
			tags = append(tags, "git:"+branch)
			logDebug("Added git branch tag: %s", branch)
		}

		var models []string
		for _, m := range strings.Split(modelsStr, ",") {
			if m = strings.TrimSpace(m); m != "" {
				models = append(models, m)
			}
		}

		p := &model.Prompt{
			Title:    title,
			Content:  content,
			Stack:    stack,
			Tags:     tags,
			Models:   models,
			Verified: verified,
		}

		// DX: Preview before adding
		if preview {
			fmt.Println()
			fmt.Println(colorPrimary + "┌" + strings.Repeat("─", 70) + "┐" + colorReset)
			fmt.Println(colorPrimary + "│" + colorReset + centerText("📋 PREVIEW", 70) + colorPrimary + "│" + colorReset)
			fmt.Println(colorPrimary + "├" + strings.Repeat("─", 70) + "┤" + colorReset)
			
			// Show content preview
			lines := strings.Split(content, "\n")
			maxLines := 12
			for i, line := range lines {
				if i >= maxLines {
					fmt.Println(colorPrimary + "│" + colorReset + fmt.Sprintf("  ... (%d more lines)", len(lines)-maxLines) + strings.Repeat(" ", 66-len(line)) + colorPrimary + "│" + colorReset)
					break
				}
				// Truncate long lines
				if len(line) > 68 {
					line = line[:65] + "..."
				}
				fmt.Println(colorPrimary + "│" + colorReset + "  " + line + strings.Repeat(" ", 68-len(line)) + colorPrimary + "│" + colorReset)
			}
			
			fmt.Println(colorPrimary + "└" + strings.Repeat("─", 70) + "┘" + colorReset)
			fmt.Println()
			
			// Show metadata
			fmt.Printf("%s Title:%s   %s\n", colorMuted, colorReset, title)
			if stack != "" {
				fmt.Printf("%s Stack:%s   %s\n", colorMuted, colorReset, stack)
			}
			if len(tags) > 0 {
				fmt.Printf("%s Tags:%s    %s\n", colorMuted, colorReset, strings.Join(tags, ", "))
			}
			if len(models) > 0 {
				fmt.Printf("%s Models:%s  %s\n", colorMuted, colorReset, strings.Join(models, ", "))
			}
			fmt.Println()
			
			// Confirm
			fmt.Printf("%s Add this prompt?%s [y/N]: ", colorInfo, colorReset)
			var confirm string
			fmt.Scanln(&confirm)
			if strings.ToLower(confirm) != "y" && confirm != "" {
				printInfo("Cancelled")
				return nil
			}
			fmt.Println()
		}

		logDebug("Adding prompt: %s (stack: %s, tags: %v)", p.Title, p.Stack, p.Tags)
		if err := database.Add(ctx, p); err != nil {
			printError("Failed to add prompt")
			fmt.Fprint(os.Stderr, suggestFix(err))
			return wrapError(err, "adding prompt")
		}

		printSuccess("Added prompt: %s", p.Title)
		logInfo("ID: %s", p.ID)
		if stack != "" {
			logInfo("Stack: %s", stack)
		}
		if len(tags) > 0 {
			logInfo("Tags: %s", strings.Join(tags, ", "))
		}
		
		return nil
	},
}

func init() {
	addCmd.Flags().StringP("title", "t", "", "Prompt title")
	addCmd.Flags().StringP("content", "c", "", "Prompt content (or pipe via stdin)")
	addCmd.Flags().StringP("stack", "s", "", "Tech stack path (e.g. frontend/react/hooks)")
	addCmd.Flags().String("tags", "", "Comma-separated tags")
	addCmd.Flags().String("models", "", "Comma-separated model names")
	addCmd.Flags().Bool("verified", false, "Mark as verified")
	// DX: Preview flag
	addCmd.Flags().Bool("preview", false, "Preview before adding")
}

// list command
var listCmd = &cobra.Command{
	Use:     "list",
	Short:   "List prompts",
	Aliases: []string{"ls", "show", "list-all"},
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		stack, _ := cmd.Flags().GetString("stack")
		short, _ := cmd.Flags().GetBool("short")
		jsonOutput, _ := cmd.Flags().GetBool("json")

		logDebug("Executing list command (stack: %s, json: %v)", stack, jsonOutput)

		// Show loading indicator
		stopLoading := showLoading("Loading prompts...")

		prompts, err := database.List(ctx, stack)
		stopLoading()

		if err != nil {
			printError("Failed to list prompts")
			return wrapError(err, "listing prompts")
		}

		if len(prompts) == 0 {
			if jsonOutput {
				fmt.Println("[]")
			} else {
				printWarning("No prompts found")
				fmt.Fprint(os.Stderr, suggestFix(fmt.Errorf("no prompts to export")))
			}
			return nil
		}

		// DX: JSON output for scripting
		if jsonOutput {
			encoder := json.NewEncoder(os.Stdout)
			encoder.SetIndent("", "  ")
			encoder.SetEscapeHTML(false)
			return encoder.Encode(prompts)
		}

		if short {
			for _, p := range prompts {
				fmt.Printf("%-36s  %s\n", p.ID[:8], p.Title)
			}
			return nil
		}

		for _, p := range prompts {
			verified := ""
			if p.Verified {
				verified = " ✓"
			}
			stackStr := ""
			if p.Stack != "" {
				stackStr = fmt.Sprintf(" [%s]", p.Stack)
			}
			fmt.Printf("%-8s  %s%s%s\n", p.ID[:8], p.Title, verified, stackStr)
		}

		fmt.Printf("\n%s %d prompt(s)\n", iconSuccess, len(prompts))
		return nil
	},
}

func init() {
	listCmd.Flags().Bool("short", false, "Short output (id + title only)")
	// DX: JSON output flag
	listCmd.Flags().Bool("json", false, "Output as JSON")
}

// get command — fetch a specific prompt and copy to clipboard
var getCmd = &cobra.Command{
	Use:     "get [id-or-title]",
	Short:   "Get a prompt by ID or title (copies to clipboard)",
	Aliases: []string{"fetch"},
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		query := args[0]
		copyFlag, _ := cmd.Flags().GetBool("copy")
		printFlag, _ := cmd.Flags().GetBool("print")

		logDebug("Executing get command (query: %s)", query)

		// Try by ID first, then search
		prompts, err := database.Search(ctx, query)
		if err != nil || len(prompts) == 0 {
			printError("No prompt found matching: %s", query)
			fmt.Fprint(os.Stderr, suggestFix(fmt.Errorf("prompt not found")))
			return fmt.Errorf("prompt not found: %s", query)
		}

		p := prompts[0]
		logDebug("Found prompt: %s (id: %s)", p.Title, p.ID)

		if printFlag || !copyFlag {
			fmt.Println(p.Content)
		}

		if copyFlag {
			if err := clipboard.WriteAll(p.Content); err != nil {
				printError("Failed to copy to clipboard")
				return wrapError(err, "copying to clipboard")
			}
			if incErr := database.IncrementUsage(ctx, p.ID); incErr != nil {
				logDebug("Failed to track usage: %v", incErr)
			}
			printSuccess("Copied '%s' to clipboard", p.Title)
		}

		return nil
	},
}

// search command
var searchCmd = &cobra.Command{
	Use:     "search [query]",
	Short:   "Full-text search across all prompts",
	Aliases: []string{"find", "query"},
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		jsonOutput, _ := cmd.Flags().GetBool("json")
		query := args[0]

		logDebug("Executing search command (query: %s)", query)

		// Show loading indicator
		stopLoading := showLoading("Searching prompts...")

		prompts, err := database.Search(ctx, query)
		stopLoading()

		if err != nil {
			printError("Search failed")
			return wrapError(err, "searching")
		}

		if len(prompts) == 0 {
			if jsonOutput {
				fmt.Println("[]")
			} else {
				printWarning("No prompts found for: %s", query)
			}
			return nil
		}

		// DX: JSON output for scripting
		if jsonOutput {
			encoder := json.NewEncoder(os.Stdout)
			encoder.SetIndent("", "  ")
			encoder.SetEscapeHTML(false)
			return encoder.Encode(prompts)
		}

		printSuccess("Found %d prompt(s):", len(prompts))
		fmt.Println()
		for _, p := range prompts {
			stack := ""
			if p.Stack != "" {
				stack = fmt.Sprintf(" [%s]", p.Stack)
			}
			fmt.Printf("  %-8s  %s%s\n", p.ID[:8], p.Title, stack)
		}
		return nil
	},
}

func init() {
	// DX: JSON output flag
	searchCmd.Flags().Bool("json", false, "Output as JSON")
}

// delete command
var deleteCmd = &cobra.Command{
	Use:     "delete [id]",
	Short:   "Delete a prompt",
	Aliases: []string{"rm", "remove", "del"},
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		force, _ := cmd.Flags().GetBool("force")

		logDebug("Executing delete command (id: %s, force: %v)", args[0], force)

		// Search to find the prompt
		prompts, err := database.Search(ctx, args[0])
		if err != nil || len(prompts) == 0 {
			// Try exact ID
			if err := database.Delete(ctx, args[0]); err != nil {
				printError("Prompt not found: %s", args[0])
				fmt.Fprint(os.Stderr, suggestFix(fmt.Errorf("prompt not found")))
				return fmt.Errorf("prompt not found: %s", args[0])
			}
			printSuccess("Deleted")
			return nil
		}

		p := prompts[0]

		if !force {
			fmt.Printf("%s Delete '%s'? [y/N] %s", colorWarning, p.Title, colorReset)
			var confirm string
			fmt.Scanln(&confirm)
			if strings.ToLower(confirm) != "y" {
				printInfo("Cancelled")
				return nil
			}
		}

		if err := database.Delete(ctx, p.ID); err != nil {
			return err
		}
		printSuccess("Deleted: %s", p.Title)
		return nil
	},
}

// export command
var exportCmd = &cobra.Command{
	Use:     "export",
	Short:   "Export prompts to various formats",
	Aliases: []string{"exp"},
	Long: `Export prompts to various formats for use with AI tools.

Available formats:
  skill.md      - Claude Code SKILL.md file
  agents.md     - OpenAI Agents format
  claude.md     - CLAUDE.md snippets
  cursorrules   - Cursor IDE rules (.cursorrules)
  windsurf      - Windsurf IDE rules (.windsurfrules)
  markdown      - Readable markdown documentation
  json          - JSON format for integrations
  text          - Plain text format

Examples:
  # Export all prompts to SKILL.md
  promptvault export --format skill.md --output SKILL.md

  # Export React prompts to Cursor rules
  promptvault export --format cursorrules --stack frontend/react

  # Export to AGENTS.md
  promptvault export --format agents.md > AGENTS.md
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		format, _ := cmd.Flags().GetString("format")
		stack, _ := cmd.Flags().GetString("stack")
		output, _ := cmd.Flags().GetString("output")

		logDebug("Executing export command (format: %s, stack: %s)", format, stack)

		// Show loading indicator
		stopLoading := showLoading("Exporting prompts...")

		prompts, err := database.List(ctx, stack)
		stopLoading()

		if err != nil {
			printError("Failed to list prompts")
			return wrapError(err, "listing prompts")
		}

		if len(prompts) == 0 {
			printWarning("No prompts to export")
			fmt.Fprint(os.Stderr, suggestFix(fmt.Errorf("no prompts to export")))
			return fmt.Errorf("no prompts to export")
		}

		e := export.New(prompts)
		result, err := e.Export(export.Format(format))
		if err != nil {
			printError("Export failed")
			return wrapError(err, "exporting")
		}

		if output == "" {
			switch format {
			case "cursorrules":
				output = ".cursorrules"
			case "windsurf":
				output = ".windsurfrules"
			case "skill.md":
				output = "SKILL.md"
			case "agents.md":
				output = "AGENTS.md"
			}
		}

		if output != "" {
			// Try to append if file exists, else create
			f, err := os.OpenFile(output, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
			if err != nil {
				printError("Failed to open file: %s", output)
				return wrapError(err, "opening file")
			}
			defer f.Close()

			if _, err := f.WriteString("\n" + result + "\n"); err != nil {
				printError("Failed to write to file")
				return wrapError(err, "writing file")
			}
			printSuccess("Exported %d prompts to %s", len(prompts), output)
			logInfo("Format: %s", format)
			logInfo("File: %s", output)
		} else {
			fmt.Print(result)
		}

		return nil
	},
}

// init command — seed vault with curated prompts
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize vault with curated starter prompts",
	Long: `Initialize your PromptVault with 15+ curated, production-grade prompts.

This command adds expertly crafted prompts for:
- React hooks and TypeScript
- FastAPI and Python backends
- Terraform and DevOps
- Docker and Kubernetes
- SQL optimization
- Code review and testing
- Security auditing
- And more!

If your vault already contains prompts, use --force to add seeds anyway.
`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		count, err := database.Count(ctx)
		if err != nil {
			return err
		}

		force, _ := cmd.Flags().GetBool("force")
		if count > 0 && !force {
			printWarning("Vault already contains %d prompts", count)
			printInfo("Use --force to add seed prompts anyway")
			return nil
		}

		// Show loading indicator
		stopLoading := showLoading("Initializing vault with curated prompts...")

		seeds := model.SeedPrompts()
		added := 0
		for _, p := range seeds {
			if err := database.Add(ctx, p); err != nil {
				stopLoading()
				printWarning("Skipping '%s': %v", p.Title, err)
				continue
			}
			added++
		}
		stopLoading()

		printSuccess("Initialized PromptVault with %d curated prompts!", added)
		fmt.Println()
		printInfo("Run 'promptvault' to open the TUI")
		printInfo("Run 'promptvault list' to see all prompts")
		return nil
	},
}

// import command — import prompts from JSON
var importCmd = &cobra.Command{
	Use:   "import [file]",
	Short: "Import prompts from a JSON file",
	Aliases: []string{"imp"},
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		filename := args[0]

		logDebug("Executing import command (file: %s)", filename)

		data, err := os.ReadFile(filename)
		if err != nil {
			printError("Failed to read file: %s", filename)
			return wrapError(err, "reading file")
		}

		var prompts []*model.Prompt
		if err := json.Unmarshal(data, &prompts); err != nil {
			printError("Failed to parse JSON")
			return wrapError(err, "parsing JSON")
		}

		printInfo("Importing %d prompts...", len(prompts))

		added := 0
		skipped := 0
		for _, p := range prompts {
			p.ID = "" // Force new ID generation
			if err := database.Add(ctx, p); err != nil {
				logDebug("Skipping '%s': %v", p.Title, err)
				skipped++
				continue
			}
			added++
		}

		printSuccess("Imported %d prompts from %s", added, filename)
		if skipped > 0 {
			printWarning("Skipped %d prompts (duplicates or errors)", skipped)
		}
		return nil
	},
}

// stats command
var statsCmd = &cobra.Command{
	Use:     "stats",
	Short:   "Show vault statistics",
	Aliases: []string{"statistics"},
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		total, stacks, err := database.Stats(ctx)
		if err != nil {
			return err
		}

		fmt.Println()
		fmt.Printf("%s %s\n", colorPrimary+iconSparkle, "PromptVault Statistics"+colorReset)
		fmt.Println(colorMuted + strings.Repeat("─", 40) + colorReset)
		fmt.Printf("  %-20s  %s%d%s\n", "Total Prompts:", colorSuccess, total, colorReset)
		fmt.Printf("  %-20s  %s%d%s\n", "Unique Stacks:", colorInfo, stacks, colorReset)
		fmt.Printf("  %-20s  %s%s%s\n", "Database Path:", colorMuted, database.Path(), colorReset)
		fmt.Println()

		return nil
	},
}

// stacks command — list all available stacks
var stacksCmd = &cobra.Command{
	Use:   "stacks",
	Short: "List all tech stack taxonomies",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Available stack paths:")
		for _, s := range model.DefaultStacks {
			fmt.Printf("  %s\n", s)
		}
		fmt.Println("\nUse: promptvault add --stack frontend/react/hooks ...")
		return nil
	},
}

// DX: Shell completion command
var completionCmd = &cobra.Command{
	Use:   "completion [bash|zsh|fish|powershell]",
	Short: "Generate shell completion scripts",
	Long: fmt.Sprintf(`Generate shell completion scripts for %s.

To load completions in your current shell:

  Bash:
    $ source <(%s completion bash)

  Zsh:
    $ source <(%s completion zsh)

  Fish:
    $ %s completion fish | source

  PowerShell:
    $ %s completion powershell | Out-String | Invoke-Expression

To load completions automatically on every shell startup:

  Bash:
    $ %s completion bash > ~/.bash_completion
    # Or: $ %s completion bash > /etc/bash_completion.d/promptvault

  Zsh:
    $ %s completion zsh > "${fpath[1]}/_promptvault"

  Fish:
    $ %s completion fish > ~/.config/fish/completions/promptvault.fish
`, rootCmd.Name(), rootCmd.Name(), rootCmd.Name(), rootCmd.Name(), rootCmd.Name(),
		rootCmd.Name(), rootCmd.Name(), rootCmd.Name(), rootCmd.Name()),
	Args:                  cobra.ExactArgs(1),
	ValidArgs:             []string{"bash", "zsh", "fish", "powershell"},
	DisableFlagsInUseLine: true,
	RunE: func(cmd *cobra.Command, args []string) error {
		switch args[0] {
		case "bash":
			return cmd.Root().GenBashCompletion(os.Stdout)
		case "zsh":
			return cmd.Root().GenZshCompletion(os.Stdout)
		case "fish":
			return cmd.Root().GenFishCompletion(os.Stdout, true)
		case "powershell":
			return cmd.Root().GenPowerShellCompletion(os.Stdout)
		}
		return fmt.Errorf("unsupported shell: %s", args[0])
	},
}

// mcp command
var mcpCmd = &cobra.Command{
	Use:   "mcp",
	Short: "Run the MCP Server over stdio for AI integration",
	RunE: func(cmd *cobra.Command, args []string) error {
		return mcp.Serve(database)
	},
}

// sync command
var syncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync prompts to/from a private GitHub Gist",
}

var syncPushCmd = &cobra.Command{
	Use:   "push",
	Short: "Backup all prompts to a private GitHub Gist",
	RunE: func(cmd *cobra.Command, args []string) error {
		token, _ := cmd.Flags().GetString("token")

		// Show loading indicator
		stopLoading := showLoading("Backing up prompts to GitHub Gist...")

		url, err := gistsync.Push(database, token)
		stopLoading()

		if err != nil {
			return err
		}
		fmt.Printf("\n✓ Successfully backed up prompts to %s\n", url)
		return nil
	},
}

var syncPullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Restore prompts from your GitHub Gist backup",
	RunE: func(cmd *cobra.Command, args []string) error {
		token, _ := cmd.Flags().GetString("token")

		// Show loading indicator
		stopLoading := showLoading("Restoring prompts from GitHub Gist...")

		added, err := gistsync.Pull(database, token)
		stopLoading()

		if err != nil {
			return err
		}
		fmt.Printf("\n✓ Successfully synced %d prompts from Gist\n", added)
		return nil
	},
}

// Execute runs the CLI
func Execute(d *db.DB) error {
	database = d
	return rootCmd.Execute()
}

func init() {
	// Note: add flags are defined inline with the command now

	// list flags
	listCmd.Flags().StringP("stack", "s", "", "Filter by stack")
	// Note: "short" and "json" flags are defined in the new init() above

	// get flags
	getCmd.Flags().BoolP("copy", "c", true, "Copy to clipboard")
	getCmd.Flags().BoolP("print", "p", false, "Print to stdout")

	// delete flags
	deleteCmd.Flags().BoolP("force", "f", false, "Skip confirmation")

	// export flags
	exportCmd.Flags().StringP("format", "f", "skill.md", "Output format: skill.md|agents.md|claude.md|cursorrules|windsurf|markdown|json|text")
	exportCmd.Flags().StringP("stack", "s", "", "Filter by stack")
	exportCmd.Flags().StringP("output", "o", "", "Output file (default: stdout)")

	// init flags
	initCmd.Flags().Bool("force", false, "Add seed prompts even if vault is not empty")

	// sync flags
	syncPushCmd.Flags().String("token", "", "GitHub Personal Access Token (or set PROMPTVAULT_GITHUB_TOKEN)")
	syncPullCmd.Flags().String("token", "", "GitHub Personal Access Token (or set PROMPTVAULT_GITHUB_TOKEN)")
	syncCmd.AddCommand(syncPushCmd, syncPullCmd)

	// Register all commands
	rootCmd.AddCommand(
		addCmd,
		listCmd,
		getCmd,
		searchCmd,
		deleteCmd,
		exportCmd,
		initCmd,
		importCmd,
		statsCmd,
		stacksCmd,
		mcpCmd,
		syncCmd,
		completionCmd, // DX: Shell completion
	)
}
