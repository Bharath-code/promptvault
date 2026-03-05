package sync

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

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

func LoadConfig() Config {
	var c Config
	data, err := os.ReadFile(configPath())
	if err == nil {
		json.Unmarshal(data, &c)
	}
	// Env fallback
	if token := os.Getenv("PROMPTVAULT_GITHUB_TOKEN"); token != "" {
		c.GitHubToken = token
	}
	if gist := os.Getenv("PROMPTVAULT_GIST_ID"); gist != "" {
		c.GistID = gist
	}
	return c
}

func SaveConfig(c Config) error {
	data, _ := json.MarshalIndent(c, "", "  ")
	return os.WriteFile(configPath(), data, 0600)
}

// Push exports the database to JSON and uploads it to a private GitHub Gist
func Push(d *db.DB, token string) (string, error) {
	if token == "" {
		return "", fmt.Errorf("GitHub token is required")
	}

	prompts, err := d.List("")
	if err != nil {
		return "", err
	}

	exporter := export.New(prompts)
	jsonContent, err := exporter.Export(export.FormatJSON)
	if err != nil {
		return "", err
	}

	cfg := LoadConfig()
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

	client := &http.Client{}
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
	cfg := LoadConfig()
	if token != "" {
		cfg.GitHubToken = token
	}

	if cfg.GistID == "" {
		return 0, fmt.Errorf("no Gist ID configured. Have you run 'sync push' yet or set PROMPTVAULT_GIST_ID?")
	}

	req, _ := http.NewRequest("GET", "https://api.github.com/gists/"+cfg.GistID, nil)
	if cfg.GitHubToken != "" {
		req.Header.Set("Authorization", "token "+cfg.GitHubToken)
	}

	client := &http.Client{}
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
		existing, err := d.Get(p.ID)
		if err == nil && existing != nil {
			d.Update(p)
		} else {
			d.Add(p)
		}
		added++
	}

	return added, nil
}
