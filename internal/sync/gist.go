package sync

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/Bharath-code/promptvault/internal/db"
	"github.com/Bharath-code/promptvault/internal/export"
	"github.com/Bharath-code/promptvault/internal/model"
)

type Config struct {
	GitHubToken string `json:"github_token"`
	GistID      string `json:"gist_id"`
}

func configPath() string {
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".promptvault", "config.json")
}

// gistIDRegex validates GitHub Gist ID format (alphanumeric only)
var gistIDRegex = regexp.MustCompile(`^[a-zA-Z0-9]+$`)

// githubTokenRegex validates GitHub Personal Access Token format (classic: ghp_, fine-grained: github_pat_)
var githubTokenRegex = regexp.MustCompile(`^(ghp_[a-zA-Z0-9]{36}|github_pat_[a-zA-Z0-9]{22}_[a-zA-Z0-9]{59})$`)

// validateGistID checks if the GistID is in valid format
func validateGistID(id string) error {
	if id == "" {
		return fmt.Errorf("GistID cannot be empty")
	}
	if len(id) > 64 {
		return fmt.Errorf("GistID too long (max 64 characters)")
	}
	if !gistIDRegex.MatchString(id) {
		return fmt.Errorf("GistID contains invalid characters (must be alphanumeric)")
	}
	return nil
}

// validateGitHubToken checks if the GitHub token is in valid format
func validateGitHubToken(token string) error {
	if token == "" {
		return fmt.Errorf("GitHub token cannot be empty")
	}
	if !githubTokenRegex.MatchString(token) {
		return fmt.Errorf("invalid GitHub token format. Token must start with 'ghp_' (classic PAT) or 'github_pat_' (fine-grained PAT)")
	}
	return nil
}

// LoadConfig loads the configuration from file
func LoadConfig() (Config, error) {
	var c Config
	data, err := os.ReadFile(configPath())
	if err == nil {
		if err := json.Unmarshal(data, &c); err != nil {
			return c, fmt.Errorf("parsing config file: %w", err)
		}
	}
	// Env fallback
	if token := os.Getenv("PROMPTVAULT_GITHUB_TOKEN"); token != "" {
		c.GitHubToken = token
	}
	if gist := os.Getenv("PROMPTVAULT_GIST_ID"); gist != "" {
		c.GistID = gist
	}
	return c, nil
}

func SaveConfig(c Config) error {
	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("marshaling config: %w", err)
	}
	if err := os.WriteFile(configPath(), data, 0600); err != nil {
		return fmt.Errorf("writing config file: %w", err)
	}
	return nil
}

// Push exports the database to JSON and uploads it to a private GitHub Gist
func Push(d *db.DB, token string) (string, error) {
	ctx := context.Background()
	if token == "" {
		return "", fmt.Errorf("GitHub token is required")
	}

	// Validate GitHub token format
	if err := validateGitHubToken(token); err != nil {
		return "", err
	}

	prompts, err := d.List(ctx, "")
	if err != nil {
		return "", err
	}

	exporter := export.New(prompts)
	jsonContent, err := exporter.Export(export.FormatJSON)
	if err != nil {
		return "", err
	}

	cfg, err := LoadConfig()
	if err != nil {
		return "", fmt.Errorf("loading config: %w", err)
	}
	if token != "" {
		cfg.GitHubToken = token
	}

	payload := map[string]interface{}{
		"description": "PromptVault Backup",
		"public":      false,
		"files": map[string]interface{}{
			"promptvault_export.json": map[string]string{
				"content": jsonContent,
			},
		},
	}

	body, _ := json.Marshal(payload)
	method := "POST"
	url := "https://api.github.com/gists"

	if cfg.GistID != "" {
		method = "PATCH"
		url = fmt.Sprintf("https://api.github.com/gists/%s", cfg.GistID)
	}

	req, _ := http.NewRequest(method, url, bytes.NewBuffer(body))
	req.Header.Set("Authorization", "token "+cfg.GitHubToken)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		b, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("GitHub API error (%d): %s", resp.StatusCode, string(b))
	}

	var result struct {
		ID      string `json:"id"`
		HTMLURL string `json:"html_url"`
	}
	json.NewDecoder(resp.Body).Decode(&result)

	cfg.GistID = result.ID
	SaveConfig(cfg)

	return result.HTMLURL, nil
}

// Pull downloads the JSON from the configured Gist and imports it
func Pull(d *db.DB, token string) (int, error) {
	ctx := context.Background()
	cfg, err := LoadConfig()
	if err != nil {
		return 0, fmt.Errorf("loading config: %w", err)
	}
	if token != "" {
		cfg.GitHubToken = token
	}

	if cfg.GistID == "" {
		return 0, fmt.Errorf("no Gist ID configured. Have you run 'sync push' yet or set PROMPTVAULT_GIST_ID?")
	}

	// Validate GistID format to prevent URL injection
	if err := validateGistID(cfg.GistID); err != nil {
		return 0, fmt.Errorf("invalid GistID: %w", err)
	}

	req, _ := http.NewRequest("GET", "https://api.github.com/gists/"+cfg.GistID, nil)
	if cfg.GitHubToken != "" {
		req.Header.Set("Authorization", "token "+cfg.GitHubToken)
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("failed to fetch gist (status %d)", resp.StatusCode)
	}

	var gist struct {
		Files map[string]struct {
			Content string `json:"content"`
		} `json:"files"`
	}
	json.NewDecoder(resp.Body).Decode(&gist)

	file, ok := gist.Files["promptvault_export.json"]
	if !ok {
		return 0, fmt.Errorf("promptvault_export.json not found in gist")
	}

	var prompts []*model.Prompt
	if err := json.Unmarshal([]byte(file.Content), &prompts); err != nil {
		return 0, fmt.Errorf("parsing gist content: %w", err)
	}

	// Just clear and reload for a pure "pull" (replace) instead of complicated merge
	// Better approach for backup restore
	added := 0
	for _, p := range prompts {
		// Only Add if it doesn't exist, or update if it does
		existing, err := d.Get(ctx, p.ID)
		if err == nil && existing != nil {
			if err := d.Update(ctx, p, "Sync from Gist", ""); err != nil {
				added++
				continue
			}
		} else {
			if err := d.Add(ctx, p); err != nil {
				added++
				continue
			}
		}
		added++
	}

	return added, nil
}
