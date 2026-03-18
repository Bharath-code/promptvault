package tui

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type SearchHistory struct {
	items  []string
	cursor int
	width  int
	height int
}

func NewSearchHistory(items []string, width, height int) *SearchHistory {
	return &SearchHistory{
		items:  items,
		cursor: 0,
		width:  width,
		height: height,
	}
}

func (sh *SearchHistory) SetItems(items []string) {
	sh.items = items
	if sh.cursor >= len(items) {
		sh.cursor = len(items) - 1
	}
	if sh.cursor < 0 {
		sh.cursor = 0
	}
}

func (sh *SearchHistory) MoveUp() {
	if sh.cursor > 0 {
		sh.cursor--
	}
}

func (sh *SearchHistory) MoveDown() {
	if sh.cursor < len(sh.items)-1 {
		sh.cursor++
	}
}

func (sh *SearchHistory) Selected() string {
	if sh.cursor >= 0 && sh.cursor < len(sh.items) {
		return sh.items[sh.cursor]
	}
	return ""
}

func (sh *SearchHistory) DeleteCurrent() string {
	if sh.cursor >= 0 && sh.cursor < len(sh.items) {
		item := sh.items[sh.cursor]
		sh.items = append(sh.items[:sh.cursor], sh.items[sh.cursor+1:]...)
		if sh.cursor >= len(sh.items) && sh.cursor > 0 {
			sh.cursor--
		}
		return item
	}
	return ""
}

func (sh *SearchHistory) Render() string {
	if len(sh.items) == 0 {
		return ""
	}

	var lines []string
	lines = append(lines, panelHeaderStyle.Render(" Recent Searches"))

	maxVisible := sh.height - 4
	if maxVisible < 0 {
		maxVisible = 0
	}

	start := 0
	if sh.cursor >= maxVisible {
		start = sh.cursor - maxVisible + 1
	}
	end := start + maxVisible
	if end > len(sh.items) {
		end = len(sh.items)
	}

	for i := start; i < end; i++ {
		prefix := "  "
		if i == sh.cursor {
			prefix = selectedItemStyle.Render("> ")
			line := lipgloss.JoinHorizontal(lipgloss.Center,
				prefix,
				lipgloss.NewStyle().Foreground(colorText).Render(sh.items[i]),
			)
			lines = append(lines, line)
		} else {
			line := lipgloss.JoinHorizontal(lipgloss.Center,
				prefix,
				lipgloss.NewStyle().Foreground(colorMuted).Render(sh.items[i]),
			)
			lines = append(lines, line)
		}
	}

	lines = append(lines, "")
	lines = append(lines, helpStyle.Render("↑/↓ select  •  Enter use  •  d delete"))

	content := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorPrimary).
		Padding(1, 2).
		Width(sh.width).
		Render(strings.Join(lines, "\n"))

	return lipgloss.Place(sh.width, sh.height,
		lipgloss.Center, lipgloss.Center,
		content)
}

func (sh *SearchHistory) RenderInline() string {
	if len(sh.items) == 0 {
		return ""
	}

	var lines []string
	maxItems := 5
	if maxItems > len(sh.items) {
		maxItems = len(sh.items)
	}

	for i := 0; i < maxItems; i++ {
		prefix := "  "
		style := helpStyle
		if i == sh.cursor {
			prefix = successStyle.Render("> ")
			style = lipgloss.NewStyle().Foreground(colorText)
		}
		lines = append(lines, lipgloss.JoinHorizontal(lipgloss.Center, prefix, style.Render(sh.items[i])))
	}

	return lipgloss.JoinVertical(lipgloss.Left, lines...)
}
