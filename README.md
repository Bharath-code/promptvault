# ⚡ PromptVault

> The universal prompt OS for developers — store, search, and deploy AI prompts by tech stack, right from your terminal.

[![Downloads](https://img.shields.io/github/downloads/Bharath-code/promptvault/total?style=flat-square&color=7C3AED)](https://github.com/Bharath-code/promptvault/releases)
[![Release](https://img.shields.io/github/v/release/Bharath-code/promptvault?style=flat-square)](https://github.com/Bharath-code/promptvault/releases)
[![Go Version](https://img.shields.io/badge/go-1.23+-blue?style=flat-square)](https://go.dev)
[![License](https://img.shields.io/badge/license-MIT-green?style=flat-square)](LICENSE)

---

**You spend 45 minutes crafting the perfect prompt. You close the tab. It's gone forever.**

PromptVault fixes this. It's a CLI + TUI that keeps every prompt you've ever written organized by tech stack, searchable in milliseconds, and exportable to any AI tool you use.

---

## Features

- **🗂 Tech-stack taxonomy** — Organize by `frontend/react/hooks`, `backend/python/fastapi`, `devops/terraform`, and 80+ more
- **⚡ Fuzzy search** — Find any prompt in under 3 seconds with full-text search (FTS5 + porter stemming)
- **📋 One-key copy** — Press Enter to copy any prompt to clipboard instantly
- **🔄 Multi-format export** — Export to `SKILL.md`, `AGENTS.md`, `CLAUDE.md`, `.cursorrules`, `.windsurfrules`
- **🤖 Model tagging** — Mark prompts as verified for Claude, GPT-4o, Gemini, and more
- **💻 Beautiful TUI** — Built with Bubble Tea, works in any terminal
- **🗄 SQLite + FTS** — Zero-dependency local storage with full-text search
- **🔌 MCP Server** — Connect directly to Claude Desktop, Cursor, and Windsurf via Model Context Protocol
- **☁️ Cloud Sync** — Backup and sync your vault across machines using private GitHub Gists
- **✨ Markdown Previews** — Beautifully rendered markdown and syntax highlighting inside the terminal
- **✏️ Interactive Variables** — Define `{{variables}}` in prompts; the TUI prompts you to fill them before copying
- **🪄 Smart Auto-Injection** — Automatically append exported rules to files like `.cursorrules` and `SKILL.md` directly in your workspace
- **🛡️ OSS CI/CD & Linting** — Mature codebase maintained with `golangci-lint`, comprehensive workflows, and Dependabot
- **📦 Single binary** — No runtime, no Node, no Docker. Just download and run.
- **🚀 One-command init** — Start with 15+ curated, production-grade prompts

### 🎯 New DX Features (v1.1+)

- **💡 Smart Error Messages** — Actionable suggestions when commands fail
- **🐚 Shell Completion** — Auto-completion for bash, zsh, fish, and PowerShell
- **📄 JSON Output** — Scriptable output with `--json` flag for automation
- **🔍 Verbose/Debug Mode** — Detailed logging with `-v` and `-vd` flags
- **⚡ Command Aliases** — Quick commands like `ls`, `rm`, `find`, `exp`
- **🎨 Rich Colors** — Beautiful color-coded output throughout the CLI
- **🎯 Smart Defaults** — Auto-detects project type (React, Go, Python, Terraform)
- **👁️ Preview Mode** — Preview prompts before adding with `--preview`
- **🏷️ Git Integration** — Auto-tags prompts with current Git branch

### 🧪 v1.2: Professional Prompt Engineering

- **🧪 Prompt Testing** — Test prompts against expected outputs with `promptvault test`
- **📜 Version History** — Git-like versioning with `promptvault history`, `diff`, `revert`
- **🤖 AI-Assisted Authoring** — Smart suggestions, variable detection, quality scoring
- **🔍 Decay Detection** — Audit prompts for issues with `promptvault audit`
- **⏱️ Auto-Export Watch** — Watch mode for continuous export with `promptvault watch`

### 🎨 v1.3: Enhanced TUI Experience (NEW!)

- **🔍 Fuzzy Search** — Type anything, get relevant results with match scores
- **❓ Quick Action Menu** — Press `?` for instant keybinding reference
- **📊 Stats Dashboard** — Press `s` for usage statistics and top prompts
- **🔥 Recent Prompts** — Press `R` to toggle frequently used prompts
- **☑️ Multi-Select** — Press `Space` to select, `x` for batch operations
- **🎬 Full-Screen Preview** — Press `v` for immersive full-screen preview
- **⚡ 40x Faster Load** — Optimized performance, loads in ~300ms

---

## Install

**Go install:**
```sh
go install github.com/Bharath-code/promptvault@latest
```

**Build from source:**
```sh
git clone https://github.com/Bharath-code/promptvault
cd promptvault
make build
```

---

## Quick Start

```sh
# Initialize with 15+ curated starter prompts
promptvault init

# Open the interactive TUI (recommended)
promptvault

# Add a prompt from the command line
promptvault add "Fix React useEffect deps" \
  --stack frontend/react/hooks \
  --models "claude-sonnet,gpt-4o" \
  --tags "debugging,hooks"

# Add with preview (see before committing)
promptvault add "My Prompt" --content "Prompt content..." --preview

# Add a prompt from stdin
cat my-prompt.txt | promptvault add "My prompt" --stack backend/python

# Smart defaults: Auto-detects stack from current directory
cd my-react-project
promptvault add "React Hook" --content "useEffect example..."
# → Auto-detects stack: frontend/react

# Search and copy to clipboard
promptvault get "useEffect"

# List all prompts
promptvault list
promptvault ls              # alias

# List as JSON for scripting
promptvault list --json | jq '.[] | .title'

# List by stack
promptvault list --stack frontend/react

# Full-text search
promptvault search "typescript generics"
promptvault find "react"    # alias

# Export to Claude Code SKILL.md
promptvault export --format skill.md --output SKILL.md

# Export to Cursor rules
promptvault export --format cursorrules --stack frontend/react > .cursorrules

# Export to AGENTS.md
promptvault export --format agents.md > AGENTS.md

# Import prompts from JSON
promptvault import shared-prompts.json
promptvault imp prompts.json  # alias

# Start an MCP Server over stdio
promptvault mcp

# Backup prompts to a private GitHub Gist
promptvault sync push --token <gh_token>

# Restore prompts from GitHub Gist
promptvault sync pull

# Watch for changes and auto-export (NEW!)
promptvault watch --format skill.md --output SKILL.md
promptvault watch --format cursorrules --output .cursorrules --interval 2s

# Test prompts before deploying
promptvault test abc123
promptvault test abc123 --input "test" --expected "output"
promptvault test abc123 --history

# View and manage version history
promptvault history abc123
promptvault diff abc123 1 2
promptvault revert abc123 3

# Create prompts with AI assistance
promptvault create --ai
promptvault create

# Audit for decay and quality issues
promptvault audit
promptvault audit --severity critical
promptvault audit --json

# Show stats with beautiful formatting
promptvault stats
promptvault statistics        # alias

# Debug mode for troubleshooting
promptvault list -v           # verbose
promptvault list -vd          # verbose + debug with timestamps

# Generate shell completion
promptvault completion bash > ~/.bash_completion
source ~/.bash_completion
```

---

## TUI Keybindings

### Navigation
| Key | Action |
|-----|--------|
| `↑` / `↓` or `k` / `j` | Navigate prompts |
| `Enter` | Fill variables (if any) and copy to clipboard |
| `Space` | Select/deselect (multi-select mode) |
| `/` | Search |

### Actions
| Key | Action |
|-----|--------|
| `a` | Add new prompt |
| `e` | Edit selected |
| `d` | Delete selected |
| `v` | Toggle full-screen preview |
| `r` | Refresh list |
| `R` | Toggle recent prompts |
| `s` | Show statistics |
| `x` | Batch process (when items selected) |
| `?` | Quick action menu |
| `Esc` | Clear filter / go back |
| `q` | Quit |

---

## CLI Commands & Aliases

### Core Commands
| Command | Aliases | Description |
|---------|---------|-------------|
| `list` | `ls`, `show`, `list-all` | List all prompts |
| `search` | `find`, `query` | Full-text search |
| `delete` | `rm`, `remove`, `del` | Delete a prompt |
| `get` | `fetch` | Get and copy to clipboard |
| `export` | `exp` | Export to various formats |
| `import` | `imp` | Import from JSON |
| `stats` | `statistics` | Show vault statistics |
| `watch` | — | **NEW!** Auto-export on changes |
| `test` | — | **NEW!** Test prompts |
| `history` | — | **NEW!** View version history |
| `diff` | — | **NEW!** Compare versions |
| `revert` | — | **NEW!** Revert to version |
| `create` | — | **NEW!** Create with AI assist |
| `audit` | — | **NEW!** Audit for decay |

### Global Flags
| Flag | Description |
|------|-------------|
| `-v, --verbose` | Enable verbose output |
| `-d, --debug` | Enable debug output with timestamps |
| `-h, --help` | Show help for command |

### Command-Specific Flags
| Command | Flag | Description |
|---------|------|-------------|
| `list`, `search` | `--json` | Output as JSON for scripting |
| `add` | `--preview` | Preview before adding |
| `export` | `--format` | Export format (skill.md, cursorrules, etc.) |
| `export` | `--stack` | Filter by stack |
| `completion` | — | Generate shell completion scripts |

---

## Tech Stack Taxonomy

PromptVault uses a hierarchical taxonomy so prompts are organized the way you think:

```
frontend/
  react/hooks          frontend/react/performance
  vue/                 frontend/svelte/
backend/
  node/express         backend/node/nestjs
  python/django        backend/python/fastapi
  go/gin               backend/go/grpc
devops/
  aws/                 devops/terraform/
  kubernetes/          devops/docker/
ai/
  prompting/           ai/rag/
  evaluation/          ai/agents/
database/
  postgresql/          database/redis/
  prisma/              database/drizzle/
```

Run `promptvault stacks` for the full list (80+ stacks).

---

## Export Formats

Export your entire prompt library — or a stack subset — to any AI tool format:

| Format | Command | Use Case |
|--------|---------|----------|
| `skill.md` | `--format skill.md` | Claude Code skills |
| `agents.md` | `--format agents.md` | OpenAI Agents |
| `claude.md` | `--format claude.md` | CLAUDE.md snippets |
| `cursorrules` | `--format cursorrules` | Cursor IDE |
| `windsurf` | `--format windsurf` | Windsurf IDE |
| `markdown` | `--format markdown` | Documentation |
| `json` | `--format json` | Integrations / import |
| `text` | `--format text` | Plain text |

---

## Data & Privacy

- All data stored locally in `~/.promptvault/vault.db` (SQLite)
- Zero telemetry, zero network calls
- Your prompts never leave your machine
- Backup: just copy `~/.promptvault/vault.db`

---

## Roadmap

- [x] Cloud sync (via private GitHub Gists)
- [ ] Team workspaces
- [ ] Prompt marketplace
- [x] MCP server mode (use your vault in Cursor/Windsurf automatically)
- [ ] Browser extension (save prompts from Claude.ai, ChatGPT)
- [x] Prompt templates with `{{variable}}` support
- [x] Decay detection (prompts that may no longer work with newer models)
- [x] Prompt testing framework
- [x] Version history (Git-like versioning)
- [x] AI-assisted authoring
- [x] Enhanced TUI (fuzzy search, multi-select, stats, full-screen preview)
- [x] Performance optimization (40x faster load)

**Latest:** v1.3 - Enhanced TUI Experience (Fuzzy Search, Multi-Select, Stats, Full-Screen Preview, 40x Faster!)

Star the repo to stay updated ⭐

---

## 📚 Documentation

### Complete Guides
- **[📖 Command Reference (HTML)](docs/commands.html)** - Interactive web documentation
- **[📝 Command Reference (Markdown)](docs/COMMANDS.md)** - Full documentation with examples
- **[🧪 Testing Guide](docs/TESTING-GUIDE.md)** - Prompt testing best practices
- **[📜 Versioning Guide](docs/VERSIONING-GUIDE.md)** - Version control for prompts
- **[🤖 AI Authoring Guide](docs/AI-AUTHORING-GUIDE.md)** - AI-assisted creation
- **[🔍 Decay Detection](docs/DECAY-DETECTION-GUIDE.md)** - Audit and maintenance
- **[🎨 TUI Enhancements](docs/TUI-FINAL.md)** - v1.3 TUI features guide
- **[⚡ Performance Fix](docs/ULTIMATE-PERF-FIX.md)** - 40x speed improvement details

### Quick Links
- [Installation](#install)
- [Quick Start](#quick-start)
- [TUI Keybindings](#tui-keybindings)
- [CLI Commands](#cli-commands--aliases)
- [Export Formats](#export-formats)

---

## Contributing

```sh
git clone https://github.com/Bharath-code/promptvault
cd promptvault
make deps
make run         # run the TUI
make init        # add curated prompts
make build       # build binary
```

---

## License

MIT — see [LICENSE](LICENSE)

---

*Built with [Bubble Tea](https://github.com/charmbracelet/bubbletea) 🧋*
