package cmd

import (
	"fmt"
	"strings"

	"github.com/Bharath-code/promptvault/internal/config"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage PromptVault configuration",
	Long:  `View and modify PromptVault configuration including themes and keybindings.`,
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}

		fmt.Println()
		fmt.Println(colorPrimary + "⚡ PromptVault Configuration" + colorReset)
		fmt.Println(strings.Repeat("─", 50))
		fmt.Println()

		fmt.Printf("%s Theme: %s\n", iconInfo, cfg.Theme.Name)
		fmt.Printf("%s Auto Copy: %t\n", iconInfo, cfg.General.AutoCopy)
		fmt.Printf("%s Confirm Delete: %t\n", iconInfo, cfg.General.ConfirmDelete)
		fmt.Printf("%s Preview Lines: %d\n", iconInfo, cfg.General.PreviewLines)
		fmt.Printf("%s Export Format: %s\n", iconInfo, cfg.General.ExportFormat)
		fmt.Printf("%s Config Path: %s\n", iconInfo, config.GetConfigPath())
		fmt.Println()

		return nil
	},
}

var configThemeCmd = &cobra.Command{
	Use:   "theme [name]",
	Short: "Get or set the color theme",
	Args:  cobra.RangeArgs(0, 1),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}

		if len(args) == 0 {
			fmt.Printf("Current theme: %s\n", cfg.Theme.Name)
			fmt.Println("\nAvailable themes:")
			config.ListPresets()
			return nil
		}

		themeName := strings.ToLower(args[0])
		colors, exists := config.PresetThemesList[themeName]
		if !exists {
			printError("Unknown theme: %s", themeName)
			fmt.Println("\nAvailable themes:")
			config.ListPresets()
			return fmt.Errorf("unknown theme: %s", themeName)
		}

		cfg.Theme.Name = themeName
		cfg.Theme.Dark = colors

		if err := config.Save(cfg); err != nil {
			return err
		}

		printSuccess("Theme set to: %s", themeName)
		return nil
	},
}

var configKeybindingsCmd = &cobra.Command{
	Use:   "keybindings",
	Short: "Show keybindings",
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}

		fmt.Println()
		fmt.Println(colorPrimary + "⚡ Keybindings" + colorReset)
		fmt.Println(strings.Repeat("─", 50))
		fmt.Println()

		fmt.Println(colorMuted + "Navigation:" + colorReset)
		for action, key := range cfg.Keybindings.Navigation {
			fmt.Printf("  %-15s %s\n", action+":", key)
		}
		fmt.Println()

		fmt.Println(colorMuted + "Actions:" + colorReset)
		for action, key := range cfg.Keybindings.Actions {
			fmt.Printf("  %-15s %s\n", action+":", key)
		}
		fmt.Println()

		fmt.Println(colorMuted + "Quick Actions:" + colorReset)
		for action, key := range cfg.Keybindings.QuickActions {
			fmt.Printf("  %-15s %s\n", action+":", key)
		}
		fmt.Println()

		return nil
	},
}

var configSetCmd = &cobra.Command{
	Use:   "set [key] [value]",
	Short: "Set a configuration value",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		cfg, err := config.Load()
		if err != nil {
			return err
		}

		key := strings.ToLower(args[0])
		value := args[1]

		set := false
		switch key {
		case "autocopy":
			cfg.General.AutoCopy = value == "true" || value == "1"
			set = true
		case "confirmdelete":
			cfg.General.ConfirmDelete = value == "true" || value == "1"
			set = true
		case "previewlines":
			fmt.Sscanf(value, "%d", &cfg.General.PreviewLines)
			set = true
		case "exportformat":
			cfg.General.ExportFormat = value
			set = true
		case "defaultstack":
			cfg.General.DefaultStack = value
			set = true
		case "editor":
			cfg.General.Editor = value
			set = true
		}

		if !set {
			printError("Unknown config key: %s", key)
			fmt.Println("\nAvailable keys:")
			fmt.Println("  autocopy        - Auto copy to clipboard")
			fmt.Println("  confirmdelete   - Confirm before delete")
			fmt.Println("  previewlines    - Lines in preview")
			fmt.Println("  exportformat    - Default export format")
			fmt.Println("  defaultstack    - Default stack filter")
			fmt.Println("  editor          - External editor")
			return fmt.Errorf("unknown key: %s", key)
		}

		if err := config.Save(cfg); err != nil {
			return err
		}

		printSuccess("Config updated: %s = %s", key, value)
		return nil
	},
}

var configResetCmd = &cobra.Command{
	Use:   "reset",
	Short: "Reset configuration to defaults",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := config.Reset(); err != nil {
			return err
		}
		printSuccess("Configuration reset to defaults")
		return nil
	},
}

var configExportCmd = &cobra.Command{
	Use:   "export [path]",
	Short: "Export configuration to a file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := config.ExportConfig(args[0]); err != nil {
			return err
		}
		printSuccess("Configuration exported to: %s", args[0])
		return nil
	},
}

var configImportCmd = &cobra.Command{
	Use:   "import [path]",
	Short: "Import configuration from a file",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := config.ImportConfig(args[0]); err != nil {
			return err
		}
		printSuccess("Configuration imported from: %s", args[0])
		return nil
	},
}

func init() {
	configCmd.AddCommand(configShowCmd)
	configCmd.AddCommand(configThemeCmd)
	configCmd.AddCommand(configKeybindingsCmd)
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configResetCmd)
	configCmd.AddCommand(configExportCmd)
	configCmd.AddCommand(configImportCmd)

	rootCmd.AddCommand(configCmd)
}
