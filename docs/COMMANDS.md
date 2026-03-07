# 📚 PromptVault Complete Command Reference

**Version:** v1.2  
**Last Updated:** March 2026

---

## Table of Contents

- [Quick Start](#quick-start)
- [Core Commands](#core-commands)
  - [Managing Prompts](#managing-prompts)
  - [Testing & Quality](#testing--quality)
  - [Version Control](#version-control)
  - [Export & Integration](#export--integration)
  - [Sync & Backup](#sync--backup)
  - [Utilities](#utilities)
- [Advanced Features](#advanced-features)
- [Real-World Workflows](#real-world-workflows)
- [Troubleshooting](#troubleshooting)

---

## Quick Start

### Installation

```bash
# Install
go install github.com/Bharath-code/promptvault@latest

# Verify installation
promptvault --version

# Initialize with starter prompts
promptvault init
```

### First Steps

```bash
# Open interactive TUI
promptvault

# Add your first prompt
promptvault add "React Hook Debugging" \
  --stack frontend/react/hooks \
  --content "Analyze this React component and fix useEffect dependencies..." \
  --tags "debugging,react,hooks" \
  --models "claude-sonnet,gpt-4o"

# Search for prompts
promptvault search "useEffect"

# Export to Cursor
promptvault export --format cursorrules --stack frontend/react > .cursorrules
```

---

## Core Commands

### Managing Prompts

#### `promptvault add` - Add a New Prompt

**Aliases:** None  
**Use Case:** Store a new prompt in your vault

```bash
# Basic usage
promptvault add "Prompt Title" --content "Prompt content here"

# With all options
promptvault add "React Component Review" \
  --stack frontend/react \
  --tags "review,quality,best-practices" \
  --models "claude-sonnet,gpt-4o" \
  --verified \
  --content "Review this React component for best practices..."

# Preview before adding
promptvault add "My Prompt" --content "..." --preview

# From stdin
cat prompt.txt | promptvault add "From File" --stack backend/python

# Auto-detect stack from current directory
cd my-react-project
promptvault add "Component" --content "..."
# → Auto-detects: frontend/react
```

**Flags:**
| Flag | Shorthand | Description |
|------|-----------|-------------|
| `--content` | `-c` | Prompt content (or pipe via stdin) |
| `--stack` | `-s` | Tech stack path |
| `--tags` | — | Comma-separated tags |
| `--models` | — | Comma-separated model names |
| `--verified` | — | Mark as verified |
| `--preview` | — | Preview before adding |

---

#### `promptvault list` - List All Prompts

**Aliases:** `ls`, `show`, `list-all`  
**Use Case:** View all prompts or filter by stack

```bash
# List all prompts
promptvault list

# List with aliases
promptvault ls
promptvault show

# Filter by stack
promptvault list --stack frontend/react
promptvault ls -s backend/python

# Short output (ID + title only)
promptvault list --short

# JSON output for scripting
promptvault list --json
promptvault ls --json | jq '.[] | select(.stack | contains("react"))'

# Verbose mode
promptvault list -v

# Debug mode
promptvault list -vd
```

**Flags:**
| Flag | Shorthand | Description |
|------|-----------|-------------|
| `--stack` | `-s` | Filter by stack |
| `--short` | — | Short output |
| `--json` | — | JSON output |

---

#### `promptvault search` - Full-Text Search

**Aliases:** `find`, `query`  
**Use Case:** Find prompts by content, title, tags, or stack

```bash
# Basic search
promptvault search "useEffect dependencies"

# With aliases
promptvault find "react hooks"
promptvault query "typescript types"

# JSON output
promptvault search "api" --json | jq 'length'

# Verbose mode
promptvault search "docker" -v

# Combine with jq for filtering
promptvault search "test" --json | \
  jq '.[] | select(.usage_count > 5) | .title'
```

**Flags:**
| Flag | Shorthand | Description |
|------|-----------|-------------|
| `--json` | — | JSON output |

---

#### `promptvault get` - Get and Copy Prompt

**Aliases:** `fetch`  
**Use Case:** Retrieve a specific prompt and copy to clipboard

```bash
# Get by title or ID (copies to clipboard)
promptvault get "useEffect"

# Fetch alias
promptvault fetch "React Hooks"

# Print to stdout without copying
promptvault get "API Design" --print

# Copy and print
promptvault get "TypeScript" --copy --print
```

**Flags:**
| Flag | Shorthand | Description |
|------|-----------|-------------|
| `--copy` | `-c` | Copy to clipboard (default: true) |
| `--print` | `-p` | Print to stdout |

---

#### `promptvault delete` - Delete a Prompt

**Aliases:** `rm`, `remove`, `del`  
**Use Case:** Remove unwanted prompts

```bash
# Delete with confirmation
promptvault delete "Old Prompt"

# Force delete (skip confirmation)
promptvault delete "Old Prompt" --force
promptvault rm "Deprecated" -f

# Delete by ID
promptvault delete abc123 --force
```

**Flags:**
| Flag | Shorthand | Description |
|------|-----------|-------------|
| `--force` | `-f` | Skip confirmation |

---

#### `promptvault create` - Create with AI Assistance

**Aliases:** None  
**Use Case:** Interactive prompt creation with AI suggestions

```bash
# Standard interactive creation
promptvault create

# AI-assisted creation (requires API key)
promptvault create --ai

# Set API key first
export ANTHROPIC_API_KEY=sk-ant-...
promptvault create --ai
```

**Features:**
- ✅ Variable detection (`{{variable}}`)
- ✅ Tag recommendations
- ✅ Stack auto-detection
- ✅ Anti-pattern warnings
- ✅ Quality scoring (0-100)

**Flags:**
| Flag | Description |
|------|-------------|
| `--ai` | Use AI-assisted creation |

---

### Testing & Quality

#### `promptvault test` - Test Prompts

**Aliases:** None  
**Use Case:** Validate prompts against expected outputs

```bash
# Interactive test session
promptvault test abc123

# Single test
promptvault test abc123 \
  --input "Convert to TypeScript" \
  --expected "TypeScript with proper types" \
  --model claude-sonnet

# View test history
promptvault test abc123 --history

# Run all saved tests
promptvault test abc123 --all
```

**Example Session:**
```bash
$ promptvault test abc123

📋 Testing prompt: React Hook Converter
   Model: claude-sonnet

Enter test input (empty line to skip, 'quit' to exit):
────────────────────────────────────────────────────────────

💬 Input: Convert this useState to useReducer
✨ Expected output: Should use useReducer with dispatch

🔍 Running test...
✅ Test PASSED

Score: 94.2/100
Latency: 234ms
Tokens: 156
```

**Flags:**
| Flag | Description |
|------|-------------|
| `--input` | Test input |
| `--expected` | Expected output |
| `--model` | Model to test against |
| `--all` | Run all saved tests |
| `--history` | Show test history |

---

#### `promptvault audit` - Audit for Decay

**Aliases:** None  
**Use Case:** Find problematic prompts (deprecated models, unused, low quality)

```bash
# Full audit
promptvault audit

# Critical issues only
promptvault audit --severity critical

# Warning and above
promptvault audit --severity warning

# JSON output for CI/CD
promptvault audit --json | jq '.issues_found'

# Filter by severity in CI
if [ $(promptvault audit --severity critical --json | jq '.issues_found') -gt 0 ]; then
  echo "❌ Critical issues found"
  exit 1
fi
```

**Example Output:**
```bash
$ promptvault audit

🔍 PromptVault Audit Report
────────────────────────────────────────────────────────────
Generated: 2026-03-07 07:33:51

📊 Summary:
   Total Prompts:   156
   Healthy:         142 (91.0%)
   Issues Found:    14

💡 Recommendations:
   🔴 Critical: Update 3 prompts using deprecated models
   🔴 Critical: Fix 2 prompts with low test success rates
   🟡 Warning: Review 7 unused prompts
   🟢 Info: Update 2 outdated prompts
```

**Detects:**
- 🔴 **Deprecated Models** (gpt-3.5-turbo, claude-2, etc.)
- 🔴 **Low Success Rate** (< 50% test pass rate)
- 🟡 **Unused Prompts** (90+ days)
- 🟢 **Old Versions** (180+ days)

**Flags:**
| Flag | Description |
|------|-------------|
| `--json` | JSON output |
| `--severity` | Filter by severity (critical/warning/info) |

---

### Version Control

#### `promptvault history` - View Version History

**Aliases:** None  
**Use Case:** See all versions of a prompt

```bash
# View history
promptvault history abc123

# Example output:
📜 Version History: React Hook Converter
────────────────────────────────────────────────────────────
▶ v5  2026-03-06 14:30  johndoe  Fixed edge case with useEffect
  v4  2026-03-06 11:15  johndoe  Added TypeScript support
  v3  2026-03-05 16:45  janedoe  Improved error handling
  v2  2026-03-05 10:00  johndoe  Updated for React 19
  v1  2026-03-04 09:00  johndoe  Initial version

Total versions: 5
```

---

#### `promptvault diff` - Compare Versions

**Aliases:** None  
**Use Case:** See differences between versions

```bash
# Compare specific versions
promptvault diff abc123 1 2

# Compare with current
promptvault diff abc123 3 current
promptvault diff abc123 3 HEAD

# Compare last two versions
promptvault diff abc123 HEAD~1 HEAD

# Example output:
📊 Diff: abc123 (v1 → v2)
────────────────────────────────────────────────────────────

ℹ  Title:
  - React Hook
  + React Hook Converter

ℹ  Content:
  Convert this to React
- function useState() {
+ function useState<T>(initial: T) {
```

**Version Specifiers:**
- `1`, `2`, `3` - Specific version numbers
- `current`, `HEAD`, `latest` - Current version
- `HEAD~1` - Previous version

---

#### `promptvault revert` - Revert to Previous Version

**Aliases:** None  
**Use Case:** Restore a prompt to an earlier version

```bash
# Revert to version 3
promptvault revert abc123 3

# With custom message
promptvault revert abc123 3 --message "Reverting broken changes"
promptvault revert abc123 3 -m "Undo last changes"

# Example output:
✓ Reverted abc123 to v3
New version: v6
```

**Flags:**
| Flag | Shorthand | Description |
|------|-----------|-------------|
| `--message` | `-m` | Commit message |

---

### Export & Integration

#### `promptvault export` - Export Prompts

**Aliases:** `exp`  
**Use Case:** Export prompts to various formats

```bash
# Export all prompts to SKILL.md
promptvault export --format skill.md --output SKILL.md

# Export specific stack
promptvault export --format cursorrules --stack frontend/react

# Export to AGENTS.md
promptvault export --format agents.md > AGENTS.md

# Export to .windsurfrules
promptvault export --format windsurf > .windsurfrules

# Export as JSON
promptvault export --format json --output prompts.json

# Verbose mode
promptvault exp --format skill.md -o SKILL.md -v
```

**Supported Formats:**
| Format | Command | Use Case |
|--------|---------|----------|
| `skill.md` | `--format skill.md` | Claude Code skills |
| `agents.md` | `--format agents.md` | OpenAI Agents |
| `claude.md` | `--format claude.md` | CLAUDE.md snippets |
| `cursorrules` | `--format cursorrules` | Cursor IDE |
| `windsurf` | `--format windsurf` | Windsurf IDE |
| `markdown` | `--format markdown` | Documentation |
| `json` | `--format json` | Integrations |
| `text` | `--format text` | Plain text |

**Flags:**
| Flag | Shorthand | Description |
|------|-----------|-------------|
| `--format` | `-f` | Export format |
| `--stack` | `-s` | Filter by stack |
| `--output` | `-o` | Output file |

---

#### `promptvault watch` - Auto-Export on Changes

**Aliases:** None  
**Use Case:** Continuously export when database changes

```bash
# Basic watch
promptvault watch --format skill.md --output SKILL.md

# Custom interval
promptvault watch --format cursorrules --output .cursorrules --interval 2s

# Filter by stack
promptvault watch --format skill.md --output REACT.md --stack frontend/react

# Verbose mode
promptvault watch -f skill.md -o SKILL.md -v

# Example output:
ℹ  Performing initial export...
✓ Watching for changes... (interval: 5s)
ℹ  Export format: skill.md → SKILL.md
ℹ  Press Ctrl+C to stop

ℹ  Detected change (#1), exporting...
✓ Exported SKILL.md
```

**Flags:**
| Flag | Shorthand | Description |
|------|-----------|-------------|
| `--format` | `-f` | Export format |
| `--output` | `-o` | Output file (required) |
| `--stack` | `-s` | Filter by stack |
| `--interval` | — | Check interval (default: 5s) |

---

#### `promptvault import` - Import Prompts

**Aliases:** `imp`  
**Use Case:** Import prompts from JSON

```bash
# Import from JSON
promptvault import prompts.json
promptvault imp shared-prompts.json

# Import and verify
promptvault import backup.json
# Shows skipped duplicates
```

---

### Sync & Backup

#### `promptvault sync push` - Backup to GitHub Gist

**Aliases:** None  
**Use Case:** Backup vault to private GitHub Gist

```bash
# Backup to Gist
promptvault sync push --token <gh_token>

# Token from environment
export PROMPTVAULT_GITHUB_TOKEN=ghp_...
promptvault sync push
```

**Required:** GitHub Personal Access Token with `gist` scope

---

#### `promptvault sync pull` - Restore from Gist

**Aliases:** None  
**Use Case:** Restore vault from GitHub Gist

```bash
# Restore from Gist
promptvault sync pull --token <gh_token>

# Token from environment
export PROMPTVAULT_GITHUB_TOKEN=ghp_...
promptvault sync pull
```

---

### Utilities

#### `promptvault stats` - Show Statistics

**Aliases:** `statistics`  
**Use Case:** View vault statistics

```bash
# Basic stats
promptvault stats
promptvault statistics

# Example output:
⚡ PromptVault Statistics
────────────────────────────────────────────────────────────
  Total Prompts:        156
  Unique Stacks:        14
  Database Path:        /Users/user/.promptvault/vault.db
```

---

#### `promptvault stacks` - List Available Stacks

**Aliases:** None  
**Use Case:** View all available tech stacks

```bash
# List all stacks
promptvault stacks

# Example output:
Available stack paths:
  frontend/react/hooks
  frontend/react/performance
  backend/python/django
  backend/python/fastapi
  devops/docker
  devops/kubernetes
  ...
```

---

#### `promptvault mcp` - Start MCP Server

**Aliases:** None  
**Use Case:** Run MCP server for AI integration

```bash
# Start MCP server
promptvault mcp

# Integrates with:
# - Claude Desktop
# - Cursor
# - Windsurf
```

---

#### `promptvault completion` - Generate Shell Completion

**Aliases:** None  
**Use Case:** Enable tab completion

```bash
# Bash
promptvault completion bash > ~/.bash_completion
source ~/.bash_completion

# Zsh
promptvault completion zsh > "${fpath[1]}/_promptvault"

# Fish
promptvault completion fish > ~/.config/fish/completions/promptvault.fish

# PowerShell
promptvault completion powershell | Out-String | Invoke-Expression
```

---

#### `promptvault init` - Initialize Vault

**Aliases:** None  
**Use Case:** Add curated starter prompts

```bash
# Initialize with starter prompts
promptvault init

# Force add even if vault exists
promptvault init --force
```

---

## Advanced Features

### Global Flags

These flags work with all commands:

```bash
# Verbose output
promptvault list -v
promptvault search "react" --verbose

# Debug output (with timestamps)
promptvault list -vd
promptvault audit -d

# Help
promptvault --help
promptvault list --help
```

| Flag | Shorthand | Description |
|------|-----------|-------------|
| `--verbose` | `-v` | Verbose output |
| `--debug` | `-d` | Debug output with timestamps |
| `--help` | `-h` | Show help |

---

### Smart Defaults

PromptVault automatically detects context:

```bash
# Auto-detect stack from project type
cd my-react-project
promptvault add "Component" --content "..."
# → Auto-detects: frontend/react

cd my-python-api
promptvault add "Endpoint" --content "..."
# → Auto-detects: backend/python/fastapi

# Auto-tag with Git branch
cd my-project
git checkout feature/new-auth
promptvault add "Auth Prompt" --content "..."
# → Auto-tags: git:feature/new-auth
```

---

### JSON Output for Scripting

```bash
# Count prompts
promptvault list --json | jq 'length'

# Filter by stack
promptvault list --json | \
  jq '.[] | select(.stack | contains("react")) | .title'

# Find unused prompts
promptvault list --json | \
  jq '.[] | select(.usage_count == 0) | .title'

# Export titles only
promptvault list --json | jq -r '.[].title'

# Count by stack
promptvault list --json | \
  jq 'group_by(.stack) | map({stack: .[0].stack, count: length})'
```

---

### CI/CD Integration

```yaml
# .github/workflows/prompt-audit.yml
name: Prompt Audit

on:
  schedule:
    - cron: '0 9 * * 1'  # Every Monday at 9 AM

jobs:
  audit:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Install PromptVault
        run: go install github.com/Bharath-code/promptvault@latest
      
      - name: Run Audit
        run: promptvault audit --json > audit-report.json
      
      - name: Check for Critical Issues
        run: |
          ISSUES=$(cat audit-report.json | jq '.issues_found')
          if [ "$ISSUES" -gt 0 ]; then
            echo "❌ Found $ISSUES prompt issues"
            exit 1
          fi
```

---

## Real-World Workflows

### Workflow 1: Daily Prompt Development

```bash
# 1. Start day - check what needs attention
promptvault audit --severity critical

# 2. Create new prompt with AI assistance
promptvault create --ai

# 3. Test the prompt
promptvault test abc123 --input "test input" --expected "expected output"

# 4. If tests fail, edit and retest
promptvault  # Open TUI, press 'e' to edit

# 5. Export to team files
promptvault watch --format cursorrules --output .cursorrules &

# 6. End day - version control
promptvault history abc123  # Review changes
```

---

### Workflow 2: Team Collaboration

```bash
# 1. Export team prompts
promptvault export --stack frontend/react --format skill.md --output TEAM.md

# 2. Share with team
git add TEAM.md
git commit -m "Update React prompts"
git push

# 3. Team members import
promptvault import TEAM.md

# 4. Set up auto-sync
promptvault watch --format skill.md --output SKILL.md &
```

---

### Workflow 3: Quality Assurance

```bash
# 1. Weekly audit
promptvault audit --json > weekly-audit.json

# 2. Fix deprecated models
promptvault list --json | \
  jq -r '.[] | select(.models[] | contains("gpt-3.5")) | .id' | \
  xargs -I {} promptvault edit {}

# 3. Test critical prompts
for id in $(promptvault list --json | jq -r '.[] | select(.verified) | .id'); do
  promptvault test $id --all
done

# 4. Remove unused prompts
promptvault audit --severity warning | grep "unused" | \
  awk '{print $NF}' | xargs -I {} promptvault delete {} --force
```

---

### Workflow 4: Migration from Other Tools

```bash
# From Notion/Cursor/ChatGPT history
# 1. Export to JSON format
# 2. Import to PromptVault
promptvault import migration.json

# 3. Verify import
promptvault list --json | jq 'length'

# 4. Organize by stack
promptvault list --json | jq -r '.[] | .id' | \
  while read id; do
    # Add stack based on content analysis
    promptvault edit $id --stack auto-detected
  done
```

---

## Troubleshooting

### Common Issues

#### "Command not found"
```bash
# Check installation
which promptvault

# Reinstall
go install github.com/Bharath-code/promptvault@latest

# Add to PATH
export PATH=$PATH:$(go env GOPATH)/bin
```

#### "Database locked"
```bash
# Close other PromptVault instances
pkill -f promptvault

# Remove WAL files
rm ~/.promptvault/vault.db-wal
rm ~/.promptvault/vault.db-shm
```

#### "API key not set"
```bash
# Set API keys
export ANTHROPIC_API_KEY=sk-ant-...
export OPENAI_API_KEY=sk-...

# Add to shell profile
echo 'export ANTHROPIC_API_KEY=sk-ant-...' >> ~/.zshrc
source ~/.zshrc
```

#### "No prompts found"
```bash
# Initialize with starter prompts
promptvault init

# Or add first prompt
promptvault add "My First Prompt" --content "..."
```

#### "Shell completion not working"
```bash
# Bash
promptvault completion bash > ~/.bash_completion
source ~/.bash_completion

# Verify
complete -p promptvault

# Zsh
promptvault completion zsh > "${fpath[1]}/_promptvault"
```

---

### Getting Help

```bash
# General help
promptvault --help

# Command-specific help
promptvault list --help
promptvault test --help

# Verbose debugging
promptvault list -vd

# Check version
promptvault --version
```

---

### Report Issues

- **GitHub Issues:** https://github.com/Bharath-code/promptvault/issues
- **Discussions:** https://github.com/Bharath-code/promptvault/discussions
- **Documentation:** https://github.com/Bharath-code/promptvault/docs

---

## Quick Reference Card

```
┌─────────────────────────────────────────────────────────────┐
│                    PROMPTVAULT QUICK REF                    │
├─────────────────────────────────────────────────────────────┤
│ MANAGE PROMPTS                                              │
│   add       Create new prompt                               │
│   list      List all prompts (ls, show)                     │
│   search    Find prompts (find, query)                      │
│   get       Get & copy (fetch)                              │
│   delete    Remove prompt (rm, remove)                      │
│   create    AI-assisted creation                            │
├─────────────────────────────────────────────────────────────┤
│ TEST & QUALITY                                              │
│   test      Test prompt                                     │
│   audit     Find issues                                     │
├─────────────────────────────────────────────────────────────┤
│ VERSION CONTROL                                             │
│   history   View versions                                   │
│   diff      Compare versions                                │
│   revert    Restore version                                 │
├─────────────────────────────────────────────────────────────┤
│ EXPORT & SYNC                                               │
│   export    Export formats (exp)                            │
│   watch     Auto-export                                     │
│   import    Import JSON (imp)                               │
│   sync      Backup/restore                                  │
├─────────────────────────────────────────────────────────────┤
│ UTILITIES                                                   │
│   stats     Statistics                                      │
│   stacks    List stacks                                     │
│   mcp       MCP server                                      │
│   init      Initialize                                      │
│   completion Shell completion                               │
├─────────────────────────────────────────────────────────────┤
│ FLAGS                                                       │
│   -v        Verbose                                         │
│   -d        Debug                                           │
│   --json    JSON output                                     │
│   --preview Preview before add                              │
└─────────────────────────────────────────────────────────────┘
```

---

**Happy Prompting! 🚀**

For more help: `promptvault --help` or visit https://github.com/Bharath-code/promptvault
