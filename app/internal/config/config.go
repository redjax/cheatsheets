package config

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/adrg/xdg"
	"github.com/knadh/koanf/parsers/dotenv"
	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/posflag"
	"github.com/knadh/koanf/v2"
	"github.com/redjax/cheatsheets/internal/constants"
	"github.com/redjax/cheatsheets/internal/utils"
	"github.com/spf13/pflag"
	"golang.org/x/term"
	yamlv3 "gopkg.in/yaml.v3"
)

const envPrefix = "CHEATSHEETS_"

var K = koanf.New(".")

// Config represents the application configuration
type Config struct {
	SheetsPath string    `koanf:"sheets_path" path:"expand"`
	Git        GitConfig `koanf:"git"`
	Debug      bool      `koanf:"debug"`
}

// GitConfig struct for app-level configuration
type GitConfig struct {
	RepoUrl       string `koanf:"repo_url"`
	ClonePath     string `koanf:"clone_path" path:"expand"`
	Token         string `koanf:"token"`
	AutoBranch    bool   `koanf:"auto_branch"`
	WorkingBranch string `koanf:"working_branch"`
	AuthorName    string `koanf:"author_name"`
	AuthorEmail   string `koanf:"author_email"`
}

// String implements the Stringer interface to mask sensitive fields when printed
func (g GitConfig) String() string {
	token := "<empty>"

	if g.Token != "" {
		token = utils.MaskToken(g.Token)
	}

	return fmt.Sprintf("GitConfig{RepoUrl: %s, ClonePath: %s, Token: %s}", g.RepoUrl, g.ClonePath, token)
}

// GetDefaultConfigPath returns the default config file path in XDG config directory
func GetDefaultConfigPath() string {
	return filepath.Join(xdg.ConfigHome, constants.AppDataDirName, "config.yml")
}

// FindConfigFile checks for a .local variant of the config file first,
// falling back to the original if .local doesn't exist.
// If configFile is empty, returns the default XDG config path.
func FindConfigFile(configFile string) string {
	if configFile == "" {
		// Check XDG config location (~/.config/cheatsheets/config.yml)
		defaultPath := GetDefaultConfigPath()
		localPath := strings.TrimSuffix(defaultPath, ".yml") + ".local.yml"

		// Prefer .local variant
		if _, err := os.Stat(localPath); err == nil {
			return localPath
		}

		// Fall back to default (may or may not exist yet)
		return defaultPath
	}

	// Check for .local variant (e.g., config.yml -> config.local.yml)
	ext := filepath.Ext(configFile)
	base := strings.TrimSuffix(configFile, ext)
	localFile := base + ".local" + ext

	if _, err := os.Stat(localFile); err == nil {
		return localFile
	}

	return configFile
}

// getDefaultClonePath returns the default clone path using XDG directories
func getDefaultClonePath() string {
	return filepath.Join(xdg.DataHome, constants.AppDataDirName)
}

// getDefaultSheetsPath returns the default sheets path (user's home directory)
func getDefaultSheetsPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return "~/.cheatsheets"
	}
	return filepath.Join(home, ".cheatsheets")
}

// LoadConfig loads configuration from file, environment variables, and CLI flags
// Returns the parsed config struct or an error
func LoadConfig(flagSet *pflag.FlagSet, configFile string) (*Config, error) {
	// Create default config file if it doesn't exist
	if configFile != "" {
		if _, err := os.Stat(configFile); os.IsNotExist(err) {
			if err := ensureConfigFile(configFile); err != nil {
				return nil, fmt.Errorf("failed to create config file: %w", err)
			}
		}
	}

	if configFile != "" {
		parser, err := parserForFile(configFile)
		if err != nil {
			return nil, fmt.Errorf("unsupported config file format: %w", err)
		}

		if err := K.Load(file.Provider(configFile), parser); err != nil {
			return nil, fmt.Errorf("error loading config file: %w", err)
		}
	}

	// Load from env var
	if err := K.Load(env.Provider(envPrefix, ".", func(s string) string {
		return strings.Replace(strings.ToLower(
			strings.TrimPrefix(s, envPrefix)), "_", ".", -1)
	}), nil); err != nil {
		return nil, fmt.Errorf("error loading env vars: %w", err)
	}

	// Load from CLI flags
	if flagSet != nil {
		if err := K.Load(posflag.Provider(flagSet, ".", K), nil); err != nil {
			return nil, fmt.Errorf("error loading flags: %w", err)
		}
	}

	// Unmarshal into config struct
	var cfg Config
	if err := K.Unmarshal("", &cfg); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	// Set defaults for empty values
	if cfg.Git.RepoUrl == "" {
		cfg.Git.RepoUrl = constants.RepoURL
	}
	if cfg.Git.ClonePath == "" {
		cfg.Git.ClonePath = getDefaultClonePath()
	}
	if cfg.SheetsPath == "" {
		cfg.SheetsPath = getDefaultSheetsPath()
	}

	// Expand paths
	cfg.expandPaths()

	return &cfg, nil
}

func parserForFile(path string) (koanf.Parser, error) {
	ext := strings.ToLower(filepath.Ext(path))

	switch ext {
	case ".yaml", ".yml":
		return yaml.Parser(), nil
	case ".json":
		return json.Parser(), nil
	case ".toml":
		return toml.Parser(), nil
	case ".env":
		return dotenv.Parser(), nil
	default:
		return nil, fmt.Errorf("unsupported config file format: %s", ext)
	}
}

// expandPaths walks the config struct and expands ~ in any field tagged with path:"expand"
func (c *Config) expandPaths() {
	expandStructPaths(reflect.ValueOf(c).Elem())
}

// expandStructPaths recursively walks a struct and expands paths in tagged fields
func expandStructPaths(v reflect.Value) {
	if v.Kind() != reflect.Struct {
		return
	}

	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)

		// Check if field has path:"expand" tag
		if tag := fieldType.Tag.Get("path"); tag == "expand" {
			if field.Kind() == reflect.String && field.CanSet() {
				field.SetString(expandPath(field.String()))
			}
		}

		// Recursively handle nested structs
		if field.Kind() == reflect.Struct {
			expandStructPaths(field)
		}
	}
}

// expandPath returns the expanded path, handling ~ for home directory and converting to absolute path
func expandPath(path string) string {
	// Handle ~ expansion
	if strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err == nil {
			path = filepath.Join(home, path[2:])
		}
	}

	// Convert to absolute path
	absPath, err := filepath.Abs(path)
	if err == nil {
		return absPath
	}

	return path // Return original if expansion fails
}

// ensureConfigFile creates the config file if it doesn't exist
func ensureConfigFile(configFile string) error {
	// Create config directory if it doesn't exist
	configDir := filepath.Dir(configFile)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %w", err)
	}

	fmt.Println("=== First Time Setup ===")
	fmt.Printf("Creating config file: %s\n\n", configFile)

	// Prompt for GitHub token if not set via env var
	token := os.Getenv("CHEATSHEETS_GIT_TOKEN")
	if token == "" {
		token = promptForToken()
	}

	// Generate default config with env vars and prompted token
	configData := createDefaultConfigWithEnvVars(token)

	// Marshal to YAML
	data, err := yamlv3.Marshal(configData)
	if err != nil {
		return fmt.Errorf("failed to marshal default config: %w", err)
	}

	// Write to file
	if err := os.WriteFile(configFile, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	fmt.Printf("\n✓ Config file created: %s\n", configFile)
	return nil
}

// createDefaultConfigWithEnvVars creates a default config map with values from env vars if available
func createDefaultConfigWithEnvVars(token string) map[string]interface{} {
	config := map[string]interface{}{
		"sheets_path": getEnvOrDefault("CHEATSHEETS_SHEETS_PATH", "~/.cheatsheets"),
		"debug":       getEnvOrDefault("CHEATSHEETS_DEBUG", "true") == "true",
		"git": map[string]interface{}{
			"repo_url":       getEnvOrDefault("CHEATSHEETS_GIT_REPO_URL", constants.RepoURL),
			"clone_path":     getEnvOrDefault("CHEATSHEETS_GIT_CLONE_PATH", ""),
			"token":          token,
			"auto_branch":    getEnvOrDefault("CHEATSHEETS_GIT_AUTO_BRANCH", "true") == "true",
			"working_branch": getEnvOrDefault("CHEATSHEETS_GIT_WORKING_BRANCH", "working"),
			"author_name":    getEnvOrDefault("CHEATSHEETS_GIT_AUTHOR_NAME", ""),
			"author_email":   getEnvOrDefault("CHEATSHEETS_GIT_AUTHOR_EMAIL", ""),
		},
	}

	return config
}

// getEnvOrDefault gets an environment variable or returns the default value
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// promptForToken prompts the user to enter a GitHub Personal Access Token
func promptForToken() string {
	fmt.Println("GitHub Personal Access Token (PAT):")
	fmt.Println("A PAT is required to clone private repositories and make commits.")
	fmt.Println("Create one at: https://github.com/settings/tokens")
	fmt.Println("Required scopes: repo (for private repos)")
	fmt.Println()
	fmt.Print("Enter your GitHub PAT (input hidden, or press Enter to skip): ")

	// Read password without echoing to terminal
	tokenBytes, err := term.ReadPassword(int(os.Stdin.Fd()))
	fmt.Println() // Add newline after hidden input

	if err != nil {
		return ""
	}

	token := strings.TrimSpace(string(tokenBytes))
	return token
}
