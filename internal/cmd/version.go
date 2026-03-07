package cmd

import (
	"context"
	"fmt"
	"os"
	"os/user"
	"strings"

	"github.com/sergi/go-diff/diffmatchpatch"
	"github.com/spf13/cobra"
)

// history command
var historyCmd = &cobra.Command{
	Use:   "history [prompt-id]",
	Short: "View version history of a prompt",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		promptID := args[0]

		// Get prompt
		prompt, err := database.Get(ctx, promptID)
		if err != nil {
			printError("Prompt not found: %s", promptID)
			return err
		}

		// Get history
		versions, err := database.GetPromptHistory(ctx, promptID)
		if err != nil {
			printError("Failed to load history: %v", err)
			return err
		}

		fmt.Println()
		fmt.Printf("📜 Version History: %s\n", prompt.Title)
		fmt.Println(strings.Repeat("─", 70))

		if len(versions) == 0 {
			fmt.Println("No versions recorded yet")
			return nil
		}

		for _, v := range versions {
			author := v.Author
			if author == "" {
				author = "unknown"
			}
			commitMsg := v.CommitMsg
			if commitMsg == "" {
				commitMsg = "No message"
			}

			marker := "  "
			if v.Version == len(versions) {
				marker = "▶ "
			}

			fmt.Printf("%s v%d  %s  %s  %s\n",
				marker, v.Version, v.CreatedAt.Format("2006-01-02 15:04"),
				author, commitMsg)
		}

		fmt.Println()
		fmt.Printf("Total versions: %d\n", len(versions))
		return nil
	},
}

// diff command
var diffCmd = &cobra.Command{
	Use:   "diff [prompt-id] [version1] [version2]",
	Short: "Compare two versions of a prompt",
	Long: `Compare two versions of a prompt and show the differences.

Examples:
  # Compare version 1 and 2
  promptvault diff abc123 1 2

  # Compare current (latest) with version 1
  promptvault diff abc123 1 current

  # Compare last two versions
  promptvault diff abc123 HEAD~1 HEAD
`,
	Args: cobra.ExactArgs(3),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		promptID := args[0]
		v1Str := args[1]
		v2Str := args[2]

		// Get current version number
		currentVersion, err := database.GetCurrentVersion(ctx, promptID)
		if err != nil {
			printError("Failed to get current version: %v", err)
			return err
		}

		// Parse version strings
		v1, err := parseVersion(v1Str, currentVersion)
		if err != nil {
			printError("Invalid version: %s", v1Str)
			return err
		}

		v2, err := parseVersion(v2Str, currentVersion)
		if err != nil {
			printError("Invalid version: %s", v2Str)
			return err
		}

		// Get versions
		ver1, err := database.GetPromptVersion(ctx, promptID, v1)
		if err != nil {
			printError("Version %d not found", v1)
			return err
		}

		ver2, err := database.GetPromptVersion(ctx, promptID, v2)
		if err != nil {
			printError("Version %d not found", v2)
			return err
		}

		// Show diff
		fmt.Println()
		fmt.Printf("📊 Diff: %s (v%d → v%d)\n", promptID, v1, v2)
		fmt.Println(strings.Repeat("─", 70))

		// Title diff
		if ver1.Title != ver2.Title {
			fmt.Printf("\n%s Title:\n", colorInfo)
			fmt.Printf("  - %s\n", ver1.Title)
			fmt.Printf("  + %s\n", ver2.Title)
		}

		// Content diff
		dmp := diffmatchpatch.New()
		diffs := dmp.DiffMain(ver1.Content, ver2.Content, false)

		fmt.Printf("\n%s Content:\n", colorInfo)
		for _, diff := range diffs {
			switch diff.Type {
			case diffmatchpatch.DiffDelete:
				fmt.Printf("%s- %s%s\n", colorError, diff.Text, colorReset)
			case diffmatchpatch.DiffInsert:
				fmt.Printf("%s+ %s%s\n", colorSuccess, diff.Text, colorReset)
			default:
				fmt.Printf("  %s\n", diff.Text)
			}
		}

		fmt.Println()
		return nil
	},
}

// revert command
var revertCmd = &cobra.Command{
	Use:   "revert [prompt-id] [version]",
	Short: "Revert a prompt to a previous version",
	Long: `Revert a prompt to a previous version.

This creates a new version with the content from the specified version.

Examples:
  # Revert to version 3
  promptvault revert abc123 3

  # Revert with custom message
  promptvault revert abc123 3 --message "Reverting broken changes"
`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		promptID := args[0]
		versionStr := args[1]

		message, _ := cmd.Flags().GetString("message")

		// Get current version
		currentVersion, err := database.GetCurrentVersion(ctx, promptID)
		if err != nil {
			printError("Failed to get current version: %v", err)
			return err
		}

		// Parse version
		targetVersion, err := parseVersion(versionStr, currentVersion)
		if err != nil {
			printError("Invalid version: %s", versionStr)
			return err
		}

		// Get target version
		target, err := database.GetPromptVersion(ctx, promptID, targetVersion)
		if err != nil {
			printError("Version %d not found", targetVersion)
			return err
		}

		// Get current prompt
		prompt, err := database.Get(ctx, promptID)
		if err != nil {
			printError("Prompt not found: %s", promptID)
			return err
		}

		// Update prompt with old content
		prompt.Title = target.Title
		prompt.Content = target.Content
		prompt.Tags = target.Tags
		prompt.Stack = target.Stack
		prompt.Models = target.Models
		prompt.Verified = target.Verified

		// Get author
		author := getAuthor()
		if message == "" {
			message = fmt.Sprintf("Reverted to v%d", targetVersion)
		}

		if err := database.Update(ctx, prompt, message, author); err != nil {
			printError("Failed to revert: %v", err)
			return err
		}

		printSuccess("Reverted %s to v%d", promptID, targetVersion)
		fmt.Printf("New version: v%d\n", currentVersion+1)
		return nil
	},
}

func init() {
	revertCmd.Flags().StringP("message", "m", "", "Commit message for the revert")
}

// parseVersion parses version string (supports "current", "HEAD", numbers)
func parseVersion(v string, current int) (int, error) {
	switch v {
	case "current", "HEAD", "latest":
		return current, nil
	case "HEAD~1":
		if current <= 1 {
			return 0, fmt.Errorf("no previous version")
		}
		return current - 1, nil
	default:
		var version int
		_, err := fmt.Sscanf(v, "%d", &version)
		if err != nil {
			return 0, err
		}
		return version, nil
	}
}

// getAuthor returns the current user's username
func getAuthor() string {
	u, err := user.Current()
	if err != nil {
		return os.Getenv("USER")
	}
	return u.Username
}
