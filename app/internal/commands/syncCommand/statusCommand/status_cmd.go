package statuscommand

import (
	"fmt"

	"github.com/redjax/cheatsheets/internal/config"
	reposervices "github.com/redjax/cheatsheets/internal/services/repoServices"
	syncservice "github.com/redjax/cheatsheets/internal/services/syncService"
	"github.com/spf13/cobra"
)

// StatusCmd represents the status command
var StatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Check sync status",
	Long:  `Checks the current state of the cheatsheets sync (symlink or copy).`,
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
			configFile = config.FindConfigFile("")
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

		fmt.Println("=== Sync Status ===")
		fmt.Printf("Source: %s\n", cfg.Git.ClonePath+"/cheatsheets")
		fmt.Printf("Destination: %s\n\n", cfg.SheetsPath)

		// Check sync status
		status := syncservice.CheckSyncStatus(cfg.Git.ClonePath, cfg.SheetsPath)

		if status.Error != nil {
			fmt.Printf("✗ Error: %v\n", status.Error)
			return nil
		}

		if !status.Exists {
			fmt.Println("✗ Sync does not exist")
			fmt.Println("  Run 'chtsht sync' to create it")
			return nil
		}

		fmt.Printf("Sync exists\n")
		fmt.Printf("  Type: %s\n", status.Type)

		if status.IsValid {
			fmt.Println("  Status: Valid")
		} else {
			fmt.Println("  Status: Invalid (points to wrong location or corrupted)")
			fmt.Println("  Run 'chtsht sync --force' to fix")
		}

		return nil
	},
}
