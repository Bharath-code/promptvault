# 🧪 Prompt Testing Guide

**Test your prompts like code — before deploying them to production.**

---

## 🎯 Why Test Prompts?

Just like you wouldn't deploy code without tests, you shouldn't deploy prompts without validating they work correctly.

**Benefits:**
- ✅ Catch broken prompts before they reach your team
- ✅ Track prompt quality over time
- ✅ Compare performance across models
- ✅ Document expected behavior
- ✅ Prevent regressions when editing prompts

---

## 🚀 Quick Start

### 1. Interactive Testing

```bash
# Start interactive test session
promptvault test <prompt-id>

# Example
promptvault test abc123
```

**What happens:**
1. Enter test input
2. Enter expected output
3. AI model runs the test
4. See pass/fail result with similarity score

### 2. Single Test

```bash
# Test with specific input/output
promptvault test abc123 \
  --input "Convert this to TypeScript" \
  --expected "Here's the TypeScript version..." \
  --model claude-sonnet
```

### 3. View Test History

```bash
# See all tests for a prompt
promptvault test abc123 --history

# Output:
📊 Test History for: React Hook Converter
────────────────────────────────────────────────────────────
Total Tests: 15
Pass Rate: 86.7%
Average Score: 82.3/100

Recent Tests:
  ✅ [claude-sonnet] Score: 94.2% - 2026-03-06 10:30
  ❌ [claude-sonnet] Score: 65.1% - 2026-03-06 09:15
  ✅ [gpt-4o] Score: 88.7% - 2026-03-05 16:45
```

---

## 📋 Command Reference

### `promptvault test [prompt-id]`

Test a prompt against AI models.

**Flags:**

| Flag | Description | Default |
|------|-------------|---------|
| `--input` | Test input | Required for non-interactive |
| `--expected` | Expected output | Required for non-interactive |
| `--model` | Model to test against | `claude-sonnet` |
| `--all` | Run all saved tests | `false` |
| `--history` | Show test history | `false` |

**Supported Models:**
- `claude-sonnet` (default)
- `claude-opus`
- `claude-haiku`
- `gpt-4o`
- `gpt-4-turbo`
- `gpt-3.5-turbo`

---

## 🔧 Configuration

### API Keys

Testing requires API keys for the models you want to test against.

**Set via environment variables:**

```bash
# Anthropic (Claude)
export ANTHROPIC_API_KEY=sk-ant-...

# OpenAI (GPT)
export OPENAI_API_KEY=sk-...
```

**Or add to `.env` file:**

```bash
# ~/.promptvault/.env
ANTHROPIC_API_KEY=sk-ant-...
OPENAI_API_KEY=sk-...
```

---

## 📊 Understanding Test Results

### Test Result Fields

| Field | Description |
|-------|-------------|
| **Passed** | ✅ if score ≥ 70%, ❌ otherwise |
| **Score** | Similarity score (0-100) |
| **Latency** | API response time in ms |
| **Token Usage** | Total tokens consumed |
| **Error Message** | Error details if test failed |

### Scoring Algorithm

Tests use **word overlap similarity**:

```
Score = (Matching Words / Expected Words) × 100
```

**Thresholds:**
- ✅ **Passed**: ≥ 70% similarity
- ❌ **Failed**: < 70% similarity

**Example:**
```
Expected: "Create a React component with useState"
Actual:   "Create React component using useState hook"

Matching words: Create, React, component, useState (4/5)
Score: 80% ✅ PASSED
```

---

## 🎓 Best Practices

### 1. Write Specific Tests

**❌ Bad:**
```
Input: "test it"
Expected: "good output"
```

**✅ Good:**
```
Input: "Convert this JavaScript function to TypeScript"
Expected: "Function should have proper type annotations for parameters and return type"
```

### 2. Test Edge Cases

```bash
# Empty input
promptvault test abc123 --input "" --expected "Error message"

# Malformed input
promptvault test abc123 --input "<<<invalid>>>" --expected "Graceful error"

# Very long input
promptvault test abc123 --input "$(cat large-file.txt)" --expected "Summary"
```

### 3. Test Multiple Models

```bash
# Test on Claude
promptvault test abc123 --model claude-sonnet --input "..." --expected "..."

# Test on GPT-4
promptvault test abc123 --model gpt-4o --input "..." --expected "..."
```

### 4. Maintain Test Suites

For critical prompts, maintain 5-10 tests covering:
- ✅ Happy path
- ✅ Edge cases
- ✅ Error conditions
- ✅ Common use cases
- ✅ Performance (latency)

---

## 🔄 CI/CD Integration

### GitHub Actions Example

```yaml
# .github/workflows/test-prompts.yml
name: Test Prompts

on:
  push:
    paths:
      - 'prompts/**'

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      
      - name: Install PromptVault
        run: |
          go install github.com/Bharath-code/promptvault@latest
      
      - name: Run Tests
        env:
          ANTHROPIC_API_KEY: ${{ secrets.ANTHROPIC_API_KEY }}
        run: |
          promptvault test abc123 --all
          
      - name: Validate Test Results
        run: |
          promptvault test abc123 --history | grep "Pass Rate" | awk '{if ($4 < 80) exit 1}'
```

---

## 🐛 Troubleshooting

### "ANTHROPIC_API_KEY not set"

**Solution:**
```bash
export ANTHROPIC_API_KEY=sk-ant-...
```

### "Prompt not found"

**Solution:**
```bash
# Get correct prompt ID
promptvault list

# Then test
promptvault test <correct-id>
```

### "Test timeout"

**Cause:** API is slow or network issue

**Solution:**
- Check internet connection
- Try again with smaller input
- Increase timeout in code (currently 60s)

### Low Scores Despite Good Output

**Cause:** Similarity algorithm is basic word overlap

**Solution:**
- Adjust expected output to match likely phrasing
- Focus on key terms that must appear
- Consider score trends, not absolute values

---

## 📈 Advanced Features (Coming Soon)

### Test Suites
```bash
# Create named test suite
promptvault test create-suite abc123 "Regression Tests"

# Run suite
promptvault test run-suite "Regression Tests"
```

### Visual Regression
```bash
# Compare outputs across versions
promptvault test diff abc123 v1.0 v2.0
```

### Performance Testing
```bash
# Load test with 100 iterations
promptvault test perf abc123 --iterations 100
```

---

## 🎯 Next Steps

1. **Start testing** your most critical prompts today
2. **Build test suites** for important prompts
3. **Integrate with CI/CD** to catch regressions
4. **Share test results** with your team
5. **Track quality trends** over time

---

**Remember:** A tested prompt is a reliable prompt. 🧪✨

For more help: `promptvault test --help`
