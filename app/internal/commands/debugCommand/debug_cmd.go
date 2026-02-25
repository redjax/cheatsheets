package debugcommand

import (
	"fmt"

	"github.com/redjax/cheatsheets/internal/config"
	debugservice "github.com/redjax/cheatsheets/internal/services/debugService"
	"github.com/spf13/cobra"
)

// DebugCmd represents the debug command
var DebugCmd = &cobra.Command{
	Use:   "debug",
	Short: "Print debug information about the configuration",
	Long:  `Displays the current configuration including paths, URLs, and other settings for debugging purposes.`,
	Run: func(command *cobra.Command, args []string) {
		// Get config file from persistent flag
		configFile, _ := command.Flags().GetString("config-file")

		var configPath string
		if configFile == "" {
			// No explicit config specified, use default with .local fallback
			configPath = config.FindConfigFile("")
		} else {
			// Explicit config specified, use it directly (no .local fallback)
			configPath = configFile
		}

		// Load configuration
		cfg, err := config.LoadConfig(nil, configPath)
		if err != nil {
			fmt.Printf("Failed to load config: %v\n", err)
			return
		}

		// Print debug info
		debugservice.PrintConfig(cfg)
	},
}
