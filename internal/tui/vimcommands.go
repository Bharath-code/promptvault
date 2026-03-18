package tui

import "strings"

type VimAction struct {
	Name    string
	Aliases []string
	Handler func(*App, []string) bool
}

var vimActions = []VimAction{
	{"quit", []string{"q"}, vimQuit},
	{"wq", []string{"wq", "x"}, vimWQ},
	{"write", []string{"w"}, vimWrite},
	{"q!", []string{"q!"}, vimForceQuit},
	{"wq!", []string{"wq!", "x!"}, vimForceWQ},
	{"e", []string{"e"}, vimRefresh},
	{"help", []string{"h", "help"}, vimHelp},
	{"set", []string{"set"}, vimSet},
}

func vimQuit(a *App, args []string) bool {
	return true
}

func vimWQ(a *App, args []string) bool {
	// Save and quit - in TUI context, save is implicit
	return true
}

func vimWrite(a *App, args []string) bool {
	a.showSuccess("Changes saved")
	return false
}

func vimForceQuit(a *App, args []string) bool {
	return true
}

func vimForceWQ(a *App, args []string) bool {
	return true
}

func vimRefresh(a *App, args []string) bool {
	a.loading = true
	return false
}

func vimHelp(a *App, args []string) bool {
	if a.state == stateHelpMenu {
		a.state = stateList
	} else {
		a.state = stateHelpMenu
	}
	return false
}

func vimSet(a *App, args []string) bool {
	if len(args) < 2 {
		a.showInfo("Usage: set [option] [value]")
		return false
	}

	option := args[0]
	value := strings.Join(args[1:], " ")

	switch option {
	case "nu", "number":
		// Toggle line numbers
		return false
	case "hlsearch", "hls":
		return false
	case "wrap":
		return false
	case "vim":
		if value == "true" || value == "1" {
			a.vimMode.SetMode(VimNormal)
		} else {
			a.vimMode.Enabled = false
		}
		return false
	default:
		a.showWarning("Unknown option: " + option)
		return false
	}
}

func parseVimCommand(input string) (string, []string) {
	input = strings.TrimSpace(input)
	if strings.HasPrefix(input, ":") {
		input = input[1:]
	}

	parts := strings.Fields(input)
	if len(parts) == 0 {
		return "", nil
	}

	cmd := strings.ToLower(parts[0])
	var args []string
	if len(parts) > 1 {
		args = parts[1:]
	}

	return cmd, args
}

func executeVimCommand(a *App, input string) bool {
	cmd, args := parseVimCommand(input)

	if cmd == "" {
		return false
	}

	// Handle repeat last command
	if cmd == "." {
		if a.vimMode.LastCommand != "" {
			cmd, args = parseVimCommand(a.vimMode.LastCommand)
		} else {
			return false
		}
	}

	for _, vc := range vimActions {
		if vc.Name == cmd || contains(vc.Aliases, cmd) {
			a.vimMode.LastCommand = input
			return vc.Handler(a, args)
		}
	}

	a.showWarning("Unknown command: " + cmd)
	return false
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
