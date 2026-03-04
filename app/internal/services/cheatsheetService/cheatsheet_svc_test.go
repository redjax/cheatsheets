package cheatsheetservice

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// Helper to create a test cheatsheets directory structure
func createTestCheatsheetsDir(t *testing.T) string {
	t.Helper()

	baseDir := t.TempDir()
	cheatsheetsDir := filepath.Join(baseDir, "cheatsheets")

	// Create type directories
	types := []string{"app", "command", "language", "system"}
	for _, typeDir := range types {
		typePath := filepath.Join(cheatsheetsDir, typeDir)
		if err := os.MkdirAll(typePath, 0755); err != nil {
			t.Fatalf("failed to create type directory: %v", err)
		}
	}

	// Create some test cheatsheets
	testCheatsheets := map[string][]string{
		"app":      {"neovim.md", "helix.md"},
		"command":  {"git.md", "docker.md"},
		"language": {"python.md", "go.md"},
		"system":   {"linux.md"},
	}

	for typeDir, files := range testCheatsheets {
		for _, file := range files {
			filePath := filepath.Join(cheatsheetsDir, typeDir, file)
			content := "# " + strings.TrimSuffix(file, ".md") + "\n\nTest content"
			if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
				t.Fatalf("failed to create cheatsheet: %v", err)
			}
		}
	}

	return baseDir
}

// TestGetCheatsheetsPath tests path detection logic
func TestGetCheatsheetsPath(t *testing.T) {
	tests := []struct {
		name      string
		setupFunc func() string
		wantPath  string // relative to base
	}{
		{
			name: "with cheatsheets subdirectory",
			setupFunc: func() string {
				return createTestCheatsheetsDir(t)
			},
			wantPath: "cheatsheets",
		},
		{
			name: "without cheatsheets subdirectory",
			setupFunc: func() string {
				baseDir := t.TempDir()
				// Create type directories directly in base
				os.MkdirAll(filepath.Join(baseDir, "app"), 0755)
				return baseDir
			},
			wantPath: "",
		},
		{
			name: "empty directory",
			setupFunc: func() string {
				return t.TempDir()
			},
			wantPath: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			basePath := tt.setupFunc()
			got := GetCheatsheetsPath(basePath)

			var expected string
			if tt.wantPath != "" {
				expected = filepath.Join(basePath, tt.wantPath)
			} else {
				expected = basePath
			}

			if got != expected {
				t.Errorf("GetCheatsheetsPath() = %q, want %q", got, expected)
			}
		})
	}
}

// TestValidateCheatsheetsDirectory tests directory validation
func TestValidateCheatsheetsDirectory(t *testing.T) {
	tests := []struct {
		name      string
		setupFunc func() string
		wantError bool
	}{
		{
			name: "valid cheatsheets directory",
			setupFunc: func() string {
				return createTestCheatsheetsDir(t)
			},
			wantError: false,
		},
		{
			name: "nonexistent directory",
			setupFunc: func() string {
				return filepath.Join(t.TempDir(), "nonexistent")
			},
			wantError: true,
		},
		{
			name: "empty directory without subdirectory",
			setupFunc: func() string {
				return t.TempDir()
			},
			wantError: false, // Directory exists, even if empty
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			basePath := tt.setupFunc()
			err := ValidateCheatsheetsDirectory(basePath)

			if (err != nil) != tt.wantError {
				t.Errorf("ValidateCheatsheetsDirectory() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

// TestGetAvailableTypes tests type directory listing
func TestGetAvailableTypes(t *testing.T) {
	tests := []struct {
		name      string
		setupFunc func() string
		want      []string
		wantError bool
	}{
		{
			name: "multiple types",
			setupFunc: func() string {
				return createTestCheatsheetsDir(t)
			},
			want:      []string{"app", "command", "language", "system"},
			wantError: false,
		},
		{
			name: "empty cheatsheets directory",
			setupFunc: func() string {
				baseDir := t.TempDir()
				os.MkdirAll(filepath.Join(baseDir, "cheatsheets"), 0755)
				return baseDir
			},
			want:      []string{},
			wantError: false,
		},
		{
			name: "nonexistent directory",
			setupFunc: func() string {
				return filepath.Join(t.TempDir(), "nonexistent")
			},
			want:      nil,
			wantError: true,
		},
		{
			name: "directory with files (no subdirs)",
			setupFunc: func() string {
				baseDir := t.TempDir()
				cheatsheetsDir := filepath.Join(baseDir, "cheatsheets")
				os.MkdirAll(cheatsheetsDir, 0755)
				// Add files but no subdirectories
				os.WriteFile(filepath.Join(cheatsheetsDir, "test.md"), []byte("test"), 0644)
				return baseDir
			},
			want:      []string{},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			basePath := tt.setupFunc()
			got, err := GetAvailableTypes(basePath)

			if (err != nil) != tt.wantError {
				t.Errorf("GetAvailableTypes() error = %v, wantError %v", err, tt.wantError)
				return
			}

			if !tt.wantError {
				if len(got) != len(tt.want) {
					t.Errorf("GetAvailableTypes() returned %d types, want %d", len(got), len(tt.want))
					t.Errorf("got: %v, want: %v", got, tt.want)
					return
				}

				for i, v := range got {
					if v != tt.want[i] {
						t.Errorf("GetAvailableTypes()[%d] = %q, want %q", i, v, tt.want[i])
					}
				}
			}
		})
	}
}

// TestGetCheatsheetsByType tests listing cheatsheets in a type
func TestGetCheatsheetsByType(t *testing.T) {
	tests := []struct {
		name      string
		setupFunc func() (basePath, typeDir string)
		want      []string
		wantError bool
	}{
		{
			name: "valid type with cheatsheets",
			setupFunc: func() (string, string) {
				return createTestCheatsheetsDir(t), "app"
			},
			want:      []string{"helix", "neovim"},
			wantError: false,
		},
		{
			name: "empty type directory",
			setupFunc: func() (string, string) {
				baseDir := createTestCheatsheetsDir(t)
				emptyType := filepath.Join(baseDir, "cheatsheets", "empty")
				os.MkdirAll(emptyType, 0755)
				return baseDir, "empty"
			},
			want:      []string{},
			wantError: false,
		},
		{
			name: "nonexistent type",
			setupFunc: func() (string, string) {
				return createTestCheatsheetsDir(t), "nonexistent"
			},
			want:      nil,
			wantError: true,
		},
		{
			name: "type with mixed files",
			setupFunc: func() (string, string) {
				baseDir := createTestCheatsheetsDir(t)
				typePath := filepath.Join(baseDir, "cheatsheets", "mixed")
				os.MkdirAll(typePath, 0755)
				// Add .md files
				os.WriteFile(filepath.Join(typePath, "file1.md"), []byte("test"), 0644)
				os.WriteFile(filepath.Join(typePath, "file2.MD"), []byte("test"), 0644)
				// Add non-.md files (should be ignored)
				os.WriteFile(filepath.Join(typePath, "readme.txt"), []byte("test"), 0644)
				os.WriteFile(filepath.Join(typePath, "config.json"), []byte("test"), 0644)
				return baseDir, "mixed"
			},
			want:      []string{"file1", "file2.MD"}, // .MD uppercase not trimmed
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			basePath, typeDir := tt.setupFunc()
			got, err := GetCheatsheetsByType(basePath, typeDir)

			if (err != nil) != tt.wantError {
				t.Errorf("GetCheatsheetsByType() error = %v, wantError %v", err, tt.wantError)
				return
			}

			if !tt.wantError {
				if len(got) != len(tt.want) {
					t.Errorf("GetCheatsheetsByType() returned %d cheatsheets, want %d", len(got), len(tt.want))
					t.Errorf("got: %v, want: %v", got, tt.want)
					return
				}

				for i, v := range got {
					if v != tt.want[i] {
						t.Errorf("GetCheatsheetsByType()[%d] = %q, want %q", i, v, tt.want[i])
					}
				}
			}
		})
	}
}

// TestValidateType tests type validation
func TestValidateType(t *testing.T) {
	tests := []struct {
		name      string
		setupFunc func() (basePath, typeFilter string)
		wantValid bool
		wantTypes []string
		wantError bool
	}{
		{
			name: "valid existing type",
			setupFunc: func() (string, string) {
				return createTestCheatsheetsDir(t), "app"
			},
			wantValid: true,
			wantTypes: []string{"app", "command", "language", "system"},
			wantError: false,
		},
		{
			name: "empty type filter (show all)",
			setupFunc: func() (string, string) {
				return createTestCheatsheetsDir(t), ""
			},
			wantValid: true,
			wantTypes: []string{"app", "command", "language", "system"},
			wantError: false,
		},
		{
			name: "nonexistent type",
			setupFunc: func() (string, string) {
				return createTestCheatsheetsDir(t), "nonexistent"
			},
			wantValid: false,
			wantTypes: []string{"app", "command", "language", "system"},
			wantError: false,
		},
		{
			name: "invalid repo path",
			setupFunc: func() (string, string) {
				return filepath.Join(t.TempDir(), "nonexistent"), "app"
			},
			wantValid: false,
			wantTypes: nil,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			basePath, typeFilter := tt.setupFunc()
			valid, types, err := ValidateType(basePath, typeFilter)

			if (err != nil) != tt.wantError {
				t.Errorf("ValidateType() error = %v, wantError %v", err, tt.wantError)
				return
			}

			if valid != tt.wantValid {
				t.Errorf("ValidateType() valid = %v, want %v", valid, tt.wantValid)
			}

			if !tt.wantError {
				if len(types) != len(tt.wantTypes) {
					t.Errorf("ValidateType() returned %d types, want %d", len(types), len(tt.wantTypes))
				}
			}
		})
	}
}

// TestGetCheatsheetPath tests path construction
func TestGetCheatsheetPath(t *testing.T) {
	tests := []struct {
		name       string
		setupFunc  func() (repoPath, typeDir, cheatName string)
		wantSuffix string
	}{
		{
			name: "standard path",
			setupFunc: func() (string, string, string) {
				return createTestCheatsheetsDir(t), "app", "neovim"
			},
			wantSuffix: filepath.Join("cheatsheets", "app", "neovim.md"),
		},
		{
			name: "path without extension",
			setupFunc: func() (string, string, string) {
				return createTestCheatsheetsDir(t), "command", "git"
			},
			wantSuffix: filepath.Join("cheatsheets", "command", "git.md"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoPath, typeDir, cheatName := tt.setupFunc()
			got := GetCheatsheetPath(repoPath, typeDir, cheatName)

			if !strings.HasSuffix(got, tt.wantSuffix) {
				t.Errorf("GetCheatsheetPath() = %q, should end with %q", got, tt.wantSuffix)
			}

			// Verify it's an absolute path
			if !filepath.IsAbs(got) {
				t.Errorf("GetCheatsheetPath() = %q, should be absolute path", got)
			}
		})
	}
}

// TestStripFrontmatter tests YAML frontmatter removal
func TestStripFrontmatter(t *testing.T) {
	tests := []struct {
		name    string
		content string
		want    string
	}{
		{
			name: "with frontmatter",
			content: `---
title: Test
tags: [test, example]
---

# Heading

Content`,
			want: `
# Heading

Content`, // Includes newline after closing delimiter
		},
		{
			name: "without frontmatter",
			content: `# Heading

Content`,
			want: `# Heading

Content`,
		},
		{
			name:    "empty content",
			content: "",
			want:    "",
		},
		{
			name: "only frontmatter",
			content: `---
title: Test
---`,
			want: "",
		},
		{
			name: "frontmatter with content on same line",
			content: `---
title: Test
---
# Immediate Content`,
			want: "# Immediate Content",
		},
		{
			name: "malformed frontmatter (no closing)",
			content: `---
title: Test

# Heading`,
			want: `---
title: Test

# Heading`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := stripFrontmatter(tt.content)

			if got != tt.want {
				t.Errorf("stripFrontmatter() = %q, want %q", got, tt.want)
			}
		})
	}
}

// Helper to create test templates
func createTestTemplates(t *testing.T, basePath string) {
	t.Helper()

	templateDir := filepath.Join(basePath, ".templates")
	if err := os.MkdirAll(templateDir, 0755); err != nil {
		t.Fatalf("failed to create template directory: %v", err)
	}

	// Create template content
	templateContent := `---
title: {{title}}
description: {{description}}
tags: [{{tags}}]
last_updated: {{last_updated}}
---

# {{title}}

{{description}}
`

	// Create templates for each type
	types := []string{"app", "command", "language", "system"}
	for _, typeDir := range types {
		templatePath := filepath.Join(templateDir, typeDir+".md")
		if err := os.WriteFile(templatePath, []byte(templateContent), 0644); err != nil {
			t.Fatalf("failed to create template: %v", err)
		}
	}
}

// TestCreateCheatsheet tests cheatsheet file creation
func TestCreateCheatsheet(t *testing.T) {
	tests := []struct {
		name      string
		setupFunc func() (repoPath string, opts CreateCheatsheetOptions)
		wantError bool
		validate  func(*testing.T, string)
	}{
		{
			name: "create new cheatsheet",
			setupFunc: func() (string, CreateCheatsheetOptions) {
				basePath := createTestCheatsheetsDir(t)
				createTestTemplates(t, basePath)
				opts := CreateCheatsheetOptions{
					Type:        "app",
					Name:        "vscode",
					Title:       "vscode",
					Description: "VS Code shortcuts",
					Tags:        "editor, ide",
				}
				return basePath, opts
			},
			wantError: false,
			validate: func(t *testing.T, filePath string) {
				if _, err := os.Stat(filePath); os.IsNotExist(err) {
					t.Errorf("file not created: %s", filePath)
				}

				content, err := os.ReadFile(filePath)
				if err != nil {
					t.Fatalf("failed to read created file: %v", err)
				}

				contentStr := string(content)
				if !strings.Contains(contentStr, "title: vscode") {
					t.Error("file should contain title in frontmatter")
				}
				if !strings.Contains(contentStr, "VS Code shortcuts") {
					t.Error("file should contain description")
				}
			},
		},
		{
			name: "create in nonexistent type (should fail)",
			setupFunc: func() (string, CreateCheatsheetOptions) {
				basePath := createTestCheatsheetsDir(t)
				createTestTemplates(t, basePath)
				opts := CreateCheatsheetOptions{
					Type:        "newtype",
					Name:        "test",
					Title:       "test",
					Description: "Test cheatsheet",
				}
				return basePath, opts
			},
			wantError: true, // Invalid type should fail
			validate:  nil,
		},
		{
			name: "file already exists",
			setupFunc: func() (string, CreateCheatsheetOptions) {
				basePath := createTestCheatsheetsDir(t)
				createTestTemplates(t, basePath)
				opts := CreateCheatsheetOptions{
					Type:        "app",
					Name:        "neovim", // Already exists
					Title:       "neovim",
					Description: "Test",
				}
				return basePath, opts
			},
			wantError: true,
			validate:  nil,
		},
		{
			name: "missing template",
			setupFunc: func() (string, CreateCheatsheetOptions) {
				basePath := createTestCheatsheetsDir(t)
				// Don't create templates
				opts := CreateCheatsheetOptions{
					Type:        "app",
					Name:        "test",
					Title:       "test",
					Description: "Test",
				}
				return basePath, opts
			},
			wantError: true,
			validate:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoPath, opts := tt.setupFunc()
			filePath, err := CreateCheatsheet(repoPath, opts)

			if (err != nil) != tt.wantError {
				t.Errorf("CreateCheatsheet() error = %v, wantError %v", err, tt.wantError)
				return
			}

			if !tt.wantError && tt.validate != nil {
				tt.validate(t, filePath)
			}
		})
	}
}

// TestDeleteCheatsheet tests cheatsheet deletion
func TestDeleteCheatsheet(t *testing.T) {
	tests := []struct {
		name      string
		setupFunc func() (repoPath, typeDir, cheatName string)
		wantError bool
	}{
		{
			name: "delete existing cheatsheet",
			setupFunc: func() (string, string, string) {
				return createTestCheatsheetsDir(t), "app", "neovim"
			},
			wantError: false,
		},
		{
			name: "delete nonexistent cheatsheet",
			setupFunc: func() (string, string, string) {
				return createTestCheatsheetsDir(t), "app", "nonexistent"
			},
			wantError: true,
		},
		{
			name: "delete from nonexistent type",
			setupFunc: func() (string, string, string) {
				return createTestCheatsheetsDir(t), "nonexistent", "test"
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoPath, typeDir, cheatName := tt.setupFunc()

			// Get the file path before deletion
			filePath := GetCheatsheetPath(repoPath, typeDir, cheatName)
			existedBefore := false
			if _, err := os.Stat(filePath); err == nil {
				existedBefore = true
			}

			err := DeleteCheatsheet(repoPath, typeDir, cheatName)

			if (err != nil) != tt.wantError {
				t.Errorf("DeleteCheatsheet() error = %v, wantError %v", err, tt.wantError)
				return
			}

			// If no error expected and file existed, verify it's deleted
			if !tt.wantError && existedBefore {
				if _, err := os.Stat(filePath); err == nil {
					t.Errorf("file still exists after deletion: %s", filePath)
				}
			}
		})
	}
}
