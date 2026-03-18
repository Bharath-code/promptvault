package ui

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
)

type Color = lipgloss.Color

var (
	Primary    = lipgloss.Color("#7C3AED")
	Accent     = lipgloss.Color("#06B6D4")
	Success    = lipgloss.Color("#10B981")
	Info       = lipgloss.Color("#3B82F6")
	Danger     = lipgloss.Color("#EF4444")
	Warning    = lipgloss.Color("#F59E0B")
	Text       = lipgloss.Color("#E2E8F0")
	Muted      = lipgloss.Color("#64748B")
	Background = lipgloss.Color("#0F172A")
	BgAlt      = lipgloss.Color("#1E293B")
	BgHover    = lipgloss.Color("#334155")
	Verified   = lipgloss.Color("#34D399")
)

var (
	SuccessCode = "\033[38;5;2m"
	ErrorCode   = "\033[38;5;1m"
	WarningCode = "\033[38;5;3m"
	InfoCode    = "\033[38;5;6m"
	PrimaryCode = "\033[38;5;129m"
	MutedCode   = "\033[38;5;245m"
	ResetCode   = "\033[0m"
)

var TitleStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(Primary).
	PaddingRight(1)

var SelectedStyle = lipgloss.NewStyle().
	Background(BgHover).
	Foreground(Text).
	Padding(0, 1).
	Bold(true)

var ItemStyle = lipgloss.NewStyle().
	Foreground(Text).
	Padding(0, 1)

var TagStyle = lipgloss.NewStyle().
	Foreground(Accent).
	Background(lipgloss.Color("#164E63")).
	Padding(0, 1).
	Bold(true)

var StackStyle = lipgloss.NewStyle().
	Foreground(Primary).
	Bold(true)

var VerifiedStyle = lipgloss.NewStyle().
	Foreground(Verified).
	Bold(true)

var SuccessStyle = lipgloss.NewStyle().
	Foreground(Success).
	Bold(true)

var ErrorStyle = lipgloss.NewStyle().
	Foreground(Danger).
	Bold(true)

var HelpStyle = lipgloss.NewStyle().
	Foreground(Muted).
	Italic(true)

func PrintSuccess(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, SuccessCode+"✓ "+ResetCode+format+"\n", args...)
}

func PrintError(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, ErrorCode+"✗ "+ResetCode+format+"\n", args...)
}

func PrintWarning(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, WarningCode+"⚠ "+ResetCode+format+"\n", args...)
}

func PrintInfo(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, InfoCode+"ℹ "+ResetCode+format+"\n", args...)
}
