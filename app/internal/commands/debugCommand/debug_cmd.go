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
	Run: func(cmd *cobra.Command, args []string) {
		// Load configuration
		configFile := config.FindConfigFile("config.yml")
		cfg, err := config.LoadConfig(nil, configFile)
		if err != nil {
			fmt.Printf("Failed to load config: %v\n", err)
			return
		}

		// Print debug info
		debugservice.PrintConfig(cfg)
	},
}
