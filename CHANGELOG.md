# Changelog

All notable changes to PromptVault will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

---

## [Unreleased]

### Planned
- Team workspaces for collaborative prompt management
- Prompt marketplace for sharing community prompts
- Browser extension for saving prompts from Claude.ai, ChatGPT
- Advanced analytics dashboard
- Custom export templates

---

## [1.3.0] - 2026-03-06

### 🎨 Added - Enhanced TUI Experience

#### Fuzzy Search
- Intelligent fuzzy matching algorithm with 0-100% scoring
- Searches across title, stack, tags, and content
- Results sorted by relevance score
- Tolerates typos and partial matches
- Match percentage displayed in results

#### Quick Action Menu
- Press `?` to show comprehensive keybinding reference
- Organized by categories (Navigation, Actions, Quick Actions, Other)
- Clean overlay design with any-key-to-close behavior
- Always accessible from any state

#### Usage Statistics Dashboard
- Press `s` to show real-time statistics
- Total prompts and usage count
- Top 5 stacks with medal emojis (🥇🥈🥉)
- Top 5 most used prompts
- Beautiful formatted overlay display

#### Recent Prompts Section
- Press `R` to toggle frequently used prompts
- Shows top 5 most used prompts
- Always visible at top of list when enabled
- Intelligent caching for instant toggle

#### Multi-Select Batch Operations
- Press `Space` to select/deselect individual prompts
- Visual checkmark (✓) on selected items
- Press `x` for batch processing
- Batch increment usage tracking
- Clear selection after processing

#### Full-Screen Preview Overlay
- Press `v` to toggle immersive full-screen preview
- 100% screen space for content
- Markdown rendering with syntax highlighting
- Scrollable viewport with ↑/↓
- Header with title and metadata
- Footer with action hints

### ⚡ Performance

#### Critical Optimizations
- **40x faster initial load** (10-12s → 0.3s)
- **200x faster navigation** (200ms → 1ms)
- **200x faster subsequent renders** with intelligent caching
- Recent prompts caching with dirty flag pattern
- Skipped preview rendering on initial load
- Plain text preview during navigation (no markdown)
- Lazy initialization of expensive operations

#### Memory Usage
- Minimal overhead (+2MB for caching)
- Total memory usage: ~52MB
- Efficient cache invalidation

### 🐛 Fixed
- Recent prompts recalculating on every render (critical performance bug)
- Preview rendering blocking initial load
- Multi-select index checking wrong cursor position
- 'v' key not expanding preview pane
- Glamour markdown rendering on every navigation
- Unnecessary map lookups in render loop

### 📚 Documentation
- Added 12 comprehensive documentation files
- Updated README.md with v1.3 features
- Updated website with v1.3 showcase
- Added performance optimization guides
- Added TUI enhancement guides

### 🔧 Technical
- Added `recentCache` and `recentDirty` fields to App struct
- Modified `renderRecentPrompts()` to use caching
- Rewrote `updatePreview()` for ultra-fast plain text
- Separated Enter and Space key handlers
- Added `renderFullScreenPreview()` function
- Optimized render loop to avoid expensive operations

---

## [1.2.0] - 2026-03-05

### 🧪 Added - Professional Prompt Engineering

#### Prompt Testing Framework
- `promptvault test <id>` - Interactive test sessions
- Test prompts against expected outputs
- Similarity scoring (0-100%)
- Track pass/fail rates
- View test history
- Support for multiple models (Claude, GPT-4o, etc.)
- Latency and token usage tracking

#### Version History
- `promptvault history <id>` - View all versions
- `promptvault diff <id> <v1> <v2>` - Compare versions
- `promptvault revert <id> <version>` - Revert to previous version
- Automatic versioning on every edit
- Commit messages support
- Author tracking
- Git-like version control for prompts

#### AI-Assisted Authoring
- `promptvault create --ai` - AI-assisted creation
- Variable detection (`{{variable}}` syntax)
- Tag recommendations based on content
- Stack auto-detection
- Anti-pattern detection (6 patterns)
- Quality scoring (0-100)
- Improvement suggestions
- Interactive creation flow

#### Decay Detection
- `promptvault audit` - Comprehensive audit
- Detects deprecated models (gpt-3.5-turbo, claude-2, etc.)
- Identifies unused prompts (90+ days)
- Finds low test success rates (<50%)
- Flags old versions (180+ days)
- JSON output for CI/CD integration
- Severity levels (critical, warning, info)
- Actionable recommendations

#### Auto-Export Watch
- `promptvault watch --format <fmt> --output <file>`
- Continuous monitoring for database changes
- Auto-export on any change
- Configurable check interval (default: 5s)
- Stack filtering support
- Verbose mode for debugging
- Graceful shutdown (Ctrl+C)

### 📚 Documentation
- Added TESTING-GUIDE.md
- Added VERSIONING-GUIDE.md
- Added AI-AUTHORING-GUIDE.md
- Added DECAY-DETECTION-GUIDE.md
- Updated COMMANDS.md with new commands
- Updated commands.html with examples

### 🔧 Technical
- Added `internal/prompttest/tester.go` - Testing engine
- Added `internal/decay/detector.go` - Decay detection
- Added `internal/ai/assistant.go` - AI assistance
- Added test_results table to database
- Added prompt_versions table to database
- Implemented similarity scoring algorithm
- Implemented deprecated model detection
- Implemented anti-pattern detection

---

## [1.1.0] - 2026-03-04

### 💡 Added - Developer Experience Improvements

#### Smart Error Messages
- Actionable suggestions for common errors
- Context-aware recommendations
- Examples included in error messages
- Covers: missing title, missing content, GitHub token issues, prompt not found

#### Shell Completion
- `promptvault completion bash` - Bash completion
- `promptvault completion zsh` - Zsh completion
- `promptvault completion fish` - Fish completion
- `promptvault completion powershell` - PowerShell completion
- Auto-completes commands, flags, and options

#### JSON Output
- `--json` flag for list and search commands
- Scriptable output for automation
- Pipe-friendly for jq processing
- Consistent JSON schema across commands

#### Verbose/Debug Mode
- `-v, --verbose` - Verbose output
- `-d, --debug` - Debug output with timestamps
- See database queries and execution times
- Helpful for troubleshooting

#### Command Aliases
- `list` → `ls`, `show`, `list-all`
- `search` → `find`, `query`
- `delete` → `rm`, `remove`, `del`
- `get` → `fetch`
- `export` → `exp`
- `import` → `imp`
- `stats` → `statistics`

#### Rich Colors & Icons
- Color-coded output throughout CLI
- Success (green ✓), Error (red ✗), Warning (yellow ⚠️)
- Info messages (cyan ℹ️)
- Primary actions (purple ⚡)

#### Smart Defaults
- Auto-detect stack from current directory
- Supports: React, Go, Python, Terraform, Docker, Kubernetes
- Auto-tag with Git branch name
- Reduces typing for common operations

#### Preview Mode
- `--preview` flag for add command
- Beautiful boxed preview before committing
- Shows metadata (title, stack, tags, models)
- Confirmation prompt

#### Git Integration
- Auto-detects Git repository
- Adds branch name as tag (e.g., `git:main`, `git:feature-xyz`)
- Track which branch created each prompt

### 📚 Documentation
- Created comprehensive COMMANDS.md
- Created interactive commands.html
- Updated README.md with all new features
- Added quick reference cards
- Added troubleshooting guides

### 🔧 Technical
- Added error suggestion engine
- Integrated cobra completion generators
- Added color code constants
- Implemented directory-based stack detection
- Added Git branch detection
- Enhanced error handling throughout

---

## [1.0.0] - 2026-02-15

### 🎉 Initial Release

#### Core Features

##### Tech-Stack Taxonomy
- Hierarchical organization (e.g., `frontend/react/hooks`)
- 80+ built-in stacks
- Custom stack support
- Inheritance model (child inherits from parent)

##### Fuzzy Search
- Full-text search across title, content, tags, stack
- SQLite FTS5 with porter stemming
- Results ranked by relevance
- Search in under 3 seconds

##### One-Key Copy
- Press Enter to copy to clipboard
- Instant clipboard integration
- Usage tracking on copy

##### Multi-Format Export
- `skill.md` - Claude Code skills
- `agents.md` - OpenAI Agents
- `claude.md` - CLAUDE.md snippets
- `cursorrules` - Cursor IDE
- `windsurf` - Windsurf IDE
- `markdown` - Documentation
- `json` - Integrations
- `text` - Plain text

##### Model Tagging
- Mark prompts as verified for specific models
- Support for Claude, GPT-4o, Gemini, and more
- Model-specific optimization tracking

##### Beautiful TUI
- Built with Bubble Tea
- Works in any terminal
- Keyboard-driven navigation
- Real-time search
- Markdown previews

##### SQLite + FTS
- Zero-dependency local storage
- Full-text search support
- Fast queries (<100ms)
- Single file database

##### MCP Server
- `promptvault mcp` command
- Connect to Claude Desktop, Cursor, Windsurf
- Model Context Protocol support
- Seamless AI integration

##### Cloud Sync
- `promptvault sync push` - Backup to GitHub Gist
- `promptvault sync pull` - Restore from Gist
- Private Gist support
- Token-based authentication

##### Markdown Previews
- Beautiful markdown rendering
- Syntax highlighting
- Glamour integration
- Readable in-terminal display

##### Interactive Variables
- Define `{{variables}}` in prompts
- TUI prompts to fill before copying
- Template-style prompts
- Reusable prompt templates

##### Smart Auto-Injection
- Auto-append to `.cursorrules`, `SKILL.md`, etc.
- Workspace integration
- Identify target files automatically

##### Single Binary
- No runtime dependencies
- No Node.js, no Docker
- Just download and run
- Cross-platform support

##### One-Command Init
- `promptvault init` - Initialize with 15+ curated prompts
- Production-grade starter prompts
- Covers multiple stacks

### 📚 Documentation
- README.md with complete feature list
- Quick start guide
- Installation instructions
- TUI keybindings reference
- CLI commands reference
- Export formats guide
- Tech stack taxonomy

### 🔧 Technical
- Go 1.23+
- SQLite3 with FTS5
- Bubble Tea for TUI
- Cobra for CLI
- Glamour for markdown rendering
- SQLite triggers for FTS index maintenance

---

## Version History Summary

| Version | Date | Codename | Key Features |
|---------|------|----------|--------------|
| [1.3.0] | 2026-03-06 | Enhanced TUI | Fuzzy search, multi-select, stats, full-screen, 40x faster |
| [1.2.0] | 2026-03-05 | Professional | Testing, versioning, AI-assist, decay detection |
| [1.1.0] | 2026-03-04 | DX | Smart errors, completion, JSON, verbose mode |
| [1.0.0] | 2026-02-15 | Initial | Core features, TUI, CLI, export, sync |

---

## Migration Guide

### From v1.2 to v1.3
- No breaking changes
- Database schema unchanged
- All existing prompts compatible
- New features available immediately

### From v1.1 to v1.2
- No breaking changes
- Database schema extended (new tables)
- Automatic migration on first run
- All existing prompts compatible

### From v1.0 to v1.1
- No breaking changes
- All features backward compatible
- Configuration unchanged

---

## Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

### Reporting Issues
- Use GitHub Issues
- Include version number
- Include steps to reproduce
- Include expected vs actual behavior

### Submitting PRs
- Fork the repository
- Create a feature branch
- Make your changes
- Add tests if applicable
- Submit a pull request

---

## Support

- **Documentation:** https://github.com/Bharath-code/promptvault/tree/main/docs
- **Issues:** https://github.com/Bharath-code/promptvault/issues
- **Discussions:** https://github.com/Bharath-code/promptvault/discussions
- **Releases:** https://github.com/Bharath-code/promptvault/releases

---

## License

PromptVault is released under the [MIT License](LICENSE).

---

**Last Updated:** March 6, 2026  
**Current Version:** v1.3.0  
**Status:** ✅ Stable
