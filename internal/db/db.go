package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"github.com/Bharath-code/promptvault/internal/model"
)

// Maximum content length for prompts (100KB)
const maxContentLength = 100 * 1024

// Maximum title length for prompts
const maxTitleLength = 500

type DB struct {
	conn *sql.DB
	path string
	mu   sync.RWMutex // Protects concurrent access to the database
}

// Open opens (or creates) the SQLite database
func Open() (*DB, error) {
	dir, err := dataDir()
	if err != nil {
		return nil, err
	}

	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("creating data dir: %w", err)
	}

	dbPath := filepath.Join(dir, "vault.db")
	return OpenPath(dbPath)
}

// OpenPath opens (or creates) the SQLite database at the specified path
func OpenPath(dbPath string) (*DB, error) {
	// Validate and clean the path to prevent path traversal
	if dbPath == "" {
		return nil, fmt.Errorf("database path cannot be empty")
	}
	
	// Clean the path to resolve any .. or . components
	dbPath = filepath.Clean(dbPath)
	
	// Ensure the path is absolute
	if !filepath.IsAbs(dbPath) {
		return nil, fmt.Errorf("database path must be absolute")
	}
	
	// Ensure the path has .db extension for safety
	if !strings.HasSuffix(dbPath, ".db") {
		return nil, fmt.Errorf("database path must have .db extension")
	}
	
	// Enable WAL mode, NORMAL sync (faster than FULL), and a large in-memory cache to speed up reads
	dsn := dbPath + "?_journal_mode=WAL&_synchronous=NORMAL&_cache_size=100000&_busy_timeout=5000"
	conn, err := sql.Open("sqlite3", dsn)
	if err != nil {
		return nil, fmt.Errorf("opening db: %w", err)
	}

	// Crucial for sqlite to prevent "database is locked" and optimize concurrent FTS scanning
	conn.SetMaxOpenConns(1)

	d := &DB{conn: conn, path: dbPath}
	if err := d.migrate(); err != nil {
		return nil, fmt.Errorf("migrating db: %w", err)
	}

	return d, nil
}

func (d *DB) Close() error {
	return d.conn.Close()
}

func (d *DB) Path() string {
	return d.path
}

// validatePrompt checks if the prompt has valid content
func validatePrompt(p *model.Prompt) error {
	if p == nil {
		return fmt.Errorf("prompt cannot be nil")
	}
	
	title := strings.TrimSpace(p.Title)
	if title == "" {
		return fmt.Errorf("title is required")
	}
	if len(title) > maxTitleLength {
		return fmt.Errorf("title exceeds maximum length of %d characters", maxTitleLength)
	}
	
	content := strings.TrimSpace(p.Content)
	if content == "" {
		return fmt.Errorf("content is required")
	}
	if len(content) > maxContentLength {
		return fmt.Errorf("content exceeds maximum length of %d bytes", maxContentLength)
	}
	
	if p.Stack != "" && len(p.Stack) > 200 {
		return fmt.Errorf("stack path exceeds maximum length of 200 characters")
	}
	
	return nil
}

// migrate creates all tables
func (d *DB) migrate() error {
	schema := `
	CREATE TABLE IF NOT EXISTS prompts (
		id          TEXT PRIMARY KEY,
		title       TEXT NOT NULL,
		content     TEXT NOT NULL,
		tags        TEXT NOT NULL DEFAULT '[]',
		stack       TEXT NOT NULL DEFAULT '',
		models      TEXT NOT NULL DEFAULT '[]',
		verified    INTEGER NOT NULL DEFAULT 0,
		usage_count INTEGER NOT NULL DEFAULT 0,
		created_at  DATETIME NOT NULL,
		updated_at  DATETIME NOT NULL,
		last_used_at DATETIME
	);

	CREATE TABLE IF NOT EXISTS test_results (
		id              TEXT PRIMARY KEY,
		prompt_id       TEXT NOT NULL,
		model           TEXT NOT NULL,
		input           TEXT NOT NULL,
		expected_output TEXT NOT NULL,
		actual_output   TEXT NOT NULL,
		passed          INTEGER NOT NULL,
		score           REAL NOT NULL,
		latency_ms      INTEGER NOT NULL,
		token_usage     INTEGER NOT NULL,
		error_message   TEXT,
		created_at      DATETIME NOT NULL,
		FOREIGN KEY (prompt_id) REFERENCES prompts(id) ON DELETE CASCADE
	);

	CREATE TABLE IF NOT EXISTS prompt_versions (
		id          TEXT PRIMARY KEY,
		prompt_id   TEXT NOT NULL,
		version     INTEGER NOT NULL,
		title       TEXT NOT NULL,
		content     TEXT NOT NULL,
		tags        TEXT NOT NULL,
		stack       TEXT NOT NULL,
		models      TEXT NOT NULL,
		verified    INTEGER NOT NULL,
		commit_msg  TEXT,
		author      TEXT,
		created_at  DATETIME NOT NULL,
		FOREIGN KEY (prompt_id) REFERENCES prompts(id) ON DELETE CASCADE
	);

	CREATE INDEX IF NOT EXISTS idx_prompts_stack ON prompts(stack);
	CREATE INDEX IF NOT EXISTS idx_prompts_updated ON prompts(updated_at DESC);
	CREATE INDEX IF NOT EXISTS idx_test_results_prompt ON test_results(prompt_id);
	CREATE INDEX IF NOT EXISTS idx_versions_prompt ON prompt_versions(prompt_id);
	CREATE INDEX IF NOT EXISTS idx_versions_version ON prompt_versions(prompt_id, version);
	CREATE VIRTUAL TABLE IF NOT EXISTS prompts_fts USING fts5(
		id UNINDEXED,
		title,
		content,
		stack,
		tags,
		tokenize='porter ascii'
	);

	CREATE TRIGGER IF NOT EXISTS prompts_ai AFTER INSERT ON prompts BEGIN
		INSERT INTO prompts_fts(id, title, content, stack, tags)
		VALUES (new.id, new.title, new.content, new.stack, new.tags);
	END;

	CREATE TRIGGER IF NOT EXISTS prompts_ad AFTER DELETE ON prompts BEGIN
		DELETE FROM prompts_fts WHERE id = old.id;
	END;

	CREATE TRIGGER IF NOT EXISTS prompts_au AFTER UPDATE ON prompts BEGIN
		DELETE FROM prompts_fts WHERE id = old.id;
		INSERT INTO prompts_fts(id, title, content, stack, tags)
		VALUES (new.id, new.title, new.content, new.stack, new.tags);
	END;
	`
	_, err := d.conn.Exec(schema)
	return err
}

// Add inserts a new prompt
func (d *DB) Add(ctx context.Context, p *model.Prompt) error {
	// Validate prompt before inserting
	if err := validatePrompt(p); err != nil {
		return fmt.Errorf("invalid prompt: %w", err)
	}
	
	d.mu.Lock()
	defer d.mu.Unlock()
	
	if p.ID == "" {
		p.ID = uuid.New().String()
	}
	now := time.Now().UTC()
	p.CreatedAt = now
	p.UpdatedAt = now

	tags, _ := json.Marshal(p.Tags)
	models, _ := json.Marshal(p.Models)

	_, err := d.conn.ExecContext(ctx, `
		INSERT INTO prompts (id, title, content, tags, stack, models, verified, usage_count, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, 0, ?, ?)`,
		p.ID, p.Title, p.Content, string(tags), p.Stack, string(models),
		boolToInt(p.Verified), now, now,
	)
	return err
}

// Get returns a prompt by ID
func (d *DB) Get(ctx context.Context, id string) (*model.Prompt, error) {
	row := d.conn.QueryRowContext(ctx, `SELECT * FROM prompts WHERE id = ?`, id)
	return scanPrompt(row)
}

// List returns all prompts, optionally filtered by stack prefix
func (d *DB) List(ctx context.Context, stackFilter string) ([]*model.Prompt, error) {
	var rows *sql.Rows
	var err error

	if stackFilter != "" {
		rows, err = d.conn.QueryContext(ctx, `
			SELECT * FROM prompts
			WHERE stack = ? OR stack LIKE ?
			ORDER BY updated_at DESC`,
			stackFilter, stackFilter+"/%",
		)
	} else {
		rows, err = d.conn.QueryContext(ctx, `SELECT * FROM prompts ORDER BY updated_at DESC`)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanPrompts(ctx, rows)
}

// Search does full-text search across title, content, tags, stack
func (d *DB) Search(ctx context.Context, query string) ([]*model.Prompt, error) {
	// Clean query for FTS5
	query = strings.ReplaceAll(query, `"`, `""`)

	rows, err := d.conn.QueryContext(ctx, `
		SELECT p.* FROM prompts p
		JOIN prompts_fts fts ON p.id = fts.id
		WHERE prompts_fts MATCH ?
		ORDER BY rank
		LIMIT 50`,
		query+"*",
	)
	if err != nil {
		// Fallback to LIKE search if FTS fails
		like := "%" + query + "%"
		rows, err = d.conn.QueryContext(ctx, `
			SELECT * FROM prompts
			WHERE title LIKE ? OR content LIKE ? OR stack LIKE ? OR tags LIKE ?
			ORDER BY updated_at DESC LIMIT 50`,
			like, like, like, like,
		)
		if err != nil {
			return nil, err
		}
	}
	defer rows.Close()
	return scanPrompts(ctx, rows)
}

// Update modifies an existing prompt
func (d *DB) Update(ctx context.Context, p *model.Prompt, commitMsg, author string) error {
	// Validate prompt before updating
	if err := validatePrompt(p); err != nil {
		return fmt.Errorf("invalid prompt: %w", err)
	}
	
	d.mu.Lock()
	defer d.mu.Unlock()
	
	// Create version snapshot before updating
	if err := d.createVersionInternal(ctx, p, commitMsg, author); err != nil {
		// Don't fail the update if versioning fails
		// logDebug("Failed to create version snapshot: %v", err)
	}
	
	p.UpdatedAt = time.Now().UTC()
	tags, _ := json.Marshal(p.Tags)
	models, _ := json.Marshal(p.Models)

	res, err := d.conn.ExecContext(ctx, `
		UPDATE prompts SET
			title = ?, content = ?, tags = ?, stack = ?,
			models = ?, verified = ?, updated_at = ?
		WHERE id = ?`,
		p.Title, p.Content, string(tags), p.Stack,
		string(models), boolToInt(p.Verified), p.UpdatedAt, p.ID,
	)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("prompt not found: %s", p.ID)
	}
	return nil
}

// createVersionInternal creates a version snapshot (internal, no lock)
func (d *DB) createVersionInternal(ctx context.Context, prompt *model.Prompt, commitMsg, author string) error {
	// Get current max version
	var maxVersion int
	err := d.conn.QueryRowContext(ctx, `
		SELECT COALESCE(MAX(version), 0) FROM prompt_versions WHERE prompt_id = ?`,
		prompt.ID,
	).Scan(&maxVersion)
	if err != nil {
		return err
	}
	
	tags, _ := json.Marshal(prompt.Tags)
	models, _ := json.Marshal(prompt.Models)
	
	_, err = d.conn.ExecContext(ctx, `
		INSERT INTO prompt_versions (id, prompt_id, version, title, content, tags, stack, models, verified, commit_msg, author, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		uuid.New().String(), prompt.ID, maxVersion+1, prompt.Title, prompt.Content,
		string(tags), prompt.Stack, string(models), boolToInt(prompt.Verified),
		commitMsg, author, time.Now().UTC(),
	)
	return err
}

// Delete removes a prompt
func (d *DB) Delete(ctx context.Context, id string) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	
	res, err := d.conn.ExecContext(ctx, `DELETE FROM prompts WHERE id = ?`, id)
	if err != nil {
		return err
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return fmt.Errorf("prompt not found: %s", id)
	}
	return nil
}

// IncrementUsage bumps the usage counter and last_used_at
func (d *DB) IncrementUsage(ctx context.Context, id string) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	
	now := time.Now().UTC()
	_, err := d.conn.ExecContext(ctx, `
		UPDATE prompts SET usage_count = usage_count + 1, last_used_at = ? WHERE id = ?`,
		now, id,
	)
	return err
}

// Stats returns aggregate statistics
func (d *DB) Stats(ctx context.Context) (total int, stacks int, err error) {
	err = d.conn.QueryRowContext(ctx, `SELECT COUNT(*) FROM prompts`).Scan(&total)
	if err != nil {
		return
	}
	err = d.conn.QueryRowContext(ctx, `SELECT COUNT(DISTINCT stack) FROM prompts WHERE stack != ''`).Scan(&stacks)
	return
}

// Count returns just the total number of prompts
func (d *DB) Count(ctx context.Context) (int, error) {
	var count int
	err := d.conn.QueryRowContext(ctx, `SELECT COUNT(*) FROM prompts`).Scan(&count)
	return count, err
}

// SaveTestResult saves a test result for a prompt
func (d *DB) SaveTestResult(ctx context.Context, result *model.TestResult) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	
	_, err := d.conn.ExecContext(ctx, `
		INSERT INTO test_results (id, prompt_id, model, input, expected_output, actual_output, passed, score, latency_ms, token_usage, error_message, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		result.ID, result.PromptID, result.Model, result.Input, result.ExpectedOutput,
		result.ActualOutput, boolToInt(result.Passed), result.Score, result.LatencyMs,
		result.TokenUsage, result.ErrorMessage, result.CreatedAt,
	)
	return err
}

// GetTestResults returns all test results for a prompt
func (d *DB) GetTestResults(ctx context.Context, promptID string) ([]*model.TestResult, error) {
	rows, err := d.conn.QueryContext(ctx, `
		SELECT * FROM test_results
		WHERE prompt_id = ?
		ORDER BY created_at DESC`,
		promptID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanTestResults(rows)
}

// GetPromptTestSuite returns aggregated test suite for a prompt
func (d *DB) GetPromptTestSuite(ctx context.Context, promptID string) (*model.TestSuite, error) {
	results, err := d.GetTestResults(ctx, promptID)
	if err != nil {
		return nil, err
	}
	
	if len(results) == 0 {
		return &model.TestSuite{
			PromptID: promptID,
			Tests:    []model.TestResult{},
		}, nil
	}
	
	passed := 0
	totalScore := 0.0
	testResults := make([]model.TestResult, len(results))
	
	for i, r := range results {
		testResults[i] = *r
		if r.Passed {
			passed++
		}
		totalScore += r.Score
	}
	
	return &model.TestSuite{
		PromptID: promptID,
		Tests:    testResults,
		PassRate: float64(passed) / float64(len(results)) * 100,
		AvgScore: totalScore / float64(len(results)),
	}, nil
}

// DeleteTestResults deletes all test results for a prompt
func (d *DB) DeleteTestResults(ctx context.Context, promptID string) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	
	_, err := d.conn.ExecContext(ctx, `DELETE FROM test_results WHERE prompt_id = ?`, promptID)
	return err
}

// CreateVersion creates a new version snapshot of a prompt
func (d *DB) CreateVersion(ctx context.Context, prompt *model.Prompt, commitMsg, author string) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	
	// Get current max version
	var maxVersion int
	err := d.conn.QueryRowContext(ctx, `
		SELECT COALESCE(MAX(version), 0) FROM prompt_versions WHERE prompt_id = ?`,
		prompt.ID,
	).Scan(&maxVersion)
	if err != nil {
		return err
	}
	
	tags, _ := json.Marshal(prompt.Tags)
	models, _ := json.Marshal(prompt.Models)
	
	_, err = d.conn.ExecContext(ctx, `
		INSERT INTO prompt_versions (id, prompt_id, version, title, content, tags, stack, models, verified, commit_msg, author, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		uuid.New().String(), prompt.ID, maxVersion+1, prompt.Title, prompt.Content,
		string(tags), prompt.Stack, string(models), boolToInt(prompt.Verified),
		commitMsg, author, time.Now().UTC(),
	)
	return err
}

// GetPromptHistory returns all versions of a prompt
func (d *DB) GetPromptHistory(ctx context.Context, promptID string) ([]*model.PromptVersion, error) {
	rows, err := d.conn.QueryContext(ctx, `
		SELECT * FROM prompt_versions
		WHERE prompt_id = ?
		ORDER BY version DESC`,
		promptID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanPromptVersions(rows)
}

// GetPromptVersion returns a specific version of a prompt
func (d *DB) GetPromptVersion(ctx context.Context, promptID string, version int) (*model.PromptVersion, error) {
	row := d.conn.QueryRowContext(ctx, `
		SELECT * FROM prompt_versions
		WHERE prompt_id = ? AND version = ?`,
		promptID, version,
	)
	return scanPromptVersion(row)
}

// GetCurrentVersion returns the latest version number for a prompt
func (d *DB) GetCurrentVersion(ctx context.Context, promptID string) (int, error) {
	var version int
	err := d.conn.QueryRowContext(ctx, `
		SELECT COALESCE(MAX(version), 0) FROM prompt_versions WHERE prompt_id = ?`,
		promptID,
	).Scan(&version)
	return version, err
}

// DeletePromptVersions deletes all versions for a prompt
func (d *DB) DeletePromptVersions(ctx context.Context, promptID string) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	
	_, err := d.conn.ExecContext(ctx, `DELETE FROM prompt_versions WHERE prompt_id = ?`, promptID)
	return err
}

// --- helpers ---

func dataDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".promptvault"), nil
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func scanPrompt(row *sql.Row) (*model.Prompt, error) {
	var p model.Prompt
	var tags, models string
	var lastUsed sql.NullTime
	var verified int

	err := row.Scan(
		&p.ID, &p.Title, &p.Content, &tags, &p.Stack,
		&models, &verified, &p.UsageCount,
		&p.CreatedAt, &p.UpdatedAt, &lastUsed,
	)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal([]byte(tags), &p.Tags); err != nil {
		return nil, fmt.Errorf("parsing tags: %w", err)
	}
	if err := json.Unmarshal([]byte(models), &p.Models); err != nil {
		return nil, fmt.Errorf("parsing models: %w", err)
	}
	p.Verified = verified == 1
	if lastUsed.Valid {
		p.LastUsedAt = &lastUsed.Time
	}
	return &p, nil
}

func scanPrompts(ctx context.Context, rows *sql.Rows) ([]*model.Prompt, error) {
	var prompts []*model.Prompt
	for rows.Next() {
		// Check for context cancellation
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}

		var p model.Prompt
		var tags, models string
		var lastUsed sql.NullTime
		var verified int

		err := rows.Scan(
			&p.ID, &p.Title, &p.Content, &tags, &p.Stack,
			&models, &verified, &p.UsageCount,
			&p.CreatedAt, &p.UpdatedAt, &lastUsed,
		)
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal([]byte(tags), &p.Tags); err != nil {
			return nil, fmt.Errorf("parsing tags: %w", err)
		}
		if err := json.Unmarshal([]byte(models), &p.Models); err != nil {
			return nil, fmt.Errorf("parsing models: %w", err)
		}
		p.Verified = verified == 1
		if lastUsed.Valid {
			p.LastUsedAt = &lastUsed.Time
		}
		prompts = append(prompts, &p)
	}
	return prompts, rows.Err()
}

func scanTestResults(rows *sql.Rows) ([]*model.TestResult, error) {
	var results []*model.TestResult
	for rows.Next() {
		var r model.TestResult
		var passed int
		var errorMessage sql.NullString

		err := rows.Scan(
			&r.ID, &r.PromptID, &r.Model, &r.Input, &r.ExpectedOutput,
			&r.ActualOutput, &passed, &r.Score, &r.LatencyMs,
			&r.TokenUsage, &errorMessage, &r.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		r.Passed = passed == 1
		if errorMessage.Valid {
			r.ErrorMessage = errorMessage.String
		}
		results = append(results, &r)
	}
	return results, rows.Err()
}

func scanPromptVersion(row *sql.Row) (*model.PromptVersion, error) {
	var v model.PromptVersion
	var tags, models string
	var verified int
	var commitMsg, author sql.NullString

	err := row.Scan(
		&v.ID, &v.PromptID, &v.Version, &v.Title, &v.Content,
		&tags, &v.Stack, &models, &verified, &commitMsg, &author, &v.CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal([]byte(tags), &v.Tags); err != nil {
		return nil, fmt.Errorf("parsing tags: %w", err)
	}
	if err := json.Unmarshal([]byte(models), &v.Models); err != nil {
		return nil, fmt.Errorf("parsing models: %w", err)
	}
	v.Verified = verified == 1
	if commitMsg.Valid {
		v.CommitMsg = commitMsg.String
	}
	if author.Valid {
		v.Author = author.String
	}
	return &v, nil
}

func scanPromptVersions(rows *sql.Rows) ([]*model.PromptVersion, error) {
	var versions []*model.PromptVersion
	for rows.Next() {
		var v model.PromptVersion
		var tags, models string
		var verified int
		var commitMsg, author sql.NullString

		err := rows.Scan(
			&v.ID, &v.PromptID, &v.Version, &v.Title, &v.Content,
			&tags, &v.Stack, &models, &verified, &commitMsg, &author, &v.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal([]byte(tags), &v.Tags); err != nil {
			return nil, fmt.Errorf("parsing tags: %w", err)
		}
		if err := json.Unmarshal([]byte(models), &v.Models); err != nil {
			return nil, fmt.Errorf("parsing models: %w", err)
		}
		v.Verified = verified == 1
		if commitMsg.Valid {
			v.CommitMsg = commitMsg.String
		}
		if author.Valid {
			v.Author = author.String
		}
		versions = append(versions, &v)
	}
	return versions, rows.Err()
}
