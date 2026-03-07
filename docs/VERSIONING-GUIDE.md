# 📜 Version History Guide

**Track changes like Git for your prompts.**

---

## 🎯 Why Version Prompts?

Just like code, prompts evolve over time. Version history lets you:
- ✅ Track what changed and when
- ✅ Compare versions with diff
- ✅ Revert broken changes
- ✅ Document why changes were made
- ✅ Experiment safely

---

## 🚀 Quick Start

### 1. View History

```bash
# See all versions of a prompt
promptvault history abc123

# Output:
📜 Version History: React Hook Converter
────────────────────────────────────────────────────────────
▶ v5  2026-03-06 14:30  johndoe  Fixed edge case with useEffect
  v4  2026-03-06 11:15  johndoe  Added TypeScript support
  v3  2026-03-05 16:45  janedoe  Improved error handling
  v2  2026-03-05 10:00  johndoe  Updated for React 19
  v1  2026-03-04 09:00  johndoe  Initial version

Total versions: 5
```

### 2. Compare Versions

```bash
# Compare version 1 and 2
promptvault diff abc123 1 2

# Compare current with version 3
promptvault diff abc123 3 current

# Compare last two versions
promptvault diff abc123 HEAD~1 HEAD
```

### 3. Revert to Previous Version

```bash
# Revert to version 3
promptvault revert abc123 3

# Revert with custom message
promptvault revert abc123 3 -m "Reverting broken changes"
```

---

## 📋 Command Reference

### `promptvault history [prompt-id]`

View version history of a prompt.

**Output:**
- `▶` marks current version
- Version number
- Timestamp
- Author
- Commit message

---

### `promptvault diff [prompt-id] [v1] [v2]`

Compare two versions and show differences.

**Version Specifiers:**
- `1`, `2`, `3` - Specific version numbers
- `current`, `HEAD`, `latest` - Current version
- `HEAD~1` - Previous version

**Output:**
- Title changes (if any)
- Content diff with `-` for removed, `+` for added

**Example:**
```bash
$ promptvault diff abc123 1 2

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

---

### `promptvault revert [prompt-id] [version]`

Revert a prompt to a previous version.

**Flags:**
| Flag | Description |
|------|-------------|
| `-m, --message` | Custom commit message |

**What it does:**
1. Creates a new version with old content
2. Preserves version history
3. Doesn't delete any versions

**Example:**
```bash
$ promptvault revert abc123 3 -m "Undo last changes"

✓ Reverted abc123 to v3
New version: v6
```

---

## 🔄 How Versioning Works

### Automatic Versioning

**Every edit creates a version:**
- Editing in TUI → Auto-versioned
- `promptvault edit` → Auto-versioned
- Sync from Gist → Auto-versioned

**Version includes:**
- Full snapshot of prompt
- Timestamp
- Author (from `$USER` or `os.Getenv("USER")`)
- Commit message

### Manual Versioning

Coming soon:
```bash
# Create version snapshot without editing
promptvault version create abc123 -m "Before major refactor"
```

---

## 🎓 Best Practices

### 1. Write Meaningful Commit Messages

**❌ Bad:**
```
promptvault edit abc123 -m "fix"
```

**✅ Good:**
```
promptvault edit abc123 -m "Add error handling for null input"
```

### 2. Version Before Risky Changes

```bash
# Before major refactor
promptvault edit abc123 -m "Starting major refactor"

# Make changes...

# If something breaks, revert easily
promptvault revert abc123 HEAD~1
```

### 3. Compare Before Reverting

```bash
# Don't revert blindly
promptvault diff abc123 3 current

# See what you'll lose, then decide
promptvault revert abc123 3
```

### 4. Use in CI/CD

```yaml
# .github/workflows/test-prompts.yml
- name: Check for regressions
  run: |
    promptvault diff abc123 HEAD~1 HEAD | grep "^-" && \
    echo "⚠️ Content was removed" && exit 1
```

---

## 🔍 Advanced Usage

### Version Aliases

| Alias | Meaning |
|-------|---------|
| `current` | Latest version |
| `HEAD` | Latest version |
| `latest` | Latest version |
| `HEAD~1` | Previous version |
| `HEAD~2` | Two versions ago |

### Workflow Examples

#### 1. Experiment Safely

```bash
# Save current state
promptvault edit abc123 -m "Before experiment"

# Make risky changes...

# If it works, great! If not:
promptvault revert abc123 HEAD~1
```

#### 2. Collaborative Editing

```bash
# See who changed what
promptvault history abc123

# Compare your changes with teammate's
promptvault diff abc123 5 6

# Revert if needed
promptvault revert abc123 5 -m "Reverting team member's changes"
```

#### 3. Audit Trail

```bash
# Who changed this prompt?
promptvault history abc123 | grep "janedoe"

# What changed in version 3?
promptvault diff abc123 2 3
```

---

## 📊 Database Schema

Versions are stored in `prompt_versions` table:

```sql
CREATE TABLE prompt_versions (
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
  FOREIGN KEY (prompt_id) REFERENCES prompts(id)
);
```

**Cascade Delete:** When a prompt is deleted, all versions are deleted.

---

## 🐛 Troubleshooting

### "No versions recorded yet"

**Cause:** Prompt was created before versioning was implemented, or never edited.

**Solution:** Edit the prompt to create first version:
```bash
promptvault edit abc123 -m "Initial version snapshot"
```

### "Version X not found"

**Cause:** Version number doesn't exist.

**Solution:** Check available versions:
```bash
promptvault history abc123
```

### "Failed to create version snapshot"

**Cause:** Database error or disk full.

**Solution:** 
1. Check disk space
2. Check database integrity
3. Update will still succeed, versioning is best-effort

---

## 📈 Future Features

Coming soon:
- [ ] `promptvault version create` - Manual versioning
- [ ] `promptvault version list --all` - List all versions across prompts
- [ ] `promptvault version tag v1.0` - Tag specific versions
- [ ] `promptvault version branch` - Experimental branching
- [ ] TUI version browser with interactive diff

---

## 🔗 Related Commands

- `promptvault test` - Test prompts before/after changes
- `promptvault edit` - Edit prompts (auto-versions)
- `promptvault export` - Export specific versions

---

**Remember:** With great versioning power comes great responsibility! 📜✨

For more help: `promptvault history --help`, `promptvault diff --help`, `promptvault revert --help`
