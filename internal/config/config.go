package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

type Config struct {
	Theme       Theme       `json:"theme"`
	Keybindings Keybindings `json:"keybindings"`
	General     General     `json:"general"`
}

type Theme struct {
	Name  string      `json:"name"`
	Dark  ThemeColors `json:"dark"`
	Light ThemeColors `json:"light"`
}

type ThemeColors struct {
	Primary    string `json:"primary"`
	Secondary  string `json:"secondary"`
	Accent     string `json:"accent"`
	Background string `json:"background"`
	Surface    string `json:"surface"`
	Text       string `json:"text"`
	TextMuted  string `json:"textMuted"`
	Success    string `json:"success"`
	Warning    string `json:"warning"`
	Error      string `json:"error"`
	Border     string `json:"border"`
	Selected   string `json:"selected"`
	Highlight  string `json:"highlight"`
}

type Keybindings struct {
	Navigation   KeyMap `json:"navigation"`
	Actions      KeyMap `json:"actions"`
	QuickActions KeyMap `json:"quickActions"`
}

type KeyMap map[string]string

type General struct {
	AutoCopy        bool   `json:"autoCopy"`
	ConfirmDelete   bool   `json:"confirmDelete"`
	ShowLineNumbers bool   `json:"showLineNumbers"`
	PreviewLines    int    `json:"previewLines"`
	DefaultStack    string `json:"defaultStack"`
	ExportFormat    string `json:"exportFormat"`
	Editor          string `json:"editor"`
	Shell           string `json:"shell"`
}

var DefaultConfig = Config{
	Theme: Theme{
		Name: "dark",
		Dark: ThemeColors{
			Primary:    "#7C3AED",
			Secondary:  "#8B5CF6",
			Accent:     "#06B6D4",
			Background: "#0F172A",
			Surface:    "#1E293B",
			Text:       "#E2E8F0",
			TextMuted:  "#64748B",
			Success:    "#10B981",
			Warning:    "#F59E0B",
			Error:      "#EF4444",
			Border:     "#334155",
			Selected:   "#334155",
			Highlight:  "#A78BFA",
		},
		Light: ThemeColors{
			Primary:    "#7C3AED",
			Secondary:  "#8B5CF6",
			Accent:     "#0891B2",
			Background: "#FFFFFF",
			Surface:    "#F1F5F9",
			Text:       "#1E293B",
			TextMuted:  "#64748B",
			Success:    "#10B981",
			Warning:    "#F59E0B",
			Error:      "#EF4444",
			Border:     "#E2E8F0",
			Selected:   "#E2E8F0",
			Highlight:  "#7C3AED",
		},
	},
	Keybindings: Keybindings{
		Navigation: KeyMap{
			"up":          "k",
			"down":        "j",
			"search":      "/",
			"select":      "enter",
			"toggleMulti": "space",
		},
		Actions: KeyMap{
			"add":    "a",
			"edit":   "e",
			"delete": "d",
			"quit":   "q",
		},
		QuickActions: KeyMap{
			"refresh":      "r",
			"stats":        "s",
			"preview":      "v",
			"recent":       "R",
			"stacks":       "t",
			"palette":      ":",
			"help":         "?",
			"batchProcess": "x",
		},
	},
	General: General{
		AutoCopy:        true,
		ConfirmDelete:   true,
		ShowLineNumbers: false,
		PreviewLines:    15,
		DefaultStack:    "",
		ExportFormat:    "skill.md",
		Editor:          os.Getenv("EDITOR"),
		Shell:           os.Getenv("SHELL"),
	},
}

type PresetThemes struct {
	Dark      ThemeColors
	Light     ThemeColors
	Cyberpunk ThemeColors
	Monokai   ThemeColors
	Dracula   ThemeColors
	Nord      ThemeColors
	GitHub    ThemeColors
}

var PresetThemesList = map[string]ThemeColors{
	"dark": {
		Primary:    "#7C3AED",
		Secondary:  "#8B5CF6",
		Accent:     "#06B6D4",
		Background: "#0F172A",
		Surface:    "#1E293B",
		Text:       "#E2E8F0",
		TextMuted:  "#64748B",
		Success:    "#10B981",
		Warning:    "#F59E0B",
		Error:      "#EF4444",
		Border:     "#334155",
		Selected:   "#334155",
		Highlight:  "#A78BFA",
	},
	"cyberpunk": {
		Primary:    "#FF00FF",
		Secondary:  "#00FFFF",
		Accent:     "#FFFF00",
		Background: "#0D0221",
		Surface:    "#1A0A2E",
		Text:       "#FFFFFF",
		TextMuted:  "#FF00FF",
		Success:    "#00FF00",
		Warning:    "#FFA500",
		Error:      "#FF0000",
		Border:     "#FF00FF",
		Selected:   "#2D1B4E",
		Highlight:  "#00FFFF",
	},
	"monokai": {
		Primary:    "#F92672",
		Secondary:  "#AE81FF",
		Accent:     "#66D9EF",
		Background: "#272822",
		Surface:    "#3E3D32",
		Text:       "#F8F8F2",
		TextMuted:  "#75715E",
		Success:    "#A6E22E",
		Warning:    "#E6DB74",
		Error:      "#F92672",
		Border:     "#49483E",
		Selected:   "#49483E",
		Highlight:  "#A6E22E",
	},
	"dracula": {
		Primary:    "#BD93F9",
		Secondary:  "#FF79C6",
		Accent:     "#8BE9FD",
		Background: "#282A36",
		Surface:    "#44475A",
		Text:       "#F8F8F2",
		TextMuted:  "#6272A4",
		Success:    "#50FA7B",
		Warning:    "#FFB86C",
		Error:      "#FF5555",
		Border:     "#44475A",
		Selected:   "#44475A",
		Highlight:  "#FF79C6",
	},
	"nord": {
		Primary:    "#5E81AC",
		Secondary:  "#81A1C1",
		Accent:     "#88C0D0",
		Background: "#2E3440",
		Surface:    "#3B4252",
		Text:       "#ECEFF4",
		TextMuted:  "#4C566A",
		Success:    "#A3BE8C",
		Warning:    "#EBCB8B",
		Error:      "#BF616A",
		Border:     "#4C566A",
		Selected:   "#434C5E",
		Highlight:  "#88C0D0",
	},
	"github": {
		Primary:    "#0366D6",
		Secondary:  "#28A745",
		Accent:     "#FFA657",
		Background: "#0D1117",
		Surface:    "#161B22",
		Text:       "#C9D1D9",
		TextMuted:  "#8B949E",
		Success:    "#3FB950",
		Warning:    "#D29922",
		Error:      "#F85149",
		Border:     "#30363D",
		Selected:   "#21262D",
		Highlight:  "#58A6FF",
	},
}

func GetConfigPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		home = "/tmp"
	}
	return filepath.Join(home, ".promptvault", "config.json")
}

func Load() (*Config, error) {
	path := GetConfigPath()

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			cfg := DefaultConfig
			if err := Save(&cfg); err != nil {
				return &cfg, nil
			}
			return &cfg, nil
		}
		return nil, err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func Save(cfg *Config) error {
	path := GetConfigPath()

	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

func Reset() error {
	return Save(&DefaultConfig)
}

func GetAvailableThemes() []string {
	themes := []string{}
	for name := range PresetThemesList {
		themes = append(themes, name)
	}
	return themes
}

func GetTheme(name string) *ThemeColors {
	colors, exists := PresetThemesList[name]
	if !exists {
		colors = PresetThemesList["dark"]
	}
	return &colors
}

func ListPresets() {
	fmt.Println("Available themes:")
	fmt.Println("-----------------")
	for name, colors := range PresetThemesList {
		fmt.Printf("  %-12s  #%s\n", name, colors.Primary)
	}
}

func ExportConfig(path string) error {
	cfg, err := Load()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

func ImportConfig(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return err
	}

	return Save(&cfg)
}
