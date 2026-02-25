package synccommand

import (
	"fmt"

	statuscommand "github.com/redjax/cheatsheets/internal/commands/syncCommand/statusCommand"
	"github.com/redjax/cheatsheets/internal/config"
	reposervices "github.com/redjax/cheatsheets/internal/services/repoServices"
	syncservice "github.com/redjax/cheatsheets/internal/services/syncService"
	"github.com/spf13/cobra"
)

var (
	force      bool
	sheetsPath string
)

// SyncCmd represents the sync command
var SyncCmd = &cobra.Command{
	Use:     "sync",
	Short:   "Sync cheatsheets to local directory. Aliases: sync, synch, synchronize.",
	Long:    `Creates a symlink (or copy on Windows) from the repository's cheatsheets directory to a local path (default: ~/.cheatsheets).`,
	Aliases: []string{"synch", "synchronize"},
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

		// Use flag value if provided, otherwise use config value
		destination := cfg.SheetsPath
		if sheetsPath != "" {
			destination = sheetsPath
		}

		fmt.Println("=== Syncing Cheatsheets ===")
		fmt.Printf("Source: %s/cheatsheets\n", cfg.Git.ClonePath)
		fmt.Printf("Destination: %s\n\n", destination)

		// Create the sync
		err = syncservice.CreateSync(cfg.Git.ClonePath, destination, force)
		if err != nil {
			return fmt.Errorf("failed to create sync: %w", err)
		}

		fmt.Println()
		fmt.Println("Sync completed successfully")
		return nil
	},
}

func init() {
	SyncCmd.Flags().BoolVarP(&force, "force", "f", false, "Force recreate sync even if it exists")
	SyncCmd.Flags().StringVarP(&sheetsPath, "sheets-path", "p", "", "Override destination path (default from config)")

	// Register subcommands
	SyncCmd.AddCommand(statuscommand.StatusCmd)
}
