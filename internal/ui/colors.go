package ui

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
)

var noColor = os.Getenv("NO_COLOR") != ""

func c(hex string) lipgloss.Color {
	if noColor {
		return lipgloss.Color("")
	}
	return lipgloss.Color(hex)
}

type Color = lipgloss.Color

var (
	Primary    = c("#7C3AED")
	Accent     = c("#06B6D4")
	Success    = c("#10B981")
	Info       = c("#3B82F6")
	Danger     = c("#EF4444")
	Warning    = c("#F59E0B")
	Text       = c("#E2E8F0")
	Muted      = c("#64748B")
	Background = c("#0F172A")
	BgAlt      = c("#1E293B")
	BgHover    = c("#334155")
	Verified   = c("#34D399")
)

var (
	SuccessCode = ansiColor(2)
	ErrorCode   = ansiColor(1)
	WarningCode = ansiColor(3)
	InfoCode    = ansiColor(6)
	PrimaryCode = ansiColor(129)
	MutedCode   = ansiColor(245)
	ResetCode   = "\033[0m"
)

func ansiColor(code int) string {
	if noColor {
		return ""
	}
	return fmt.Sprintf("\033[38;5;%dm", code)
}

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
	Background(c("#164E63")).
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
