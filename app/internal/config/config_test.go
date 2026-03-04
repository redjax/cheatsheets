package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/knadh/koanf/v2"
	"github.com/redjax/cheatsheets/internal/constants"
)

// Helper to create a test config file
func createTestConfigFile(t *testing.T, content string) string {
	t.Helper()
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yml")

	if err := os.WriteFile(configPath, []byte(content), 0644); err != nil {
		t.Fatalf("failed to create test config: %v", err)
	}

	return configPath
}

// TestFindConfigFile tests config file discovery logic
func TestFindConfigFile(t *testing.T) {
	tests := []struct {
		name          string
		setupFunc     func() (string, string) // Returns (configFile, expectedFile)
		wantLocalFile bool
	}{
		{
			name: "prefers local variant when it exists",
			setupFunc: func() (string, string) {
				tempDir := t.TempDir()
				configPath := filepath.Join(tempDir, "config.yml")
				localPath := filepath.Join(tempDir, "config.local.yml")

				// Create both files
				os.WriteFile(configPath, []byte("original"), 0644)
				os.WriteFile(localPath, []byte("local"), 0644)

				return configPath, localPath
			},
			wantLocalFile: true,
		},
		{
			name: "uses original when local doesn't exist",
			setupFunc: func() (string, string) {
				tempDir := t.TempDir()
				configPath := filepath.Join(tempDir, "config.yml")
				os.WriteFile(configPath, []byte("original"), 0644)

				return configPath, configPath
			},
			wantLocalFile: false,
		},
		{
			name: "handles empty string (default path)",
			setupFunc: func() (string, string) {
				expected := GetDefaultConfigPath()
				return "", expected
			},
			wantLocalFile: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			configFile, expectedFile := tt.setupFunc()
			result := FindConfigFile(configFile)

			if tt.wantLocalFile {
				if !strings.Contains(result, ".local.") {
					t.Errorf("FindConfigFile() = %q, should contain '.local.'", result)
				}
			}

			if configFile != "" && result != expectedFile {
				// Only check exact match when we created specific files
				if _, err := os.Stat(expectedFile); err == nil {
					if result != expectedFile {
						t.Errorf("FindConfigFile() = %q, want %q", result, expectedFile)
					}
				}
			}
		})
	}
}

// TestExpandPath tests path expansion logic
func TestExpandPath(t *testing.T) {
	home, err := os.UserHomeDir()
	if err != nil {
		t.Skip("cannot get home directory")
	}

	tests := []struct {
		name      string
		input     string
		wantStart string // Check if result starts with this
	}{
		{
			name:      "expands tilde",
			input:     "~/test/path",
			wantStart: home,
		},
		{
			name:      "handles absolute path",
			input:     "/absolute/path",
			wantStart: "/",
		},
		{
			name:      "handles relative path",
			input:     "relative/path",
			wantStart: "", // Will be absolute after expansion
		},
		{
			name:      "handles empty string",
			input:     "",
			wantStart: "", // Empty string becomes current dir after filepath.Abs
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := expandPath(tt.input)

			if tt.input == "" {
				// Empty string gets converted to current dir by filepath.Abs
				if !filepath.IsAbs(result) {
					t.Errorf("expandPath(%q) = %q, should be absolute path", tt.input, result)
				}
				return
			}

			// Check if result is absolute (after expansion)
			if !filepath.IsAbs(result) && tt.input != "" {
				t.Errorf("expandPath(%q) = %q, should be absolute path", tt.input, result)
			}

			// Check starting path for tilde expansion
			if tt.wantStart != "" && !strings.HasPrefix(result, tt.wantStart) {
				t.Errorf("expandPath(%q) = %q, should start with %q", tt.input, result, tt.wantStart)
			}
		})
	}
}

// TestParserForFile tests file extension parser detection
func TestParserForFile(t *testing.T) {
	tests := []struct {
		name      string
		filename  string
		wantError bool
	}{
		{
			name:      "yaml extension",
			filename:  "config.yaml",
			wantError: false,
		},
		{
			name:      "yml extension",
			filename:  "config.yml",
			wantError: false,
		},
		{
			name:      "json extension",
			filename:  "config.json",
			wantError: false,
		},
		{
			name:      "toml extension",
			filename:  "config.toml",
			wantError: false,
		},
		{
			name:      "env extension",
			filename:  ".env",
			wantError: false,
		},
		{
			name:      "uppercase extension",
			filename:  "config.YAML",
			wantError: false,
		},
		{
			name:      "unsupported extension",
			filename:  "config.txt",
			wantError: true,
		},
		{
			name:      "no extension",
			filename:  "config",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parser, err := parserForFile(tt.filename)

			if tt.wantError {
				if err == nil {
					t.Errorf("parserForFile(%q) expected error, got nil", tt.filename)
				}
				if parser != nil {
					t.Errorf("parserForFile(%q) expected nil parser on error", tt.filename)
				}
			} else {
				if err != nil {
					t.Errorf("parserForFile(%q) unexpected error: %v", tt.filename, err)
				}
				if parser == nil {
					t.Errorf("parserForFile(%q) expected parser, got nil", tt.filename)
				}
			}
		})
	}
}

// TestLoadConfig tests configuration loading
func TestLoadConfig(t *testing.T) {
	tests := []struct {
		name      string
		content   string
		wantError bool
		validate  func(*testing.T, *Config)
	}{
		{
			name: "valid minimal config",
			content: `
git:
  repo_url: https://github.com/test/repo.git
  clone_path: /tmp/test
  token: test_token
`,
			wantError: false,
			validate: func(t *testing.T, cfg *Config) {
				if cfg.Git.RepoUrl != "https://github.com/test/repo.git" {
					t.Errorf("RepoUrl = %q, want %q", cfg.Git.RepoUrl, "https://github.com/test/repo.git")
				}
				if cfg.Git.Token != "test_token" {
					t.Errorf("Token = %q, want %q", cfg.Git.Token, "test_token")
				}
			},
		},
		{
			name: "config with defaults applied",
			content: `
debug: true
`,
			wantError: false,
			validate: func(t *testing.T, cfg *Config) {
				// Should apply default repo URL
				if cfg.Git.RepoUrl != constants.RepoURL {
					t.Errorf("RepoUrl = %q, want default %q", cfg.Git.RepoUrl, constants.RepoURL)
				}
				// Should apply default clone path
				if cfg.Git.ClonePath == "" {
					t.Error("ClonePath should not be empty after defaults")
				}
			},
		},
		{
			name: "empty config gets all defaults",
			content: `
`,
			wantError: false,
			validate: func(t *testing.T, cfg *Config) {
				if cfg.Git.RepoUrl == "" {
					t.Error("RepoUrl should have default value")
				}
				if cfg.Git.ClonePath == "" {
					t.Error("ClonePath should have default value")
				}
				if cfg.SheetsPath == "" {
					t.Error("SheetsPath should have default value")
				}
			},
		},
		{
			name: "invalid YAML",
			content: `
this is not: valid: yaml: structure
  bad indentation
`,
			wantError: true,
			validate:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset koanf instance for each test
			K = koanf.New(".")

			configPath := createTestConfigFile(t, tt.content)
			cfg, err := LoadConfig(nil, configPath)

			if tt.wantError {
				if err == nil {
					t.Error("LoadConfig() expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("LoadConfig() unexpected error: %v", err)
			}

			if cfg == nil {
				t.Fatal("LoadConfig() returned nil config")
			}

			if tt.validate != nil {
				tt.validate(t, cfg)
			}
		})
	}
}

// TestLoadConfigWithEnvVars tests environment variable overrides
func TestLoadConfigWithEnvVars(t *testing.T) {
	// Reset koanf instance
	K = koanf.New(".")

	// Create minimal config
	configPath := createTestConfigFile(t, `
git:
  repo_url: https://github.com/original/repo.git
  token: original_token
`)

	// Set env var to override
	t.Setenv("CHEATSHEETS_GIT_TOKEN", "env_token")

	cfg, err := LoadConfig(nil, configPath)
	if err != nil {
		t.Fatalf("LoadConfig() error: %v", err)
	}

	// Env var should override config file
	if cfg.Git.Token != "env_token" {
		t.Errorf("Token = %q, want %q (from env)", cfg.Git.Token, "env_token")
	}

	// Original value should still be present
	if cfg.Git.RepoUrl != "https://github.com/original/repo.git" {
		t.Errorf("RepoUrl = %q, want original value", cfg.Git.RepoUrl)
	}
}

// TestGitConfigString tests token masking in string representation
func TestGitConfigString(t *testing.T) {
	tests := []struct {
		name      string
		config    GitConfig
		wantMask  bool
		wantEmpty bool
	}{
		{
			name: "masks non-empty token",
			config: GitConfig{
				RepoUrl:   "https://github.com/test/repo.git",
				ClonePath: "/tmp/test",
				Token:     "ghp_1234567890abcdefghijklmnopqrstuvwxyz",
			},
			wantMask:  true,
			wantEmpty: false,
		},
		{
			name: "shows empty for no token",
			config: GitConfig{
				RepoUrl:   "https://github.com/test/repo.git",
				ClonePath: "/tmp/test",
				Token:     "",
			},
			wantMask:  false,
			wantEmpty: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.config.String()

			if tt.wantMask {
				if !strings.Contains(result, "*") {
					t.Errorf("String() = %q, should contain masked token", result)
				}
				if strings.Contains(result, tt.config.Token) {
					t.Errorf("String() = %q, should not contain raw token", result)
				}
			}

			if tt.wantEmpty {
				if !strings.Contains(result, "<empty>") {
					t.Errorf("String() = %q, should show <empty>", result)
				}
			}

			// Should always contain RepoUrl and ClonePath
			if !strings.Contains(result, tt.config.RepoUrl) {
				t.Errorf("String() = %q, should contain RepoUrl", result)
			}
			if !strings.Contains(result, tt.config.ClonePath) {
				t.Errorf("String() = %q, should contain ClonePath", result)
			}
		})
	}
}

// TestGetDefaultConfigPath tests default config path generation
func TestGetDefaultConfigPath(t *testing.T) {
	path := GetDefaultConfigPath()

	if path == "" {
		t.Error("GetDefaultConfigPath() returned empty string")
	}

	if !strings.Contains(path, constants.AppDataDirName) {
		t.Errorf("GetDefaultConfigPath() = %q, should contain %q", path, constants.AppDataDirName)
	}

	if !strings.HasSuffix(path, "config.yml") {
		t.Errorf("GetDefaultConfigPath() = %q, should end with 'config.yml'", path)
	}

	// Should be an absolute path
	if !filepath.IsAbs(path) {
		t.Errorf("GetDefaultConfigPath() = %q, should be absolute path", path)
	}
}

// TestGetEnvOrDefault tests environment variable fallback logic
func TestGetEnvOrDefault(t *testing.T) {
	tests := []struct {
		name         string
		envKey       string
		envValue     string
		defaultValue string
		want         string
	}{
		{
			name:         "returns env value when set",
			envKey:       "TEST_ENV_VAR",
			envValue:     "env_value",
			defaultValue: "default_value",
			want:         "env_value",
		},
		{
			name:         "returns default when env not set",
			envKey:       "TEST_NONEXISTENT_VAR",
			envValue:     "",
			defaultValue: "default_value",
			want:         "default_value",
		},
		{
			name:         "returns empty string if that's the env value",
			envKey:       "TEST_EMPTY_VAR",
			envValue:     "",
			defaultValue: "default_value",
			want:         "default_value",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.envValue != "" {
				t.Setenv(tt.envKey, tt.envValue)
			}

			result := getEnvOrDefault(tt.envKey, tt.defaultValue)

			if result != tt.want {
				t.Errorf("getEnvOrDefault(%q, %q) = %q, want %q",
					tt.envKey, tt.defaultValue, result, tt.want)
			}
		})
	}
}

// TestConfigPathExpansion tests that paths are expanded after loading
func TestConfigPathExpansion(t *testing.T) {
	// Reset koanf instance
	K = koanf.New(".")

	configPath := createTestConfigFile(t, `
sheets_path: ~/test/sheets
git:
  clone_path: ~/test/clone
`)

	cfg, err := LoadConfig(nil, configPath)
	if err != nil {
		t.Fatalf("LoadConfig() error: %v", err)
	}

	// Paths should be expanded (no ~ in final result)
	if strings.Contains(cfg.SheetsPath, "~") {
		t.Errorf("SheetsPath = %q, should be expanded (no ~)", cfg.SheetsPath)
	}

	if strings.Contains(cfg.Git.ClonePath, "~") {
		t.Errorf("ClonePath = %q, should be expanded (no ~)", cfg.Git.ClonePath)
	}

	// Should be absolute paths
	if !filepath.IsAbs(cfg.SheetsPath) {
		t.Errorf("SheetsPath = %q, should be absolute", cfg.SheetsPath)
	}

	if !filepath.IsAbs(cfg.Git.ClonePath) {
		t.Errorf("ClonePath = %q, should be absolute", cfg.Git.ClonePath)
	}
}

// TestLoadConfigNonexistentFile tests creating config file on first run
func TestLoadConfigNonexistentFile(t *testing.T) {
	// Reset koanf instance
	K = koanf.New(".")

	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yml")

	// Set token via env to avoid interactive prompt
	t.Setenv("CHEATSHEETS_GIT_TOKEN", "test_token")

	// File doesn't exist, should be created
	cfg, err := LoadConfig(nil, configPath)
	if err != nil {
		t.Fatalf("LoadConfig() should create missing config, got error: %v", err)
	}

	if cfg == nil {
		t.Fatal("LoadConfig() returned nil config")
	}

	// Verify file was created
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		t.Errorf("LoadConfig() should have created config file at %q", configPath)
	}

	// Verify token from env was used
	if cfg.Git.Token != "test_token" {
		t.Errorf("Token = %q, want %q (from env)", cfg.Git.Token, "test_token")
	}

	// Should have defaults
	if cfg.Git.RepoUrl == "" {
		t.Error("RepoUrl should have default value")
	}
}

// TestCreateDefaultConfigWithEnvVars tests default config generation
func TestCreateDefaultConfigWithEnvVars(t *testing.T) {
	t.Setenv("CHEATSHEETS_GIT_REPO_URL", "https://custom.repo.url")
	t.Setenv("CHEATSHEETS_DEBUG", "false")

	config := createDefaultConfigWithEnvVars("test_token")

	// Check token
	gitConfig := config["git"].(map[string]interface{})
	if gitConfig["token"] != "test_token" {
		t.Errorf("token = %q, want %q", gitConfig["token"], "test_token")
	}

	// Check env override
	if gitConfig["repo_url"] != "https://custom.repo.url" {
		t.Errorf("repo_url = %q, want env value", gitConfig["repo_url"])
	}

	// Check debug flag
	if config["debug"] != false {
		t.Errorf("debug = %v, want false (from env)", config["debug"])
	}
}
