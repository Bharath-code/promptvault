package tui

import "github.com/charmbracelet/lipgloss"

// ── Color Palette ─────────────────────────────────────────────
var (
	colorPrimary  = lipgloss.Color("#7C3AED") // Vibrant purple
	colorAccent   = lipgloss.Color("#06B6D4") // Cyan accent
	colorSuccess  = lipgloss.Color("#10B981") // Emerald green
	colorDanger   = lipgloss.Color("#EF4444") // Red
	colorWarning  = lipgloss.Color("#F59E0B") // Amber
	colorText     = lipgloss.Color("#E2E8F0") // Slate-200
	colorMuted    = lipgloss.Color("#64748B") // Slate-500
	colorBg       = lipgloss.Color("#0F172A") // Slate-900
	colorBgAlt    = lipgloss.Color("#1E293B") // Slate-800
	colorBgHover  = lipgloss.Color("#334155") // Slate-700
	colorVerified = lipgloss.Color("#34D399") // Emerald-400
)

// ── Header ────────────────────────────────────────────────────
var titleStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(colorPrimary).
	PaddingRight(1)

var subtitleStyle = lipgloss.NewStyle().
	Foreground(colorMuted)

// ── Search ────────────────────────────────────────────────────
var searchStyle = lipgloss.NewStyle().
	Background(colorBgHover).
	Foreground(colorText).
	Padding(0, 1)

// ── Help / Keybindings ────────────────────────────────────────
var helpStyle = lipgloss.NewStyle().
	Foreground(colorMuted).
	Italic(true)

// ── Tags / Stack / Metadata ───────────────────────────────────
var tagStyle = lipgloss.NewStyle().
	Foreground(colorAccent).
	Background(lipgloss.Color("#164E63")).
	Padding(0, 1).
	Bold(true)

var stackStyle = lipgloss.NewStyle().
	Foreground(colorPrimary).
	Bold(true)

var verifiedStyle = lipgloss.NewStyle().
	Foreground(colorVerified).
	Bold(true)

var usageStyle = lipgloss.NewStyle().
	Foreground(colorMuted).
	Italic(true)

var scoreStyle = lipgloss.NewStyle().
	Foreground(colorAccent).
	Bold(true)

// ── List Items ────────────────────────────────────────────────
var selectedItemStyle = lipgloss.NewStyle().
	Background(colorBgHover).
	Foreground(colorText).
	Padding(0, 1).
	Bold(true)

var itemStyle = lipgloss.NewStyle().
	Foreground(colorText).
	Padding(0, 1)

// ── Preview Panel ─────────────────────────────────────────────
var panelHeaderStyle = lipgloss.NewStyle().
	Bold(true).
	Foreground(colorPrimary).
	BorderBottom(true).
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(colorMuted).
	PaddingBottom(1).
	MarginBottom(1)

var previewBorderStyle = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderForeground(colorMuted).
	Padding(1, 2)

// ── Divider ───────────────────────────────────────────────────
var dividerStyle = lipgloss.NewStyle().
	Foreground(colorMuted)

// ── Status Bar ────────────────────────────────────────────────
var statusBarStyle = lipgloss.NewStyle().
	Background(colorBgAlt).
	Foreground(colorPrimary).
	Padding(0, 1).
	Bold(true)

var statusBarMutedStyle = lipgloss.NewStyle().
	Background(colorBgAlt).
	Foreground(colorMuted).
	Padding(0, 1)

// ── Feedback ──────────────────────────────────────────────────
var successStyle = lipgloss.NewStyle().
	Foreground(colorSuccess).
	Bold(true)

var errorStyle = lipgloss.NewStyle().
	Foreground(colorDanger).
	Bold(true)

// ── Form ──────────────────────────────────────────────────────
var formLabelStyle = lipgloss.NewStyle().
	Foreground(colorAccent).
	Bold(true).
	Width(10)

var formInputStyle = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderForeground(colorMuted).
	Padding(0, 1)

var formActiveInputStyle = lipgloss.NewStyle().
	Border(lipgloss.RoundedBorder()).
	BorderForeground(colorPrimary).
	Padding(0, 1)
