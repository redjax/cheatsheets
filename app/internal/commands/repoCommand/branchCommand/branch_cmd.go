package branchcommand

import (
	"fmt"

	"github.com/redjax/cheatsheets/internal/config"
	reposervices "github.com/redjax/cheatsheets/internal/services/repoServices"
	"github.com/spf13/cobra"
)

// BranchCmd represents the branch command
var BranchCmd = &cobra.Command{
	Use:   "branch",
	Short: "Manage git branches",
	Long:  `View and manage git branches. Without arguments, shows the current branch.`,
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

		// Check if repository is cloned
		cloned, err := reposervices.IsRepositoryCloned(cfg.Git.ClonePath)
		if err != nil {
			return fmt.Errorf("error checking repository: %w", err)
		}

		if !cloned {
			return fmt.Errorf("repository is not cloned at %s. Run 'chtsht repo clone' first", cfg.Git.ClonePath)
		}

		// Get current branch
		branch, err := reposervices.GetCurrentBranch(cfg.Git.ClonePath)
		if err != nil {
			return fmt.Errorf("failed to get current branch: %w", err)
		}

		fmt.Printf("Current branch: %s\n", branch)

		return nil
	},
}

// ListCmd lists all branches
var ListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all branches",
	Long:  `List all local branches, marking the current branch.`,
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

		// Check if repository is cloned
		cloned, err := reposervices.IsRepositoryCloned(cfg.Git.ClonePath)
		if err != nil {
			return fmt.Errorf("error checking repository: %w", err)
		}

		if !cloned {
			return fmt.Errorf("repository is not cloned at %s. Run 'chtsht repo clone' first", cfg.Git.ClonePath)
		}

		// Get current branch
		currentBranch, err := reposervices.GetCurrentBranch(cfg.Git.ClonePath)
		if err != nil {
			return fmt.Errorf("failed to get current branch: %w", err)
		}

		// List all branches
		branches, err := reposervices.ListBranches(cfg.Git.ClonePath)
		if err != nil {
			return fmt.Errorf("failed to list branches: %w", err)
		}

		fmt.Println("Local branches:")
		for _, branch := range branches {
			if branch == currentBranch {
				fmt.Printf("* %s (current)\n", branch)
			} else {
				fmt.Printf("  %s\n", branch)
			}
		}

		return nil
	},
}

// EnsureCmd ensures the machine branch exists and is checked out
var EnsureCmd = &cobra.Command{
	Use:   "ensure",
	Short: "Ensure machine-specific branch",
	Long:  `Ensure the machine-specific branch (e.g., local-hostname) exists and is checked out.`,
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

		// Check if repository is cloned
		cloned, err := reposervices.IsRepositoryCloned(cfg.Git.ClonePath)
		if err != nil {
			return fmt.Errorf("error checking repository: %w", err)
		}

		if !cloned {
			return fmt.Errorf("repository is not cloned at %s. Run 'chtsht repo clone' first", cfg.Git.ClonePath)
		}

		// Check if working tree is clean
		clean, err := reposervices.IsWorkingTreeClean(cfg.Git.ClonePath)
		if err != nil {
			return fmt.Errorf("failed to check working tree: %w", err)
		}

		if !clean {
			return fmt.Errorf("you have uncommitted changes. Commit or stash them before switching branches")
		}

		// Get branch prefix from config or use default
		prefix := cfg.Git.BranchPrefix
		if prefix == "" {
			prefix = "local"
		}

		// Ensure machine branch
		branchName, err := reposervices.EnsureMachineBranch(cfg.Git.ClonePath, prefix)
		if err != nil {
			return fmt.Errorf("failed to ensure machine branch: %w", err)
		}

		fmt.Printf("✓ Now on branch '%s'\n", branchName)

		return nil
	},
}

// SwitchCmd switches to a different branch
var SwitchCmd = &cobra.Command{
	Use:   "switch <branch>",
	Short: "Switch to a different branch",
	Long:  `Switch to an existing branch.`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		branchName := args[0]

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

		// Check if repository is cloned
		cloned, err := reposervices.IsRepositoryCloned(cfg.Git.ClonePath)
		if err != nil {
			return fmt.Errorf("error checking repository: %w", err)
		}

		if !cloned {
			return fmt.Errorf("repository is not cloned at %s. Run 'chtsht repo clone' first", cfg.Git.ClonePath)
		}

		// Check if working tree is clean
		clean, err := reposervices.IsWorkingTreeClean(cfg.Git.ClonePath)
		if err != nil {
			return fmt.Errorf("failed to check working tree: %w", err)
		}

		if !clean {
			return fmt.Errorf("you have uncommitted changes. Commit or stash them before switching branches")
		}

		// Check if branch exists
		exists, err := reposervices.BranchExists(cfg.Git.ClonePath, branchName)
		if err != nil {
			return fmt.Errorf("failed to check branch: %w", err)
		}

		if !exists {
			return fmt.Errorf("branch '%s' does not exist locally", branchName)
		}

		// Switch to branch
		err = reposervices.CheckoutBranch(cfg.Git.ClonePath, branchName)
		if err != nil {
			return fmt.Errorf("failed to switch branch: %w", err)
		}

		fmt.Printf("✓ Switched to branch '%s'\n", branchName)

		return nil
	},
}

func init() {
	BranchCmd.AddCommand(ListCmd)
	BranchCmd.AddCommand(EnsureCmd)
	BranchCmd.AddCommand(SwitchCmd)
}
