package statuscommand

import (
	"fmt"

	"github.com/redjax/cheatsheets/internal/config"
	reposervices "github.com/redjax/cheatsheets/internal/services/repoServices"
	"github.com/spf13/cobra"
)

// StatusCmd represents the status command
var StatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show git repository status",
	Long:  `Displays detailed git status including branch, staged/unstaged files, and sync status with remote.`,
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

		// Check if repository is cloned
		cloned, err := reposervices.IsRepositoryCloned(cfg.Git.ClonePath)
		if err != nil {
			return fmt.Errorf("error checking repository: %w", err)
		}

		if !cloned {
			return fmt.Errorf("repository is not cloned at %s. Run 'chtsht repo clone' first", cfg.Git.ClonePath)
		}

		// Get detailed status
		status, err := reposervices.GetDetailedStatus(cfg.Git.ClonePath, cfg.Git.Token)
		if err != nil {
			return fmt.Errorf("failed to get repository status: %w", err)
		}

		// Display status
		fmt.Println("=== Repository Status ===")
		fmt.Printf("Location: %s\n", cfg.Git.ClonePath)
		fmt.Printf("Branch: %s\n", status.CurrentBranch)

		// Show last commit
		if status.LastCommit != nil {
			fmt.Printf("Last commit: %s - %s (%s)\n",
				status.LastCommit.Hash,
				status.LastCommit.Message,
				status.LastCommit.Date.Format("2006-01-02 15:04"))
		}

		// Show branch divergence
		if status.AheadBy > 0 || status.BehindBy > 0 {
			fmt.Print("Remote status: ")
			if status.AheadBy > 0 {
				fmt.Printf("ahead by %d commit(s)", status.AheadBy)
			}
			if status.AheadBy > 0 && status.BehindBy > 0 {
				fmt.Print(", ")
			}
			if status.BehindBy > 0 {
				fmt.Printf("behind by %d commit(s)", status.BehindBy)
			}
			fmt.Println()

			if status.BehindBy > 0 {
				fmt.Println("  → Run 'chtsht repo pull' to update")
			}
			if status.AheadBy > 0 {
				fmt.Println("  → Run 'chtsht repo push' to sync changes")
			}
		} else {
			fmt.Println("Remote status: up to date ✓")
		}

		fmt.Println()

		// Show file changes
		totalChanges := len(status.StagedFiles) + len(status.UnstagedFiles) + len(status.UntrackedFiles)

		if totalChanges == 0 {
			fmt.Println("Working tree clean ✓")
			fmt.Println("No changes to commit")
		} else {
			fmt.Printf("Changes detected (%d files):\n", totalChanges)
			fmt.Println()

			if len(status.StagedFiles) > 0 {
				fmt.Println("Staged for commit:")
				for _, file := range status.StagedFiles {
					fmt.Printf("  %s\n", file)
				}
				fmt.Println()
			}

			if len(status.UnstagedFiles) > 0 {
				fmt.Println("Modified (not staged):")
				for _, file := range status.UnstagedFiles {
					fmt.Printf("  • %s\n", file)
				}
				fmt.Println()
			}

			if len(status.UntrackedFiles) > 0 {
				fmt.Println("Untracked files:")
				for _, file := range status.UntrackedFiles {
					fmt.Printf("  ? %s\n", file)
				}
				fmt.Println()
			}

			// Show next steps
			if len(status.StagedFiles) > 0 {
				fmt.Println("Next: chtsht repo commit -m \"your message\"")
			} else if len(status.UnstagedFiles) > 0 || len(status.UntrackedFiles) > 0 {
				fmt.Println("Next: chtsht repo stage --all")
			}
		}

		return nil
	},
}
