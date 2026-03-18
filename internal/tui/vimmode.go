package tui

import "github.com/charmbracelet/lipgloss"

type VimModeState string

const (
	VimNormal  VimModeState = "NORMAL"
	VimInsert  VimModeState = "INSERT"
	VimVisual  VimModeState = "VISUAL"
	VimCommand VimModeState = "COMMAND"
)

type VimModeHandler struct {
	Mode                VimModeState
	Enabled             bool
	CommandBuffer       string
	CommandHistory      []string
	CommandHistoryIndex int
	LastCommand         string
}

func NewVimModeHandler() *VimModeHandler {
	return &VimModeHandler{
		Mode:                VimNormal,
		Enabled:             true,
		CommandHistory:      []string{},
		CommandHistoryIndex: -1,
	}
}

func (v *VimModeHandler) Toggle() {
	v.Enabled = !v.Enabled
	if v.Enabled {
		v.Mode = VimNormal
	} else {
		v.Mode = VimNormal
	}
}

func (v *VimModeHandler) SetMode(mode VimModeState) {
	v.Mode = mode
	v.CommandBuffer = ""
	v.CommandHistoryIndex = -1
}

func (v *VimModeHandler) EnterInsert() {
	v.Mode = VimInsert
}

func (v *VimModeHandler) EnterNormal() {
	v.Mode = VimNormal
}

func (v *VimModeHandler) EnterVisual() {
	v.Mode = VimVisual
}

func (v *VimModeHandler) EnterCommand() {
	v.Mode = VimCommand
	v.CommandBuffer = ""
}

func (v *VimModeHandler) AddCommandHistory(cmd string) {
	if cmd == "" {
		return
	}
	// Don't duplicate consecutive commands
	if len(v.CommandHistory) > 0 && v.CommandHistory[len(v.CommandHistory)-1] == cmd {
		return
	}
	v.CommandHistory = append(v.CommandHistory, cmd)
	v.CommandHistoryIndex = -1
}

func (v *VimModeHandler) HistoryUp() string {
	if len(v.CommandHistory) == 0 {
		return ""
	}
	if v.CommandHistoryIndex < len(v.CommandHistory)-1 {
		v.CommandHistoryIndex++
	}
	return v.CommandHistory[len(v.CommandHistory)-1-v.CommandHistoryIndex]
}

func (v *VimModeHandler) HistoryDown() string {
	if v.CommandHistoryIndex <= 0 {
		v.CommandHistoryIndex = -1
		return ""
	}
	v.CommandHistoryIndex--
	return v.CommandHistory[len(v.CommandHistory)-1-v.CommandHistoryIndex]
}

func (v *VimModeHandler) RenderModeIndicator() string {
	style := lipgloss.NewStyle()

	switch v.Mode {
	case VimNormal:
		style = style.Foreground(lipgloss.Color("82")).Bold(true)
	case VimInsert:
		style = style.Foreground(lipgloss.Color("208")).Bold(true)
	case VimVisual:
		style = style.Foreground(lipgloss.Color("226")).Bold(true)
	case VimCommand:
		style = style.Foreground(lipgloss.Color("147")).Bold(true)
	}

	return style.Render(string(v.Mode))
}
