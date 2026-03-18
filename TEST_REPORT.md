# PromptVault Test Report

## Build Configuration

- **Build tag required**: `fts5` (SQLite FTS5 module)
- **Build command**: `go build -tags fts5 .`
- **Test command**: `go test -tags fts5 ./...`
- **Vet**: `go vet ./...` (clean)

## Bugs Found & Fixed

### 1. Panic: close of closed channel in `showLoading` (`root.go:268`)

**Severity**: High (crashes on `--force` init when prompts skip)

**Cause**: The `stopLoading` closure calls `close(done)` each time it's invoked. During `init --force`, prompts that fail (FTS5 unavailable) call `stopLoading()` inside the loop, then again after the loop — double-closing the channel.

**Fix**: Added `sync.Once` to make `stopLoading` idempotent:

```go
var once sync.Once
return func() {
    once.Do(func() { close(done) })
    wg.Wait()
}
```

### 2. Duplicate `history` command conflict (`root.go`)

**Severity**: Medium (search history command shadowed by version history)

**Cause**: `searchHistoryCmd` used `Use: "history [subcommand]"` which conflicted with `historyCmd` (`Use: "history"`). Both registered to root command — version history won.

**Fix**: Renamed to `Use: "search-history [subcommand]"` so commands are distinct.

### 3. Missing FTS5 build tags on test files

**Severity**: Low (tests fail on non-FTS5 environments without clear indication)

**Cause**: `version_test.go` and `detector_test.go` lacked `//go:build fts5` tags, causing them to attempt DB operations requiring FTS5 even without the build tag.

**Fix**: Added `//go:build fts5` to both files.

## Test Results

### CLI Commands — All PASS ✅

| Command | Flags | Status |
|---------|-------|--------|
| `init` | `--force` | ✅ Adds 15 curated prompts |
| `add` | `[title]` + stdin | ✅ Reads from stdin, auto-detects stack |
| `add` | `--title`, `--content`, `--stack`, `--tags`, `--models`, `--verified` | ✅ All flags work |
| `add` | `--preview` | ✅ Shows box preview, prompts for confirmation |
| `list` | (none) | ✅ Lists all prompts with ID, title, stack |
| `list` | `--short` | ✅ Compact ID + title output |
| `list` | `--json` | ✅ Valid JSON array |
| `list` | `--stack frontend/react` | ✅ Filters correctly |
| `get` | `[id-or-title]` | ✅ Copies to clipboard, increments usage |
| `get` | `--print` | ✅ Prints to stdout |
| `search` | `[query]` | ✅ FTS5 search returns matching prompts |
| `search` | `--json` | ✅ JSON output |
| `search` | (no results) | ✅ Shows warning |
| `delete` | `[id] --force` | ✅ Deletes without confirmation |
| `export` | `--format text` | ✅ Plain text output |
| `export` | `--format json` | ✅ Valid JSON |
| `export` | `--format markdown` | ✅ Markdown with headers |
| `export` | `--format skill.md` | ✅ Claude Code SKILL.md format |
| `export` | `--format bulk --output ./dir/` | ✅ Creates directory + per-prompt files |
| `export` | `--id <full-uuid>` | ✅ Exports specific prompt |
| `export` | `--stack frontend/react` | ✅ Filters by stack |
| `import` | `[file] --dry-run` (JSON) | ✅ Shows preview without importing |
| `import` | `[file]` (JSON) | ✅ Imports successfully |
| `import` | `[file] --dry-run` (Markdown) | ✅ Parses ## headers as titles |
| `stats` | | ✅ Shows total prompts, unique stacks, DB path |
| `stacks` | | ✅ Lists all 57 stack paths |
| `history` | `[prompt-id]` | ✅ Shows version history |
| `diff` | `[id] [v1] [v2]` | ✅ Shows diff between versions |
| `revert` | `[id] [version]` | ✅ Reverts to specified version |
| `config show` | | ✅ Shows current config (theme, autocopy, etc.) |
| `config theme` | | ✅ Lists available themes |
| `config theme monokai` | | ✅ Sets theme, persists to config.json |
| `config set autocopy false` | | ✅ Updates config |
| `config set previewlines 20` | | ✅ Updates config |
| `config keybindings` | | ✅ Shows navigation/actions/quick actions |
| `config reset` | | ✅ Resets to defaults |
| `search-history list` | | ✅ Shows recent searches |
| `search-history clear` | | ✅ Clears search history |
| `completion bash` | | ✅ Generates bash completions |
| `completion zsh` | | ✅ Generates zsh completions |
| `completion fish` | | ✅ Generates fish completions |
| `completion powershell` | | ✅ Generates powershell completions |
| `audit` | | ✅ Audits prompts for decay |
| `create` | | ✅ AI-assisted prompt creation |
| `test` | `[prompt-id]` | ✅ Tests prompts against models |
| `watch` | `--format --output --stack --interval` | ✅ Watches for DB changes |

### Unit Tests — All PASS ✅

| Package | With FTS5 | Without FTS5 |
|---------|-----------|--------------|
| `internal/db` | ✅ PASS | ⏭️ Skipped (build tag) |
| `internal/decay` | ✅ PASS | ⏭️ Skipped (build tag) |
| `internal/export` | ✅ PASS | ✅ PASS (no FTS5 dep) |
| `internal/model` | ✅ PASS | ✅ PASS (no FTS5 dep) |

### `go vet` — Clean ✅

No issues found.

## Build Status

- `go build -tags fts5 .` ✅
- `go test -tags fts5 ./...` ✅
- `go vet ./...` ✅

## Known Environment Notes

1. **FTS5 Required**: Build and tests require `go build -tags fts5` on this environment (macOS with CGO-enabled SQLite). Without the tag, prompts fail to add with `no such module: fts5`.

2. **No `PROMPTVAULT_DB` env var**: The database path is always `~/.promptvault/vault.db` (hardcoded in `db.go:dataDir()`). There is no environment variable override.

3. **Database auto-creates**: Running any command creates `~/.promptvault/` and `vault.db` automatically.

## Command Summary

```
promptvault              # Open interactive TUI
promptvault add          # Add prompt (stdin, --title, --content, --stack, --tags, --models, --verified, --preview)
promptvault list         # List prompts (--short, --json, --stack)
promptvault get          # Get prompt, copy to clipboard (--print, --copy)
promptvault search       # Full-text search (--json)
promptvault delete       # Delete prompt (--force)
promptvault export       # Export (--format, --stack, --output, --id)
promptvault import       # Import JSON/Markdown (--dry-run)
promptvault init         # Seed with 15 curated prompts (--force)
promptvault stats        # Show statistics
promptvault stacks       # List stack taxonomies
promptvault history      # Version history (requires prompt-id)
promptvault diff         # Compare versions
promptvault revert       # Revert to version
promptvault config       # Config (show, theme, set, keybindings, reset, export, import)
promptvault search-history # Search history (list, clear)
promptvault audit        # Prompt decay audit
promptvault create       # AI-assisted prompt creation
promptvault test         # Test prompts against models
promptvault watch        # Auto-export on DB changes
promptvault mcp          # MCP server over stdio
promptvault sync         # GitHub Gist sync (push, pull)
promptvault completion   # Shell completions (bash, zsh, fish, powershell)
```
