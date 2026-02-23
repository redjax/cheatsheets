package showcommand

import (
	"fmt"

	"github.com/redjax/cheatsheets/internal/config"
	cheatsheetservice "github.com/redjax/cheatsheets/internal/services/cheatsheetService"
	"github.com/spf13/cobra"
)

var (
	typeFilter string
)

// ShowCmd represents the show command
var ShowCmd = &cobra.Command{
	Use:     "show [cheatsheet-name]",
	Aliases: []string{"view", "read"},
	Short:   "Display a cheatsheet. Aliases: view, read",
	Long: `Display a cheatsheet with beautiful markdown rendering.

Examples:
  # Show a specific cheatsheet by type and name
  chtsht show -t language python
  chtsht show --type language python

  # Show interactive selector for all cheatsheets
  chtsht show

  # Show selector filtered by type
  chtsht show -t language`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get the config file path from the persistent flag
		configFile, err := cmd.Flags().GetString("config-file")
		if err != nil {
			return fmt.Errorf("error getting config-file flag: %w", err)
		}

		// Load config
		var cfg *config.Config
		if configFile == "" {
			configFile = config.FindConfigFile("config.yml")
			cfg, err = config.LoadConfig(nil, configFile)
		} else {
			cfg, err = config.LoadConfig(nil, configFile)
		}

		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		// Use SheetsPath from config (defaults to ~/.cheatsheets)
		sheetsPath := cfg.SheetsPath

		// Validate cheatsheets directory exists
		if err := cheatsheetservice.ValidateCheatsheetsDirectory(sheetsPath); err != nil {
			return fmt.Errorf("cheatsheets not found at %s. Run 'chtsht sync' first", sheetsPath)
		}

		// Determine cheatsheet name from args
		var cheatsheetName string
		if len(args) > 0 {
			cheatsheetName = args[0]
		}

		// Handle different scenarios
		if typeFilter != "" && cheatsheetName != "" {
			// Scenario 1: Type and name provided - show specific cheatsheet
			return cheatsheetservice.ShowCheatsheet(sheetsPath, typeFilter, cheatsheetName)
		} else if typeFilter != "" && cheatsheetName == "" {
			// Scenario 2: Only type provided - show selector for that type
			return cheatsheetservice.ShowCheatsheetSelector(sheetsPath, typeFilter)
		} else if typeFilter == "" && cheatsheetName != "" {
			// Scenario 3: Only name provided - search all types for that name
			return cheatsheetservice.ShowCheatsheetByName(sheetsPath, cheatsheetName)
		} else {
			// Scenario 4: No type or name - show selector for all cheatsheets
			return cheatsheetservice.ShowCheatsheetSelector(sheetsPath, "")
		}
	},
}

func init() {
	ShowCmd.Flags().StringVarP(&typeFilter, "type", "t", "", "Filter by cheatsheet type (app, command, language, system, etc.)")
}
