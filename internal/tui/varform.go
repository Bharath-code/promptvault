package tui

import (
	"regexp"
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var templateVarRe = regexp.MustCompile(`\{\{([^}]+)\}\}`)

func ExtractVars(content string) []string {
	matches := templateVarRe.FindAllStringSubmatch(content, -1)
	vars := []string{}
	seen := map[string]bool{}
	for _, m := range matches {
		name := strings.TrimSpace(m[1])
		if !seen[name] {
			vars = append(vars, name)
			seen[name] = true
		}
	}
	return vars
}

type VarFormResult struct {
	Submitted bool
	Content   string
	Form      *VarForm
}

type VarForm struct {
	vars    []string
	inputs  []textinput.Model
	active  int
	content string
}

func NewVarForm(content string, vars []string) *VarForm {
	inputs := make([]textinput.Model, len(vars))
	for i, v := range vars {
		ti := textinput.New()
		ti.Placeholder = "Value for " + v
		ti.Prompt = formLabelStyle.Render(v + ": ")
		ti.Width = 60
		if i == 0 {
			ti.Focus()
		}
		inputs[i] = ti
	}
	return &VarForm{
		vars:    vars,
		inputs:  inputs,
		content: content,
	}
}

func (f *VarForm) Init() tea.Cmd {
	return textinput.Blink
}

func (f *VarForm) Update(msg tea.KeyMsg) (VarFormResult, tea.Cmd) {
	switch msg.String() {
	case "ctrl+s", "enter":
		if f.active == len(f.inputs)-1 || msg.String() == "ctrl+s" {
			// Submit
			final := f.content
			for i, v := range f.vars {
				val := f.inputs[i].Value()
				if val == "" {
					val = f.inputs[i].Placeholder
				}
				final = strings.ReplaceAll(final, "{{"+v+"}}", val)
				final = strings.ReplaceAll(final, "{{ "+v+" }}", val)
			}
			return VarFormResult{Submitted: true, Content: final, Form: f}, nil
		}
		// next field
		f.inputs[f.active].Blur()
		f.active++
		f.inputs[f.active].Focus()
		return VarFormResult{Form: f}, nil

	case "tab", "down":
		f.inputs[f.active].Blur()
		f.active = (f.active + 1) % len(f.inputs)
		f.inputs[f.active].Focus()
		return VarFormResult{Form: f}, nil

	case "shift+tab", "up":
		f.inputs[f.active].Blur()
		f.active = (f.active - 1 + len(f.inputs)) % len(f.inputs)
		f.inputs[f.active].Focus()
		return VarFormResult{Form: f}, nil
	}

	var cmd tea.Cmd
	f.inputs[f.active], cmd = f.inputs[f.active].Update(msg)
	return VarFormResult{Form: f}, cmd
}

func (f *VarForm) View(width, height int) string {
	var rows []string
	rows = append(rows, lipgloss.NewStyle().Foreground(colorPrimary).Bold(true).Render("⚡ Fill Variables Before Copying"))
	rows = append(rows, "")

	for i := range f.inputs {
		style := formInputStyle
		if i == f.active {
			style = formActiveInputStyle
		}
		rows = append(rows, style.Width(width-20).Render(f.inputs[i].View()))
		rows = append(rows, "")
	}

	rows = append(rows, helpStyle.Render("TAB/ENTER next  •  ESC cancel  •  CTRL+S finish"))

	return lipgloss.NewStyle().Padding(2, 4).Render(lipgloss.JoinVertical(lipgloss.Left, rows...))
}
