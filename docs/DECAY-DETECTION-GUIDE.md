# 🔍 Decay Detection Guide

**Keep your prompts fresh and effective.**

---

## 🎯 Why Decay Detection?

Prompts decay over time:
- 🤖 Models get deprecated
- 📉 Test success rates drop
- 🕐 Prompts become obsolete
- 📚 Framework best practices change

Decay detection helps you:
- ✅ Identify problematic prompts before they cause issues
- ✅ Keep your vault up-to-date
- ✅ Maintain high quality standards
- ✅ Save time on manual review

---

## 🚀 Quick Start

```bash
# Run full audit
promptvault audit

# See only critical issues
promptvault audit --severity critical

# JSON output for scripting
promptvault audit --json
```

---

## 📋 What Gets Detected

### 1. Unused Prompts (🟡 Warning)

**Detection:** Not used in 90+ days OR never used after 30+ days

**Why it matters:**
- Wastes storage
- Clutters search results
- May be obsolete

**Suggestion:** Test or remove

**Example Output:**
```
🟡 [unused] React v17 Component Pattern
   Not used in 120 days
   💡 Review if this prompt is still relevant
   📅 Last used: 120 days ago
```

---

### 2. Deprecated Models (🔴 Critical)

**Detection:** Uses deprecated AI models

**Deprecated Models:**
| Model | Replacement |
|-------|-------------|
| `gpt-3.5-turbo` | `gpt-4o` or `gpt-4-turbo` |
| `gpt-4-turbo` | `gpt-4o` |
| `claude-2` | `claude-3-sonnet` or `claude-3-opus` |
| `claude-instant` | `claude-3-haiku` |
| `text-davinci-003` | `gpt-4o` |
| `text-davinci-002` | `gpt-4o` |
| `code-davinci-002` | `gpt-4o` |

**Why it matters:**
- Deprecated models may be removed
- Newer models perform better
- Cost efficiency

**Suggestion:** Update to recommended model

**Example Output:**
```
🔴 [deprecated] Legacy Code Converter
   Uses deprecated model: gpt-3.5-turbo
   💡 Use gpt-4o or gpt-4-turbo instead
   🤖 Deprecated model: gpt-3.5-turbo
```

---

### 3. Low Test Success Rate (🔴 Critical)

**Detection:** Test pass rate < 50% (minimum 3 tests)

**Why it matters:**
- Prompt isn't working reliably
- May produce incorrect outputs
- Wastes API credits

**Suggestion:** Review test failures and update prompt

**Example Output:**
```
🔴 [low_success] TypeScript Type Generator
   Low test success rate: 33.3%
   💡 Review test failures and update prompt
   📊 Pass rate: 33.3%
```

---

### 4. Old Versions (🟢 Info)

**Detection:** Not updated in 180+ days

**Why it matters:**
- Best practices may have changed
- Framework updates may break it
- Missing recent improvements

**Suggestion:** Review and update if needed

**Example Output:**
```
🟢 [old_version] Express.js Middleware
   Not updated in 200 days
   💡 Review and update if needed
```

---

## 📊 Audit Report

### Sample Output

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

📋 Issues Found (14):
────────────────────────────────────────────────────────────

🔴 [deprecated] Legacy Code Converter
   Uses deprecated model: gpt-3.5-turbo
   💡 Use gpt-4o or gpt-4-turbo instead
   🤖 Deprecated model: gpt-3.5-turbo

🔴 [low_success] TypeScript Type Generator
   Low test success rate: 33.3%
   💡 Review test failures and update prompt
   📊 Pass rate: 33.3%

🟡 [unused] React v17 Component Pattern
   Not used in 120 days
   💡 Review if this prompt is still relevant
   📅 Last used: 120 days ago

... and 11 more issues

Showing 10 of 14 issues

Run 'promptvault audit --severity critical' to see only critical issues
```

---

## 🎯 Usage Examples

### 1. Quick Health Check

```bash
# Just show critical issues
promptvault audit --severity critical
```

### 2. CI/CD Integration

```bash
# JSON output for parsing
promptvault audit --json | jq '.issues_found'

# Fail CI if critical issues found
if [ $(promptvault audit --severity critical --json | jq '.issues_found') -gt 0 ]; then
  echo "❌ Critical prompt issues found"
  exit 1
fi
```

### 3. Regular Maintenance

```bash
# Add to weekly cron
0 9 * * 1 promptvault audit --json >> /var/log/promptvault-audit.json
```

### 4. Filter by Issue Type

```bash
# See only deprecated models
promptvault audit | grep "deprecated"

# See only unused prompts
promptvault audit | grep "unused"
```

---

## 🔧 Decay Heuristics

### Configuration

| Issue Type | Threshold | Severity |
|------------|-----------|----------|
| Unused | 90 days | Warning |
| Unused (never) | 30 days | Warning |
| Deprecated Model | Any usage | Critical |
| Low Success Rate | < 50% | Critical |
| Old Version | 180 days | Info |

### Customization

Coming soon:
```yaml
# ~/.promptvault/config.yaml
decay:
  unused_threshold_days: 90
  old_version_threshold_days: 180
  min_test_success_rate: 50
  exclude_stacks:
    - experimental
```

---

## 📈 Quality Metrics

### Health Score Calculation

```
Health Score = (Healthy Prompts / Total Prompts) × 100
```

**Interpretation:**
- **90-100%**: Excellent maintenance
- **70-89%**: Good, some attention needed
- **50-69%**: Needs work
- **< 50%**: Critical, immediate attention required

### Recommended Audit Frequency

| Vault Size | Frequency |
|------------|-----------|
| < 50 prompts | Monthly |
| 50-200 prompts | Bi-weekly |
| 200-500 prompts | Weekly |
| 500+ prompts | Daily (CI/CD) |

---

## 🎓 Best Practices

### 1. Fix Critical Issues First

**Priority Order:**
1. 🔴 Deprecated models (may stop working anytime)
2. 🔴 Low success rates (producing bad outputs)
3. 🟡 Unused prompts (cleaning up clutter)
4. 🟢 Old versions (preventive maintenance)

### 2. Test Before Removing

```bash
# Test a prompt before deciding to remove it
promptvault test abc123

# If it still works, keep it
# If it fails, update or remove
```

### 3. Update in Batches

```bash
# Group similar updates
# Update all deprecated gpt-3.5-turbo prompts at once
# Update all claude-2 prompts together
```

### 4. Track Improvements

```bash
# Run audit before and after cleanup
promptvault audit --json > audit-before.json

# Make updates...

promptvault audit --json > audit-after.json

# Compare improvements
```

---

## 🐛 Troubleshooting

### "No issues found" but prompts seem broken

**Cause:** Decay detection only checks specific heuristics

**Solution:**
```bash
# Manually test suspicious prompts
promptvault test abc123

# Check test history
promptvault test abc123 --history
```

### "Too many false positives"

**Cause:** Thresholds may be too aggressive for your use case

**Solution:**
- Some prompts are intentionally rarely used (reference prompts)
- Consider tagging them as `reference` or `archive`
- Future: Exclude by tag

### Audit is slow

**Cause:** Large vault or many tests to analyze

**Solution:**
```bash
# Filter by severity for faster results
promptvault audit --severity critical

# Use JSON for programmatic access
promptvault audit --json
```

---

## 📊 Integration Examples

### GitHub Actions

```yaml
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

### Pre-commit Hook

```bash
#!/bin/bash
# .git/hooks/pre-commit

# Check for deprecated models in changed prompts
if promptvault audit --severity critical | grep -q "deprecated"; then
  echo "❌ Critical prompt issues found"
  echo "Run 'promptvault audit' to see details"
  exit 1
fi
```

---

## 📈 Future Features

Coming soon:
- [ ] Custom decay thresholds
- [ ] Exclude by tag/stack
- [ ] Auto-fix for deprecated models
- [ ] Trend analysis (improving vs declining)
- [ ] Email notifications for critical issues
- [ ] Slack/Teams integration

---

## 🔗 Related Commands

- `promptvault test` - Test individual prompts
- `promptvault history` - View version history
- `promptvault stats` - Overall statistics

---

**Remember:** Regular audits keep your prompt vault healthy and effective! 🔍✨

For more help: `promptvault audit --help`
