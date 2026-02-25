package pullcommand

import (
	"fmt"

	"github.com/redjax/cheatsheets/internal/config"
	reposervices "github.com/redjax/cheatsheets/internal/services/repoServices"
	"github.com/spf13/cobra"
)

var PullCmd = &cobra.Command{
	Use:   "pull",
	Short: "Pull changes from remote repository",
	Long:  "Pull the latest changes from the remote repository for the current branch.",
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

		// Check what branch we're on
		currentBranch, err := reposervices.GetCurrentBranch(repoPath)
		if err != nil {
			return fmt.Errorf("failed to get current branch: %w", err)
		}

		fmt.Printf("Pulling branch '%s' from remote\n", currentBranch)

		// Pull changes
		err = reposervices.UpdateRepository(repoPath, cfg.Git.Token)
		if err != nil {
			return fmt.Errorf("failed to pull changes: %w", err)
		}

		// Get status after pull
		status, err := reposervices.GetDetailedStatus(repoPath, cfg.Git.Token)
		if err == nil {
			if status.AheadBy > 0 {
				fmt.Printf("Pull complete. Your branch is %d commit(s) ahead of remote.\n", status.AheadBy)
			} else if status.BehindBy > 0 {
				fmt.Printf("Pull complete. Fast-forwarded %d commit(s).\n", status.BehindBy)
			} else {
				fmt.Println("Already up to date.")
			}

			if len(status.StagedFiles) > 0 || len(status.UnstagedFiles) > 0 || len(status.UntrackedFiles) > 0 {
				fmt.Println("\nNote: You have uncommitted local changes.")
			}
		} else {
			fmt.Println("Pull complete.")
		}

		return nil
	},
}
