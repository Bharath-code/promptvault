# PromptVault

> The universal prompt OS for developers — store, search, and deploy AI prompts by tech stack, right from your terminal.

[![Version](https://img.shields.io/npm/v/promptvault.svg)](https://www.npmjs.com/package/promptvault)
[![Downloads](https://img.shields.io/npm/dm/promptvault.svg)](https://www.npmjs.com/package/promptvault)
[![License](https://img.shields.io/npm/l/promptvault.svg)](https://github.com/Bharath-code/promptvault/blob/main/LICENSE)

## 🚀 Quick Start

```bash
# Install globally
npm install -g promptvault

# Verify installation
promptvault --version

# Get started
promptvault --help
```

## 📦 Installation

### NPM (Recommended)

```bash
npm install -g promptvault
```

This will automatically download the appropriate binary for your platform.

### Manual Installation

If the automatic installation fails, you can download the binary manually:

1. Go to [GitHub Releases](https://github.com/Bharath-code/promptvault/releases)
2. Download the binary for your platform
3. Extract and add to your PATH

### Platform Support

| Platform | Architecture | Binary |
|----------|--------------|--------|
| macOS | ARM64 (M1/M2) | `promptvault-darwin-arm64` |
| macOS | Intel x64 | `promptvault-darwin-amd64` |
| Linux | ARM64 | `promptvault-linux-arm64` |
| Linux | x64 | `promptvault-linux-amd64` |
| Windows | x64 | `promptvault-windows-amd64.exe` |

## 🎯 Features

- **🗂 Tech-stack taxonomy** — Organize by `frontend/react/hooks`, `backend/python/fastapi`, `devops/terraform`, and 80+ more
- **⚡ Fuzzy search** — Find any prompt in under 3 seconds
- **📋 One-key copy** — Press Enter to copy any prompt to clipboard instantly
- **🔄 Multi-format export** — Export to `SKILL.md`, `AGENTS.md`, `.cursorrules`, `.windsurfrules`
- **🤖 Model tagging** — Mark prompts as verified for Claude, GPT-4o, Gemini
- **💻 Beautiful TUI** — Built with Bubble Tea, works in any terminal
- **🗄 SQLite + FTS** — Zero-dependency local storage with full-text search
- **🔌 MCP Server** — Connect to Claude Desktop, Cursor, and Windsurf
- **☁️ Cloud Sync** — Backup and sync using private GitHub Gists

## 📖 Documentation

For complete documentation, visit the [GitHub repository](https://github.com/Bharath-code/promptvault).

### Quick Commands

```bash
# Initialize with starter prompts
promptvault init

# Open interactive TUI
promptvault

# Add a prompt
promptvault add "My prompt" --stack frontend/react --content "..."

# Search prompts
promptvault search "react hooks"

# Export to Cursor rules
promptvault export --format cursorrules --stack frontend/react > .cursorrules

# View help
promptvault --help
```

## 🔧 Troubleshooting

### Installation Failed

If installation fails, check:

1. **Internet connection** - The postinstall script needs to download the binary
2. **Platform support** - Ensure your OS/architecture is supported (see table above)
3. **Firewall/Proxy** - Some corporate networks block GitHub downloads

**Manual installation workaround:**
```bash
# Download binary from GitHub Releases
# https://github.com/Bharath-code/promptvault/releases

# Extract and move to PATH
mv promptvault /usr/local/bin/
```

### Binary Not Found

If you get "binary not found" error:

```bash
# Reinstall
npm uninstall -g promptvault
npm install -g promptvault

# Or verify installation
which promptvault
```

### Permission Issues (Unix/Linux/macOS)

```bash
# Fix permissions
sudo chown -R $(whoami) $(npm config get prefix)/{lib/node_modules,bin}

# Or reinstall with sudo
sudo npm install -g promptvault --unsafe-perm
```

## 📝 License

MIT License - see [LICENSE](https://github.com/Bharath-code/promptvault/blob/main/LICENSE) for details.

## 🔗 Links

- [GitHub Repository](https://github.com/Bharath-code/promptvault)
- [NPM Package](https://www.npmjs.com/package/promptvault)
- [Full Documentation](https://github.com/Bharath-code/promptvault#readme)
- [Report Issues](https://github.com/Bharath-code/promptvault/issues)

---

**Built with ❤️ for developers**
