package selfcommand

import (
	upgradecommand "github.com/redjax/cheatsheets/internal/commands/selfCommand/upgradeCommand"
	versioncommand "github.com/redjax/cheatsheets/internal/commands/selfCommand/versionCommand"
	"github.com/spf13/cobra"
)

// SelfCmd represents the self command
var SelfCmd = &cobra.Command{
	Use:   "self",
	Short: "Manage the chtsht application itself",
	Long:  `Commands for managing the chtsht application, including version information and upgrades.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Show help when no subcommand is provided
		cmd.Help()
	},
}

func init() {
	// Register subcommands
	SelfCmd.AddCommand(versioncommand.VersionCmd)
	SelfCmd.AddCommand(upgradecommand.UpgradeCmd)
}
