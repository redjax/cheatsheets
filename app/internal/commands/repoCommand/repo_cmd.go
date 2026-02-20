package repocommand

import (
	clonecommand "github.com/redjax/cheatsheets/internal/commands/repoCommand/cloneCommand"
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
	RepoCmd.AddCommand(clonecommand.CloneCmd)
	RepoCmd.AddCommand(validatecommand.ValidateCmd)
}
