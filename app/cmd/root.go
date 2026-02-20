package cmd

import (
	"fmt"
	"os"

	debugcommand "github.com/redjax/cheatsheets/internal/commands/debugCommand"
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
	rootCmd.AddCommand(debugcommand.DebugCmd)

	// Global persistent flags
	// Empty default means "use default with .local fallback"
	rootCmd.PersistentFlags().StringVarP(&configFile, "config-file", "c", "", "config file path (default: config.yml with .local fallback)")
}
