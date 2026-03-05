package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
	"github.com/Bharath-code/promptvault/internal/db"
	"github.com/Bharath-code/promptvault/internal/model"
)

// viewState tracks which panel is active
type viewState int

const (
	stateList viewState = iota
	stateSearch
	stateDetail
	stateAdd
	stateEdit
	stateDeleteConfirm
	stateCopied
	stateFillVars
)

// App is the root Bubble Tea model
type App struct {
	db       *db.DB
	width    int
	height   int
	state    viewState
	prevState viewState

	// Data
	prompts  []*model.Prompt
	filtered []*model.Prompt
	cursor   int

	// Sub-components
	search        textinput.Model
	preview       viewport.Model
	cachedPreview string

	// Renderer cache
	glamourRenderer *glamour.TermRenderer
	lastWrapWidth   int

	// Add/Edit form
	form     *Form
	varForm  *VarForm

	// Feedback
	statusMsg   string
	statusIsErr bool
	statusTimer time.Time

	// Stack filter
	stackFilter string

	// Flash message timer
	flashMsg  string
	flashTime time.Time
}

type tickMsg time.Time
type promptsLoadedMsg []*model.Prompt

// New creates a new App instance
func New(database *db.DB) *App {
	search := textinput.New()
	search.Placeholder = "Search prompts... (title, content, tags, stack)"
	search.CharLimit = 256

	preview := viewport.New(0, 0)

	return &App{
		db:      database,
		search:  search,
		preview: preview,
	}
}

// Init implements tea.Model
func (a *App) Init() tea.Cmd {
	return tea.Batch(
		a.loadPrompts(),
		tea.EnterAltScreen,
	)
}

// loadPrompts fetches prompts from the db
func (a *App) loadPrompts() tea.Cmd {
	return func() tea.Msg {
		prompts, err := a.db.List(a.stackFilter)
		if err != nil {
			return promptsLoadedMsg(nil)
		}
		return promptsLoadedMsg(prompts)
	}
}

// Update implements tea.Model
func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height
		a.preview.Width = a.previewWidth()
		a.preview.Height = a.contentHeight()
		// Calculate safe width for search box avoiding overflow: 20 is min, cap it to 40 max or percentage
		sw := a.width/3
		if sw < 20 { sw = 20 }
		a.search.Width = sw
		a.updatePreview()

	case promptsLoadedMsg:
		a.prompts = msg
		a.applyFilter()
		a.updatePreview()

	case tickMsg:
		// Clear status after 2 seconds
		if !a.statusTimer.IsZero() && time.Since(a.statusTimer) > 2*time.Second {
			a.statusMsg = ""
			a.statusTimer = time.Time{}
		}
		if !a.flashTime.IsZero() && time.Since(a.flashTime) > 1500*time.Millisecond {
			a.flashMsg = ""
			a.flashTime = time.Time{}
		}
		if !a.statusTimer.IsZero() || !a.flashTime.IsZero() {
			cmds = append(cmds, tick())
		}
		return a, tea.Batch(cmds...)

	case tea.KeyMsg:
		if msg.Type == tea.KeyCtrlC {
			return a, tea.Quit
		}
		var cmd tea.Cmd
		_, cmd = a.handleKey(msg)
		if cmd != nil {
			cmds = append(cmds, cmd)
		}
	}

	// Update sub-components
	if a.state == stateSearch {
		// Only update search with non-key messages here, key messages are handled in handleSearchKey
		if _, ok := msg.(tea.KeyMsg); !ok {
			var cmd tea.Cmd
			a.search, cmd = a.search.Update(msg)
			cmds = append(cmds, cmd)
		}
	}

	if a.state == stateDetail {
		var cmd tea.Cmd
		a.preview, cmd = a.preview.Update(msg)
		cmds = append(cmds, cmd)
	}

	return a, tea.Batch(cmds...)
}

func (a *App) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch a.state {

	case stateList, stateDetail:
		return a.handleListKey(msg)

	case stateSearch:
		return a.handleSearchKey(msg)

	case stateAdd, stateEdit:
		return a.handleFormKey(msg)

	case stateFillVars:
		return a.handleVarFormKey(msg)

	case stateDeleteConfirm:
		return a.handleDeleteKey(msg)
	}

	return a, nil
}

func (a *App) handleListKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {

	case "q", "ctrl+c":
		return a, tea.Quit

	case "up", "k":
		if a.state == stateDetail {
			break // Let viewport handle it
		}
		if a.cursor > 0 {
			a.cursor--
			a.updatePreview()
		}

	case "down", "j":
		if a.state == stateDetail {
			break // Let viewport handle it
		}
		if a.cursor < len(a.filtered)-1 {
			a.cursor++
			a.updatePreview()
		}

	case "/":
		a.state = stateSearch
		a.search.Focus()
		return a, textinput.Blink

	case "enter", " ":
		// Check for variables and fill if present
		if p := a.selectedPrompt(); p != nil {
			vars := ExtractVars(p.Content)
			if len(vars) > 0 {
				a.state = stateFillVars
				a.varForm = NewVarForm(p.Content, vars)
				return a, a.varForm.Init()
			}

			// Copy prompt directly to clipboard
			if err := clipboard.WriteAll(p.Content); err == nil {
				_ = a.db.IncrementUsage(p.ID)
				a.flashMsg = "✓ Copied to clipboard!"
				a.flashTime = time.Now()
				return a, tick()
			} else {
				a.setStatus("Failed to copy: "+err.Error(), true)
				return a, tick()
			}
		}

	case "a":
		a.state = stateAdd
		a.form = NewForm(nil)
		return a, a.form.Init()

	case "e":
		if p := a.selectedPrompt(); p != nil {
			a.state = stateEdit
			a.form = NewForm(p)
			return a, a.form.Init()
		}

	case "d", "delete":
		if a.selectedPrompt() != nil {
			a.state = stateDeleteConfirm
		}

	case "v":
		if a.state == stateList {
			a.state = stateDetail
		} else {
			a.state = stateList
		}
		a.updatePreview()

	case "r":
		return a, a.loadPrompts()

	case "esc":
		a.state = stateList
		a.stackFilter = ""
		return a, a.loadPrompts()
	}

	return a, nil
}

func (a *App) handleSearchKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "esc":
		a.state = stateList
		a.search.SetValue("")
		a.search.Blur()
		a.applyFilter()
		return a, nil

	case "enter":
		a.state = stateList
		a.search.Blur()
		return a, nil

	case "up", "down":
		// These shouldn't act as input but rather navigation
		return a.handleListKey(msg)
	}

	var cmd tea.Cmd
	a.search, cmd = a.search.Update(msg)
	a.applyFilter()
	a.cursor = 0
	a.updatePreview()
	return a, cmd
}

func (a *App) handleFormKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if msg.String() == "esc" {
		a.state = stateList
		a.form = nil
		return a, nil
	}

	result, cmd := a.form.Update(msg)
	if result.Submitted {
		p := result.Prompt
		var err error
		if a.state == stateAdd {
			err = a.db.Add(p)
			if err == nil {
				a.setStatus("✓ Prompt added!", false)
			}
		} else {
			err = a.db.Update(p)
			if err == nil {
				a.setStatus("✓ Prompt updated!", false)
			}
		}
		if err != nil {
			a.setStatus("Error: "+err.Error(), true)
			return a, tea.Batch(cmd, a.loadPrompts(), tick())
		}
		a.state = stateList
		a.form = nil
		return a, tea.Batch(cmd, a.loadPrompts(), tick())
	}

	a.form = result.Form
	return a, cmd
}

func (a *App) handleVarFormKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	if msg.String() == "esc" {
		a.state = stateList
		a.varForm = nil
		return a, nil
	}

	res, cmd := a.varForm.Update(msg)
	if res.Submitted {
		if err := clipboard.WriteAll(res.Content); err == nil {
			if p := a.selectedPrompt(); p != nil {
				_ = a.db.IncrementUsage(p.ID)
			}
			a.flashMsg = "✓ Filled & Copied to clipboard!"
			a.flashTime = time.Now()
			a.state = stateList
			a.varForm = nil
			return a, tea.Batch(cmd, tick())
		} else {
			a.setStatus("Failed to copy: "+err.Error(), true)
			a.state = stateList
			a.varForm = nil
			return a, tea.Batch(cmd, tick())
		}
	}
	a.varForm = res.Form
	return a, cmd
}

func (a *App) handleDeleteKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch msg.String() {
	case "y", "Y":
		if p := a.selectedPrompt(); p != nil {
			if err := a.db.Delete(p.ID); err == nil {
				a.setStatus("Prompt deleted", false)
				if a.cursor > 0 {
					a.cursor--
				}
				a.state = stateList
				return a, tea.Batch(a.loadPrompts(), tick())
			}
		}
		a.state = stateList
		return a, a.loadPrompts()

	case "n", "N", "esc":
		a.state = stateList
	}
	return a, nil
}

// View implements tea.Model
func (a *App) View() string {
	if a.width == 0 {
		return "Loading..."
	}

	switch a.state {
	case stateAdd, stateEdit:
		if a.form != nil {
			return a.renderForm()
		}
	case stateDeleteConfirm:
		return a.renderDeleteConfirm()
	case stateFillVars:
		if a.varForm != nil {
			return a.varForm.View(a.width, a.height)
		}
	}

	return a.renderMain()
}

func (a *App) renderMain() string {
	header := a.renderHeader()
	body := a.renderBody()
	statusBar := a.renderStatusBar()

	return lipgloss.JoinVertical(lipgloss.Left,
		header,
		body,
		statusBar,
	)
}

func (a *App) renderHeader() string {
	title := titleStyle.Render("⚡ PromptVault")

	total := fmt.Sprintf("%d prompts", len(a.prompts))
	if a.stackFilter != "" {
		total += " in " + stackStyle.Render(a.stackFilter)
	}
	count := subtitleStyle.Render(total)

	searchBox := ""
	if a.state == stateSearch {
		searchBox = searchStyle.Render(a.search.View())
	} else {
		searchBox = helpStyle.Render("/ to search")
	}

	left := lipgloss.JoinHorizontal(lipgloss.Center, title, "  ", count)
	
	// Ensure searchBox doesn't push layout offscreen. Give priority to searchBox.
	actualLeftWidth := lipgloss.Width(left)
	actualSearchWidth := lipgloss.Width(searchBox)
	
	// If the terminal is incredibly narrow, hide the left side completely
	if actualLeftWidth + actualSearchWidth + 4 > a.width {
		left = title
		actualLeftWidth = lipgloss.Width(left)
		if actualLeftWidth + actualSearchWidth + 4 > a.width {
			left = "" // hide title to make room for search box
			actualLeftWidth = 0
		}
	}
	
	gapW := a.width - actualLeftWidth - actualSearchWidth - 4
	if gapW < 0 {
		gapW = 0
	}
	gap := lipgloss.NewStyle().Width(gapW).Render("")

	header := lipgloss.NewStyle().
		Width(a.width).
		BorderBottom(true).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(colorMuted).
		Padding(0, 1).
		Render(lipgloss.JoinHorizontal(lipgloss.Center,
			left,
			gap,
			searchBox,
		))

	return header
}

func (a *App) renderBody() string {
	listWidth := a.listWidth()
	previewWidth := a.previewWidth()
	height := a.contentHeight()

	list := a.renderList(listWidth, height)

	var preview string
	if a.state == stateDetail {
		preview = a.renderPreview(previewWidth, height)
	} else {
		preview = a.renderPreviewPane(previewWidth, height)
	}

	divider := dividerStyle.Render(strings.Repeat("│\n", height))

	return lipgloss.JoinHorizontal(lipgloss.Top, list, divider, preview)
}

func (a *App) renderList(width, height int) string {
	var items []string

	if len(a.filtered) == 0 {
		empty := lipgloss.NewStyle().
			Width(width).
			Height(height).
			Align(lipgloss.Center, lipgloss.Center).
			Foreground(colorMuted).
			Render("No prompts found.\nPress 'a' to add one.\n\nOr run: promptvault init")
		return empty
	}

	maxVisible := height - 2
	start := 0
	if a.cursor >= maxVisible {
		start = a.cursor - maxVisible + 1
	}
	end := start + maxVisible
	if end > len(a.filtered) {
		end = len(a.filtered)
	}

	for i := start; i < end; i++ {
		p := a.filtered[i]
		item := a.renderListItem(p, i == a.cursor, width)
		items = append(items, item)
	}

	// Flash message overlay
	if a.flashMsg != "" {
		items = append([]string{successStyle.PaddingLeft(2).Render(a.flashMsg)}, items...)
	}

	content := strings.Join(items, "\n")
	return lipgloss.NewStyle().
		Width(width).
		Height(height).
		Render(content)
}

func (a *App) renderListItem(p *model.Prompt, selected bool, width int) string {
	verified := ""
	if p.Verified {
		verified = verifiedStyle.Render(" ✓")
	}

	stack := ""
	if p.Stack != "" {
		parts := strings.Split(p.Stack, "/")
		leaf := parts[len(parts)-1]
		stack = " " + tagStyle.Render(leaf)
	}

	usage := ""
	if p.UsageCount > 0 {
		usage = usageStyle.Render(fmt.Sprintf(" (%dx)", p.UsageCount))
	}

	title := p.Title + verified + usage
	meta := stack

	line := lipgloss.JoinHorizontal(lipgloss.Center,
		title,
		lipgloss.NewStyle().Width(width-lipgloss.Width(title)-lipgloss.Width(meta)-4).Render(""),
		meta,
	)

	if selected {
		return selectedItemStyle.Width(width - 2).Render(line)
	}
	return itemStyle.Width(width).Render(line)
}

func (a *App) renderPreviewPane(width, height int) string {
	p := a.selectedPrompt()
	if p == nil {
		return lipgloss.NewStyle().
			Width(width).
			Height(height).
			Align(lipgloss.Center, lipgloss.Center).
			Foreground(colorMuted).
			Render("Select a prompt to preview")
	}

	header := panelHeaderStyle.Width(width - 4).Render(p.Title)

	// Stack + models
	meta := ""
	if p.Stack != "" {
		meta += stackStyle.Render(p.Stack) + "  "
	}
	for _, m := range p.Models {
		meta += tagStyle.Render(m) + " "
	}

	contentStyle := lipgloss.NewStyle().
		Foreground(colorText).
		Width(width - 4).
		Margin(1, 0)

	tags := ""
	for _, t := range p.Tags {
		tags += tagStyle.Render("#"+t) + " "
	}

	footer := helpStyle.Render("ENTER copy  •  v expand  •  e edit  •  d delete")

	body := lipgloss.JoinVertical(lipgloss.Left,
		header,
		meta,
		contentStyle.Render(a.cachedPreview),
		tags,
		footer,
	)

	return previewBorderStyle.Width(width - 2).Height(height - 2).Render(body)
}

func (a *App) renderPreview(width, height int) string {
	p := a.selectedPrompt()
	if p == nil {
		return ""
	}

	header := panelHeaderStyle.Width(width-4).Render("▶ " + p.Title)

	body := lipgloss.JoinVertical(lipgloss.Left,
		header,
		a.preview.View(),
	)

	return previewBorderStyle.Width(width - 2).Height(height - 2).Render(body)
}

func (a *App) renderStatusBar() string {
	left := ""
	if a.statusMsg != "" {
		if a.statusIsErr {
			left = errorStyle.Render("✗ " + a.statusMsg)
		} else {
			left = successStyle.Render(a.statusMsg)
		}
	} else {
		left = statusBarStyle.Render("PromptVault")
	}

	keys := statusBarMutedStyle.Render("a add  •  e edit  •  d del  •  / search  •  v preview  •  q quit")

	gap := lipgloss.NewStyle().
		Background(colorBgAlt).
		Width(a.width - lipgloss.Width(left) - lipgloss.Width(keys)).
		Render("")

	return lipgloss.JoinHorizontal(lipgloss.Bottom, left, gap, keys)
}

func (a *App) renderForm() string {
	title := "Add Prompt"
	if a.state == stateEdit {
		title = "Edit Prompt"
	}

	header := lipgloss.NewStyle().
		Bold(true).
		Foreground(colorPrimary).
		PaddingBottom(1).
		Render("⚡ PromptVault — " + title)

	formView := a.form.View(a.width, a.height-6)

	return lipgloss.NewStyle().
		Padding(1, 2).
		Render(lipgloss.JoinVertical(lipgloss.Left,
			header,
			formView,
			helpStyle.Render("TAB next field  •  ESC cancel  •  CTRL+S save"),
		))
}

func (a *App) renderDeleteConfirm() string {
	p := a.selectedPrompt()
	if p == nil {
		return ""
	}

	msg := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorDanger).
		Padding(2, 4).
		Render(lipgloss.JoinVertical(lipgloss.Center,
			errorStyle.Render("Delete Prompt?"),
			"",
			lipgloss.NewStyle().Foreground(colorText).Render(p.Title),
			"",
			helpStyle.Render("y confirm  •  n / ESC cancel"),
		))

	return lipgloss.NewStyle().
		Width(a.width).
		Height(a.height).
		Align(lipgloss.Center, lipgloss.Center).
		Render(msg)
}

// --- helpers ---

func (a *App) selectedPrompt() *model.Prompt {
	if len(a.filtered) == 0 || a.cursor >= len(a.filtered) {
		return nil
	}
	return a.filtered[a.cursor]
}

func (a *App) applyFilter() {
	q := strings.TrimSpace(a.search.Value())
	if q == "" {
		a.filtered = a.prompts
		return
	}

	var results []*model.Prompt
	prompts, err := a.db.Search(q)
	if err == nil {
		results = prompts
	}
	a.filtered = results
}

func (a *App) updatePreview() {
	p := a.selectedPrompt()
	if p == nil {
		a.cachedPreview = ""
		a.preview.SetContent("")
		return
	}
	
	w := a.preview.Width - 4
	if w < 20 {
		w = 80
	}
	
	if a.glamourRenderer == nil || a.lastWrapWidth != w {
		renderer, err := glamour.NewTermRenderer(
			glamour.WithAutoStyle(),
			glamour.WithWordWrap(w),
		)
		if err == nil {
			a.glamourRenderer = renderer
			a.lastWrapWidth = w
		}
	}
	
	// Fast path for preview pane items during list scrolling
	lines := strings.Split(p.Content, "\n")
	paneText := p.Content
	if len(lines) > 20 {
		paneText = strings.Join(lines[:20], "\n") + "\n\n..."
	}

	if a.glamourRenderer != nil {
		if str, err := a.glamourRenderer.Render(paneText); err == nil {
			paneText = str
		}
	}
	a.cachedPreview = paneText

	// Only full-render when we specifically expand to detail view
	fullRenderedContent := p.Content
	if a.state == stateDetail {
		// Glamour's regex engine is extremely CPU intensive. Avoid rendering if >2500 bytes (blocks for seconds)
		if a.glamourRenderer != nil && len(p.Content) < 2500 {
			if str, err := a.glamourRenderer.Render(p.Content); err == nil {
				fullRenderedContent = str
			}
		}
	}

	// Create the full text for viewport
	meta := ""
	if p.Stack != "" {
		meta += stackStyle.Render(p.Stack) + "  "
	}
	if p.Verified {
		meta += verifiedStyle.Render("✓ Verified") + "  "
	}
	for _, m := range p.Models {
		meta += tagStyle.Render(m) + " "
	}

	fullContent := lipgloss.JoinVertical(lipgloss.Left,
		meta,
		"",
		fullRenderedContent,
		"",
		usageStyle.Render(fmt.Sprintf("Used %d times", p.UsageCount)),
	)
	
	a.preview.SetContent(fullContent)
	// Optionally clear scroll state to top when changing prompt
	if a.preview.YOffset > 0 {
		a.preview.GotoTop()
	}
}

func (a *App) setStatus(msg string, isErr bool) {
	a.statusMsg = msg
	a.statusIsErr = isErr
	a.statusTimer = time.Now()
}

func (a *App) listWidth() int {
	return a.width / 2
}

func (a *App) previewWidth() int {
	return a.width - a.listWidth() - 1
}

func (a *App) contentHeight() int {
	return a.height - 4 // header + status bar
}

func tick() tea.Cmd {
	return tea.Tick(300*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}

// Run starts the TUI
func Run(database *db.DB) error {
	app := New(database)
	p := tea.NewProgram(app, tea.WithAltScreen(), tea.WithMouseCellMotion())
	_, err := p.Run()
	return err
}
