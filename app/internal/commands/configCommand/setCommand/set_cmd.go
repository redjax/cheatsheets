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
	
Keys can use either camelCase or snake_case - both will work:
  - git.workingbranch or git.working_branch
  - git.authorname or git.author_name
	
Examples:
  chtsht config set git.token ghp_xxxxx
  chtsht config set git.auto_branch true
  chtsht config set git.author_name "John Doe"
  chtsht config set git.workingbranch main`,
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
		part := normalizeKey(parts[i])

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
	finalKey := normalizeKey(parts[len(parts)-1])
	current[finalKey] = inferType(value)

	return nil
}

// normalizeKey converts a key to snake_case to match koanf struct tags
// Examples: workingbranch -> working_branch, authorname -> author_name
func normalizeKey(key string) string {
	// Common config key mappings (add more as needed)
	keyMap := map[string]string{
		"sheetspath":    "sheets_path",
		"repourl":       "repo_url",
		"clonepath":     "clone_path",
		"autobranch":    "auto_branch",
		"workingbranch": "working_branch",
		"authorname":    "author_name",
		"authoremail":   "author_email",
	}

	// Convert to lowercase first for case-insensitive matching
	lowerKey := strings.ToLower(key)

	// Check if we have a direct mapping
	if normalized, ok := keyMap[lowerKey]; ok {
		return normalized
	}

	// If the key already has underscores, return as-is
	if strings.Contains(key, "_") {
		return strings.ToLower(key)
	}

	// Return the key as-is (might already be correct)
	return strings.ToLower(key)
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
