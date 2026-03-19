package tui

import (
	"strings"

	"github.com/Bharath-code/promptvault/internal/config"
	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type HelpOverlay struct {
	width       int
	height      int
	state       viewState
	vim         *VimModeHandler
	keybindings config.Keybindings
}

func NewHelpOverlay(width, height int, state viewState, vim *VimModeHandler, kb config.Keybindings) *HelpOverlay {
	return &HelpOverlay{
		width:       width,
		height:      height,
		state:       state,
		vim:         vim,
		keybindings: kb,
	}
}

func (h *HelpOverlay) Init() tea.Cmd {
	return nil
}

func (h *HelpOverlay) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return h, nil
}

func (h *HelpOverlay) View() string {
	colWidth := (h.width - 10) / 2
	if colWidth < 30 {
		colWidth = h.width - 10
	}

	leftCol := h.renderLeftColumn(colWidth)
	rightCol := h.renderRightColumn(colWidth)

	panel := lipgloss.JoinHorizontal(lipgloss.Top, leftCol, rightCol)

	return lipgloss.NewStyle().
		Width(h.width).
		Height(h.height).
		Align(lipgloss.Center, lipgloss.Center).
		Render(
			lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(colorPrimary).
				Padding(1, 2).
				Render(
					lipgloss.JoinVertical(lipgloss.Left,
						h.renderHeader(),
						lipgloss.NewStyle().
							Foreground(colorMuted).
							Render(strings.Repeat("─", max(20, h.width-12))),
						panel,
						lipgloss.NewStyle().
							Foreground(colorMuted).
							Render(strings.Repeat("─", max(20, h.width-12))),
						h.renderFooter(),
					),
				),
		)
}

func (h *HelpOverlay) renderHeader() string {
	title := lipgloss.NewStyle().Bold(true).Foreground(colorPrimary).Render("⚡ PromptVault Help")

	var hint string
	switch h.state {
	case stateOnboarding:
		hint = lipgloss.NewStyle().Foreground(colorWarning).Render(" [Onboarding Tour]")
	case stateStackTree:
		hint = lipgloss.NewStyle().Foreground(colorInfo).Render(" [Stack Tree]")
	case stateSearch:
		hint = lipgloss.NewStyle().Foreground(colorInfo).Render(" [Search]")
	case stateConfig:
		hint = lipgloss.NewStyle().Foreground(colorInfo).Render(" [Config]")
	}

	subtitle := lipgloss.NewStyle().Foreground(colorMuted).Render("Keyboard shortcuts reference")

	var parts []string
	if hint != "" {
		parts = []string{title, hint, "  ", subtitle}
	} else {
		parts = []string{title, "  ", subtitle}
	}
	return lipgloss.JoinHorizontal(lipgloss.Left, parts...)
}

func (h *HelpOverlay) renderFooter() string {
	vimHint := ""
	if h.vim != nil && h.vim.Enabled {
		vimHint = "  •  Vim: " + h.vim.RenderModeIndicator()
	}
	return lipgloss.NewStyle().
		Foreground(colorMuted).
		Render("↑↓ scroll  •  Esc / ? to close" + vimHint)
}

func (h *HelpOverlay) renderLeftColumn(width int) string {
	var lines []string

	lines = append(lines, h.section("Navigation"))
	lines = append(lines, h.kl("↑ / ↓   or   k / j", "Navigate prompts"))
	lines = append(lines, h.kl("Enter", "Copy to clipboard"))
	lines = append(lines, h.kl("Space", "Multi-select"))
	lines = append(lines, h.kl("Tab", "Quick actions panel"))
	lines = append(lines, h.kl("/   or   s", "Search prompts"))
	lines = append(lines, h.kl("Click", "Select prompt"))
	lines = append(lines, h.kl("Scroll", "Navigate/scroll"))
	lines = append(lines, h.kl("Right-click", "Delete prompt"))

	lines = append(lines, "")
	lines = append(lines, h.section("Actions"))
	lines = append(lines, h.kl("a", "Add new prompt"))
	lines = append(lines, h.kl("e", "Edit selected"))
	lines = append(lines, h.kl("d", "Delete selected"))
	lines = append(lines, h.kl("v", "Toggle preview"))
	lines = append(lines, h.kl("c", "Copy selected"))
	lines = append(lines, h.kl("r", "Refresh list"))
	lines = append(lines, h.kl("R", "Recent prompts"))

	lines = append(lines, "")
	lines = append(lines, h.section("Quick Actions"))
	lines = append(lines, h.kl("t", "Stack tree view"))
	lines = append(lines, h.kl("g", "Theme picker"))
	lines = append(lines, h.kl("s", "Statistics"))
	lines = append(lines, h.kl("x", "Batch process"))
	lines = append(lines, h.kl("h", "Search history"))
	lines = append(lines, h.kl("0", "Go to top"))
	lines = append(lines, h.kl("$", "Go to bottom"))

	if h.state == stateSearch {
		lines = append(lines, "")
		lines = append(lines, h.section("Search"))
		lines = append(lines, h.kl("↑ / ↓", "History navigation"))
		lines = append(lines, h.kl("Ctrl+U", "Clear search"))
		lines = append(lines, h.kl("Enter", "Execute search"))
		lines = append(lines, h.kl("Esc", "Close search"))
	}

	if h.state == stateStackTree {
		lines = append(lines, "")
		lines = append(lines, h.section("Stack Tree"))
		lines = append(lines, h.kl("↑ / ↓   or   k / j", "Navigate nodes"))
		lines = append(lines, h.kl("→ / l", "Expand node"))
		lines = append(lines, h.kl("← / h", "Collapse node"))
		lines = append(lines, h.kl("Enter", "Filter by stack"))
		lines = append(lines, h.kl("Esc", "Close tree"))
	}

	if h.state == stateOnboarding {
		lines = append(lines, "")
		lines = append(lines, h.section("Onboarding Tour"))
		lines = append(lines, h.kl("← / h", "Previous step"))
		lines = append(lines, h.kl("→ / l", "Next step"))
		lines = append(lines, h.kl("Enter / Space", "Advance"))
		lines = append(lines, h.kl("Esc", "Skip tour"))
	}

	return lipgloss.NewStyle().Width(width).Render(strings.Join(lines, "\n"))
}

func (h *HelpOverlay) renderRightColumn(width int) string {
	var lines []string

	if h.vim != nil && h.vim.Enabled {
		lines = append(lines, h.section("Vim Mode"))
		lines = append(lines, h.kl(string(h.vim.Mode), "Current mode"))
		lines = append(lines, "")
		lines = append(lines, h.kl("i / a / A / o / O", "Insert mode"))
		lines = append(lines, h.kl("h / j / k / l", "Navigate (←↓↑→)"))
		lines = append(lines, h.kl("yy", "Yank (copy)"))
		lines = append(lines, h.kl("dd", "Delete prompt"))
		lines = append(lines, h.kl("p", "Put (paste)"))
		lines = append(lines, h.kl("v / V", "Visual mode"))
		lines = append(lines, h.kl("y / d", "Yank / Delete"))
		lines = append(lines, h.kl("/ : g G", "Search / Cmd / Top / Bot"))
		lines = append(lines, h.kl("Esc", "Normal mode"))
		lines = append(lines, "")
		lines = append(lines, h.section("Command Palette"))
	} else {
		lines = append(lines, h.section("Command Palette"))
	}

	lines = append(lines, h.kl(":", "Open command palette"))
	lines = append(lines, h.kl("?", "This help overlay"))

	lines = append(lines, "")
	lines = append(lines, h.section("Global"))
	lines = append(lines, h.kl("Esc", "Clear / Cancel / Back"))
	lines = append(lines, h.kl("q / Ctrl+C", "Quit PromptVault"))

	lines = append(lines, "")
	lines = append(lines, h.section("CLI Flags"))
	lines = append(lines, h.kl("--json", "Machine-readable output"))
	lines = append(lines, h.kl("-v / --verbose", "Verbose output"))
	lines = append(lines, h.kl("-d / --debug", "Debug mode"))
	lines = append(lines, h.kl("--stack", "Filter by stack"))

	lines = append(lines, "")
	lines = append(lines, h.section("Shell / Piping"))
	lines = append(lines, h.kl("stdout", "Data (pipeable → grep/jq)"))
	lines = append(lines, h.kl("stderr", "Errors & logs (colored)"))

	return lipgloss.NewStyle().Width(width).Render(strings.Join(lines, "\n"))
}

func (h *HelpOverlay) section(title string) string {
	return "\n" + lipgloss.NewStyle().
		Bold(true).
		Foreground(colorPrimary).
		Render("▸ "+title)
}

func (h *HelpOverlay) kl(key, desc string) string {
	colWidth := (h.width - 10) / 2
	if colWidth < 30 {
		colWidth = h.width - 10
	}
	keyWidth := colWidth/3 + 1

	keyStyle := lipgloss.NewStyle().
		Foreground(colorAccent).
		Bold(true).
		Width(keyWidth)

	descStyle := lipgloss.NewStyle().Foreground(colorText)

	return keyStyle.Render(key) + "  " + descStyle.Render(desc)
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
