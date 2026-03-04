package syncservice

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	cheatsheetservice "github.com/redjax/cheatsheets/internal/services/cheatsheetService"
)

// Helper to create a test cheatsheets directory structure
func createTestCheatsheetsDir(t *testing.T) string {
	t.Helper()

	baseDir := t.TempDir()
	cheatsheetsDir := filepath.Join(baseDir, "cheatsheets")

	// Create type directories
	types := []string{"app", "command"}
	for _, typeDir := range types {
		typePath := filepath.Join(cheatsheetsDir, typeDir)
		if err := os.MkdirAll(typePath, 0755); err != nil {
			t.Fatalf("failed to create type directory: %v", err)
		}
	}

	// Create some test cheatsheets
	testFiles := map[string]string{
		"app/test.md":    "# Test App\n\nContent",
		"command/git.md": "# Git\n\nGit commands",
	}

	for relPath, content := range testFiles {
		filePath := filepath.Join(cheatsheetsDir, relPath)
		if err := os.WriteFile(filePath, []byte(content), 0644); err != nil {
			t.Fatalf("failed to create test file: %v", err)
		}
	}

	return baseDir
}

// TestCheckSyncStatus tests sync status detection
func TestCheckSyncStatus(t *testing.T) {
	tests := []struct {
		name      string
		setupFunc func() (repoPath, sheetsPath string)
		wantType  SyncType
		wantValid bool
		wantError bool
	}{
		{
			name: "nonexistent destination",
			setupFunc: func() (string, string) {
				repoPath := createTestCheatsheetsDir(t)
				sheetsPath := filepath.Join(t.TempDir(), "sheets")
				return repoPath, sheetsPath
			},
			wantType:  SyncTypeNone,
			wantValid: false,
			wantError: false,
		},
		{
			name: "valid symlink",
			setupFunc: func() (string, string) {
				if runtime.GOOS == "windows" {
					t.Skip("Symlink test skipped on Windows")
				}
				repoPath := createTestCheatsheetsDir(t)
				sheetsPath := filepath.Join(t.TempDir(), "sheets")
				sourcePath := cheatsheetservice.GetCheatsheetsPath(repoPath)
				if err := os.Symlink(sourcePath, sheetsPath); err != nil {
					t.Fatalf("failed to create symlink: %v", err)
				}
				return repoPath, sheetsPath
			},
			wantType:  SyncTypeSymlink,
			wantValid: true,
			wantError: false,
		},
		{
			name: "invalid symlink (wrong target)",
			setupFunc: func() (string, string) {
				if runtime.GOOS == "windows" {
					t.Skip("Symlink test skipped on Windows")
				}
				repoPath := createTestCheatsheetsDir(t)
				sheetsPath := filepath.Join(t.TempDir(), "sheets")
				wrongPath := filepath.Join(t.TempDir(), "wrong")
				os.MkdirAll(wrongPath, 0755)
				if err := os.Symlink(wrongPath, sheetsPath); err != nil {
					t.Fatalf("failed to create symlink: %v", err)
				}
				return repoPath, sheetsPath
			},
			wantType:  SyncTypeSymlink,
			wantValid: false,
			wantError: false,
		},
		{
			name: "copied directory",
			setupFunc: func() (string, string) {
				repoPath := createTestCheatsheetsDir(t)
				sheetsPath := filepath.Join(t.TempDir(), "sheets")
				sourcePath := cheatsheetservice.GetCheatsheetsPath(repoPath)
				if err := CopyDirectory(sourcePath, sheetsPath); err != nil {
					t.Fatalf("failed to copy directory: %v", err)
				}
				return repoPath, sheetsPath
			},
			wantType:  SyncTypeCopy,
			wantValid: true,
			wantError: false,
		},
		{
			name: "file instead of directory (mismatch)",
			setupFunc: func() (string, string) {
				repoPath := createTestCheatsheetsDir(t)
				sheetsPath := filepath.Join(t.TempDir(), "sheets")
				if err := os.WriteFile(sheetsPath, []byte("test"), 0644); err != nil {
					t.Fatalf("failed to create file: %v", err)
				}
				return repoPath, sheetsPath
			},
			wantType:  SyncTypeMismatch,
			wantValid: false,
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoPath, sheetsPath := tt.setupFunc()
			status := CheckSyncStatus(repoPath, sheetsPath)

			if (status.Error != nil) != tt.wantError {
				t.Errorf("CheckSyncStatus() error = %v, wantError %v", status.Error, tt.wantError)
			}

			if status.Type != tt.wantType {
				t.Errorf("CheckSyncStatus() Type = %v, want %v", status.Type, tt.wantType)
			}

			if status.IsValid != tt.wantValid {
				t.Errorf("CheckSyncStatus() IsValid = %v, want %v", status.IsValid, tt.wantValid)
			}

			// Check Exists flag
			expectedExists := tt.wantType != SyncTypeNone
			if status.Exists != expectedExists {
				t.Errorf("CheckSyncStatus() Exists = %v, want %v", status.Exists, expectedExists)
			}
		})
	}
}

// TestCopyFile tests single file copying
func TestCopyFile(t *testing.T) {
	tests := []struct {
		name      string
		setupFunc func() (src, dst string)
		wantError bool
		validate  func(*testing.T, string)
	}{
		{
			name: "copy regular file",
			setupFunc: func() (string, string) {
				src := filepath.Join(t.TempDir(), "source.txt")
				dst := filepath.Join(t.TempDir(), "dest.txt")
				if err := os.WriteFile(src, []byte("test content"), 0644); err != nil {
					t.Fatalf("failed to create source file: %v", err)
				}
				return src, dst
			},
			wantError: false,
			validate: func(t *testing.T, dst string) {
				content, err := os.ReadFile(dst)
				if err != nil {
					t.Fatalf("failed to read destination: %v", err)
				}
				if string(content) != "test content" {
					t.Errorf("content = %q, want %q", string(content), "test content")
				}
			},
		},
		{
			name: "copy preserves permissions",
			setupFunc: func() (string, string) {
				src := filepath.Join(t.TempDir(), "source.sh")
				dst := filepath.Join(t.TempDir(), "dest.sh")
				if err := os.WriteFile(src, []byte("#!/bin/bash"), 0755); err != nil {
					t.Fatalf("failed to create source file: %v", err)
				}
				return src, dst
			},
			wantError: false,
			validate: func(t *testing.T, dst string) {
				info, err := os.Stat(dst)
				if err != nil {
					t.Fatalf("failed to stat destination: %v", err)
				}
				if info.Mode().Perm() != 0755 {
					t.Errorf("permissions = %o, want 0755", info.Mode().Perm())
				}
			},
		},
		{
			name: "nonexistent source",
			setupFunc: func() (string, string) {
				src := filepath.Join(t.TempDir(), "nonexistent.txt")
				dst := filepath.Join(t.TempDir(), "dest.txt")
				return src, dst
			},
			wantError: true,
			validate:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			src, dst := tt.setupFunc()
			err := CopyFile(src, dst)

			if (err != nil) != tt.wantError {
				t.Errorf("CopyFile() error = %v, wantError %v", err, tt.wantError)
				return
			}

			if !tt.wantError && tt.validate != nil {
				tt.validate(t, dst)
			}
		})
	}
}

// TestCopyDirectory tests recursive directory copying
func TestCopyDirectory(t *testing.T) {
	tests := []struct {
		name      string
		setupFunc func() (src, dst string)
		wantError bool
		validate  func(*testing.T, string)
	}{
		{
			name: "copy directory with files",
			setupFunc: func() (string, string) {
				src := filepath.Join(t.TempDir(), "source")
				dst := filepath.Join(t.TempDir(), "dest")
				os.MkdirAll(src, 0755)
				os.WriteFile(filepath.Join(src, "file1.txt"), []byte("content1"), 0644)
				os.WriteFile(filepath.Join(src, "file2.txt"), []byte("content2"), 0644)
				return src, dst
			},
			wantError: false,
			validate: func(t *testing.T, dst string) {
				// Check files exist
				if _, err := os.Stat(filepath.Join(dst, "file1.txt")); os.IsNotExist(err) {
					t.Error("file1.txt not copied")
				}
				if _, err := os.Stat(filepath.Join(dst, "file2.txt")); os.IsNotExist(err) {
					t.Error("file2.txt not copied")
				}

				// Check content
				content, _ := os.ReadFile(filepath.Join(dst, "file1.txt"))
				if string(content) != "content1" {
					t.Errorf("file1 content = %q, want %q", string(content), "content1")
				}
			},
		},
		{
			name: "copy nested directories",
			setupFunc: func() (string, string) {
				src := filepath.Join(t.TempDir(), "source")
				dst := filepath.Join(t.TempDir(), "dest")
				os.MkdirAll(filepath.Join(src, "subdir1", "subdir2"), 0755)
				os.WriteFile(filepath.Join(src, "root.txt"), []byte("root"), 0644)
				os.WriteFile(filepath.Join(src, "subdir1", "sub1.txt"), []byte("sub1"), 0644)
				os.WriteFile(filepath.Join(src, "subdir1", "subdir2", "sub2.txt"), []byte("sub2"), 0644)
				return src, dst
			},
			wantError: false,
			validate: func(t *testing.T, dst string) {
				// Check nested structure
				paths := []string{
					"root.txt",
					"subdir1/sub1.txt",
					"subdir1/subdir2/sub2.txt",
				}
				for _, p := range paths {
					fullPath := filepath.Join(dst, p)
					if _, err := os.Stat(fullPath); os.IsNotExist(err) {
						t.Errorf("%s not copied", p)
					}
				}
			},
		},
		{
			name: "copy empty directory",
			setupFunc: func() (string, string) {
				src := filepath.Join(t.TempDir(), "source")
				dst := filepath.Join(t.TempDir(), "dest")
				os.MkdirAll(src, 0755)
				return src, dst
			},
			wantError: false,
			validate: func(t *testing.T, dst string) {
				info, err := os.Stat(dst)
				if os.IsNotExist(err) {
					t.Error("destination not created")
				}
				if !info.IsDir() {
					t.Error("destination is not a directory")
				}
			},
		},
		{
			name: "nonexistent source",
			setupFunc: func() (string, string) {
				src := filepath.Join(t.TempDir(), "nonexistent")
				dst := filepath.Join(t.TempDir(), "dest")
				return src, dst
			},
			wantError: true,
			validate:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			src, dst := tt.setupFunc()
			err := CopyDirectory(src, dst)

			if (err != nil) != tt.wantError {
				t.Errorf("CopyDirectory() error = %v, wantError %v", err, tt.wantError)
				return
			}

			if !tt.wantError && tt.validate != nil {
				tt.validate(t, dst)
			}
		})
	}
}

// TestCreateSync tests sync creation
func TestCreateSync(t *testing.T) {
	tests := []struct {
		name      string
		setupFunc func() (repoPath, sheetsPath string, force bool)
		wantError bool
		validate  func(*testing.T, string, string)
	}{
		{
			name: "create new sync",
			setupFunc: func() (string, string, bool) {
				repoPath := createTestCheatsheetsDir(t)
				sheetsPath := filepath.Join(t.TempDir(), "sheets")
				return repoPath, sheetsPath, false
			},
			wantError: false,
			validate: func(t *testing.T, repoPath, sheetsPath string) {
				status := CheckSyncStatus(repoPath, sheetsPath)
				if !status.Exists {
					t.Error("sync not created")
				}
				if !status.IsValid {
					t.Error("sync not valid")
				}
			},
		},
		{
			name: "create with existing valid sync (no force)",
			setupFunc: func() (string, string, bool) {
				repoPath := createTestCheatsheetsDir(t)
				sheetsPath := filepath.Join(t.TempDir(), "sheets")
				// Create initial sync
				sourcePath := cheatsheetservice.GetCheatsheetsPath(repoPath)
				CopyDirectory(sourcePath, sheetsPath)
				return repoPath, sheetsPath, false
			},
			wantError: false, // Should succeed but do nothing
			validate: func(t *testing.T, repoPath, sheetsPath string) {
				status := CheckSyncStatus(repoPath, sheetsPath)
				if !status.IsValid {
					t.Error("sync should still be valid")
				}
			},
		},
		{
			name: "force recreate existing sync",
			setupFunc: func() (string, string, bool) {
				repoPath := createTestCheatsheetsDir(t)
				sheetsPath := filepath.Join(t.TempDir(), "sheets")
				// Create initial sync
				sourcePath := cheatsheetservice.GetCheatsheetsPath(repoPath)
				CopyDirectory(sourcePath, sheetsPath)
				return repoPath, sheetsPath, true
			},
			wantError: false,
			validate: func(t *testing.T, repoPath, sheetsPath string) {
				status := CheckSyncStatus(repoPath, sheetsPath)
				if !status.IsValid {
					t.Error("recreated sync should be valid")
				}
			},
		},
		{
			name: "replace invalid sync",
			setupFunc: func() (string, string, bool) {
				repoPath := createTestCheatsheetsDir(t)
				sheetsPath := filepath.Join(t.TempDir(), "sheets")
				// Create a file instead of directory (invalid)
				os.WriteFile(sheetsPath, []byte("invalid"), 0644)
				return repoPath, sheetsPath, false
			},
			wantError: false,
			validate: func(t *testing.T, repoPath, sheetsPath string) {
				status := CheckSyncStatus(repoPath, sheetsPath)
				if !status.IsValid {
					t.Error("replaced sync should be valid")
				}
				// Should be a directory now
				if status.Type == SyncTypeMismatch {
					t.Error("should not be mismatch after replacement")
				}
			},
		},
		{
			name: "invalid source",
			setupFunc: func() (string, string, bool) {
				repoPath := filepath.Join(t.TempDir(), "nonexistent")
				sheetsPath := filepath.Join(t.TempDir(), "sheets")
				return repoPath, sheetsPath, false
			},
			wantError: true,
			validate:  nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoPath, sheetsPath, force := tt.setupFunc()
			err := CreateSync(repoPath, sheetsPath, force)

			if (err != nil) != tt.wantError {
				t.Errorf("CreateSync() error = %v, wantError %v", err, tt.wantError)
				return
			}

			if !tt.wantError && tt.validate != nil {
				tt.validate(t, repoPath, sheetsPath)
			}
		})
	}
}

// TestUpdateSync tests sync updates
func TestUpdateSync(t *testing.T) {
	tests := []struct {
		name      string
		setupFunc func() (repoPath, sheetsPath string)
		wantError bool
	}{
		{
			name: "update nonexistent sync",
			setupFunc: func() (string, string) {
				repoPath := createTestCheatsheetsDir(t)
				sheetsPath := filepath.Join(t.TempDir(), "sheets")
				return repoPath, sheetsPath
			},
			wantError: true,
		},
		{
			name: "update symlink (no-op)",
			setupFunc: func() (string, string) {
				if runtime.GOOS == "windows" {
					t.Skip("Symlink test skipped on Windows")
				}
				repoPath := createTestCheatsheetsDir(t)
				sheetsPath := filepath.Join(t.TempDir(), "sheets")
				sourcePath := cheatsheetservice.GetCheatsheetsPath(repoPath)
				os.Symlink(sourcePath, sheetsPath)
				return repoPath, sheetsPath
			},
			wantError: false,
		},
		{
			name: "update copied sync",
			setupFunc: func() (string, string) {
				repoPath := createTestCheatsheetsDir(t)
				sheetsPath := filepath.Join(t.TempDir(), "sheets")
				sourcePath := cheatsheetservice.GetCheatsheetsPath(repoPath)
				CopyDirectory(sourcePath, sheetsPath)
				return repoPath, sheetsPath
			},
			wantError: false,
		},
		{
			name: "update mismatch type",
			setupFunc: func() (string, string) {
				repoPath := createTestCheatsheetsDir(t)
				sheetsPath := filepath.Join(t.TempDir(), "sheets")
				os.WriteFile(sheetsPath, []byte("file"), 0644)
				return repoPath, sheetsPath
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoPath, sheetsPath := tt.setupFunc()
			err := UpdateSync(repoPath, sheetsPath)

			if (err != nil) != tt.wantError {
				t.Errorf("UpdateSync() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}
