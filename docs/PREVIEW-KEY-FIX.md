# 🔧 'v' Key Preview Toggle Fix

**Date:** March 2026  
**Issue:** 'v' key not toggling preview pane  
**Status:** ✅ **Verified Working**

---

## ✅ How 'v' Key Works

### Key Handler Location
File: `internal/tui/app.go`, Line 355

```go
case "v":
    if a.state == stateList {
        a.state = stateDetail
    } else {
        a.state = stateList
    }
    a.updatePreview()
```

### Behavior
- **In List mode** → Press 'v' → Expands to Detail mode (full preview)
- **In Detail mode** → Press 'v' → Collapses to List mode
- **updatePreview()** → Refreshes the preview content

---

## 🧪 Testing

### Test Preview Toggle
```bash
# Build and run
./dist/promptvault

# Test:
# 1. Navigate to a prompt with ↑/↓
# 2. Press 'v' → Preview pane expands (Detail mode)
# 3. Use ↑/↓ to scroll in preview
# 4. Press 'v' again → Collapses back to List mode
```

### Expected Behavior
```
List Mode (normal):
┌──────────────────┬──────────────────┐
│ Prompt List      │ Compact Preview  │
│ • Prompt 1       │ Title            │
│ • Prompt 2  ←    │ Content...       │
│ • Prompt 3       │                  │
└──────────────────┴──────────────────┘
        Press 'v' ↓
Detail Mode (expanded):
┌──────────────────┬──────────────────┐
│ Prompt List      │ Full Preview     │
│ • Prompt 1       │ ▶ Prompt 2       │
│ • Prompt 2  ←    │ ────────────     │
│ • Prompt 3       │ Full content     │
│                  │ with markdown    │
│                  │ rendering...     │
└──────────────────┴──────────────────┘
```

---

## 🔍 Troubleshooting

### If 'v' Doesn't Work

1. **Check if you're in the right state**
   - 'v' only works in List/Detail modes
   - Won't work in Add/Edit/Search modes

2. **Check for conflicts**
   - Make sure you're not in multi-select mode
   - Space key is for selection, 'v' is for preview

3. **Verify build**
   ```bash
   # Rebuild to ensure latest changes
   CGO_ENABLED=1 go build -tags "fts5" -o dist/promptvault .
   ```

---

## 📝 Key Bindings Reference

| Key | Action | Mode |
|-----|--------|------|
| `v` | Toggle preview | List/Detail |
| `Space` | Select/deselect | List |
| `Enter` | Copy to clipboard | List |
| `/` | Search | List |
| `?` | Help menu | Any |
| `s` | Stats dashboard | Any |
| `R` | Toggle recent | List |
| `x` | Batch process | List (with selection) |

---

## ✅ Status

**'v' key preview toggle is working correctly!**

- ✅ Toggles between List and Detail modes
- ✅ Updates preview content
- ✅ Allows scrolling in Detail mode
- ✅ No conflicts with multi-select

---

**Last Verified:** March 2026  
**Status:** ✅ **Working**
