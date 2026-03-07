package cmd

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
	"github.com/Bharath-code/promptvault/internal/db"
	"github.com/Bharath-code/promptvault/internal/model"
	"github.com/Bharath-code/promptvault/internal/prompttest"
)

// test command
var testCmd = &cobra.Command{
	Use:   "test [prompt-id]",
	Short: "Test a prompt against AI models",
	Long: `Test a prompt against AI models to verify it produces expected outputs.

This command helps you validate that your prompts work correctly before deploying them.

Examples:
  # Interactive test mode
  promptvault test abc123

  # Test with specific input and expected output
  promptvault test abc123 --input "test input" --expected "expected output"

  # Test against specific model
  promptvault test abc123 --model claude-sonnet

  # Run all tests for a prompt
  promptvault test abc123 --all

  # View test history
  promptvault test abc123 --history
`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := context.Background()
		promptID := args[0]
		
		// Flags
		input, _ := cmd.Flags().GetString("input")
		expected, _ := cmd.Flags().GetString("expected")
		model, _ := cmd.Flags().GetString("model")
		all, _ := cmd.Flags().GetBool("all")
		history, _ := cmd.Flags().GetBool("history")
		
		// Get the prompt
		prompt, err := database.Get(ctx, promptID)
		if err != nil {
			printError("Prompt not found: %s", promptID)
			return fmt.Errorf("prompt not found: %s", promptID)
		}
		
		// Show test history if requested
		if history {
			return showTestHistory(ctx, database, prompt)
		}
		
		// Run all tests if requested
		if all {
			return runAllTests(ctx, database, prompt)
		}
		
		// Interactive mode if no input/expected provided
		if input == "" || expected == "" {
			return runInteractiveTest(ctx, database, prompt, model)
		}
		
		// Single test with provided input/expected
		return runSingleTest(ctx, database, prompt, input, expected, model)
	},
}

func init() {
	testCmd.Flags().String("input", "", "Test input")
	testCmd.Flags().String("expected", "", "Expected output")
	testCmd.Flags().String("model", "claude-sonnet", "Model to test against")
	testCmd.Flags().Bool("all", false, "Run all saved tests for this prompt")
	testCmd.Flags().Bool("history", false, "Show test history")
}

// runInteractiveTest runs an interactive test session
func runInteractiveTest(ctx context.Context, db *db.DB, prompt *model.Prompt, model string) error {
	fmt.Println()
	fmt.Printf("📋 Testing prompt: %s\n", prompt.Title)
	fmt.Printf("   Model: %s\n", model)
	fmt.Println()
	fmt.Println("Enter test input (empty line to skip, 'quit' to exit):")
	fmt.Println(strings.Repeat("─", 60))
	
	scanner := bufio.NewScanner(os.Stdin)
	
	for {
		fmt.Print("\n💬 Input: ")
		if !scanner.Scan() {
			break
		}
		input := strings.TrimSpace(scanner.Text())
		
		if input == "quit" || input == "q" {
			printInfo("Test session ended")
			break
		}
		
		if input == "" {
			continue
		}
		
		fmt.Print("✨ Expected output: ")
		if !scanner.Scan() {
			break
		}
		expected := strings.TrimSpace(scanner.Text())
		
		if expected == "" {
			continue
		}
		
		// Run the test
		fmt.Println()
		logInfo("Running test...")
		
		tester := prompttest.NewTester()
		result, err := tester.TestPrompt(ctx, prompt, input, expected, model)
		if err != nil {
			printError("Test failed: %v", err)
			continue
		}
		
		// Save test result
		if err := db.SaveTestResult(ctx, result); err != nil {
			logDebug("Failed to save test result: %v", err)
		}
		
		// Display results
		displayTestResult(result)
	}
	
	return nil
}

// runSingleTest runs a single test with provided parameters
func runSingleTest(ctx context.Context, db *db.DB, prompt *model.Prompt, input, expected, model string) error {
	logInfo("Testing prompt: %s", prompt.Title)
	logInfo("Model: %s", model)
	fmt.Println()
	
	tester := prompttest.NewTester()
	result, err := tester.TestPrompt(ctx, prompt, input, expected, model)
	if err != nil {
		printError("Test failed: %v", err)
		return err
	}
	
	// Save test result
	if err := db.SaveTestResult(ctx, result); err != nil {
		logDebug("Failed to save test result: %v", err)
	}
	
	displayTestResult(result)
	return nil
}

// runAllTests runs all saved tests for a prompt
func runAllTests(ctx context.Context, db *db.DB, prompt *model.Prompt) error {
	printError("--all flag not yet implemented")
	printInfo("This will run all previously saved tests for this prompt")
	return nil
}

// showTestHistory shows test history for a prompt
func showTestHistory(ctx context.Context, db *db.DB, prompt *model.Prompt) error {
	suite, err := db.GetPromptTestSuite(ctx, prompt.ID)
	if err != nil {
		printError("Failed to load test history: %v", err)
		return err
	}
	
	fmt.Println()
	fmt.Printf("📊 Test History for: %s\n", prompt.Title)
	fmt.Println(strings.Repeat("─", 60))
	
	if len(suite.Tests) == 0 {
		fmt.Println("No tests run yet for this prompt")
		return nil
	}
	
	fmt.Printf("Total Tests: %d\n", len(suite.Tests))
	fmt.Printf("Pass Rate: %.1f%%\n", suite.PassRate)
	fmt.Printf("Average Score: %.1f/100\n", suite.AvgScore)
	fmt.Println()
	
	// Show last 10 tests
	limit := 10
	if len(suite.Tests) < limit {
		limit = len(suite.Tests)
	}
	
	fmt.Println("Recent Tests:")
	for i := 0; i < limit; i++ {
		test := suite.Tests[i]
		status := "✅"
		if !test.Passed {
			status = "❌"
		}
		fmt.Printf("  %s [%s] Score: %.1f%% - %s\n", 
			status, test.Model, test.Score, test.CreatedAt.Format("2006-01-02 15:04"))
	}
	
	return nil
}

// displayTestResult displays a test result
func displayTestResult(result *model.TestResult) {
	fmt.Println()
	if result.Passed {
		printSuccess("✅ Test PASSED")
	} else {
		printError("❌ Test FAILED")
	}
	
	fmt.Println()
	fmt.Printf("Score: %.1f/100\n", result.Score)
	fmt.Printf("Latency: %dms\n", result.LatencyMs)
	fmt.Printf("Tokens: %d\n", result.TokenUsage)
	fmt.Println()
	
	if result.ErrorMessage != "" {
		printError("Error: %s", result.ErrorMessage)
		fmt.Println()
	}
	
	fmt.Println("Actual Output:")
	fmt.Println(strings.Repeat("─", 60))
	fmt.Println(result.ActualOutput)
	fmt.Println(strings.Repeat("─", 60))
}
