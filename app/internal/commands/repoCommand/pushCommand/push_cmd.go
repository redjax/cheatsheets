package pushcommand

import (
	"fmt"

	"github.com/redjax/cheatsheets/internal/config"
	"github.com/redjax/cheatsheets/internal/guards"
	reposervices "github.com/redjax/cheatsheets/internal/services/repoServices"
	"github.com/spf13/cobra"
)

var (
	setUpstream bool
)

// PushCmd represents the push command
var PushCmd = &cobra.Command{
	Use:   "push",
	Short: "Push commits to remote",
	Long: `Push the current branch to the remote repository.
	
Examples:
  chtsht repo push              # Push current branch
  chtsht repo push --set-upstream  # Push and set upstream tracking`,
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

		// Pre-flight check
		guardCtx := guards.NewGuardContext(cfg)
		if err := guards.CheckAll(guardCtx, guards.RepoCloned); err != nil {
			return err
		}

		// Check if working tree is clean
		clean, err := reposervices.IsWorkingTreeClean(cfg.Git.ClonePath)
		if err != nil {
			return fmt.Errorf("failed to check working tree: %w", err)
		}

		if !clean {
			return fmt.Errorf("you have uncommitted changes. Commit them first with 'chtsht repo commit'")
		}

		// Get current branch
		branch, err := reposervices.GetCurrentBranch(cfg.Git.ClonePath)
		if err != nil {
			return fmt.Errorf("failed to get current branch: %w", err)
		}

		// Check if we need to set upstream (first push)
		// We'll try to get divergence - if remote doesn't exist, we need upstream
		_, _, err = reposervices.GetBranchDivergence(cfg.Git.ClonePath, branch, cfg.Git.Token)
		needsUpstream := err != nil

		if needsUpstream && !setUpstream {
			setUpstream = true
			fmt.Printf("Setting upstream for new branch '%s'\n", branch)
		}

		fmt.Printf("Pushing branch '%s' to origin\n", branch)

		// Push
		err = reposervices.PushBranch(cfg.Git.ClonePath, branch, cfg.Git.Token, setUpstream)
		if err != nil {
			return fmt.Errorf("failed to push: %w", err)
		}

		fmt.Printf("Successfully pushed '%s' to origin\n", branch)

		return nil
	},
}

func init() {
	PushCmd.Flags().BoolVarP(&setUpstream, "set-upstream", "u", false, "Set upstream tracking for the branch")
}
