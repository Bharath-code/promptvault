# 🤖 AI-Assisted Authoring Guide

**Create better prompts with AI-powered recommendations.**

---

## 🎯 Why AI-Assisted Authoring?

Writing effective prompts is hard. AI-assisted authoring helps you:
- ✅ Detect variables automatically
- ✅ Get smart tag and stack suggestions
- ✅ Catch anti-patterns before saving
- ✅ Improve prompt quality with AI feedback
- ✅ Learn prompt engineering best practices

---

## 🚀 Quick Start

### 1. AI-Assisted Creation

```bash
# Start interactive AI-assisted creation
promptvault create --ai

# Or standard interactive creation
promptvault create
```

### 2. What You Get

**During creation, AI analyzes:**
- Variables (`{{variable}}` syntax)
- Suggested tags based on content
- Recommended tech stack
- Anti-patterns to fix
- Quality score (0-100)
- Improvement suggestions

---

## 📋 Features

### Variable Detection

Automatically detects `{{variable}}` patterns:

**Example:**
```
Convert this {{language}} code to {{target_language}}.
Use {{framework}} framework.

Variables detected: language, target_language, framework
```

**Benefits:**
- Makes prompts reusable
- Clear parameter documentation
- Easy to fill in TUI

---

### Tag Recommendation

AI suggests relevant tags based on content analysis:

**Content Keywords → Suggested Tags:**
| Keywords | Tags |
|----------|------|
| react, component, hooks | `react`, `frontend`, `hooks` |
| typescript, type, interface | `typescript`, `types`, `javascript` |
| docker, container, image | `docker`, `containers`, `devops` |
| test, unit, coverage | `testing`, `unit-test`, `quality` |
| debug, error, fix | `debugging`, `troubleshooting`, `fix` |

---

### Stack Auto-Detection

Suggests appropriate tech stack:

**Detection Rules:**
```
"react" → frontend/react
"nextjs" → frontend/react/nextjs
"typescript" → frontend/typescript
"python" → backend/python
"django" → backend/python/django
"fastapi" → backend/python/fastapi
"go" → backend/go
"docker" → devops/docker
"kubernetes" → devops/kubernetes
"terraform" → devops/terraform
"postgresql" → database/postgresql
```

---

### Anti-Pattern Detection

Catches common prompt mistakes:

| Anti-Pattern | Detection | Suggestion |
|--------------|-----------|------------|
| Too short | < 20 words | "Add more context and specifics" |
| No examples | Missing "example" keyword | "Consider adding examples" |
| No output format | Missing "output/format/return" | "Specify expected output format" |
| ALL CAPS | All uppercase | "Use normal casing" |
| Too nested | > 5 if/then/else | "Simplify the logic" |
| No constraints | Missing don't/avoid/never | "Add constraints" |

---

### Quality Scoring

Calculates 0-100 score based on:

| Factor | Points |
|--------|--------|
| Base score | 50 |
| > 50 words | +10 |
| > 100 words | +10 |
| > 200 words | +5 |
| Has structure (numbered/bulleted) | +10 |
| Has examples | +10 |
| Has variables | +5 |
| Has constraints | +5 |
| **Maximum** | **100** |

**Score Interpretation:**
- **90-100**: Excellent prompt
- **70-89**: Good prompt, minor improvements possible
- **50-69**: Decent prompt, needs work
- **< 50**: Poor prompt, significant revision needed

---

## 🎓 Usage Examples

### Example 1: Creating a React Prompt

```bash
$ promptvault create --ai

🤖 AI-Assisted Prompt Creation
────────────────────────────────────────────────────────────

📝 Prompt Title: React Component with Accessibility

📄 Enter prompt content (type 'DONE' on a new line to finish):
────────────────────────────────────────────────────────────
Create a React component that follows WCAG 2.1 AA guidelines.
Include proper ARIA labels, keyboard navigation, and focus management.
The component should accept {{label}} and {{onChange}} props.
DONE

🔍 Analyzing prompt...

📊 Analysis Results:
────────────────────────────────────────────────────────────
🏷️  Variables detected: label, onChange
📚 Suggested stack: frontend/react
🏷️  Suggested tags: react, frontend, accessibility
⭐ Quality score: 75/100

⚠️  Anti-patterns found:
   💡 Consider adding examples for clarity
   📝 Specify the expected output format

💡 Suggested improvements:
   • Add a role definition (e.g., "You are an expert React developer")
   • Provide more context about the use case

Continue with these suggestions? [Y/n]: y

📚 Tech stack (e.g., frontend/react/hooks): frontend/react
🏷️  Tags (comma-separated): accessibility, a11y, wcag
🤖 Models (comma-separated): claude-sonnet

✓ Created prompt: React Component with Accessibility
ℹ  ID: abc123
ℹ  Stack: frontend/react
ℹ  Tags: accessibility, a11y, wcag
```

---

### Example 2: Quick Creation Without AI

```bash
$ promptvault create

📝 Create New Prompt
────────────────────────────────────────────────────────────

📝 Title: Python FastAPI Endpoint

📄 Content (type 'DONE' on a new line to finish):
────────────────────────────────────────────────────────────
Create a FastAPI endpoint with:
1. Pydantic validation
2. Error handling
3. Proper status codes
DONE

📚 Tech stack: backend/python/fastapi
🏷️  Tags: api, validation, pydantic
🤖 Models: claude-sonnet, gpt-4o

✓ Created prompt: Python FastAPI Endpoint
```

---

## 🎯 Best Practices

### 1. Use Variables for Reusability

**❌ Bad:**
```
Convert this JavaScript to TypeScript
```

**✅ Good:**
```
Convert this {{source_language}} code to {{target_language}}
```

### 2. Provide Context

**❌ Bad:**
```
Fix the bug
```

**✅ Good:**
```
You are a senior {{language}} developer. Review this code for {{issue_type}} bugs.
Focus on {{specific_concern}} issues.
```

### 3. Specify Output Format

**❌ Bad:**
```
Explain the code
```

**✅ Good:**
```
Explain the code in this format:
1. What it does (1-2 sentences)
2. Key functions
3. Potential issues
4. Suggested improvements
```

### 4. Add Examples

**❌ Bad:**
```
Create a function
```

**✅ Good:**
```
Create a function. For example:

Input: [1, 2, 3]
Output: [3, 2, 1]

Input: "hello"
Output: "olleh"
```

### 5. Set Constraints

**❌ Bad:**
```
Write the code
```

**✅ Good:**
```
Write the code. Don't use external libraries.
Avoid recursion. Never mutate the input array.
```

---

## 🔧 Configuration

### API Keys (for AI Analysis)

AI-powered suggestions require API keys:

```bash
# Set Anthropic API key
export ANTHROPIC_API_KEY=sk-ant-...

# Set OpenAI API key
export OPENAI_API_KEY=sk-...
```

**Without API keys:**
- ✅ Variable detection works (rule-based)
- ✅ Tag suggestions work (keyword-based)
- ✅ Stack detection works (keyword-based)
- ✅ Anti-pattern detection works (rule-based)
- ✅ Quality scoring works (rule-based)
- ❌ AI-powered improvements unavailable

---

## 📊 Quality Score Breakdown

### How to Improve Your Score

| Current Score | Action | Potential Gain |
|---------------|--------|----------------|
| < 50 | Add more content (aim for 100+ words) | +20 |
| < 60 | Add numbered/bulleted structure | +10 |
| < 70 | Add examples | +10 |
| < 80 | Add variable placeholders | +5 |
| < 90 | Add constraints (don't/avoid/never) | +5 |

### Score Examples

**Score: 50 (Base)**
```
Fix this code.
```

**Score: 75 (Base + Length + Structure)**
```
Fix this code:
1. Check for null values
2. Validate input types
3. Handle edge cases
```

**Score: 95 (All bonuses)**
```
You are a senior developer. Fix this code:

1. Check for null values
2. Validate input types  
3. Handle edge cases

For example:
Input: null
Output: Error message

Don't use external libraries.
Avoid mutating the input.
Use {{language}} syntax.
```

---

## 🐛 Troubleshooting

### "AI analysis failed"

**Cause:** API key not set or API error

**Solution:**
```bash
# Check API key
echo $ANTHROPIC_API_KEY

# If empty, set it
export ANTHROPIC_API_KEY=sk-ant-...
```

**Note:** Creation continues with rule-based analysis.

### "No suggestions shown"

**Cause:** Content too short for analysis

**Solution:** Add more content (aim for 50+ words)

### Variables not detected

**Cause:** Wrong syntax

**Solution:** Use `{{variable}}` with double braces:
```
✅ {{variable}}
❌ {variable}
❌ <variable>
```

---

## 📈 Future Features

Coming soon:
- [ ] Real-time AI feedback as you type
- [ ] Auto-rewrite suggestions
- [ ] Model-specific optimization tips
- [ ] Prompt template library
- [ ] A/B testing for prompts
- [ ] Performance benchmarking across models

---

## 🔗 Related Commands

- `promptvault add` - Quick add without AI
- `promptvault test` - Test your prompt
- `promptvault edit` - Edit with version history

---

**Remember:** AI assistance is a tool, not a replacement for thinking. Always review suggestions critically! 🤖✨

For more help: `promptvault create --help`
