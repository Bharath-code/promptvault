package tui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type ActionCategory string

const (
	CategoryNavigation ActionCategory = "Navigation"
	CategoryClipboard  ActionCategory = "Clipboard"
	CategoryEdit       ActionCategory = "Edit"
	CategoryView       ActionCategory = "View"
	CategoryExport     ActionCategory = "Export"
	CategoryMeta       ActionCategory = "Metadata"
)

type QuickAction struct {
	Label       string
	Key         string
	Icon        string
	Category    ActionCategory
	Shortcut    string
	Description string
}

var (
	actionPanelStyle = lipgloss.NewStyle().
				Border(lipgloss.NormalBorder()).
				BorderForeground(color("#334155")).
				Padding(1, 1)

	actionCategoryStyle = lipgloss.NewStyle().
				Foreground(color("#7C3AED")).
				Bold(true).
				PaddingBottom(1)

	actionItemStyle = lipgloss.NewStyle().
			Foreground(color("#94A3B8"))

	actionSelectedStyle = lipgloss.NewStyle().
				Background(color("#334155")).
				Foreground(color("#E2E8F0"))

	actionKeyStyle = lipgloss.NewStyle().
			Foreground(color("#06B6D4")).
			Background(color("#164E63")).
			Padding(0, 1).
			Bold(true)

	actionIconStyle = lipgloss.NewStyle().
			Foreground(color("#7C3AED"))
)

type QuickActionsPanel struct {
	actions []QuickAction
	cursor  int
	width   int
	height  int
	visible bool
}

func NewQuickActionsPanel(width, height int) *QuickActionsPanel {
	return &QuickActionsPanel{
		width:   width,
		height:  height,
		visible: false,
	}
}

func (qap *QuickActionsPanel) Toggle() {
	qap.visible = !qap.visible
}

func (qap *QuickActionsPanel) Show() {
	qap.visible = true
}

func (qap *QuickActionsPanel) Hide() {
	qap.visible = false
}

func (qap *QuickActionsPanel) IsVisible() bool {
	return qap.visible
}

func (qap *QuickActionsPanel) SetActions(actions []QuickAction) {
	qap.actions = actions
	qap.cursor = 0
}

func (qap *QuickActionsPanel) MoveUp() {
	if qap.cursor > 0 {
		qap.cursor--
	}
}

func (qap *QuickActionsPanel) MoveDown() {
	if qap.cursor < len(qap.actions)-1 {
		qap.cursor++
	}
}

func (qap *QuickActionsPanel) CurrentAction() *QuickAction {
	if qap.cursor >= 0 && qap.cursor < len(qap.actions) {
		return &qap.actions[qap.cursor]
	}
	return nil
}

func (qap *QuickActionsPanel) Execute() string {
	action := qap.CurrentAction()
	if action != nil {
		return action.Key
	}
	return ""
}

func (qap *QuickActionsPanel) Render() string {
	if !qap.visible || len(qap.actions) == 0 {
		return ""
	}

	var lines []string

	header := lipgloss.NewStyle().
		Bold(true).
		Foreground(color("#7C3AED")).
		PaddingBottom(1).
		Width(qap.width).
		Render("⚡ Actions")

	lines = append(lines, header)

	currentCategory := ActionCategory("")
	for i, action := range qap.actions {
		if action.Category != currentCategory {
			if currentCategory != "" {
				lines = append(lines, "")
			}
			lines = append(lines, actionCategoryStyle.Render(string(action.Category)))
			currentCategory = action.Category
		}

		selected := i == qap.cursor
		line := qap.renderAction(action, selected)
		lines = append(lines, line)
	}

	footer := lipgloss.NewStyle().
		Foreground(color("#64748B")).
		Italic(true).
		Width(qap.width).
		Render("↑/↓ navigate  •  Enter execute")

	lines = append(lines, "")
	lines = append(lines, footer)

	content := strings.Join(lines, "\n")

	return actionPanelStyle.
		Width(qap.width).
		Height(qap.height).
		Render(content)
}

func (qap *QuickActionsPanel) renderAction(action QuickAction, selected bool) string {
	icon := actionIconStyle.Render(action.Icon)
	key := actionKeyStyle.Render("[" + action.Shortcut + "]")
	label := action.Label

	if selected {
		return " " + icon + " " + key + " " + lipgloss.NewStyle().
			Background(color("#334155")).
			Foreground(color("#E2E8F0")).
			Render(label)
	}

	return " " + icon + " " + key + " " + actionItemStyle.Render(label)
}

type ActionsBuilder struct {
	actions []QuickAction
}

func NewActionsBuilder() *ActionsBuilder {
	return &ActionsBuilder{}
}

func (ab *ActionsBuilder) WithNavigation() *ActionsBuilder {
	ab.actions = append(ab.actions,
		QuickAction{
			Label:    "Search Prompts",
			Key:      "search",
			Icon:     "🔍",
			Category: CategoryNavigation,
			Shortcut: "/",
		},
		QuickAction{
			Label:    "Refresh List",
			Key:      "refresh",
			Icon:     "🔄",
			Category: CategoryNavigation,
			Shortcut: "r",
		},
		QuickAction{
			Label:    "Stack Browser",
			Key:      "stacks",
			Icon:     "📁",
			Category: CategoryNavigation,
			Shortcut: "t",
		},
	)
	return ab
}

func (ab *ActionsBuilder) WithClipboard() *ActionsBuilder {
	ab.actions = append(ab.actions,
		QuickAction{
			Label:    "Copy to Clipboard",
			Key:      "copy",
			Icon:     "📋",
			Category: CategoryClipboard,
			Shortcut: "Enter",
		},
		QuickAction{
			Label:    "Copy as JSON",
			Key:      "copy-json",
			Icon:     "📄",
			Category: CategoryClipboard,
			Shortcut: "J",
		},
	)
	return ab
}

func (ab *ActionsBuilder) WithEdit(hasSelection bool) *ActionsBuilder {
	ab.actions = append(ab.actions,
		QuickAction{
			Label:    "Add New Prompt",
			Key:      "add",
			Icon:     "➕",
			Category: CategoryEdit,
			Shortcut: "a",
		},
	)

	if hasSelection {
		ab.actions = append(ab.actions,
			QuickAction{
				Label:    "Edit Prompt",
				Key:      "edit",
				Icon:     "✏️",
				Category: CategoryEdit,
				Shortcut: "e",
			},
			QuickAction{
				Label:    "Delete Prompt",
				Key:      "delete",
				Icon:     "🗑️",
				Category: CategoryEdit,
				Shortcut: "d",
			},
			QuickAction{
				Label:    "Duplicate Prompt",
				Key:      "duplicate",
				Icon:     "📑",
				Category: CategoryEdit,
				Shortcut: "D",
			},
		)
	}
	return ab
}

func (ab *ActionsBuilder) WithView() *ActionsBuilder {
	ab.actions = append(ab.actions,
		QuickAction{
			Label:    "Toggle Preview",
			Key:      "preview",
			Icon:     "👁️",
			Category: CategoryView,
			Shortcut: "v",
		},
		QuickAction{
			Label:    "Toggle Recent",
			Key:      "recent",
			Icon:     "🔥",
			Category: CategoryView,
			Shortcut: "R",
		},
		QuickAction{
			Label:    "Statistics",
			Key:      "stats",
			Icon:     "📊",
			Category: CategoryView,
			Shortcut: "s",
		},
	)
	return ab
}

func (ab *ActionsBuilder) WithExport() *ActionsBuilder {
	ab.actions = append(ab.actions,
		QuickAction{
			Label:    "Export to skill.md",
			Key:      "export-skill",
			Icon:     "📤",
			Category: CategoryExport,
			Shortcut: "1",
		},
		QuickAction{
			Label:    "Export to .cursorrules",
			Key:      "export-cursor",
			Icon:     "📤",
			Category: CategoryExport,
			Shortcut: "2",
		},
		QuickAction{
			Label:    "Export to JSON",
			Key:      "export-json",
			Icon:     "📤",
			Category: CategoryExport,
			Shortcut: "3",
		},
	)
	return ab
}

func (ab *ActionsBuilder) WithMeta() *ActionsBuilder {
	ab.actions = append(ab.actions,
		QuickAction{
			Label:    "Toggle Verified",
			Key:      "verify",
			Icon:     "✓",
			Category: CategoryMeta,
			Shortcut: "V",
		},
		QuickAction{
			Label:    "View History",
			Key:      "history",
			Icon:     "📜",
			Category: CategoryMeta,
			Shortcut: "H",
		},
		QuickAction{
			Label:    "View Tests",
			Key:      "tests",
			Icon:     "🧪",
			Category: CategoryMeta,
			Shortcut: "T",
		},
	)
	return ab
}

func (ab *ActionsBuilder) WithAll(hasSelection bool) *ActionsBuilder {
	return ab.WithNavigation().
		WithClipboard().
		WithEdit(hasSelection).
		WithView().
		WithExport().
		WithMeta()
}

func (ab *ActionsBuilder) Build() []QuickAction {
	return ab.actions
}

func (ab *ActionsBuilder) Simple() []QuickAction {
	actions := []QuickAction{
		{Icon: "🔍", Label: "Search", Key: "search", Shortcut: "/", Category: CategoryNavigation},
		{Icon: "➕", Label: "Add", Key: "add", Shortcut: "a", Category: CategoryEdit},
		{Icon: "📋", Label: "Copy", Key: "copy", Shortcut: "Enter", Category: CategoryClipboard},
		{Icon: "✏️", Label: "Edit", Key: "edit", Shortcut: "e", Category: CategoryEdit},
		{Icon: "🗑️", Label: "Delete", Key: "delete", Shortcut: "d", Category: CategoryEdit},
		{Icon: "👁️", Label: "Preview", Key: "preview", Shortcut: "v", Category: CategoryView},
		{Icon: "📊", Label: "Stats", Key: "stats", Shortcut: "s", Category: CategoryView},
		{Icon: "📁", Label: "Stacks", Key: "stacks", Shortcut: "t", Category: CategoryNavigation},
		{Icon: "🔄", Label: "Refresh", Key: "refresh", Shortcut: "r", Category: CategoryNavigation},
		{Icon: "⚡", Label: "Themes", Key: "themes", Shortcut: "g", Category: CategoryView},
		{Icon: "❓", Label: "Help", Key: "help", Shortcut: "?", Category: CategoryNavigation},
	}
	return actions
}

func FormatActionResult(key, result string) string {
	return fmt.Sprintf("%s %s", actionIconStyle.Render("→"), result)
}
