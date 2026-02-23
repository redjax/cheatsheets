package listcommand

import (
	"fmt"
	"strings"

	"github.com/redjax/cheatsheets/internal/config"
	cheatsheetservice "github.com/redjax/cheatsheets/internal/services/cheatsheetService"
	reposervices "github.com/redjax/cheatsheets/internal/services/repoServices"
	"github.com/spf13/cobra"
)

var (
	typeFilter string
)

// ListCmd represents the list command
var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List available cheatsheets",
	Long:  `Lists all cheatsheets in the repository, optionally filtered by type (app, command, language, system, etc.).`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get the config file path from the persistent flag
		configFile, err := cmd.Flags().GetString("config-file")
		if err != nil {
			return fmt.Errorf("error getting config-file flag: %w", err)
		}

		// Load config
		var cfg *config.Config
		if configFile == "" {
			// Use FindConfigFile for .local fallback
			configFile = config.FindConfigFile("config.yml")
			cfg, err = config.LoadConfig(nil, configFile)
		} else {
			// Use explicit file
			cfg, err = config.LoadConfig(nil, configFile)
		}

		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		// Check if repository is cloned
		cloned, err := reposervices.IsRepositoryCloned(cfg.Git.ClonePath)
		if err != nil {
			return fmt.Errorf("error checking repository: %w", err)
		}

		if !cloned {
			return fmt.Errorf("repository is not cloned at %s. Run 'chtsht repo clone' first", cfg.Git.ClonePath)
		}

		// Validate cheatsheets directory exists
		if err := cheatsheetservice.ValidateCheatsheetsDirectory(cfg.Git.ClonePath); err != nil {
			return err
		}

		// Get available types
		availableTypes, err := cheatsheetservice.GetAvailableTypes(cfg.Git.ClonePath)
		if err != nil {
			return fmt.Errorf("error getting available types: %w", err)
		}

		if len(availableTypes) == 0 {
			fmt.Println("No cheatsheet types found in repository")
			return nil
		}

		// If type filter provided, validate it
		if typeFilter != "" {
			valid := false
			for _, t := range availableTypes {
				if t == typeFilter {
					valid = true
					break
				}
			}
			if !valid {
				fmt.Printf("⚠ Warning: Type '%s' not found in repository\n", typeFilter)
				fmt.Printf("Available types: %s\n\n", strings.Join(availableTypes, ", "))
				return nil
			}
		}

		// List cheatsheets
		result, err := cheatsheetservice.ListCheatsheets(cfg.Git.ClonePath, typeFilter)
		if err != nil {
			return fmt.Errorf("error listing cheatsheets: %w", err)
		}

		// Display results
		fmt.Println("=== Cheatsheets ===")
		if typeFilter != "" {
			fmt.Printf("Type: %s\n\n", typeFilter)
		} else {
			fmt.Printf("Available types: %s\n\n", strings.Join(availableTypes, ", "))
		}

		if result.TotalCount == 0 {
			fmt.Println("No cheatsheets found")
			return nil
		}

		// Sort types for consistent display
		var types []string
		for t := range result.TypesWithSheets {
			types = append(types, t)
		}
		// Use a simple sort to maintain alphabetical order
		sortedTypes := make([]string, len(types))
		copy(sortedTypes, types)
		for i := 0; i < len(sortedTypes); i++ {
			for j := i + 1; j < len(sortedTypes); j++ {
				if sortedTypes[i] > sortedTypes[j] {
					sortedTypes[i], sortedTypes[j] = sortedTypes[j], sortedTypes[i]
				}
			}
		}

		for _, t := range sortedTypes {
			sheets := result.TypesWithSheets[t]
			fmt.Printf("[%s]\n", t)
			for _, sheet := range sheets {
				fmt.Printf("  - %s\n", sheet)
			}
			fmt.Println()
		}

		fmt.Printf("Total: %d cheatsheet(s)\n", result.TotalCount)
		return nil
	},
}

func init() {
	ListCmd.Flags().StringVarP(&typeFilter, "type", "t", "", "Filter by type (app, command, language, system, etc.)")
}
