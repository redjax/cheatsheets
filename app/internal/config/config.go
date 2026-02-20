package config

import (
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"github.com/knadh/koanf/parsers/dotenv"
	"github.com/knadh/koanf/parsers/json"
	"github.com/knadh/koanf/parsers/toml"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/env"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/providers/posflag"
	"github.com/knadh/koanf/v2"
	"github.com/spf13/pflag"
)

const envPrefix = "CHEATSHEETS_"

var K = koanf.New(".")

// Config represents the application configuration
type Config struct {
	Git   GitConfig `koanf:"git"`
	Debug bool      `koanf:"debug"`
}

// GitConfig struct for app-level configuration
type GitConfig struct {
	ClonePath string `koanf:"clone_path" path:"expand"`
	Token     string `koanf:"token"`
}

// FindConfigFile checks for a .local variant of the config file first,
// falling back to the original if .local doesn't exist
func FindConfigFile(configFile string) string {
	if configFile == "" {
		return ""
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

// LoadConfig loads configuration from file, environment variables, and CLI flags
// Returns the parsed config struct or an error
func LoadConfig(flagSet *pflag.FlagSet, configFile string) (*Config, error) {
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
