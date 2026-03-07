package prompttest

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/Bharath-code/promptvault/internal/model"
)

// Tester handles prompt testing against AI models
type Tester struct {
	claudeAPIKey  string
	openaiAPIKey  string
	httpClient    *http.Client
}

// NewTester creates a new Tester instance
func NewTester() *Tester {
	return &Tester{
		claudeAPIKey: os.Getenv("ANTHROPIC_API_KEY"),
		openaiAPIKey: os.Getenv("OPENAI_API_KEY"),
		httpClient:   &http.Client{Timeout: 60 * time.Second},
	}
}

// TestPrompt tests a prompt against a model and returns the result
func (t *Tester) TestPrompt(ctx context.Context, prompt *model.Prompt, input, expectedOutput, modelName string) (*model.TestResult, error) {
	startTime := time.Now()
	
	// Execute the prompt
	actualOutput, tokenUsage, err := t.executePrompt(ctx, prompt.Content, input, modelName)
	latency := time.Since(startTime).Milliseconds()
	
	if err != nil {
		return &model.TestResult{
			ID:           uuid.New().String(),
			PromptID:     prompt.ID,
			Model:        modelName,
			Input:        input,
			ExpectedOutput: expectedOutput,
			ActualOutput: "",
			Passed:       false,
			Score:        0,
			LatencyMs:    int(latency),
			TokenUsage:   tokenUsage,
			ErrorMessage: err.Error(),
			CreatedAt:    time.Now().UTC(),
		}, nil
	}
	
	// Calculate similarity score
	score := calculateSimilarity(expectedOutput, actualOutput)
	passed := score >= 70.0 // 70% similarity threshold
	
	return &model.TestResult{
		ID:           uuid.New().String(),
		PromptID:     prompt.ID,
		Model:        modelName,
		Input:        input,
		ExpectedOutput: expectedOutput,
		ActualOutput: actualOutput,
		Passed:       passed,
		Score:        score,
		LatencyMs:    int(latency),
		TokenUsage:   tokenUsage,
		ErrorMessage: "",
		CreatedAt:    time.Now().UTC(),
	}, nil
}

// executePrompt sends the prompt to the AI model and returns the response
func (t *Tester) executePrompt(ctx context.Context, promptContent, input, modelName string) (string, int, error) {
	// Combine prompt with input
	fullPrompt := promptContent
	if input != "" {
		fullPrompt = promptContent + "\n\nInput:\n" + input
	}
	
	// Route to appropriate model
	if strings.Contains(modelName, "claude") {
		return t.callClaude(ctx, fullPrompt, modelName)
	} else if strings.Contains(modelName, "gpt") || strings.Contains(modelName, "chatgpt") {
		return t.callOpenAI(ctx, fullPrompt, modelName)
	}
	
	return "", 0, fmt.Errorf("unsupported model: %s", modelName)
}

// callClaude calls Anthropic's Claude API
func (t *Tester) callClaude(ctx context.Context, prompt, model string) (string, int, error) {
	if t.claudeAPIKey == "" {
		return "", 0, fmt.Errorf("ANTHROPIC_API_KEY not set")
	}
	
	// Default to claude-sonnet-4-20250514 if no specific version
	if model == "claude-sonnet" || model == "" {
		model = "claude-sonnet-4-20250514"
	}
	
	payload := map[string]interface{}{
		"model": model,
		"max_tokens": 4096,
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
	}
	
	body, _ := json.Marshal(payload)
	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.anthropic.com/v1/messages", bytes.NewBuffer(body))
	if err != nil {
		return "", 0, err
	}
	
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", t.claudeAPIKey)
	req.Header.Set("anthropic-version", "2023-06-01")
	
	resp, err := t.httpClient.Do(req)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()
	
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", 0, err
	}
	
	if resp.StatusCode != http.StatusOK {
		return "", 0, fmt.Errorf("Claude API error (%d): %s", resp.StatusCode, string(respBody))
	}
	
	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", 0, err
	}
	
	// Extract content
	content, _ := result["content"].([]interface{})
	if len(content) == 0 {
		return "", 0, fmt.Errorf("no content in Claude response")
	}
	
	text, _ := content[0].(map[string]interface{})["text"].(string)
	
	// Extract usage
	usage, _ := result["usage"].(map[string]interface{})
	inputTokens, _ := usage["input_tokens"].(float64)
	outputTokens, _ := usage["output_tokens"].(float64)
	
	return text, int(inputTokens + outputTokens), nil
}

// callOpenAI calls OpenAI's API
func (t *Tester) callOpenAI(ctx context.Context, prompt, model string) (string, int, error) {
	if t.openaiAPIKey == "" {
		return "", 0, fmt.Errorf("OPENAI_API_KEY not set")
	}
	
	// Default to gpt-4o if no specific version
	if model == "gpt-4o" || model == "" {
		model = "gpt-4o"
	}
	
	payload := map[string]interface{}{
		"model": model,
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
	}
	
	body, _ := json.Marshal(payload)
	req, err := http.NewRequestWithContext(ctx, "POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer(body))
	if err != nil {
		return "", 0, err
	}
	
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+t.openaiAPIKey)
	
	resp, err := t.httpClient.Do(req)
	if err != nil {
		return "", 0, err
	}
	defer resp.Body.Close()
	
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", 0, err
	}
	
	if resp.StatusCode != http.StatusOK {
		return "", 0, fmt.Errorf("OpenAI API error (%d): %s", resp.StatusCode, string(respBody))
	}
	
	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return "", 0, err
	}
	
	// Extract content
	choices, _ := result["choices"].([]interface{})
	if len(choices) == 0 {
		return "", 0, fmt.Errorf("no choices in OpenAI response")
	}
	
	message, _ := choices[0].(map[string]interface{})["message"].(map[string]interface{})
	text, _ := message["content"].(string)
	
	// Extract usage
	usage, _ := result["usage"].(map[string]interface{})
	inputTokens, _ := usage["prompt_tokens"].(float64)
	outputTokens, _ := usage["completion_tokens"].(float64)
	
	return text, int(inputTokens + outputTokens), nil
}

// calculateSimilarity calculates similarity between expected and actual output
func calculateSimilarity(expected, actual string) float64 {
	// Simple word overlap similarity (can be improved with more sophisticated algorithms)
	expectedWords := strings.Fields(strings.ToLower(expected))
	actualWords := strings.Fields(strings.ToLower(actual))
	
	if len(expectedWords) == 0 {
		return 0
	}
	
	// Count matching words
	matches := 0
	for _, ew := range expectedWords {
		for _, aw := range actualWords {
			if ew == aw {
				matches++
				break
			}
		}
	}
	
	// Jaccard-like similarity
	similarity := float64(matches) / float64(len(expectedWords)) * 100
	
	// Cap at 100
	if similarity > 100 {
		similarity = 100
	}
	
	return similarity
}
