package synccommand

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/redjax/cheatsheets/internal/config"
	reposervices "github.com/redjax/cheatsheets/internal/services/repoServices"
	"github.com/spf13/cobra"
)

var message string

var SyncCmd = &cobra.Command{
	Use:   "sync",
	Short: "Sync repository (pull, commit, push)",
	Long: `Convenience command to sync your working branch with remote.
	
This command will:
1. Pull the latest changes from remote
2. Stage all changes (if any)
3. Prompt for a commit message (if changes exist)
4. Commit the changes
5. Push to remote

This is equivalent to running:
  chtsht repo pull
  chtsht repo stage --all
  chtsht repo commit -m "message"
  chtsht repo push`,
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

		// Warn if on main
		if currentBranch == "main" || currentBranch == "master" {
			fmt.Println("⚠️  Warning: You are on the main branch.")
			fmt.Println("Consider using 'chtsht repo branch ensure' to switch to your working branch.")
		}

		fmt.Printf("Syncing branch '%s'\n\n", currentBranch)

		// Step 1: Pull
		fmt.Println("1. Pulling from remote")
		err = reposervices.UpdateRepository(repoPath, cfg.Git.Token)
		if err != nil {
			return fmt.Errorf("failed to pull: %w", err)
		}
		fmt.Println("   Pull complete")

		// Step 2: Check for local changes
		status, err := reposervices.GetDetailedStatus(repoPath, cfg.Git.Token)
		if err != nil {
			return fmt.Errorf("failed to get status: %w", err)
		}

		hasChanges := len(status.StagedFiles) > 0 || len(status.UnstagedFiles) > 0 || len(status.UntrackedFiles) > 0
		if !hasChanges {
			fmt.Println("\nNo local changes to sync. Already up to date.")
			return nil
		}

		// Step 3: Stage all changes
		fmt.Println("\n2. Staging all changes")
		stagedFiles, err := reposervices.GetStagedFiles(repoPath)
		if err != nil {
			return fmt.Errorf("failed to get staged files: %w", err)
		}

		unstagedFiles, err := reposervices.GetUnstagedFiles(repoPath)
		if err != nil {
			return fmt.Errorf("failed to get unstaged files: %w", err)
		}

		untrackedFiles, err := reposervices.GetUntrackedFiles(repoPath)
		if err != nil {
			return fmt.Errorf("failed to get untracked files: %w", err)
		}

		// Only stage if there are unstaged or untracked files
		if len(unstagedFiles) > 0 || len(untrackedFiles) > 0 {
			err = reposervices.StageAll(repoPath)
			if err != nil {
				return fmt.Errorf("failed to stage changes: %w", err)
			}
			fmt.Printf("   Staged %d file(s)\n", len(unstagedFiles)+len(untrackedFiles))
		} else if len(stagedFiles) > 0 {
			fmt.Printf("   %d file(s) already staged\n", len(stagedFiles))
		}

		// Step 4: Get commit message
		fmt.Println("\n3. Creating commit")
		commitMessage := message
		if commitMessage == "" {
			// Prompt for commit message
			fmt.Print("   Enter commit message: ")
			reader := bufio.NewReader(os.Stdin)
			commitMessage, err = reader.ReadString('\n')
			if err != nil {
				return fmt.Errorf("failed to read commit message: %w", err)
			}
			commitMessage = strings.TrimSpace(commitMessage)

			if commitMessage == "" {
				return fmt.Errorf("commit message cannot be empty")
			}
		}

		// Get author info
		authorName, authorEmail, err := reposervices.GetGitAuthor(cfg.Git.AuthorName, cfg.Git.AuthorEmail)
		if err != nil {
			return fmt.Errorf("failed to get git author: %w\n\nYou can set author in config.yml or run:\n  git config --global user.name \"Your Name\"\n  git config --global user.email \"you@example.com\"", err)
		}

		// Commit changes
		hash, err := reposervices.CommitChanges(repoPath, commitMessage, authorName, authorEmail)
		if err != nil {
			return fmt.Errorf("failed to commit: %w", err)
		}

		shortHash := hash
		if len(hash) > 7 {
			shortHash = hash[:7]
		}
		fmt.Printf("   Committed as %s\n", shortHash)

		// Step 5: Push
		fmt.Println("\n4. Pushing to remote")
		err = reposervices.PushBranch(repoPath, currentBranch, cfg.Git.Token, false)
		if err != nil {
			return fmt.Errorf("failed to push: %w", err)
		}
		fmt.Println("   Push complete")

		fmt.Println("\nSync complete!")

		return nil
	},
}

func init() {
	SyncCmd.Flags().StringVarP(&message, "message", "m", "", "Commit message (if not provided, will prompt)")
}
