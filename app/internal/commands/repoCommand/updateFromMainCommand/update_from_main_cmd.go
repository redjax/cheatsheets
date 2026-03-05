package updatefrommaincommand

import (
	"fmt"
	"os"
	"strings"

	"github.com/redjax/cheatsheets/internal/config"
	"github.com/redjax/cheatsheets/internal/guards"
	reposervices "github.com/redjax/cheatsheets/internal/services/repoServices"
	"github.com/spf13/cobra"
)

var (
	useRebase  bool
	showManual bool
)

var UpdateFromMainCmd = &cobra.Command{
	Use:           "update-from-main",
	Aliases:       []string{"merge-from-main"},
	Short:         "Merge or rebase changes from main into current working branch",
	SilenceUsage:  true,
	SilenceErrors: true,
	Long: `Merge or rebase the latest changes from the main branch into your current branch.

NOTE: This command is only useful if you're working on a separate branch.
If you're working directly on main, you don't need this command - just use 'chtsht repo pull'.

This command will:
  - Ensure you are not on the main branch
  - Ensure you have a clean working tree
  - Fetch and update the main branch
  - Merge or rebase main into your current branch (based on --rebase flag)
  - Push the updated branch to remote

Use --rebase for a cleaner history (but requires force-push if already published).
Use --manual to see manual git commands instead of automatic execution.

This is useful when main has been updated and you want to incorporate
those changes into your branch before continuing work.`,
	Example: `  chtsht repo update-from-main              # Merge main into current branch
  chtsht repo update-from-main --rebase     # Rebase current branch onto main
  chtsht repo update-from-main --manual     # Show manual git commands`,
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

		// If --manual flag is set, show manual instructions
		if showManual {
			return showManualInstructions(repoPath, useRebase)
		}

		// Pre-flight checks using guards
		guardCtx := guards.NewGuardContext(cfg)
		if err := guards.CheckAll(guardCtx, guards.RepoCloned, guards.CleanWorkingTree, guards.OnWorkingBranch); err != nil {
			return err
		}

		// Get current branch (guards already verified we're on working branch)
		currentBranch, err := reposervices.GetCurrentBranch(repoPath)
		if err != nil {
			return fmt.Errorf("failed to get current branch: %w", err)
		}

		operation := "Merging"
		if useRebase {
			operation = "Rebasing"
		}
		fmt.Printf("%s 'main' → '%s'\n\n", operation, currentBranch)

		// Step 1: Fetch latest from remote
		fmt.Println("1. Fetching latest changes from remote")
		err = reposervices.UpdateRepositoryQuiet(repoPath, cfg.Git.Token, true)
		if err != nil {
			return formatError("fetch", err, repoPath, currentBranch, useRebase)
		}

		// Step 2: Checkout main temporarily to update it
		fmt.Println("2. Updating main branch")
		originalBranch := currentBranch
		err = reposervices.CheckoutBranch(repoPath, "main")
		if err != nil {
			return formatError("checkout main", err, repoPath, currentBranch, useRebase)
		}

		// Pull main
		err = reposervices.UpdateRepositoryQuiet(repoPath, cfg.Git.Token, true)
		if err != nil {
			// Try to switch back to original branch
			_ = reposervices.CheckoutBranch(repoPath, originalBranch)
			return formatError("pull main", err, repoPath, currentBranch, useRebase)
		}

		// Step 3: Switch back to working branch
		fmt.Printf("3. Switching back to '%s'\n", originalBranch)
		err = reposervices.CheckoutBranch(repoPath, originalBranch)
		if err != nil {
			return formatError("checkout working branch", err, repoPath, currentBranch, useRebase)
		}

		// Step 4: Merge or Rebase main into current branch
		if useRebase {
			fmt.Printf("4. Rebasing '%s' onto 'main'\n", originalBranch)
			err = reposervices.RebaseBranch(repoPath, "main")
		} else {
			fmt.Printf("4. Merging 'main' into '%s'\n", originalBranch)
			err = reposervices.MergeBranch(repoPath, "main")
		}

		if err != nil {
			return formatError(operation, err, repoPath, currentBranch, useRebase)
		}

		// Step 5: Push the updated branch
		fmt.Printf("5. Pushing '%s' to remote\n", originalBranch)

		// Use force push for rebase, normal push for merge
		pushOpts := reposervices.PushOptions{
			Force: useRebase,
		}
		err = reposervices.PushBranchWithOptions(repoPath, originalBranch, cfg.Git.Token, pushOpts)
		if err != nil {
			return formatError("push", err, repoPath, currentBranch, useRebase)
		}

		successMsg := "merged"
		if useRebase {
			successMsg = "rebased onto"
		}
		fmt.Printf("\nSuccessfully %s 'main' and pushed to remote.\n", successMsg)
		fmt.Println("Your working branch is now up to date with main.")

		return nil
	},
}

func init() {
	// Run function wraps the command to print errors ourselves (since SilenceErrors is true)
	originalRunE := UpdateFromMainCmd.RunE
	UpdateFromMainCmd.RunE = func(cmd *cobra.Command, args []string) error {
		err := originalRunE(cmd, args)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		}
		return err
	}

	UpdateFromMainCmd.Flags().BoolVar(&useRebase, "rebase", false, "Rebase instead of merge (creates cleaner history)")
	UpdateFromMainCmd.Flags().BoolVar(&showManual, "manual", false, "Show manual git commands instead of executing automatically")
}

func showManualInstructions(repoPath string, rebase bool) error {
	fmt.Println("=== Manual Update from Main Instructions ===")
	fmt.Printf("Repository: %s\n\n", repoPath)

	if rebase {
		fmt.Println("To rebase your working branch onto main manually:")
		fmt.Println("  1. Navigate to repository:")
		fmt.Printf("     cd %s (or chtsht cd)\n\n", repoPath)
		fmt.Println("  2. Fetch latest changes:")
		fmt.Println("     git fetch origin")
		fmt.Println()
		fmt.Println("  3. Update main branch:")
		fmt.Println("     git checkout main")
		fmt.Println("     git pull origin main")
		fmt.Println()
		fmt.Println("  4. Return to your working branch:")
		fmt.Println("     git checkout working  # or your branch name")
		fmt.Println()
		fmt.Println("  5. Rebase onto main:")
		fmt.Println("     git rebase main")
		fmt.Println()
		fmt.Println("  6. If conflicts occur:")
		fmt.Println("     - Edit conflicting files")
		fmt.Println("     - git add <resolved-files>")
		fmt.Println("     - git rebase --continue")
		fmt.Println("     - Repeat until rebase completes")
		fmt.Println()
		fmt.Println("  7. Push changes (requires force):")
		fmt.Println("     git push --force-with-lease origin working")
		fmt.Println()
		fmt.Println("Note: --force-with-lease is safer than --force as it checks")
		fmt.Println("      that you're not overwriting someone else's changes.")
	} else {
		fmt.Println("To merge main into your working branch manually:")
		fmt.Println("  1. Navigate to repository:")
		fmt.Printf("     cd %s (or chtsht cd)\n\n", repoPath)
		fmt.Println("  2. Fetch latest changes:")
		fmt.Println("     git fetch origin")
		fmt.Println()
		fmt.Println("  3. Update main branch:")
		fmt.Println("     git checkout main")
		fmt.Println("     git pull origin main")
		fmt.Println()
		fmt.Println("  4. Return to your working branch:")
		fmt.Println("     git checkout working  # or your branch name")
		fmt.Println()
		fmt.Println("  5. Merge main:")
		fmt.Println("     git merge main")
		fmt.Println()
		fmt.Println("  6. If conflicts occur:")
		fmt.Println("     - Edit conflicting files")
		fmt.Println("     - git add <resolved-files>")
		fmt.Println("     - git commit")
		fmt.Println()
		fmt.Println("  7. Push changes:")
		fmt.Println("     git push origin working")
	}

	fmt.Println("\n===========================================")
	return nil
}

func formatError(operation string, err error, repoPath, currentBranch string, rebase bool) error {
	var sb strings.Builder

	// Check if this is a divergence/conflict error
	errMsg := err.Error()
	isDivergence := strings.Contains(errMsg, "diverged") || strings.Contains(errMsg, "conflict")

	if isDivergence {
		sb.WriteString(fmt.Sprintf("Cannot automatically %s - branches have diverged.\n\n", operation))
		sb.WriteString("Manual steps to resolve:\n\n")
		sb.WriteString(fmt.Sprintf("  1. cd %s\n", repoPath))

		if rebase {
			sb.WriteString("  2. git rebase main\n")
			sb.WriteString("  3. Resolve conflicts:\n")
			sb.WriteString("     - Edit conflicting files\n")
			sb.WriteString("     - git add <resolved-files>\n")
			sb.WriteString("     - git rebase --continue\n")
			sb.WriteString("  4. git push --force-with-lease origin " + currentBranch + "\n")
		} else {
			sb.WriteString("  2. git merge main\n")
			sb.WriteString("  3. Resolve conflicts:\n")
			sb.WriteString("     - Edit conflicting files\n")
			sb.WriteString("     - git add <resolved-files>\n")
			sb.WriteString("     - git commit\n")
			sb.WriteString("  4. git push origin " + currentBranch + "\n")
		}
		sb.WriteString("\nFor detailed instructions: chtsht repo update-from-main --manual")
	} else {
		// Generic error
		sb.WriteString(fmt.Sprintf("Failed to %s: %v\n\n", operation, err))
		sb.WriteString("Manual steps:\n\n")
		sb.WriteString(fmt.Sprintf("  1. cd %s\n", repoPath))
		sb.WriteString("  2. git status  # Check repository state\n")
		if operation == "push" {
			sb.WriteString("  3. chtsht repo push  # Retry push\n")
		} else {
			sb.WriteString("  3. Resolve the issue manually\n")
		}
	}

	return fmt.Errorf("%s", sb.String())
}
