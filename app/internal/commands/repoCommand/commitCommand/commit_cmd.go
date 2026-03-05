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
	autoPush   bool
)

// CommitCmd represents the commit command
var CommitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Commit staged changes",
	Long: `Commit staged changes with a message.
	
Examples:
  chtsht repo commit -m "Updated Linux cheatsheet"
  chtsht repo commit --auto-commit                     # Stage all + auto-generate message
  chtsht repo commit --auto-commit --push              # Stage all + commit + push
  chtsht repo commit --auto-commit -m "Custom" --push  # Stage all + custom message + push`,
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

		// If auto-commit with auto-push, use the all-in-one function
		if autoCommit && autoPush {
			// Check if there are changes
			clean, err := reposervices.IsWorkingTreeClean(cfg.Git.ClonePath)
			if err != nil {
				return fmt.Errorf("failed to check working tree: %w", err)
			}

			if clean {
				fmt.Println("No changes to commit")
				return nil
			}

			// Generate message if not provided
			if message == "" {
				// Get list of changed files for the message
				fileStatuses, err := reposervices.GetDetailedFileStatus(cfg.Git.ClonePath)
				if err == nil && len(fileStatuses) > 0 {
					// Count files and generate message
					if len(fileStatuses) == 1 {
						message = fmt.Sprintf("update %s", fileStatuses[0].Path)
					} else if len(fileStatuses) <= 3 {
						message = fmt.Sprintf("update %d files", len(fileStatuses))
					} else {
						message = fmt.Sprintf("update %d files", len(fileStatuses))
					}
				} else {
					message = "update multiple files"
				}
				fmt.Printf("Auto-generated message: %s\n", message)
			}

			// Use the all-in-one StageCommitAndPush function
			return reposervices.StageCommitAndPush(
				cfg.Git.ClonePath,
				message,
				cfg.Git.AuthorName,
				cfg.Git.AuthorEmail,
				cfg.Git.Token,
			)
		}

		// Original logic for commit-only

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

			// If no message provided with --auto-commit, generate one
			if message == "" {
				// Get list of changed files for the message
				staged, err := reposervices.GetStagedFiles(cfg.Git.ClonePath)
				if err == nil && len(staged) > 0 {
					// Generate message based on files
					if len(staged) == 1 {
						message = fmt.Sprintf("update %s", staged[0])
					} else if len(staged) <= 3 {
						message = fmt.Sprintf("update %d files", len(staged))
					} else {
						message = fmt.Sprintf("update %d files", len(staged))
					}
				} else {
					message = "update multiple files"
				}
				fmt.Printf("Auto-generated message: %s\n", message)
			}
		} else {
			// Not using auto-commit, message is required
			if message == "" {
				return fmt.Errorf("commit message is required. Use -m flag: chtsht repo commit -m \"your message\"")
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

		// If --push flag is set, push after committing
		if autoPush {
			fmt.Println("\nPushing to remote...")
			err = reposervices.PushBranch(cfg.Git.ClonePath, branch, cfg.Git.Token, false)
			if err != nil {
				return fmt.Errorf("commit succeeded but push failed: %w", err)
			}
			fmt.Printf("Pushed changes to remote branch '%s'\n", branch)
		} else {
			fmt.Println("\nNext: chtsht repo push")
		}

		return nil
	},
}

func init() {
	CommitCmd.Flags().StringVarP(&message, "message", "m", "", "Commit message (optional with --auto-commit)")
	CommitCmd.Flags().BoolVarP(&autoCommit, "auto-commit", "a", false, "Automatically stage all changes and generate message if not provided")
	CommitCmd.Flags().BoolVarP(&autoPush, "push", "p", false, "Push to remote after committing")
}
