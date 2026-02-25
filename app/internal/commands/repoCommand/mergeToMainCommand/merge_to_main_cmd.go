package mergetomaincommand

import (
	"fmt"

	"github.com/redjax/cheatsheets/internal/config"
	reposervices "github.com/redjax/cheatsheets/internal/services/repoServices"
	"github.com/spf13/cobra"
)

var MergeToMainCmd = &cobra.Command{
	Use:   "merge-to-main",
	Short: "Merge working branch into main",
	Long: `Merge your working branch into the main branch and push it to remote.
	
This command will:
  - Ensure you have a clean working tree
  - Checkout the main branch
  - Pull the latest changes
  - Merge your working branch into main
  - Push main to remote
  - Switch back to your working branch`,
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

		repoPath := cfg.Git.ClonePath

		// Get current branch
		currentBranch, err := reposervices.GetCurrentBranch(repoPath)
		if err != nil {
			return fmt.Errorf("failed to get current branch: %w", err)
		}

		// Check if already on main
		if currentBranch == "main" || currentBranch == "master" {
			return fmt.Errorf("you are already on the %s branch. Switch to your working branch first", currentBranch)
		}

		workingBranch := currentBranch

		// Check for uncommitted changes
		clean, err := reposervices.IsWorkingTreeClean(repoPath)
		if err != nil {
			return fmt.Errorf("failed to check working tree: %w", err)
		}

		if !clean {
			return fmt.Errorf("you have uncommitted changes. Commit or stash them first.\nRun 'chtsht repo status' to see your changes")
		}

		fmt.Printf("Merging '%s' → 'main'\n\n", workingBranch)

		// Step 1: Checkout main
		fmt.Println("1. Switching to main branch")
		err = reposervices.CheckoutBranch(repoPath, "main")
		if err != nil {
			return fmt.Errorf("failed to checkout main: %w", err)
		}

		// Step 2: Pull latest main
		fmt.Println("2. Pulling latest main")
		err = reposervices.UpdateRepository(repoPath, cfg.Git.Token)
		if err != nil {
			// Try to switch back to working branch
			_ = reposervices.CheckoutBranch(repoPath, workingBranch)
			return fmt.Errorf("failed to pull main: %w", err)
		}

		// Step 3: Merge working branch
		fmt.Printf("3. Merging '%s' into main\n", workingBranch)
		err = reposervices.MergeBranch(repoPath, workingBranch)
		if err != nil {
			// Try to switch back to working branch
			_ = reposervices.CheckoutBranch(repoPath, workingBranch)
			return fmt.Errorf("failed to merge: %w\n\nYou may need to resolve conflicts manually", err)
		}

		// Step 4: Push main
		fmt.Println("4. Pushing main to remote")
		err = reposervices.PushBranch(repoPath, "main", cfg.Git.Token, false)
		if err != nil {
			// Try to switch back to working branch
			_ = reposervices.CheckoutBranch(repoPath, workingBranch)
			return fmt.Errorf("failed to push main: %w\n\nThe merge was completed locally but not pushed", err)
		}

		// Step 5: Switch back to working branch
		fmt.Printf("5. Switching back to '%s'\n", workingBranch)
		err = reposervices.CheckoutBranch(repoPath, workingBranch)
		if err != nil {
			return fmt.Errorf("merge successful but failed to switch back to %s: %w", workingBranch, err)
		}

		fmt.Printf("\nSuccessfully merged '%s' into 'main' and pushed to remote.\n", workingBranch)
		fmt.Printf("Back on branch '%s'.\n", workingBranch)

		return nil
	},
}
