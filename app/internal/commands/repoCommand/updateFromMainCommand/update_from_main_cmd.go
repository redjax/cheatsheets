package updatefrommaincommand

import (
	"fmt"

	"github.com/redjax/cheatsheets/internal/config"
	"github.com/redjax/cheatsheets/internal/guards"
	reposervices "github.com/redjax/cheatsheets/internal/services/repoServices"
	"github.com/spf13/cobra"
)

var UpdateFromMainCmd = &cobra.Command{
	Use:     "update-from-main",
	Aliases: []string{"merge-from-main"},
	Short:   "Merge changes from main into current working branch",
	Long: `Merge the latest changes from the main branch into your current working branch.
	
This command will:
  - Ensure you are not on the main branch
  - Ensure you have a clean working tree
  - Fetch and update the main branch
  - Merge main into your current branch
  - Push the updated branch to remote

This is useful when main has been updated and you want to incorporate
those changes into your working branch before continuing work.`,
	Example: `  chtsht repo update-from-main
  chtsht repo merge-from-main`,
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

		// Pre-flight checks using guards
		guardCtx := guards.NewGuardContext(cfg)
		if err := guards.CheckAll(guardCtx, guards.RepoCloned, guards.CleanWorkingTree, guards.OnWorkingBranch); err != nil {
			return err
		}

		repoPath := cfg.Git.ClonePath

		// Get current branch (guards already verified we're on working branch)
		currentBranch, err := reposervices.GetCurrentBranch(repoPath)
		if err != nil {
			return fmt.Errorf("failed to get current branch: %w", err)
		}

		fmt.Printf("Merging 'main' → '%s'\n\n", currentBranch)

		// Step 1: Fetch latest from remote
		fmt.Println("1. Fetching latest changes from remote")
		err = reposervices.UpdateRepository(repoPath, cfg.Git.Token)
		if err != nil {
			return fmt.Errorf("failed to fetch: %w", err)
		}

		// Step 2: Checkout main temporarily to update it
		fmt.Println("2. Updating main branch")
		originalBranch := currentBranch
		err = reposervices.CheckoutBranch(repoPath, "main")
		if err != nil {
			return fmt.Errorf("failed to checkout main: %w", err)
		}

		// Pull main
		err = reposervices.UpdateRepository(repoPath, cfg.Git.Token)
		if err != nil {
			// Try to switch back to original branch
			_ = reposervices.CheckoutBranch(repoPath, originalBranch)
			return fmt.Errorf("failed to pull main: %w", err)
		}

		// Step 3: Switch back to working branch
		fmt.Printf("3. Switching back to '%s'\n", originalBranch)
		err = reposervices.CheckoutBranch(repoPath, originalBranch)
		if err != nil {
			return fmt.Errorf("failed to switch back to %s: %w", originalBranch, err)
		}

		// Step 4: Merge main into current branch
		fmt.Printf("4. Merging main into '%s'\n", originalBranch)
		err = reposervices.MergeBranch(repoPath, "main")
		if err != nil {
			return fmt.Errorf("failed to merge main: %w\n\nYou may need to resolve conflicts manually.\nAfter resolving, run:\n  cd %s\n  git add .\n  git commit\n  git push", err, repoPath)
		}

		// Step 5: Push the updated branch
		fmt.Printf("5. Pushing '%s' to remote\n", originalBranch)
		err = reposervices.PushBranch(repoPath, originalBranch, cfg.Git.Token, false)
		if err != nil {
			return fmt.Errorf("failed to push: %w\n\nThe merge was completed locally but not pushed.\nYou can manually push with:\n  chtsht repo push", err)
		}

		fmt.Printf("\nSuccessfully merged 'main' into '%s' and pushed to remote.\n", originalBranch)
		fmt.Println("Your working branch is now up to date with main.")

		return nil
	},
}
