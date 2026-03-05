package cleanupcommand

import (
	"fmt"

	"github.com/redjax/cheatsheets/internal/config"
	reposervices "github.com/redjax/cheatsheets/internal/services/repoServices"
	"github.com/spf13/cobra"
)

var CleanupCmd = &cobra.Command{
	Use:   "cleanup",
	Short: "Clean up local branches that have been deleted on remote",
	Long: `Remove local branches that have been deleted on the remote repository.
Only deletes branches that were previously pushed (have upstream tracking).
Branches that were never pushed to remote are kept safe.`,
	Example: "  chtsht repo cleanup",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get the config file path from the persistent flag
		configFile, err := cmd.Flags().GetString("config-file")
		if err != nil {
			return fmt.Errorf("error getting config-file flag: %w", err)
		}

		// Load config
		var cfg *config.Config
		if configFile == "" {
			configFile = config.FindConfigFile("")
			cfg, err = config.LoadConfig(nil, configFile)
		} else {
			cfg, err = config.LoadConfig(nil, configFile)
		}

		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		repoPath := cfg.Git.ClonePath

		// Check if repository exists
		exists, err := reposervices.IsRepositoryCloned(repoPath)
		if err != nil {
			return fmt.Errorf("error checking repository: %w", err)
		}
		if !exists {
			return fmt.Errorf("repository not cloned at %s\nRun: chtsht repo clone", repoPath)
		}

		fmt.Println("Cleaning up merged branches...")
		deletedBranches, err := reposervices.CleanupMergedBranches(repoPath, cfg.Git.Token)
		if err != nil {
			return fmt.Errorf("failed to cleanup branches: %w", err)
		}

		if len(deletedBranches) == 0 {
			fmt.Println("\nNo branches to clean up. ✓")
		} else {
			fmt.Printf("\nCleaned up %d branch(es):\n", len(deletedBranches))
			for _, branch := range deletedBranches {
				fmt.Printf("  ✓ %s\n", branch)
			}
		}

		return nil
	},
}
