package syncservice

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"

	cheatsheetservice "github.com/redjax/cheatsheets/internal/services/cheatsheetService"
)

// SyncType represents how the sheets are synced
type SyncType string

const (
	SyncTypeSymlink  SyncType = "symlink"
	SyncTypeCopy     SyncType = "copy"
	SyncTypeNone     SyncType = "none"
	SyncTypeMismatch SyncType = "mismatch"
)

// SyncStatus represents the current sync state
type SyncStatus struct {
	Exists      bool
	Type        SyncType
	Source      string
	Destination string
	IsValid     bool
	Error       error
}

// CheckSyncStatus checks the current state of the sheets sync
func CheckSyncStatus(repoPath, sheetsPath string) *SyncStatus {
	status := &SyncStatus{
		Source:      cheatsheetservice.GetCheatsheetsPath(repoPath),
		Destination: sheetsPath,
	}

	// Check if destination exists
	destInfo, err := os.Lstat(sheetsPath)
	if os.IsNotExist(err) {
		status.Exists = false
		status.Type = SyncTypeNone
		return status
	}
	if err != nil {
		status.Error = fmt.Errorf("error checking destination: %w", err)
		return status
	}

	status.Exists = true

	// Check if it's a symlink
	if destInfo.Mode()&os.ModeSymlink != 0 {
		status.Type = SyncTypeSymlink

		// Verify symlink points to correct location
		linkTarget, err := os.Readlink(sheetsPath)
		if err != nil {
			status.Error = fmt.Errorf("error reading symlink: %w", err)
			return status
		}

		// Resolve to absolute path for comparison
		absLinkTarget, err := filepath.Abs(linkTarget)
		if err != nil {
			absLinkTarget = linkTarget
		}

		absSource, err := filepath.Abs(status.Source)
		if err != nil {
			absSource = status.Source
		}

		status.IsValid = absLinkTarget == absSource
		return status
	}

	// Check if it's a directory
	if destInfo.IsDir() {
		status.Type = SyncTypeCopy
		// For copied directories, we consider it valid if it exists
		// User can force re-sync if needed
		status.IsValid = true
		return status
	}

	// It's a file or something else - this is a problem
	status.Type = SyncTypeMismatch
	status.IsValid = false
	return status
}

// CreateSync creates the sync (symlink or copy) from source to destination
func CreateSync(repoPath, sheetsPath string, force bool) error {
	sourcePath := cheatsheetservice.GetCheatsheetsPath(repoPath)

	// Validate source exists
	if err := cheatsheetservice.ValidateCheatsheetsDirectory(repoPath); err != nil {
		return fmt.Errorf("source cheatsheets directory not found: %w", err)
	}

	// Check current status
	status := CheckSyncStatus(repoPath, sheetsPath)

	// If exists and valid, nothing to do (unless force)
	if status.Exists && status.IsValid && !force {
		fmt.Println("✓ Sync already exists and is valid")
		return nil
	}

	// If exists, remove it (either force, invalid, or wrong type)
	if status.Exists {
		if force {
			fmt.Printf("Removing existing destination (--force): %s\n", sheetsPath)
		} else {
			fmt.Printf("Removing invalid destination: %s\n", sheetsPath)
		}
		if err := os.RemoveAll(sheetsPath); err != nil {
			return fmt.Errorf("error removing existing destination: %w", err)
		}
	}

	// Ensure parent directory exists
	parentDir := filepath.Dir(sheetsPath)
	if err := os.MkdirAll(parentDir, 0755); err != nil {
		return fmt.Errorf("error creating parent directory: %w", err)
	}

	// Try symlink first on Unix-like systems
	if runtime.GOOS != "windows" {
		err := os.Symlink(sourcePath, sheetsPath)
		if err == nil {
			fmt.Printf("Created symlink: %s -> %s\n", sheetsPath, sourcePath)
			return nil
		}
		fmt.Printf("⚠ Failed to create symlink: %v\n", err)
		fmt.Println("Falling back to copy")
	}

	// Fall back to copying
	return CopyDirectory(sourcePath, sheetsPath)
}

// CopyDirectory recursively copies a directory tree
func CopyDirectory(src, dst string) error {
	fmt.Printf("Copying %s to %s\n", src, dst)

	// Get source directory info
	srcInfo, err := os.Stat(src)
	if err != nil {
		return fmt.Errorf("error stating source: %w", err)
	}

	// Create destination directory
	if err := os.MkdirAll(dst, srcInfo.Mode()); err != nil {
		return fmt.Errorf("error creating destination: %w", err)
	}

	// Read source directory
	entries, err := os.ReadDir(src)
	if err != nil {
		return fmt.Errorf("error reading source directory: %w", err)
	}

	// Copy each entry
	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			// Recursively copy subdirectory
			if err := CopyDirectory(srcPath, dstPath); err != nil {
				return err
			}
		} else {
			// Copy file
			if err := CopyFile(srcPath, dstPath); err != nil {
				return fmt.Errorf("error copying file %s: %w", entry.Name(), err)
			}
		}
	}

	fmt.Printf("Copied directory: %s\n", dst)
	return nil
}

// CopyFile copies a single file
func CopyFile(src, dst string) error {
	// Open source file
	srcFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	// Get source file info for permissions
	srcInfo, err := srcFile.Stat()
	if err != nil {
		return err
	}

	// Create destination file
	dstFile, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, srcInfo.Mode())
	if err != nil {
		return err
	}
	defer dstFile.Close()

	// Copy contents
	if _, err := io.Copy(dstFile, srcFile); err != nil {
		return err
	}

	return nil
}

// UpdateSync updates the sync (re-copies if needed)
func UpdateSync(repoPath, sheetsPath string) error {
	status := CheckSyncStatus(repoPath, sheetsPath)

	if !status.Exists {
		return fmt.Errorf("sync does not exist. Run sync command first")
	}

	// If it's a symlink, it's always up to date
	if status.Type == SyncTypeSymlink {
		fmt.Println("Symlink is always up to date")
		return nil
	}

	// If it's a copy, we need to re-copy
	if status.Type == SyncTypeCopy {
		fmt.Println("Updating copied cheatsheets")
		return CreateSync(repoPath, sheetsPath, true)
	}

	return fmt.Errorf("destination is not a valid sync (type: %s)", status.Type)
}
