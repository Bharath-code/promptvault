package tui

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Bharath-code/promptvault/internal/config"
	"github.com/Bharath-code/promptvault/internal/db"
	"github.com/Bharath-code/promptvault/internal/model"
	"github.com/atotto/clipboard"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
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
	stateHelpMenu
	stateStats
	stateCommandPalette
	stateOnboarding
	stateStackTree
	stateConfig
	stateThemePreview
)

// Command represents a available command in palette
type Command struct {
	Name        string
	Description string
	Shortcut    string
	Action      func(*App) (tea.Model, tea.Cmd)
}

// CommandPalette holds commands for fuzzy searching
type CommandPalette struct {
	search   textinput.Model
	commands []Command
	filtered []Command
	cursor   int
}

// App is the root Bubble Tea model
type App struct {
	db        *db.DB
	width     int
	height    int
	state     viewState
	prevState viewState

	// Data
	prompts     []*model.Prompt
	filtered    []*model.Prompt
	cursor      int
	scores      []int           // Fuzzy search scores
	showRecent  bool            // Toggle recent prompts section
	selected    map[int]bool    // Multi-select indices
	recentCache []*model.Prompt // Cached recent prompts
	recentDirty bool            // Cache invalidation flag

	// Sub-components
	search        textinput.Model
	preview       viewport.Model
	cachedPreview string

	// Renderer cache
	glamourRenderer *glamour.TermRenderer
	lastWrapWidth   int

	// Add/Edit form
	form    *Form
	varForm *VarForm

	// Feedback
	statusMsg   string
	statusIsErr bool
	statusTimer time.Time

	// Stack filter
	stackFilter string

	// Flash message timer
	flashMsg  string
	flashTime time.Time

	// Loading state
	spinner spinner.Model
	loading bool

	// Command palette
	commandPalette *CommandPalette

	// Toast notifications
	toastManager *ToastManager

	// Onboarding tour
	onboarding     *OnboardingTour
	showOnboarding bool

	// Stack tree navigation
	stackTree *StackTree

	// Config view
	themePreview *ThemePreview
	configTab    int
}

type tickMsg time.Time
type promptsLoadedMsg []*model.Prompt
type startLoadingMsg struct{}
type stopLoadingMsg struct{}

// New creates a new App instance
func New(database *db.DB) *App {
	search := textinput.New()
	search.Placeholder = "Search prompts... (title, content, tags, stack)"
	search.CharLimit = 256

	preview := viewport.New(0, 0)

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return &App{
		db:           database,
		search:       search,
		preview:      preview,
		spinner:      s,
		selected:     make(map[int]bool),
		recentCache:  nil,
		recentDirty:  true,
		toastManager: &ToastManager{maxCount: 5},
	}
}

// Init checks if first run and starts onboarding if needed
func (a *App) Init() tea.Cmd {
	ctx := context.Background()
	count, _ := a.db.Count(ctx)
	a.showOnboarding = count == 0

	if a.showOnboarding {
		a.onboarding = NewOnboardingTour()
		a.state = stateOnboarding
	}

	return tea.Batch(
		tea.EnterAltScreen,
		a.loadPrompts(),
	)
}

// NewCommandPalette creates a new command palette with available commands
func NewCommandPalette() *CommandPalette {
	search := textinput.New()
	search.Placeholder = "Type a command..."
	search.CharLimit = 50

	commands := []Command{
		{"Add Prompt", "Create a new prompt", "a", nil},
		{"Edit Prompt", "Edit the selected prompt", "e", nil},
		{"Delete Prompt", "Delete the selected prompt", "d", nil},
		{"Search", "Search prompts by text", "/", nil},
		{"Toggle Preview", "Toggle full-screen preview", "v", nil},
		{"Refresh", "Reload prompts from database", "r", nil},
		{"Toggle Recent", "Show/hide recently used prompts", "R", nil},
		{"Statistics", "View usage statistics", "s", nil},
		{"Help", "Show keyboard shortcuts", "?", nil},
		{"Quit", "Exit PromptVault", "q", nil},
	}

	return &CommandPalette{
		search:   search,
		commands: commands,
		filtered: commands,
	}
}

// filterCommands filters commands by search query
func (cp *CommandPalette) filterCommands(query string) {
	if query == "" {
		cp.filtered = cp.commands
		return
	}

	query = strings.ToLower(query)
	var filtered []Command
	for _, cmd := range cp.commands {
		if strings.Contains(strings.ToLower(cmd.Name), query) ||
			strings.Contains(strings.ToLower(cmd.Description), query) {
			filtered = append(filtered, cmd)
		}
	}
	cp.filtered = filtered
}

// loadPrompts fetches prompts from the db
func (a *App) loadPrompts() tea.Cmd {
	return func() tea.Msg {
		ctx := context.Background()
		prompts, err := a.db.List(ctx, a.stackFilter)
		if err != nil {
			return promptsLoadedMsg(nil)
		}
		return promptsLoadedMsg(prompts)
	}
}

// startLoading shows the loading spinner
func (a *App) startLoading() tea.Cmd {
	a.loading = true
	return nil
}

// stopLoading hides the loading spinner
func (a *App) stopLoading() tea.Cmd {
	a.loading = false
	return nil
}

// Update implements tea.Model
func (a *App) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	// Update spinner when loading
	if a.loading {
		var cmd tea.Cmd
		a.spinner, cmd = a.spinner.Update(msg)
		cmds = append(cmds, cmd)
	}

	switch msg := msg.(type) {

	case tea.WindowSizeMsg:
		a.width = msg.Width
		a.height = msg.Height
		a.preview.Width = a.previewWidth()
		a.preview.Height = a.contentHeight()
		// Calculate safe width for search box avoiding overflow: 20 is min, cap it to 40 max or percentage
		sw := a.width / 3
		if sw < 20 {
			sw = 20
		}
		a.search.Width = sw
		a.updatePreview()

	case promptsLoadedMsg:
		a.prompts = msg
		a.applyFilter()
		// CRITICAL: Do NOT update preview on initial load!
		// This causes expensive markdown rendering on startup
		// Preview will be updated when user navigates with arrow keys
		// Mark recent cache as dirty (needs recalculation)
		a.recentDirty = true
		// Stop loading when prompts are loaded
		if a.loading {
			a.loading = false
		}

	case startLoadingMsg:
		a.loading = true
		return a, nil

	case stopLoadingMsg:
		a.loading = false
		return a, nil

	case tickMsg:
		if !a.statusTimer.IsZero() && time.Since(a.statusTimer) > 2*time.Second {
			a.statusMsg = ""
			a.statusTimer = time.Time{}
		}
		if !a.flashTime.IsZero() && time.Since(a.flashTime) > 1500*time.Millisecond {
			a.flashMsg = ""
			a.flashTime = time.Time{}
		}
		a.toastManager.RemoveExpired()
		if !a.statusTimer.IsZero() || !a.flashTime.IsZero() || a.toastManager.IsActive() {
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
		// Handle viewport scrolling in full-screen mode
		var cmd tea.Cmd
		a.preview, cmd = a.preview.Update(msg)
		cmds = append(cmds, cmd)
	}

	return a, tea.Batch(cmds...)
}

func (a *App) handleKey(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Handle states that close on any key press
	if a.state == stateHelpMenu || a.state == stateStats || a.state == stateCommandPalette {
		switch msg.String() {
		case "q", "ctrl+c":
			return a, tea.Quit
		default:
			a.state = stateList
			return a, nil
		}
	}

	// Onboarding tour keyboard handling
	if a.state == stateOnboarding {
		switch msg.String() {
		case "q", "ctrl+c":
			return a, tea.Quit
		case "esc":
			a.state = stateList
			a.onboarding = nil
			return a, nil
		case "enter", " ":
			if a.onboarding.IsLast() {
				a.state = stateList
				a.onboarding = nil
				if a.showOnboarding {
					a.showSuccess("Press a to add your first prompt!")
				}
				return a, nil
			}
			a.onboarding.Next()
			return a, nil
		case "left", "h":
			a.onboarding.Previous()
			return a, nil
		case "right", "l":
			if !a.onboarding.IsLast() {
				a.onboarding.Next()
			}
			return a, nil
		}
	}

	// Stack tree keyboard handling
	if a.state == stateStackTree {
		switch msg.String() {
		case "esc":
			a.state = stateList
			a.stackTree = nil
			return a, nil
		case "q", "ctrl+c":
			return a, tea.Quit
		case "up", "k":
			if a.stackTree != nil {
				a.stackTree.MoveUp()
			}
			return a, nil
		case "down", "j":
			if a.stackTree != nil {
				a.stackTree.MoveDown()
			}
			return a, nil
		case "left", "h":
			if a.stackTree != nil {
				a.stackTree.Collapse()
			}
			return a, nil
		case "right", "l":
			if a.stackTree != nil {
				a.stackTree.Expand()
			}
			return a, nil
		case "enter":
			if a.stackTree != nil && a.stackTree.IsSelectable() {
				node := a.stackTree.Current()
				a.stackFilter = node.Path
				a.state = stateList
				a.stackTree = nil
				a.loading = true
				return a, a.loadPrompts()
			}
			return a, nil
		case " ":
			if a.stackTree != nil {
				a.stackTree.ToggleExpand()
			}
			return a, nil
		}
	}

	// Theme preview keyboard handling
	if a.state == stateThemePreview {
		switch msg.String() {
		case "esc":
			a.state = stateList
			a.themePreview = nil
			return a, nil
		case "q", "ctrl+c":
			return a, tea.Quit
		case "up", "k":
			if a.themePreview != nil {
				a.themePreview.MoveUp()
			}
			return a, nil
		case "down", "j":
			if a.themePreview != nil {
				a.themePreview.MoveDown()
			}
			return a, nil
		case "enter":
			if a.themePreview != nil {
				selected := a.themePreview.Select()
				cfg, _ := config.Load()
				cfg.Theme.Name = selected
				cfg.Theme.Dark = *config.GetTheme(selected)
				config.Save(cfg)
				a.showSuccess("Theme changed to: " + selected)
				a.state = stateList
				a.themePreview = nil
			}
			return a, nil
		}
	}

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

	case "enter":
		if p := a.selectedPrompt(); p != nil {
			vars := ExtractVars(p.Content)
			if len(vars) > 0 {
				a.state = stateFillVars
				a.varForm = NewVarForm(p.Content, vars)
				return a, a.varForm.Init()
			}

			if err := clipboard.WriteAll(p.Content); err == nil {
				ctx := context.Background()
				if incErr := a.db.IncrementUsage(ctx, p.ID); incErr != nil {
					a.showWarning("Copied (usage tracking failed)")
				} else {
					a.showSuccess("Copied to clipboard!")
				}
				return a, tick()
			} else {
				a.showError("Failed to copy: " + err.Error())
				return a, tick()
			}
		}

	case " ":
		// Space for multi-select toggle
		if p := a.selectedPrompt(); p != nil {
			// Toggle current item selection
			if a.selected[a.cursor] {
				delete(a.selected, a.cursor)
			} else {
				a.selected[a.cursor] = true
			}
			return a, nil
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

	case "x":
		// Batch operations
		if len(a.selected) > 0 {
			return a, a.performBatchOperation()
		}

	case "v":
		// Toggle full-screen preview overlay
		if a.state == stateList {
			a.state = stateDetail
			a.updatePreview()
		} else if a.state == stateDetail {
			a.state = stateList
		}
		return a, nil

	case "r":
		a.loading = true
		return a, a.loadPrompts()

	case "R":
		// Toggle recent prompts section
		a.showRecent = !a.showRecent
		return a, nil

	case "s":
		// Toggle stats dashboard
		if a.state == stateStats {
			a.state = stateList
		} else {
			a.state = stateStats
		}
		return a, nil

	case ":":
		a.state = stateCommandPalette
		a.commandPalette = NewCommandPalette()
		a.commandPalette.search.Focus()
		return a, textinput.Blink

	case "t":
		// Open stack tree navigation
		a.openStackTree()
		return a, nil

	case "?":
		if a.state == stateHelpMenu {
			a.state = stateList
		} else {
			a.state = stateHelpMenu
		}
		return a, nil

	case "g":
		a.openThemePreview()
		return a, nil

	case "esc":
		if a.state == stateStats || a.state == stateHelpMenu {
			a.state = stateList
			return a, nil
		}
		a.state = stateList
		a.stackFilter = ""
		a.loading = true
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
	a.loading = true
	a.applyFilter()
	a.cursor = 0
	a.updatePreview()
	a.loading = false
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
		ctx := context.Background()
		p := result.Prompt
		var err error
		if a.state == stateAdd {
			err = a.db.Add(ctx, p)
			if err == nil {
				a.showSuccess("Prompt added!")
			}
		} else {
			author := os.Getenv("USER")
			if author == "" {
				author = "anonymous"
			}
			err = a.db.Update(ctx, p, "Edited in TUI", author)
			if err == nil {
				a.showSuccess("Prompt updated!")
			}
		}
		if err != nil {
			a.showError("Error: " + err.Error())
			a.loading = true
			return a, tea.Batch(cmd, a.loadPrompts(), tick())
		}
		a.state = stateList
		a.form = nil
		a.loading = true
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
		ctx := context.Background()
		if err := clipboard.WriteAll(res.Content); err == nil {
			if p := a.selectedPrompt(); p != nil {
				if incErr := a.db.IncrementUsage(ctx, p.ID); incErr != nil {
					a.showWarning("Copied (usage tracking failed)")
				} else {
					a.showSuccess("Filled & Copied to clipboard!")
				}
			} else {
				a.flashMsg = "✓ Filled & Copied to clipboard!"
				a.flashTime = time.Now()
			}
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
		ctx := context.Background()
		if p := a.selectedPrompt(); p != nil {
			if err := a.db.Delete(ctx, p.ID); err == nil {
				a.showSuccess("Prompt deleted")
				if a.cursor > 0 {
					a.cursor--
				}
				a.state = stateList
				a.loading = true
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

func (a *App) performBatchOperation() tea.Cmd {
	ctx := context.Background()
	processed := 0

	for idx := range a.selected {
		if idx < len(a.filtered) {
			p := a.filtered[idx]
			_ = a.db.IncrementUsage(ctx, p.ID)
			processed++
		}
	}

	a.showSuccess(fmt.Sprintf("Processed %d prompts", processed))
	a.selected = make(map[int]bool)
	return tick()
}

// View implements tea.Model
func (a *App) View() string {
	if a.width == 0 {
		return "Loading..."
	}

	// Show loading overlay
	if a.loading {
		return a.renderLoading()
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
	case stateHelpMenu:
		return a.renderHelpMenu()
	case stateStats:
		return a.renderStats()
	case stateCommandPalette:
		return a.renderCommandPalette()
	case stateOnboarding:
		return a.renderOnboarding()
	case stateStackTree:
		return a.renderStackTree()
	case stateThemePreview:
		return a.renderThemePreview()
	}

	return a.renderMain()
}

func (a *App) renderMain() string {
	header := a.renderHeader()
	body := a.renderBody()
	statusBar := a.renderStatusBar()
	toastBar := a.toastManager.Render(a.width)

	return lipgloss.JoinVertical(lipgloss.Left,
		header,
		body,
		toastBar,
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
	if actualLeftWidth+actualSearchWidth+4 > a.width {
		left = title
		actualLeftWidth = lipgloss.Width(left)
		if actualLeftWidth+actualSearchWidth+4 > a.width {
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
	// Full-screen preview overlay in detail mode
	if a.state == stateDetail {
		return a.renderFullScreenPreview()
	}

	listWidth := a.listWidth()
	previewWidth := a.previewWidth()
	height := a.contentHeight()

	list := a.renderList(listWidth, height)

	var preview string
	preview = a.renderPreviewPane(previewWidth, height)

	divider := dividerStyle.Render(strings.Repeat("│\n", height))

	return lipgloss.JoinHorizontal(lipgloss.Top, list, divider, preview)
}

func (a *App) renderFullScreenPreview() string {
	p := a.selectedPrompt()
	if p == nil {
		return lipgloss.NewStyle().
			Width(a.width).
			Height(a.height).
			Align(lipgloss.Center, lipgloss.Center).
			Foreground(colorMuted).
			Render("No prompt selected\n\nPress ↑/↓ to select a prompt")
	}

	// Full-screen preview with header
	header := lipgloss.NewStyle().
		Bold(true).
		Foreground(colorPrimary).
		Width(a.width).
		Render(fmt.Sprintf("⚡ %s", p.Title))

	// Metadata
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

	// Full content with markdown rendering
	content := a.preview.View()

	// Footer with instructions
	footer := helpStyle.Render("v close  •  ↑/↓ scroll  •  Enter copy  •  t stacks")

	body := lipgloss.JoinVertical(lipgloss.Left,
		header,
		meta,
		"",
		content,
		"",
		footer,
	)

	return lipgloss.NewStyle().
		Width(a.width).
		Height(a.height).
		Render(body)
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

	// Show recent prompts section if enabled and not searching
	if a.showRecent && a.search.Value() == "" && len(a.prompts) > 0 {
		recentSection := a.renderRecentPrompts(width)
		items = append(items, recentSection)
		items = append(items, "") // Spacer
	}

	maxVisible := height - 2
	if a.showRecent && a.search.Value() == "" {
		maxVisible -= 6 // Reserve space for recent section
	}

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
		item := a.renderListItem(p, i == a.cursor, i, width)
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

func (a *App) renderRecentPrompts(width int) string {
	// Use cached recent prompts if available
	if !a.recentDirty && a.recentCache != nil && len(a.recentCache) > 0 {
		// Build recent section from cache
		var lines []string
		lines = append(lines, panelHeaderStyle.Render(" 🔥 Recently Used"))

		for _, p := range a.recentCache {
			title := p.Title
			if len(title) > 50 {
				title = title[:47] + "..."
			}
			lines = append(lines, fmt.Sprintf("  • %-45s %dx", title, p.UsageCount))
		}

		return lipgloss.NewStyle().Render(strings.Join(lines, "\n"))
	}

	// Calculate recent prompts (expensive - only do when cache is dirty)
	type recentPrompt struct {
		prompt *model.Prompt
		score  int
	}
	var recents []recentPrompt

	for _, p := range a.prompts {
		if p.UsageCount > 0 {
			recents = append(recents, recentPrompt{p, p.UsageCount})
		}
	}

	// Sort by usage count (descending)
	for i := 0; i < len(recents)-1; i++ {
		for j := i + 1; j < len(recents); j++ {
			if recents[j].score > recents[i].score {
				recents[i], recents[j] = recents[j], recents[i]
			}
		}
	}

	// Take top 5 and cache
	if len(recents) > 5 {
		recents = recents[:5]
	}

	// Update cache
	a.recentCache = make([]*model.Prompt, len(recents))
	for i, rp := range recents {
		a.recentCache[i] = rp.prompt
	}
	a.recentDirty = false

	if len(recents) == 0 {
		return ""
	}

	// Build recent section
	var lines []string
	lines = append(lines, panelHeaderStyle.Render(" 🔥 Recently Used"))

	for _, rp := range recents {
		p := rp.prompt
		title := p.Title
		if len(title) > 50 {
			title = title[:47] + "..."
		}
		lines = append(lines, fmt.Sprintf("  • %-45s %dx", title, p.UsageCount))
	}

	return lipgloss.NewStyle().Render(strings.Join(lines, "\n"))
}

func (a *App) renderListItem(p *model.Prompt, selected bool, index int, width int) string {
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

	// Show match score if searching
	score := ""
	if a.scores != nil && len(a.scores) > 0 {
		for i, fp := range a.filtered {
			if fp == p && i < len(a.scores) && a.scores[i] < 100 {
				score = scoreStyle.Render(fmt.Sprintf(" %d%%", a.scores[i]))
				break
			}
		}
	}

	// Show selection indicator (check this specific item's index)
	selectIndicator := "  "
	if a.selected[index] {
		selectIndicator = successStyle.Render("✓ ")
	}

	meta := stack

	line := lipgloss.JoinHorizontal(lipgloss.Center,
		selectIndicator,
		title,
		lipgloss.NewStyle().Width(width-lipgloss.Width(title)-lipgloss.Width(meta)-lipgloss.Width(score)-6).Render(""),
		meta,
		score,
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
		Width(width-4).
		Margin(1, 0)

	tags := ""
	for _, t := range p.Tags {
		tags += tagStyle.Render("#"+t) + " "
	}

	footer := helpStyle.Render("Enter copy  •  v expand  •  t stacks  •  e edit  •  d delete")

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

	header := panelHeaderStyle.Width(width - 4).Render("▶ " + p.Title)

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

	keys := statusBarMutedStyle.Render("a add  •  e edit  •  t stacks  •  / search  •  ? help  •  q quit")

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

func (a *App) renderHelpMenu() string {
	helpItems := []struct {
		key     string
		desc    string
		section string
	}{
		// Navigation
		{"↑/↓ or k/j", "Navigate prompts", "Navigation"},
		{"/", "Search prompts", ""},
		{"Space", "Select/deselect", ""},
		{"Enter", "Copy to clipboard", ""},

		// Actions
		{"a", "Add new prompt", "Actions"},
		{"e", "Edit selected", ""},
		{"d", "Delete selected", ""},
		{"v", "Toggle preview", ""},

		// Quick Actions
		{"c", "Copy selected", "Quick Actions"},
		{"r", "Refresh list", ""},
		{"R", "Toggle recent", ""},
		{"s", "Show stats", ""},
		{"x", "Batch process", ""},
		{"t", "Stack tree", ""},
		{"g", "Themes", ""},

		// Other
		{":", "Command palette", "Other"},
		{"Esc", "Go back / Clear search", ""},
		{"q", "Quit", ""},
		{"Ctrl+C", "Exit", ""},
	}

	var sections []string
	currentSection := ""
	var items []string

	for _, item := range helpItems {
		if item.section != "" && item.section != currentSection {
			if len(items) > 0 {
				sections = append(sections, strings.Join(items, "\n"))
				items = nil
			}
			currentSection = item.section
			items = append(items, panelHeaderStyle.Render(" "+currentSection))
		}
		items = append(items, fmt.Sprintf("  %-16s %s",
			tagStyle.Render(item.key),
			item.desc))
	}

	if len(items) > 0 {
		sections = append(sections, strings.Join(items, "\n"))
	}

	content := strings.Join(sections, "\n\n")

	msg := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorPrimary).
		Padding(2, 4).
		Width(60).
		Render(lipgloss.JoinVertical(lipgloss.Left,
			lipgloss.NewStyle().
				Bold(true).
				Foreground(colorPrimary).
				PaddingBottom(1).
				Render("⚡ Quick Actions & Keybindings"),
			content,
			"",
			helpStyle.Render("Press any key to close"),
		))

	return lipgloss.NewStyle().
		Width(a.width).
		Height(a.height).
		Align(lipgloss.Center, lipgloss.Center).
		Render(msg)
}

func (a *App) renderStats() string {
	// Get stats from database
	total := len(a.prompts)
	totalUsage := 0
	stackCounts := make(map[string]int)

	for _, p := range a.prompts {
		totalUsage += p.UsageCount
		if p.Stack != "" {
			stackCounts[p.Stack]++
		}
	}

	// Get top stacks
	type stackCount struct {
		stack string
		count int
	}
	var stacks []stackCount
	for s, c := range stackCounts {
		stacks = append(stacks, stackCount{s, c})
	}
	// Sort by count
	for i := 0; i < len(stacks)-1; i++ {
		for j := i + 1; j < len(stacks); j++ {
			if stacks[j].count > stacks[i].count {
				stacks[i], stacks[j] = stacks[j], stacks[i]
			}
		}
	}
	// Take top 5
	if len(stacks) > 5 {
		stacks = stacks[:5]
	}

	// Get most used prompts
	type promptUsage struct {
		title string
		count int
	}
	var usage []promptUsage
	for _, p := range a.prompts {
		if p.UsageCount > 0 {
			usage = append(usage, promptUsage{p.Title, p.UsageCount})
		}
	}
	// Sort by usage
	for i := 0; i < len(usage)-1; i++ {
		for j := i + 1; j < len(usage); j++ {
			if usage[j].count > usage[i].count {
				usage[i], usage[j] = usage[j], usage[i]
			}
		}
	}
	// Take top 5
	if len(usage) > 5 {
		usage = usage[:5]
	}

	// Build stats display
	var lines []string
	lines = append(lines, panelHeaderStyle.Render(" 📊 PromptVault Statistics"))
	lines = append(lines, "")
	lines = append(lines, fmt.Sprintf("  %-20s  %d", "Total Prompts:", total))
	lines = append(lines, fmt.Sprintf("  %-20s  %d", "Total Usage:", totalUsage))
	lines = append(lines, "")

	// Top stacks
	lines = append(lines, panelHeaderStyle.Render(" Top Stacks"))
	for i, s := range stacks {
		medal := "  "
		if i == 0 {
			medal = "🥇"
		} else if i == 1 {
			medal = "🥈"
		} else if i == 2 {
			medal = "🥉"
		}
		lines = append(lines, fmt.Sprintf("  %s %-25s %d", medal, s.stack, s.count))
	}
	lines = append(lines, "")

	// Most used
	lines = append(lines, panelHeaderStyle.Render(" Most Used Prompts"))
	for i, u := range usage {
		medal := "  "
		if i == 0 {
			medal = "🥇"
		} else if i == 1 {
			medal = "🥈"
		} else if i == 2 {
			medal = "🥉"
		}
		lines = append(lines, fmt.Sprintf("  %s %-30s %dx", medal, u.title, u.count))
	}

	content := strings.Join(lines, "\n")

	msg := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorPrimary).
		Padding(2, 4).
		Width(70).
		Render(content)

	return lipgloss.NewStyle().
		Width(a.width).
		Height(a.height).
		Align(lipgloss.Center, lipgloss.Center).
		Render(msg)
}

func (a *App) renderLoading() string {
	// Show loading overlay with spinner
	spinnerView := a.spinner.View() + " Loading prompts..."

	msg := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorPrimary).
		Padding(2, 4).
		Render(lipgloss.JoinVertical(lipgloss.Center,
			spinnerView,
			"",
			helpStyle.Render("Please wait..."),
		))

	return lipgloss.NewStyle().
		Width(a.width).
		Height(a.height).
		Align(lipgloss.Center, lipgloss.Center).
		Render(msg)
}

func (a *App) renderCommandPalette() string {
	if a.commandPalette == nil {
		return ""
	}

	paletteWidth := 60
	paletteHeight := 15

	var items []string
	items = append(items, panelHeaderStyle.Render(" Command Palette"))

	searchView := a.commandPalette.search.View()
	items = append(items, searchView)
	items = append(items, "")

	maxItems := paletteHeight - 6
	for i := 0; i < maxItems && i < len(a.commandPalette.filtered); i++ {
		cmd := a.commandPalette.filtered[i]
		shortcut := tagStyle.Render(cmd.Shortcut)
		if i == a.commandPalette.cursor {
			item := lipgloss.JoinHorizontal(lipgloss.Center,
				selectedItemStyle.Render("> "),
				lipgloss.NewStyle().Foreground(colorText).Render(cmd.Name),
				lipgloss.NewStyle().Width(paletteWidth-25).Render(""),
				shortcut,
			)
			items = append(items, item)
		} else {
			item := lipgloss.JoinHorizontal(lipgloss.Center,
				lipgloss.NewStyle().Width(2).Render(" "),
				lipgloss.NewStyle().Foreground(colorMuted).Render(cmd.Name),
				lipgloss.NewStyle().Width(paletteWidth-25).Render(""),
				lipgloss.NewStyle().Foreground(colorMuted).Render(cmd.Shortcut),
			)
			items = append(items, item)
		}
	}

	items = append(items, "")
	items = append(items, helpStyle.Render("↑/↓ navigate  •  Enter execute  •  Esc close"))

	content := strings.Join(items, "\n")

	msg := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(colorAccent).
		Padding(1, 2).
		Width(paletteWidth).
		Render(content)

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
		a.scores = nil
		return
	}

	// Use fuzzy search for better results
	filtered, scores := FuzzySearch(q, a.prompts)
	a.filtered = filtered
	a.scores = scores
}

func (a *App) updatePreview() {
	p := a.selectedPrompt()
	if p == nil {
		a.cachedPreview = ""
		a.preview.SetContent("")
		return
	}

	lines := strings.Split(p.Content, "\n")
	paneText := p.Content
	if len(lines) > 15 {
		paneText = strings.Join(lines[:15], "\n") + "\n..."
	}

	highlightedText := HighlightPromptContent(paneText, 15)
	a.cachedPreview = highlightedText

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
		highlightedText,
		"",
		usageStyle.Render(fmt.Sprintf("Used %d times", p.UsageCount)),
	)

	a.preview.SetContent(fullContent)
	// Reset scroll to top when changing prompt
	if a.preview.YOffset > 0 {
		a.preview.GotoTop()
	}
}

func (a *App) setStatus(msg string, isErr bool) {
	a.statusMsg = msg
	a.statusIsErr = isErr
	a.statusTimer = time.Now()
}

// Toast notification helpers
func (a *App) showToast(msg string, toastType ToastType) {
	a.toastManager.Add(msg, toastType, 3*time.Second)
}

func (a *App) showSuccess(msg string) {
	a.showToast(msg, ToastSuccess)
}

func (a *App) showError(msg string) {
	a.showToast(msg, ToastError)
}

func (a *App) showWarning(msg string) {
	a.showToast(msg, ToastWarning)
}

func (a *App) showInfo(msg string) {
	a.showToast(msg, ToastInfo)
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

func (a *App) renderOnboarding() string {
	if a.onboarding == nil {
		return a.renderMain()
	}
	return a.onboarding.Render(a.width, a.height)
}

func (a *App) renderStackTree() string {
	if a.stackTree == nil {
		return a.renderMain()
	}

	var stacks []string
	stackCounts := make(map[string]int)

	for _, p := range a.prompts {
		if p.Stack != "" {
			stacks = append(stacks, p.Stack)
			stackCounts[p.Stack] = 1
		}
	}

	a.stackTree = NewStackTree(stacks, 40)
	a.stackTree.UpdateCounts(stackCounts)

	content := a.stackTree.Render()

	return lipgloss.Place(a.width, a.height,
		lipgloss.Center, lipgloss.Center,
		lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7C3AED")).
			Padding(1, 2).
			Render(content))
}

func (a *App) openStackTree() {
	var stacks []string
	for _, p := range a.prompts {
		if p.Stack != "" {
			stacks = append(stacks, p.Stack)
		}
	}
	a.stackTree = NewStackTree(stacks, 40)
	a.state = stateStackTree
}

func (a *App) renderThemePreview() string {
	if a.themePreview == nil {
		cfg, _ := config.Load()
		a.themePreview = NewThemePreview(cfg.Theme.Name, a.width, a.height)
	}

	content := a.themePreview.Render()

	return lipgloss.Place(a.width, a.height,
		lipgloss.Center, lipgloss.Center,
		lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#7C3AED")).
			Padding(1, 2).
			Render(content))
}

func (a *App) openThemePreview() {
	cfg, _ := config.Load()
	a.themePreview = NewThemePreview(cfg.Theme.Name, a.width, a.height)
	a.state = stateThemePreview
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
