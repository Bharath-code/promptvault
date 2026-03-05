package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/atotto/clipboard"
	"github.com/spf13/cobra"
	"github.com/Bharath-code/promptvault/internal/db"
	"github.com/Bharath-code/promptvault/internal/export"
	"github.com/Bharath-code/promptvault/internal/model"
	"github.com/Bharath-code/promptvault/internal/tui"
)

var database *db.DB

// Root command
var rootCmd = &cobra.Command{
	Use:   "promptvault",
	Short: "⚡ PromptVault — The universal prompt OS for developers",
	Long: `
⚡ PromptVault — Manage AI prompts by tech stack, right from your terminal.

Store, search, version, and deploy prompts across every AI tool.
Run without arguments to open the interactive TUI.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return tui.Run(database)
	},
}

// add command
var addCmd = &cobra.Command{
	Use:   "add [title]",
	Short: "Add a new prompt",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		title, _ := cmd.Flags().GetString("title")
		content, _ := cmd.Flags().GetString("content")
		stack, _ := cmd.Flags().GetString("stack")
		tagsStr, _ := cmd.Flags().GetString("tags")
		modelsStr, _ := cmd.Flags().GetString("models")
		verified, _ := cmd.Flags().GetBool("verified")

		if len(args) > 0 && title == "" {
			title = args[0]
		}

		if title == "" {
			return fmt.Errorf("title is required (use --title or pass as argument)")
		}
		if content == "" {
			// Try reading from stdin
			stat, _ := os.Stdin.Stat()
			if (stat.Mode() & os.ModeCharDevice) == 0 {
				buf := make([]byte, 1024*64)
				n, _ := os.Stdin.Read(buf)
				content = string(buf[:n])
			}
		}
		if content == "" {
			return fmt.Errorf("content is required (use --content or pipe via stdin)")
		}

		var tags []string
		for _, t := range strings.Split(tagsStr, ",") {
			if t = strings.TrimSpace(t); t != "" {
				tags = append(tags, t)
			}
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

		if err := database.Add(p); err != nil {
			return fmt.Errorf("adding prompt: %w", err)
		}

		fmt.Printf("✓ Added prompt: %s (id: %s)\n", p.Title, p.ID)
		return nil
	},
}

// list command
var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List prompts",
	Aliases: []string{"ls"},
	RunE: func(cmd *cobra.Command, args []string) error {
		stack, _ := cmd.Flags().GetString("stack")
		short, _ := cmd.Flags().GetBool("short")

		prompts, err := database.List(stack)
		if err != nil {
			return err
		}

		if len(prompts) == 0 {
			fmt.Println("No prompts found. Use 'promptvault add' or 'promptvault init' to get started.")
			return nil
		}

		for _, p := range prompts {
			if short {
				fmt.Printf("%-36s  %s\n", p.ID[:8], p.Title)
				continue
			}

			verified := ""
			if p.Verified {
				verified = " ✓"
			}
			stack := ""
			if p.Stack != "" {
				stack = fmt.Sprintf(" [%s]", p.Stack)
			}
			fmt.Printf("%-8s  %s%s%s\n", p.ID[:8], p.Title, verified, stack)
		}

		fmt.Printf("\n%d prompt(s)\n", len(prompts))
		return nil
	},
}

// get command — fetch a specific prompt and copy to clipboard
var getCmd = &cobra.Command{
	Use:   "get [id-or-title]",
	Short: "Get a prompt by ID or title (copies to clipboard)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		query := args[0]
		copyFlag, _ := cmd.Flags().GetBool("copy")
		printFlag, _ := cmd.Flags().GetBool("print")

		// Try by ID first, then search
		prompts, err := database.Search(query)
		if err != nil || len(prompts) == 0 {
			return fmt.Errorf("no prompt found matching: %s", query)
		}

		p := prompts[0]

		if printFlag || !copyFlag {
			fmt.Println(p.Content)
		}

		if copyFlag {
			if err := clipboard.WriteAll(p.Content); err != nil {
				return fmt.Errorf("copying to clipboard: %w", err)
			}
			_ = database.IncrementUsage(p.ID)
			fmt.Fprintf(os.Stderr, "✓ Copied '%s' to clipboard\n", p.Title)
		}

		return nil
	},
}

// search command
var searchCmd = &cobra.Command{
	Use:   "search [query]",
	Short: "Full-text search across all prompts",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		prompts, err := database.Search(args[0])
		if err != nil {
			return err
		}

		if len(prompts) == 0 {
			fmt.Printf("No prompts found for: %s\n", args[0])
			return nil
		}

		fmt.Printf("Found %d prompt(s):\n\n", len(prompts))
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

// delete command
var deleteCmd = &cobra.Command{
	Use:     "delete [id]",
	Short:   "Delete a prompt",
	Aliases: []string{"rm"},
	Args:    cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		force, _ := cmd.Flags().GetBool("force")

		// Search to find the prompt
		prompts, err := database.Search(args[0])
		if err != nil || len(prompts) == 0 {
			// Try exact ID
			if err := database.Delete(args[0]); err != nil {
				return fmt.Errorf("prompt not found: %s", args[0])
			}
			fmt.Println("✓ Deleted")
			return nil
		}

		p := prompts[0]

		if !force {
			fmt.Printf("Delete '%s'? [y/N] ", p.Title)
			var confirm string
			fmt.Scanln(&confirm)
			if strings.ToLower(confirm) != "y" {
				fmt.Println("Cancelled")
				return nil
			}
		}

		if err := database.Delete(p.ID); err != nil {
			return err
		}
		fmt.Printf("✓ Deleted: %s\n", p.Title)
		return nil
	},
}

// export command
var exportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export prompts to various formats (skill.md, agents.md, cursorrules, etc.)",
	RunE: func(cmd *cobra.Command, args []string) error {
		format, _ := cmd.Flags().GetString("format")
		stack, _ := cmd.Flags().GetString("stack")
		output, _ := cmd.Flags().GetString("output")

		prompts, err := database.List(stack)
		if err != nil {
			return err
		}

		if len(prompts) == 0 {
			return fmt.Errorf("no prompts to export")
		}

		e := export.New(prompts)
		result, err := e.Export(export.Format(format))
		if err != nil {
			return err
		}

		if output != "" {
			if err := os.WriteFile(output, []byte(result), 0644); err != nil {
				return fmt.Errorf("writing file: %w", err)
			}
			fmt.Fprintf(os.Stderr, "✓ Exported %d prompts to %s\n", len(prompts), output)
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
	RunE: func(cmd *cobra.Command, args []string) error {
		count, err := database.Count()
		if err != nil {
			return err
		}

		force, _ := cmd.Flags().GetBool("force")
		if count > 0 && !force {
			fmt.Printf("Vault already contains %d prompts. Use --force to add seed prompts anyway.\n", count)
			return nil
		}

		seeds := model.SeedPrompts()
		added := 0
		for _, p := range seeds {
			if err := database.Add(p); err != nil {
				fmt.Fprintf(os.Stderr, "⚠ Skipping '%s': %v\n", p.Title, err)
				continue
			}
			added++
		}

		fmt.Printf("⚡ Initialized PromptVault with %d curated prompts!\n", added)
		fmt.Println("\nRun 'promptvault' to open the TUI, or 'promptvault list' to see them.")
		return nil
	},
}

// import command — import prompts from JSON
var importCmd = &cobra.Command{
	Use:   "import [file]",
	Short: "Import prompts from a JSON file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := os.ReadFile(args[0])
		if err != nil {
			return fmt.Errorf("reading file: %w", err)
		}

		var prompts []*model.Prompt
		if err := json.Unmarshal(data, &prompts); err != nil {
			return fmt.Errorf("parsing JSON: %w", err)
		}

		added := 0
		for _, p := range prompts {
			p.ID = "" // Force new ID generation
			if err := database.Add(p); err != nil {
				fmt.Fprintf(os.Stderr, "⚠ Skipping '%s': %v\n", p.Title, err)
				continue
			}
			added++
		}

		fmt.Printf("✓ Imported %d prompts from %s\n", added, args[0])
		return nil
	},
}

// stats command
var statsCmd = &cobra.Command{
	Use:   "stats",
	Short: "Show vault statistics",
	RunE: func(cmd *cobra.Command, args []string) error {
		total, stacks, err := database.Stats()
		if err != nil {
			return err
		}
		fmt.Printf("⚡ PromptVault Stats\n")
		fmt.Printf("   Prompts : %d\n", total)
		fmt.Printf("   Stacks  : %d unique\n", stacks)
		fmt.Printf("   DB Path : %s\n", database.Path())
		return nil
	},
}

// stacks command — list all available stacks
var stacksCmd = &cobra.Command{
	Use:   "stacks",
	Short: "List all tech stack taxonomies",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("Available stack paths:\n")
		for _, s := range model.DefaultStacks {
			fmt.Printf("  %s\n", s)
		}
		fmt.Println("\nUse: promptvault add --stack frontend/react/hooks ...")
		return nil
	},
}

// Execute runs the CLI
func Execute(d *db.DB) error {
	database = d
	return rootCmd.Execute()
}

func init() {
	// add flags
	addCmd.Flags().StringP("title", "t", "", "Prompt title")
	addCmd.Flags().StringP("content", "c", "", "Prompt content (or pipe via stdin)")
	addCmd.Flags().StringP("stack", "s", "", "Tech stack path (e.g. frontend/react/hooks)")
	addCmd.Flags().String("tags", "", "Comma-separated tags")
	addCmd.Flags().String("models", "", "Comma-separated model names")
	addCmd.Flags().Bool("verified", false, "Mark as verified")

	// list flags
	listCmd.Flags().StringP("stack", "s", "", "Filter by stack")
	listCmd.Flags().Bool("short", false, "Short output (id + title only)")

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
	)
}
