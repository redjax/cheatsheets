package repocommand

import (
	branchcommand "github.com/redjax/cheatsheets/internal/commands/repoCommand/branchCommand"
	clonecommand "github.com/redjax/cheatsheets/internal/commands/repoCommand/cloneCommand"
	commitcommand "github.com/redjax/cheatsheets/internal/commands/repoCommand/commitCommand"
	pushcommand "github.com/redjax/cheatsheets/internal/commands/repoCommand/pushCommand"
	stagecommand "github.com/redjax/cheatsheets/internal/commands/repoCommand/stageCommand"
	statuscommand "github.com/redjax/cheatsheets/internal/commands/repoCommand/statusCommand"
	validatecommand "github.com/redjax/cheatsheets/internal/commands/repoCommand/validateCommand"
	"github.com/spf13/cobra"
)

// RepoCmd represents the repo command
var RepoCmd = &cobra.Command{
	Use:   "repo",
	Short: "Manage the cheatsheets repository",
	Long:  `Commands for managing the cloned cheatsheets repository, including validation, updates, and status checks.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Show help when no subcommand is provided
		cmd.Help()
	},
}

func init() {
	// Register subcommands
	RepoCmd.AddCommand(branchcommand.BranchCmd)
	RepoCmd.AddCommand(clonecommand.CloneCmd)
	RepoCmd.AddCommand(commitcommand.CommitCmd)
	RepoCmd.AddCommand(pushcommand.PushCmd)
	RepoCmd.AddCommand(stagecommand.StageCmd)
	RepoCmd.AddCommand(statuscommand.StatusCmd)
	RepoCmd.AddCommand(validatecommand.ValidateCmd)
}
