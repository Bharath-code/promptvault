package db

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"github.com/Bharath-code/promptvault/internal/model"
)

type DB struct {
	conn *sql.DB
	path string
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
	conn, err := sql.Open("sqlite3", dbPath+"?_journal=WAL&_timeout=5000")
	if err != nil {
		return nil, fmt.Errorf("opening db: %w", err)
	}

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

	CREATE INDEX IF NOT EXISTS idx_prompts_stack ON prompts(stack);
	CREATE INDEX IF NOT EXISTS idx_prompts_updated ON prompts(updated_at DESC);
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
func (d *DB) Add(p *model.Prompt) error {
	if p.ID == "" {
		p.ID = uuid.New().String()
	}
	now := time.Now().UTC()
	p.CreatedAt = now
	p.UpdatedAt = now

	tags, _ := json.Marshal(p.Tags)
	models, _ := json.Marshal(p.Models)

	_, err := d.conn.Exec(`
		INSERT INTO prompts (id, title, content, tags, stack, models, verified, usage_count, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, 0, ?, ?)`,
		p.ID, p.Title, p.Content, string(tags), p.Stack, string(models),
		boolToInt(p.Verified), now, now,
	)
	return err
}

// Get returns a prompt by ID
func (d *DB) Get(id string) (*model.Prompt, error) {
	row := d.conn.QueryRow(`SELECT * FROM prompts WHERE id = ?`, id)
	return scanPrompt(row)
}

// List returns all prompts, optionally filtered by stack prefix
func (d *DB) List(stackFilter string) ([]*model.Prompt, error) {
	var rows *sql.Rows
	var err error

	if stackFilter != "" {
		rows, err = d.conn.Query(`
			SELECT * FROM prompts 
			WHERE stack = ? OR stack LIKE ?
			ORDER BY updated_at DESC`,
			stackFilter, stackFilter+"/%",
		)
	} else {
		rows, err = d.conn.Query(`SELECT * FROM prompts ORDER BY updated_at DESC`)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanPrompts(rows)
}

// Search does full-text search across title, content, tags, stack
func (d *DB) Search(query string) ([]*model.Prompt, error) {
	// Clean query for FTS5
	query = strings.ReplaceAll(query, `"`, `""`)

	rows, err := d.conn.Query(`
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
		rows, err = d.conn.Query(`
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
	return scanPrompts(rows)
}

// Update modifies an existing prompt
func (d *DB) Update(p *model.Prompt) error {
	p.UpdatedAt = time.Now().UTC()
	tags, _ := json.Marshal(p.Tags)
	models, _ := json.Marshal(p.Models)

	res, err := d.conn.Exec(`
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

// Delete removes a prompt
func (d *DB) Delete(id string) error {
	res, err := d.conn.Exec(`DELETE FROM prompts WHERE id = ?`, id)
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
func (d *DB) IncrementUsage(id string) error {
	now := time.Now().UTC()
	_, err := d.conn.Exec(`
		UPDATE prompts SET usage_count = usage_count + 1, last_used_at = ? WHERE id = ?`,
		now, id,
	)
	return err
}

// Stats returns aggregate statistics
func (d *DB) Stats() (total int, stacks int, err error) {
	err = d.conn.QueryRow(`SELECT COUNT(*) FROM prompts`).Scan(&total)
	if err != nil {
		return
	}
	err = d.conn.QueryRow(`SELECT COUNT(DISTINCT stack) FROM prompts WHERE stack != ''`).Scan(&stacks)
	return
}

// Count returns just the total number of prompts
func (d *DB) Count() (int, error) {
	var count int
	err := d.conn.QueryRow(`SELECT COUNT(*) FROM prompts`).Scan(&count)
	return count, err
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

	_ = json.Unmarshal([]byte(tags), &p.Tags)
	_ = json.Unmarshal([]byte(models), &p.Models)
	p.Verified = verified == 1
	if lastUsed.Valid {
		p.LastUsedAt = &lastUsed.Time
	}
	return &p, nil
}

func scanPrompts(rows *sql.Rows) ([]*model.Prompt, error) {
	var prompts []*model.Prompt
	for rows.Next() {
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

		_ = json.Unmarshal([]byte(tags), &p.Tags)
		_ = json.Unmarshal([]byte(models), &p.Models)
		p.Verified = verified == 1
		if lastUsed.Valid {
			p.LastUsedAt = &lastUsed.Time
		}
		prompts = append(prompts, &p)
	}
	return prompts, rows.Err()
}
