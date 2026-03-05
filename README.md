# ⚡ PromptVault

> The universal prompt OS for developers — store, search, and deploy AI prompts by tech stack, right from your terminal.

[![Downloads](https://img.shields.io/github/downloads/Bharath-code/promptvault/total?style=flat-square&color=7C3AED)](https://github.com/Bharath-code/promptvault/releases)
[![Release](https://img.shields.io/github/v/release/Bharath-code/promptvault?style=flat-square)](https://github.com/Bharath-code/promptvault/releases)
[![Go Version](https://img.shields.io/badge/go-1.20+-blue?style=flat-square)](https://go.dev)
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
- **📦 Single binary** — No runtime, no Node, no Docker. Just download and run.
- **🚀 One-command init** — Start with 15+ curated, production-grade prompts

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

# Add a prompt from stdin
cat my-prompt.txt | promptvault add "My prompt" --stack backend/python

# Search and copy to clipboard
promptvault get "useEffect"

# List all prompts
promptvault list

# List by stack
promptvault list --stack frontend/react

# Full-text search
promptvault search "typescript generics"

# Export to Claude Code SKILL.md
promptvault export --format skill.md --output SKILL.md

# Export to Cursor rules
promptvault export --format cursorrules --stack frontend/react > .cursorrules

# Export to AGENTS.md
promptvault export --format agents.md > AGENTS.md

# Import prompts from JSON
promptvault import shared-prompts.json

# Show stats
promptvault stats
```

---

## TUI Keybindings

| Key | Action |
|-----|--------|
| `↑` / `↓` or `k` / `j` | Navigate prompts |
| `Enter` or `Space` | Copy to clipboard |
| `/` | Search |
| `a` | Add new prompt |
| `e` | Edit selected |
| `d` | Delete selected |
| `v` | Toggle preview pane |
| `r` | Refresh |
| `Esc` | Clear filter / go back |
| `q` | Quit |

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

- [ ] Cloud sync (end-to-end encrypted)
- [ ] Team workspaces
- [ ] Prompt marketplace
- [ ] MCP server mode (use your vault in Cursor/Windsurf automatically)
- [ ] Browser extension (save prompts from Claude.ai, ChatGPT)
- [ ] Prompt templates with `{{variable}}` support
- [ ] Decay detection (prompts that may no longer work with newer models)

Star the repo to stay updated ⭐

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
