package configcommand

import (
	editcommand "github.com/redjax/cheatsheets/internal/commands/configCommand/editCommand"
	setcommand "github.com/redjax/cheatsheets/internal/commands/configCommand/setCommand"
	showcommand "github.com/redjax/cheatsheets/internal/commands/configCommand/showCommand"
	"github.com/spf13/cobra"
)

// ConfigCmd represents the config command
var ConfigCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage configuration",
	Long:  `View and modify the cheatsheets configuration file.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Show help when no subcommand is provided
		cmd.Help()
	},
}

func init() {
	// Register subcommands
	ConfigCmd.AddCommand(editcommand.EditCmd)
	ConfigCmd.AddCommand(setcommand.SetCmd)
	ConfigCmd.AddCommand(showcommand.ShowCmd)
}
