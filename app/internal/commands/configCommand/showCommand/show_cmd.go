package showcommand

import (
	"fmt"

	"github.com/redjax/cheatsheets/internal/config"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var asJson bool

// ShowCmd represents the config show command
var ShowCmd = &cobra.Command{
	Use:   "show [key]",
	Short: "Show configuration",
	Long: `Display the current configuration or a specific value.
	
Examples:
  chtsht config show              # Show all config
  chtsht config show git.token    # Show specific value
  chtsht config show git          # Show git section`,
	RunE: runShow,
}

func init() {
	ShowCmd.Flags().BoolVar(&asJson, "json", false, "Output as JSON instead of YAML")
}

func runShow(cmd *cobra.Command, args []string) error {
	// Get the config file path
	configFile, err := cmd.Flags().GetString("config-file")
	if err != nil {
		return fmt.Errorf("error getting config-file flag: %w", err)
	}

	// Resolve config file path
	if configFile == "" {
		configFile = config.FindConfigFile("")
	}

	// Load config
	cfg, err := config.LoadConfig(nil, configFile)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Print config file location
	fmt.Printf("# Config file: %s\n\n", configFile)

	// If no key specified, show entire config
	if len(args) == 0 {
		return showFullConfig(cfg)
	}

	// Show specific key
	key := args[0]
	return showKey(cfg, key)
}

func showFullConfig(cfg *config.Config) error {
	// Marshal to YAML for display
	data, err := yaml.Marshal(cfg)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	fmt.Print(string(data))
	return nil
}

func showKey(cfg *config.Config, key string) error {
	// Use koanf to get the value by key path
	value := config.K.Get(key)

	if value == nil {
		return fmt.Errorf("key not found: %s", key)
	}

	// If it's a complex type (map/struct), show as YAML
	switch v := value.(type) {
	case map[string]interface{}:
		data, err := yaml.Marshal(v)
		if err != nil {
			return fmt.Errorf("failed to marshal value: %w", err)
		}
		fmt.Print(string(data))
	default:
		// Simple value, just print it
		fmt.Println(value)
	}

	return nil
}
