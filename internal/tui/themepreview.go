package tui

import (
	"fmt"
	"strings"

	"github.com/Bharath-code/promptvault/internal/config"
	"github.com/charmbracelet/lipgloss"
)

type ThemePreview struct {
	themes   []string
	cursor   int
	selected string
	width    int
	height   int
}

func NewThemePreview(currentTheme string, width, height int) *ThemePreview {
	themes := []string{}
	for name := range config.PresetThemesList {
		themes = append(themes, name)
	}

	return &ThemePreview{
		themes:   themes,
		cursor:   0,
		selected: currentTheme,
		width:    width,
		height:   height,
	}
}

func (tp *ThemePreview) MoveUp() {
	if tp.cursor > 0 {
		tp.cursor--
	}
}

func (tp *ThemePreview) MoveDown() {
	if tp.cursor < len(tp.themes)-1 {
		tp.cursor++
	}
}

func (tp *ThemePreview) Select() string {
	tp.selected = tp.themes[tp.cursor]
	return tp.selected
}

func (tp *ThemePreview) Current() string {
	return tp.themes[tp.cursor]
}

func (tp *ThemePreview) Render() string {
	var lines []string

	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(color("#7C3AED")).
		PaddingBottom(1)

	lines = append(lines, headerStyle.Render("🎨 Theme Selector"))
	lines = append(lines, "")
	lines = append(lines, lipgloss.NewStyle().
		Foreground(color("#64748B")).
		Render("Current: "+tp.selected))
	lines = append(lines, "")

	for i, themeName := range tp.themes {
		colors := config.PresetThemesList[themeName]
		selected := i == tp.cursor

		preview := fmt.Sprintf(" %s %s %s %s %s ",
			wrapColor(colors.Primary, "■"),
			wrapColor(colors.Success, "●"),
			wrapColor(colors.Warning, "◆"),
			wrapColor(colors.Error, "▲"),
			wrapColor(colors.Accent, "★"),
		)

		line := ""
		if selected {
			line = lipgloss.JoinHorizontal(lipgloss.Center,
				lipgloss.NewStyle().
					Foreground(color("#7C3AED")).
					Bold(true).
					Render("►"),
				lipgloss.NewStyle().
					Background(color("#334155")).
					Foreground(color("#E2E8F0")).
					Padding(0, 1).
					Render(themeName),
				lipgloss.NewStyle().
					Foreground(color("#64748B")).
					Width(8).
					Render(""),
				lipgloss.NewStyle().
					Render(preview),
			)
		} else {
			line = lipgloss.JoinHorizontal(lipgloss.Center,
				lipgloss.NewStyle().
					Foreground(color("#334155")).
					Width(2).
					Render(""),
				lipgloss.NewStyle().
					Foreground(color("#94A3B8")).
					Width(12).
					Render(themeName),
				lipgloss.NewStyle().
					Foreground(color("#64748B")).
					Width(8).
					Render(""),
				lipgloss.NewStyle().
					Foreground(color("#475569")).
					Render(preview),
			)
		}

		lines = append(lines, line)
	}

	lines = append(lines, "")
	lines = append(lines, lipgloss.NewStyle().
		Foreground(color("#64748B")).
		Italic(true).
		Render("↑/↓ navigate  •  Enter select  •  Esc cancel"))

	content := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(color("#7C3AED")).
		Padding(1, 2).
		Width(70).
		Render(strings.Join(lines, "\n"))

	return lipgloss.NewStyle().
		Width(tp.width).
		Height(tp.height).
		Align(lipgloss.Center, lipgloss.Center).
		Render(content)
}

func wrapColor(color, char string) string {
	return lipgloss.NewStyle().
		Foreground(lipgloss.Color(color)).
		Render(char)
}

type KeybindingEditor struct {
	bindings map[string]string
	cursor   int
	keys     []string
	width    int
}

func NewKeybindingEditor(bindings config.Keybindings, width int) *KeybindingEditor {
	keys := []string{}
	for action, key := range bindings.Navigation {
		keys = append(keys, "nav:"+action+":"+key)
	}
	for action, key := range bindings.Actions {
		keys = append(keys, "action:"+action+":"+key)
	}
	for action, key := range bindings.QuickActions {
		keys = append(keys, "quick:"+action+":"+key)
	}

	return &KeybindingEditor{
		bindings: make(map[string]string),
		cursor:   0,
		keys:     keys,
		width:    width,
	}
}

func (ke *KeybindingEditor) Render() string {
	var lines []string

	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(color("#7C3AED")).
		PaddingBottom(1)

	lines = append(lines, headerStyle.Render("⌨️  Keybindings"))
	lines = append(lines, "")
	lines = append(lines, lipgloss.NewStyle().
		Foreground(color("#64748B")).
		Render("↑/↓ navigate  •  Esc done"))

	content := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(color("#7C3AED")).
		Padding(1, 2).
		Width(ke.width).
		Render(strings.Join(lines, "\n"))

	return lipgloss.NewStyle().
		Width(ke.width+10).
		Align(lipgloss.Center, lipgloss.Center).
		Render(content)
}
