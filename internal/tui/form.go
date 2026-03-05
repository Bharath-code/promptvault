package tui

import (
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/Bharath-code/promptvault/internal/model"
)

// formField tracks which field is active
type formField int

const (
	fieldTitle formField = iota
	fieldContent
	fieldStack
	fieldTags
	fieldModels
	fieldCount // sentinel
)

// FormResult is returned when the form is submitted or cancelled
type FormResult struct {
	Submitted bool
	Prompt    *model.Prompt
	Form      *Form
}

// Form is a multi-field prompt editor
type Form struct {
	title   textinput.Model
	content textarea.Model
	stack   textinput.Model
	tags    textinput.Model
	models  textinput.Model

	active   formField
	editing  *model.Prompt // nil for add, existing for edit
	verified bool
}

// NewForm creates a form, optionally pre-filled with a prompt
func NewForm(p *model.Prompt) *Form {
	ti := textinput.New()
	ti.Placeholder = "Prompt title"
	ti.CharLimit = 200
	ti.Width = 60

	ta := textarea.New()
	ta.Placeholder = "Write your prompt here...\n\nUse {{variable}} for template variables."
	ta.CharLimit = 10000
	ta.SetWidth(60)
	ta.SetHeight(8)

	si := textinput.New()
	si.Placeholder = "e.g. frontend/react/hooks"
	si.CharLimit = 100
	si.Width = 60

	tg := textinput.New()
	tg.Placeholder = "Comma-separated tags (e.g. debugging, hooks)"
	tg.CharLimit = 200
	tg.Width = 60

	mi := textinput.New()
	mi.Placeholder = "e.g. claude-sonnet, gpt-4o"
	mi.CharLimit = 200
	mi.Width = 60

	f := &Form{
		title:   ti,
		content: ta,
		stack:   si,
		tags:    tg,
		models:  mi,
		active:  fieldTitle,
	}

	// Pre-fill if editing
	if p != nil {
		f.editing = p
		f.title.SetValue(p.Title)
		f.content.SetValue(p.Content)
		f.stack.SetValue(p.Stack)
		f.tags.SetValue(strings.Join(p.Tags, ", "))
		f.models.SetValue(strings.Join(p.Models, ", "))
		f.verified = p.Verified
	}

	f.focusCurrent()
	return f
}

// Init returns the initial command
func (f *Form) Init() tea.Cmd {
	return textinput.Blink
}

// Update handles input and returns a FormResult
func (f *Form) Update(msg tea.KeyMsg) (FormResult, tea.Cmd) {
	switch msg.String() {
	case "ctrl+s":
		// Submit
		p := f.toPrompt()
		if p.Title == "" || p.Content == "" {
			// Don't submit empty
			return FormResult{Form: f}, nil
		}
		return FormResult{Submitted: true, Prompt: p, Form: f}, nil

	case "tab", "shift+tab":
		f.blurAll()
		if msg.String() == "tab" {
			f.active = (f.active + 1) % fieldCount
		} else {
			f.active = (f.active - 1 + fieldCount) % fieldCount
		}
		f.focusCurrent()
		return FormResult{Form: f}, nil

	case "ctrl+v":
		// Toggle verified
		f.verified = !f.verified
		return FormResult{Form: f}, nil
	}

	// Update active field
	var cmd tea.Cmd
	switch f.active {
	case fieldTitle:
		f.title, cmd = f.title.Update(msg)
	case fieldContent:
		f.content, cmd = f.content.Update(msg)
	case fieldStack:
		f.stack, cmd = f.stack.Update(msg)
	case fieldTags:
		f.tags, cmd = f.tags.Update(msg)
	case fieldModels:
		f.models, cmd = f.models.Update(msg)
	}

	return FormResult{Form: f}, cmd
}

// View renders the form
func (f *Form) View(width, height int) string {
	fields := []struct {
		label string
		view  string
		idx   formField
	}{
		{"Title", f.title.View(), fieldTitle},
		{"Content", f.content.View(), fieldContent},
		{"Stack", f.stack.View(), fieldStack},
		{"Tags", f.tags.View(), fieldTags},
		{"Models", f.models.View(), fieldModels},
	}

	var rows []string
	for _, field := range fields {
		label := formLabelStyle.Render(field.label)
		inputStyle := formInputStyle
		if field.idx == f.active {
			inputStyle = formActiveInputStyle
		}
		input := inputStyle.Width(width - 16).Render(field.view)
		rows = append(rows, lipgloss.JoinHorizontal(lipgloss.Top, label, input))
	}

	// Verified toggle
	verifiedLabel := formLabelStyle.Render("Verified")
	verifiedVal := "[ ]"
	if f.verified {
		verifiedVal = verifiedStyle.Render("[✓]")
	}
	rows = append(rows, lipgloss.JoinHorizontal(lipgloss.Center,
		verifiedLabel,
		lipgloss.NewStyle().Padding(0, 1).Render(verifiedVal+" (CTRL+V to toggle)"),
	))

	return lipgloss.JoinVertical(lipgloss.Left, rows...)
}

func (f *Form) focusCurrent() {
	switch f.active {
	case fieldTitle:
		f.title.Focus()
	case fieldContent:
		f.content.Focus()
	case fieldStack:
		f.stack.Focus()
	case fieldTags:
		f.tags.Focus()
	case fieldModels:
		f.models.Focus()
	}
}

func (f *Form) blurAll() {
	f.title.Blur()
	f.content.Blur()
	f.stack.Blur()
	f.tags.Blur()
	f.models.Blur()
}

func (f *Form) toPrompt() *model.Prompt {
	p := &model.Prompt{
		Title:    strings.TrimSpace(f.title.Value()),
		Content:  strings.TrimSpace(f.content.Value()),
		Stack:    strings.TrimSpace(f.stack.Value()),
		Verified: f.verified,
	}

	// Parse tags
	for _, t := range strings.Split(f.tags.Value(), ",") {
		if t = strings.TrimSpace(t); t != "" {
			p.Tags = append(p.Tags, t)
		}
	}

	// Parse models
	for _, m := range strings.Split(f.models.Value(), ",") {
		if m = strings.TrimSpace(m); m != "" {
			p.Models = append(p.Models, m)
		}
	}

	// Preserve ID if editing
	if f.editing != nil {
		p.ID = f.editing.ID
		p.CreatedAt = f.editing.CreatedAt
		p.UsageCount = f.editing.UsageCount
		p.LastUsedAt = f.editing.LastUsedAt
	}

	return p
}
