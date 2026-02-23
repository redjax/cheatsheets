package cmd

import (
	"fmt"
	"os"

	cdcommand "github.com/redjax/cheatsheets/internal/commands/cdCommand"
	cleanupcommand "github.com/redjax/cheatsheets/internal/commands/cleanupCommand"
	debugcommand "github.com/redjax/cheatsheets/internal/commands/debugCommand"
	editcommand "github.com/redjax/cheatsheets/internal/commands/editCommand"
	listcommand "github.com/redjax/cheatsheets/internal/commands/listCommand"
	repocommand "github.com/redjax/cheatsheets/internal/commands/repoCommand"
	showcommand "github.com/redjax/cheatsheets/internal/commands/showCommand"
	synccommand "github.com/redjax/cheatsheets/internal/commands/syncCommand"
	"github.com/spf13/cobra"
)

var (
	configFile string
)

var rootCmd = &cobra.Command{
	Use:   "chtsht",
	Short: "A cheatsheet management CLI",
	Long:  `chtsht is a CLI tool for managing and accessing your personal cheatsheets stored in a Git repository.`,
	Run: func(cmd *cobra.Command, args []string) {
		// Default behavior when no subcommand is provided
		cmd.Help()
	},
}

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	// Register subcommands
	rootCmd.AddCommand(cdcommand.CdCmd)
	rootCmd.AddCommand(cleanupcommand.CleanupCmd)
	rootCmd.AddCommand(debugcommand.DebugCmd)
	rootCmd.AddCommand(editcommand.EditCmd)
	rootCmd.AddCommand(listcommand.ListCmd)
	rootCmd.AddCommand(repocommand.RepoCmd)
	rootCmd.AddCommand(showcommand.ShowCmd)
	rootCmd.AddCommand(synccommand.SyncCmd)

	// Global persistent flags
	// Empty default means "use default with .local fallback"
	rootCmd.PersistentFlags().StringVarP(&configFile, "config-file", "c", "", "config file path (default: config.yml with .local fallback)")
}
