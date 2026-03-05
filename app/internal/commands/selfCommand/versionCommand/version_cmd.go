package versioncommand

import (
	"fmt"

	"github.com/redjax/cheatsheets/internal/version"
	"github.com/spf13/cobra"
)

// VersionCmd represents the version command
var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show the current version of chtsht",
	Long:  `Display the current version, build commit, and build date of the chtsht application.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("chtsht version %s (commit: %s, built: %s)\n",
			version.GetVersion(), version.GetCommit(), version.GetDate())
	},
}
