package commitcommand

import (
	"fmt"

	"github.com/redjax/cheatsheets/internal/config"
	reposervices "github.com/redjax/cheatsheets/internal/services/repoServices"
	"github.com/spf13/cobra"
)

var (
	message    string
	autoCommit bool
)

// CommitCmd represents the commit command
var CommitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Commit staged changes",
	Long: `Commit staged changes with a message.
	
Examples:
  chtsht repo commit -m "Updated Linux cheatsheet"
  chtsht repo commit --auto-commit -m "Quick fix"  # Stage all + commit`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Require message flag
		if message == "" {
			return fmt.Errorf("commit message is required. Use -m flag: chtsht repo commit -m \"your message\"")
		}

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

		// If auto-commit, stage all changes first
		if autoCommit {
			clean, err := reposervices.IsWorkingTreeClean(cfg.Git.ClonePath)
			if err != nil {
				return fmt.Errorf("failed to check working tree: %w", err)
			}

			if !clean {
				fmt.Println("Staging all changes")
				err = reposervices.StageAll(cfg.Git.ClonePath)
				if err != nil {
					return fmt.Errorf("failed to stage changes: %w", err)
				}
			}
		}

		// Check if there are staged files
		staged, err := reposervices.GetStagedFiles(cfg.Git.ClonePath)
		if err != nil {
			return fmt.Errorf("failed to check staged files: %w", err)
		}

		if len(staged) == 0 {
			return fmt.Errorf("no changes staged for commit. Run 'chtsht repo stage' first, or use --auto-commit flag")
		}

		// Get git author info (from config or system git config)
		authorName, authorEmail, err := reposervices.GetGitAuthor(cfg.Git.AuthorName, cfg.Git.AuthorEmail)
		if err != nil {
			return fmt.Errorf("failed to get git author: %w\n\nYou can set author in config.yml or run:\n  git config --global user.name \"Your Name\"\n  git config --global user.email \"you@example.com\"", err)
		}

		// Get current branch
		branch, err := reposervices.GetCurrentBranch(cfg.Git.ClonePath)
		if err != nil {
			return fmt.Errorf("failed to get current branch: %w", err)
		}

		fmt.Printf("Committing %d file(s) on branch '%s'\n", len(staged), branch)

		// Create commit
		hash, err := reposervices.CommitChanges(cfg.Git.ClonePath, message, authorName, authorEmail)
		if err != nil {
			return fmt.Errorf("failed to create commit: %w", err)
		}

		fmt.Printf("Created commit %s\n", hash)
		fmt.Printf("  Author: %s <%s>\n", authorName, authorEmail)
		fmt.Printf("  Message: %s\n", message)
		fmt.Println("\nNext: chtsht repo push")

		return nil
	},
}

func init() {
	CommitCmd.Flags().StringVarP(&message, "message", "m", "", "Commit message (required)")
	CommitCmd.Flags().BoolVar(&autoCommit, "auto-commit", false, "Automatically stage all changes before committing")
}
