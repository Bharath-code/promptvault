package tui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"
)

type OnboardingStep struct {
	Title       string
	Description string
	Shortcut    string
	Icon        string
}

var onboardingSteps = []OnboardingStep{
	{
		Title:       "Welcome to PromptVault!",
		Description: "Your universal prompt OS for developers. Store, search, and deploy AI prompts by tech stack.",
		Shortcut:    "",
		Icon:        "⚡",
	},
	{
		Title:       "Navigate Prompts",
		Description: "Use ↑/↓ arrow keys or k/j to navigate through your prompts. The preview updates automatically.",
		Shortcut:    "↑↓ or k/j",
		Icon:        "📋",
	},
	{
		Title:       "Quick Search",
		Description: "Press / to search prompts by title, content, tags, or stack. Fuzzy matching finds results as you type.",
		Shortcut:    "/",
		Icon:        "🔍",
	},
	{
		Title:       "Copy to Clipboard",
		Description: "Press Enter to copy any prompt directly to your clipboard. Variables will be filled interactively.",
		Shortcut:    "Enter",
		Icon:        "📋",
	},
	{
		Title:       "Add New Prompts",
		Description: "Press a to create a new prompt. Fill in title, content, stack, and tags. Ctrl+S to save.",
		Shortcut:    "a",
		Icon:        "➕",
	},
	{
		Title:       "Edit & Delete",
		Description: "Press e to edit the selected prompt, or d to delete it. Confirm with y/n.",
		Shortcut:    "e or d",
		Icon:        "✏️",
	},
	{
		Title:       "Multi-Select",
		Description: "Press Space to select multiple prompts. Press x to batch process them.",
		Shortcut:    "Space, x",
		Icon:        "☑️",
	},
	{
		Title:       "Full-Screen Preview",
		Description: "Press v to toggle full-screen preview mode for immersive reading.",
		Shortcut:    "v",
		Icon:        "🖥️",
	},
	{
		Title:       "Command Palette",
		Description: "Press : to open the command palette for quick access to all commands.",
		Shortcut:    ":",
		Icon:        "🎯",
	},
	{
		Title:       "Statistics",
		Description: "Press s to view usage statistics. See your most used prompts and stack distribution.",
		Shortcut:    "s",
		Icon:        "📊",
	},
	{
		Title:       "Help",
		Description: "Press ? anytime to see all keyboard shortcuts. Happy prompting!",
		Shortcut:    "?",
		Icon:        "❓",
	},
}

type OnboardingTour struct {
	currentStep int
	totalSteps  int
}

func NewOnboardingTour() *OnboardingTour {
	return &OnboardingTour{
		currentStep: 0,
		totalSteps:  len(onboardingSteps),
	}
}

func (ot *OnboardingTour) Next() bool {
	if ot.currentStep < ot.totalSteps-1 {
		ot.currentStep++
		return true
	}
	return false
}

func (ot *OnboardingTour) Previous() bool {
	if ot.currentStep > 0 {
		ot.currentStep--
		return true
	}
	return false
}

func (ot *OnboardingTour) IsLast() bool {
	return ot.currentStep == ot.totalSteps-1
}

func (ot *OnboardingTour) IsFirst() bool {
	return ot.currentStep == 0
}

func (ot *OnboardingTour) CurrentStep() OnboardingStep {
	return onboardingSteps[ot.currentStep]
}

func (ot *OnboardingTour) Progress() string {
	return fmt.Sprintf("%d/%d", ot.currentStep+1, ot.totalSteps)
}

func (ot *OnboardingTour) Render(width, height int) string {
	step := ot.CurrentStep()

	progress := ot.Progress()
	progressStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#64748B")).
		Align(lipgloss.Right)

	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#7C3AED")).
		PaddingBottom(1)

	descStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#94A3B8")).
		Width(width - 20)

	shortcutStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#06B6D4")).
		Background(lipgloss.Color("#164E63")).
		Padding(0, 1).
		Bold(true)

	navigationStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#64748B"))

	skipStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("#475569")).
		Italic(true)

	dots := ot.renderProgressDots()

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		lipgloss.JoinHorizontal(lipgloss.Center,
			lipgloss.NewStyle().
				Foreground(lipgloss.Color("#7C3AED")).
				Bold(true).
				Render("⚡ Welcome to PromptVault"),
			progressStyle.Width(20).Render(progress),
		),
		"",
		lipgloss.NewStyle().
			Width(width).
			Align(lipgloss.Center).
			Render(step.Icon),
		"",
		headerStyle.Width(width).Render(step.Title),
		"",
		descStyle.Render(step.Description),
		"",
	)

	if step.Shortcut != "" {
		content = lipgloss.JoinVertical(
			lipgloss.Left,
			content,
			lipgloss.JoinHorizontal(lipgloss.Center,
				lipgloss.NewStyle().Width(10).Render(""),
				shortcutStyle.Render(step.Shortcut),
				lipgloss.NewStyle().Width(10).Render(""),
			),
		)
	}

	content = lipgloss.JoinVertical(
		lipgloss.Left,
		content,
		"",
		dots,
		"",
		lipgloss.JoinHorizontal(lipgloss.Center,
			navigationStyle.Render("← Back"),
			lipgloss.NewStyle().Width(width/3).Render(""),
			navigationStyle.Render("Next →"),
		),
		"",
		lipgloss.NewStyle().
			Width(width).
			Align(lipgloss.Center).
			Render(skipStyle.Render("Press Esc to skip tutorial")),
	)

	container := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("#7C3AED")).
		BorderBackground(lipgloss.Color("#1E293B")).
		Padding(2, 4).
		Width(60).
		Background(lipgloss.Color("#0F172A"))

	box := container.Render(content)

	return lipgloss.NewStyle().
		Width(width).
		Height(height).
		Align(lipgloss.Center, lipgloss.Center).
		Render(box)
}

func (ot *OnboardingTour) renderProgressDots() string {
	var dots []string
	for i := 0; i < ot.totalSteps; i++ {
		if i == ot.currentStep {
			dots = append(dots, lipgloss.NewStyle().
				Foreground(lipgloss.Color("#7C3AED")).
				Bold(true).
				Render("●"))
		} else if i < ot.currentStep {
			dots = append(dots, lipgloss.NewStyle().
				Foreground(lipgloss.Color("#7C3AED")).
				Render("●"))
		} else {
			dots = append(dots, lipgloss.NewStyle().
				Foreground(lipgloss.Color("#334155")).
				Render("○"))
		}
	}
	return lipgloss.JoinHorizontal(lipgloss.Center, dots...)
}
