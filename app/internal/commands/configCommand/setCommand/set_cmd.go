package setcommand

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/redjax/cheatsheets/internal/config"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// SetCmd represents the config set command
var SetCmd = &cobra.Command{
	Use:   "set <key> <value>",
	Short: "Set a configuration value",
	Long: `Set a configuration value by key path.
	
Examples:
  chtsht config set git.token ghp_xxxxx
  chtsht config set git.auto_branch true
  chtsht config set git.author_name "John Doe"`,
	Args: cobra.ExactArgs(2),
	RunE: runSet,
}

func runSet(cmd *cobra.Command, args []string) error {
	key := args[0]
	value := args[1]

	// Get the config file path
	configFile, err := cmd.Flags().GetString("config-file")
	if err != nil {
		return fmt.Errorf("error getting config-file flag: %w", err)
	}

	// Resolve config file path
	if configFile == "" {
		configFile = config.FindConfigFile("")
	}

	// Ensure config directory exists
	configDir := filepath.Dir(configFile)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	// Load existing config or create empty map
	var configData map[string]interface{}
	
	if _, err := os.Stat(configFile); err == nil {
		// File exists, load it
		data, err := os.ReadFile(configFile)
		if err != nil {
			return fmt.Errorf("failed to read config file: %w", err)
		}
		
		if err := yaml.Unmarshal(data, &configData); err != nil {
			return fmt.Errorf("failed to parse config file: %w", err)
		}
	} else {
		// File doesn't exist, start with empty config
		configData = make(map[string]interface{})
	}

	// Set the value in the config map
	if err := setNestedValue(configData, key, value); err != nil {
		return err
	}

	// Marshal back to YAML
	data, err := yaml.Marshal(configData)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	// Write back to file
	if err := os.WriteFile(configFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	fmt.Printf("✓ Set %s = %s\n", key, value)
	fmt.Printf("  Config file: %s\n", configFile)

	return nil
}

func setNestedValue(data map[string]interface{}, key, value string) error {
	parts := strings.Split(key, ".")
	
	// Navigate to the parent map
	current := data
	for i := 0; i < len(parts)-1; i++ {
		part := parts[i]
		
		if next, ok := current[part]; ok {
			// Key exists
			if nextMap, ok := next.(map[string]interface{}); ok {
				current = nextMap
			} else {
				return fmt.Errorf("cannot set nested key: %s is not a map", strings.Join(parts[:i+1], "."))
			}
		} else {
			// Create new nested map
			newMap := make(map[string]interface{})
			current[part] = newMap
			current = newMap
		}
	}

	// Set the final value with type inference
	finalKey := parts[len(parts)-1]
	current[finalKey] = inferType(value)

	return nil
}

func inferType(value string) interface{} {
	// Try to infer the type from the string value
	
	// Boolean
	if value == "true" {
		return true
	}
	if value == "false" {
		return false
	}

	// Integer
	if i, err := strconv.ParseInt(value, 10, 64); err == nil {
		return i
	}

	// Float
	if f, err := strconv.ParseFloat(value, 64); err == nil {
		return f
	}

	// Default to string
	return value
}
