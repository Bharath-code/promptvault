# TUI/CLI Research: Industry Best Practices for UX, DX, and UI

## The Golden Standards: Tools That Define Excellence

These tools are universally revered for their UX, DX, and UI. They're the reference points every TUI/CLI developer should study.

### The Hall of Fame

| Tool | Stars | Language | Purpose | Why It's Great |
|------|-------|----------|---------|---------------|
| **lazygit** | 57.5K | Go | Git TUI | King of git TUIs; 5-star UX; discoverable keybindings; built with Bubble Tea |
| **lazydocker** | 42.8K | Go | Docker TUI | Same dev as lazygit; live logs, stats, container management |
| **htop** | 23.5K | C | System monitor | The gold standard for real-time TUI; discoverable (`F1` help), intuitive |
| **btop** | 23.5K | C++ | System monitor | Modern htop replacement; gorgeous UI, GPU-accelerated rendering |
| **yazi** | 23K | Rust | File manager | Fastest terminal file manager; async I/O; vim-like keys |
| **ranger** | 9K+ | Python | File manager | Six-pane preview; vim keybindings; extensible |
| **superfile** | 13K | Go | File manager | Pretty, modern; mouse support |
| **gh-dash** | 8K | Go | GitHub CLI | Beautiful dashboard for PRs/issues |
| **oha** | 7.9K | Rust | HTTP load tester | TUI animation during load testing |
| **posting** | 8.3K | Python | API client | Like Postman, in terminal |
| **harlequin** | 4.3K | Python | SQL IDE | Full PostgreSQL IDE in terminal |
| **trippy** | 4.3K | Rust | Network diag | Modern traceroute with TUI |
| **soft-serve** | 5.7K | Go | Git server | Self-hostable git server over SSH; built with Charm |
| **delta** | 17K | Rust | Git pager | Beautiful git diffs with syntax highlighting |

### The Reference Documents

1. **[CLI Guidelines](https://clig.dev/)** — The definitive guide. By the creators of Docker Compose. Covers philosophy (human-first, consistency, conversation) and concrete guidelines. **Must-read.**

2. **[UX Patterns for CLI Tools](https://lucasfcosta.com/blog/ux-patterns-cli-tools)** — 8 patterns: getting started, interactive mode, input validation, human errors, colors/emojis, loading indicators, context-awareness, exit codes, streams, consistent command trees.

3. **[terminal-apps.dev](https://terminal-apps.dev/)** — Curated catalog of best TUI apps with screenshots.

4. **[awesome-tuis](https://github.com/rothgar/awesome-tuis)** — 17K stars; 330+ contributors; comprehensive list of every TUI app.

---

## Core UX Principles

### 1. Human-First Design (from clig.dev)

UNIX commands were machine-first. Modern CLIs must be human-first. This doesn't mean dumbing down — it means optimizing for human cognitive patterns first, scriptability second.

### 2. Time to Value

GUIs lead users by the hand. CLIs should do the same:
- **Getting started**: Show useful first commands, not a wall of docs
- **Examples over documentation**: Lead with examples in help text
- **Onboarding**: PromptVault's 11-step onboarding tour is exactly right

### 3. Discoverability

GUIs win at discoverability because everything is visible. CLIs can replicate this:
- Comprehensive `--help` with examples
- Suggest next commands (`git status` tells you `git add` and `git restore`)
- Fuzzy suggestions on typos (Cobra has this built-in)
- **Keybinding cheatsheet** (`?` in lazygit shows all keybindings)

### 4. Conversation as Norm

CLI usage is a dialogue:
- Invalid input → suggest corrections (`git` does this: "did you mean `git commit`?")
- Multi-step workflows → confirm intermediate state before scary actions
- Dry-run mode → show what will happen before doing it (`--dry-run` flag)
- Suggestions after success → what to do next

### 5. Error Messages That Teach

Bad: `Error: operation failed`
Good: `Error: title is required. Pass title as argument: promptvault add "My prompt"` (already implemented!)

NPM's error format is the gold standard:
```
npm ERR! code ECONNRESET
npm ERR! network This is most likely not a problem with npm itself
npm ERR! If you are behind a proxy, please set proxy config...
npm ERR! See: 'npm help config'
```

### 6. Progress & Loading States

Spinner is minimum; progress bar is better; granular progress (like Docker's layer download) is best. Never leave users wondering if the command is frozen.

---

## DX (Developer Experience) Patterns

### 1. Consistent Command Trees

`kubectl` sets the gold standard. Once you learn `kubectl get pods`, you can guess `kubectl get deployments`. Consistent patterns make CLIs guessable.

### 2. Subcommand Hierarchy

```
promptvault config show      # ✓ clear hierarchy
promptvault search-history list  # ✓ clear hierarchy
```

### 3. Flag Conventions (follow standards)

| Flag | Meaning |
|------|---------|
| `-h`, `--help` | Help |
| `-v`, `--verbose` | Verbose output |
| `-d`, `--debug` | Debug mode |
| `--json` | JSON output |
| `--dry-run` | Preview without executing |
| `-f`, `--force` | Skip confirmation |
| `-o`, `--output` | Output file/directory |
| `-s`, `--stack` | Filter by stack |

### 4. stdin/stdout/stderr Discipline

- **stdout**: machine-readable output (the actual result)
- **stderr**: log messages, errors, prompts
- **stdin**: accept piped input
- This makes piping work: `cat file.json | pv import -`

### 5. Exit Codes

- **0**: success
- **non-zero**: failure with specific meaning
- Scripts depend on this — always return 0 on success

### 6. Context-Awareness

NPM is the gold standard:
- Detects `package.json` in current directory
- Auto-adds git branch as tag
- Project-local configs via `.npmrc`

PromptVault already does stack auto-detection — this is excellent.

---

## UI Patterns (Visual Design)

### 1. Color Usage

From lazygit, npm, yarn — colors communicate **semantics**:
- **Green** = success / added
- **Red** = error / deleted
- **Yellow** = warning / changed
- **Cyan/Blue** = info / neutral
- **Muted/Gray** = secondary info

Rules:
1. Check terminal capability (`NO_COLOR` env, `TERM=dumb`)
2. Don't overuse — if everything is colored, nothing stands out
3. Output must remain `grep`able
4. Use bold/inverse before reaching for color

### 2. Unicode + Emoji (Sparingly)

Good use: ✓ ✗ ⚠ ℹ ⚡ (single-char, grepable)
Bad use: 🚀 🔥 💯 ✨ (distracting, not grepable, platform-varied)

### 3. Box-Drawing for Layout

```
┌─────────────────────────┐
│  Header                 │
├─────────────────────────┤
│  Content                │
└─────────────────────────┘
```

Used in preview boxes, diffs, stats. Already used in PromptVault's `--preview` mode.

### 4. Tables for Structured Data

For list outputs, formatted tables beat raw lines. Most CLIs don't do this enough.

### 5. Spinners & Progress Bars

- Spinner frames: `⠋⠙⠹⠸⠼⠴⠦⠧⠇⠏` (already in PromptVault)
- Show **what** is happening, not just **that** something is happening
- Docker's layer-by-layer progress is the gold standard

---

## TUI-Specific Patterns

### 1. Lazygit UX — The Gold Standard

What makes lazygit exceptional:
- **Keybinding overlay**: Press `?` → full cheatsheet appears
- **Context panels**: Files, commits, stash, branches all visible
- **Inline actions**: `space` to stage, `enter` to expand/diff
- **Command palette**: `:` opens vim-like command palette
- **Status indicators**: Colored indicators for staged/unstaged/conflicts
- **Rich diffs**: Inline diff with + green / - red
- **Branching UI**: Visual branch graph
- **Undo**: `z Z` to undo almost anything
- **Fresh/reset**: Easy state recovery

### 2. Htop UX — The System Monitor Standard

- **Help overlay**: `F1` shows all keybindings
- **Tree view**: `t` shows process tree
- **Sortable columns**: `P M N T` for CPU/Memory/PGID/Time
- **Kill/resize**: `F9`, `F5` etc.
- **Search**: `/` to filter processes
- **Color schemes**: `Z` to customize

### 3. Common TUI Keybinding Conventions

From analyzing lazygit, htop, ranger, yazi:

| Action | Keys |
|--------|------|
| Help | `?`, `F1` |
| Quit | `q`, `Esc`, `Ctrl+C` |
| Navigate up/down | `k/j`, `↑/↓`, `Ctrl+P/N` |
| Select/confirm | `Enter`, `l` |
| Back/escape | `←`, `h`, `Esc` |
| Search | `/` |
| Refresh | `r`, `Ctrl+R` |
| Command palette | `:` |
| Toggle panel | `Tab`, `Ctrl+I` |
| Delete | `d`, `dd` (vim-style) |
| Copy | `yy`, `y` |
| Preview | `v`, `p` |

### 4. Vim Mode in TUIs

Many modern TUIs adopt vim keys because devs already know them:
- `j/k` navigation
- `/` search
- `:` command palette
- `dd` delete, `yy` copy, `p` paste
- `i` insert mode, `v` visual mode

PromptVault's vim mode implementation is on the right track.

---

## Best Practices Summary

### Must-Have (Critical)

1. ✅ `stdout`/`stderr` separation (PromptVault does this)
2. ✅ Color-coded output (✓ ✗ ⚠ ℹ) (PromptVault does this)
3. ✅ `--help` with examples (PromptVault does this)
4. ✅ Human-readable error messages with suggestions (PromptVault does this)
5. ✅ Loading indicators (spinner) (PromptVault does this)
6. ✅ Exit codes (0=success, non-zero=failure)
7. ✅ `--json` flag for machine-readable output (PromptVault does this)

### Should-Have (Important)

8. 🔲 Interactive prompts when stdin is a TTY
9. 🔲 `--dry-run` for preview mode (PromptVault does this for import)
10. 🔲 Fuzzy suggestions on typos (Cobra supports this)
11. 🔲 `--no-color` / respect `NO_COLOR` env
12. 🔲 Context-awareness (auto-detect env/project)
13. 🔲 Keybinding cheatsheet (`?` overlay)
14. 🔲 Confirmation before destructive actions
15. 🔲 stdin pipe support

### Nice-to-Have (Delight)

16. 🔲 `?` help overlay showing all keybindings
17. 🔲 Command palette (vim `:command` style)
18. 🔲 Undo support
19. 🔲 Fresh/reset state
20. 🔲 ASCII art / box-drawing for visual structure
21. 🔲 Progress bars instead of spinners for long operations
22. 🔲 Quiet mode (`-q`) for scripting
23. 🔲 `--plain` flag to disable formatting for `grep`-compatibility

---

## Inspiration: What PromptVault Already Does Right

Based on the codebase analysis:

- ✅ Toast notifications (4 types, auto-expiry)
- ✅ Onboarding tour (11 steps, auto-triggers on first run)
- ✅ Syntax highlighting (keywords, strings, comments, templates)
- ✅ Stack tree navigation (hierarchical, `t` to open)
- ✅ Config system (6 themes, persisted)
- ✅ Quick actions panel (`Tab` to toggle)
- ✅ Mouse support (click, wheel, hover)
- ✅ Search history (persisted, `h` in search)
- ✅ Bulk export (individual files)
- ✅ `--dry-run` for import preview
- ✅ Color-coded output with icons
- ✅ Auto-detect stack from `package.json`, `go.mod`, etc.
- ✅ Git branch auto-tagging
- ✅ Shell completion (bash, zsh, fish, powershell)
- ✅ Configurable keybindings
- ✅ Vim mode (NORMAL, INSERT, VISUAL, COMMAND)
- ✅ Spinner loading indicator

---

## Gap Analysis: What Could Be Improved

### High Priority

1. **`?` keybinding overlay** — Every great TUI (lazygit, htop, ranger) shows a full cheatsheet when `?` is pressed. This is the single biggest discoverability improvement.

2. **Respect `NO_COLOR` env var** — The color codes should be disabled when `NO_COLOR` is set or `TERM=dumb`.

3. **stdin pipe for `add`** — The `add` command reads from stdin, but the stdin check has a bug: it checks `ModeCharDevice` which fails interactively. Should check `!os.Stdin.IsNil()` first.

4. **Interactive mode** — When running `promptvault add` with no args in a TTY, prompt interactively (like `npm init`).

5. **Status bar / HUD** — Show current stack filter, count, vim mode indicator. Inspired by lazygit's bottom status.

### Medium Priority

6. **`-q` / `--quiet` flag** — Suppress non-error output for scripting. Currently verbose by default.

7. **`--plain` flag** — For `grep`-compatible output (no ANSI codes, no box-drawing).

8. **Fresh command** — `promptvault fresh` to reset the TUI state, like lazygit's `r` key.

9. **Undo support** — Undo recent actions (delete, revert). Lazygit has excellent undo (`z Z`).

10. **Command palette** — `:` opens a fuzzy command palette (like `:palette` in vim). Already has vim mode but needs palette UI.

### Nice to Have

11. **Progress bars** — Replace spinner with progress bars for long operations (bulk export, init).

12. **Man pages** — `man promptvault` for offline docs.

13. **Copilot hints** — Suggest next actions inline (e.g., "Run `promptvault list` to see all prompts").

14. **Plugin system** — Extensible commands (the planned low-priority feature).

15. **Timestamps on list** — Show `created_at`/`updated_at` in list view with toggle.

---

## Resources

- [clig.dev](https://clig.dev/) — The definitive CLI design guide
- [lucasfcosta.com/blog/ux-patterns-cli-tools](https://lucasfcosta.com/blog/ux-patterns-cli-tools) — UX patterns
- [terminal-apps.dev](https://terminal-apps.dev/) — Curated TUI catalog
- [github.com/rothgar/awesome-tuis](https://github.com/rothgar/awesome-tuis) — 17K star list
- [evilmartians.com/chronicles/cli-ux-best-practices](https://evilmartians.com/chronicles/cli-ux-best-practices-3-patterns-for-improving-progress-displays) — Progress display patterns
- [no-color.org](https://no-color.org/) — Color disabling standard
- [github.com/jesseduffield/lazygit](https://github.com/jesseduffield/lazygit) — Study the source code
- [github.com/aristocratos/btop](https://github.com/aristocratos/btop) — Beautiful modern TUI
