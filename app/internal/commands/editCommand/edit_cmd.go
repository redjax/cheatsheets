package editcommand

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"time"

	"github.com/redjax/cheatsheets/internal/config"
	cheatsheetservice "github.com/redjax/cheatsheets/internal/services/cheatsheetService"
	reposervices "github.com/redjax/cheatsheets/internal/services/repoServices"
	"github.com/spf13/cobra"
)

var typeFilter string

// EditCmd represents the edit command
var EditCmd = &cobra.Command{
	Use:     "edit [name]",
	Short:   "Edit a cheatsheet in your default editor",
	Long:    `Opens a cheatsheet markdown file from the git repository in your default editor ($EDITOR, $VISUAL, or system default).`,
	Example: "  chtsht edit git\n  chtsht edit -t command git\n  chtsht edit -t app neovim",
	Aliases: []string{"e"},
	Args:    cobra.MaximumNArgs(1),
	RunE:    runEdit,
}

func init() {
	EditCmd.Flags().StringVarP(&typeFilter, "type", "t", "", "Filter by cheatsheet type (app, command, language, system)")
}

func runEdit(cmd *cobra.Command, args []string) error {
	// Get the config file path from the persistent flag
	configFile, err := cmd.Flags().GetString("config-file")
	if err != nil {
		return fmt.Errorf("error getting config-file flag: %w", err)
	}

	// Load config
	var cfg *config.Config
	if configFile == "" {
		configFile = config.FindConfigFile("config.yml")
		cfg, err = config.LoadConfig(nil, configFile)
	} else {
		cfg, err = config.LoadConfig(nil, configFile)
	}

	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Validate git repository path is configured
	if cfg.Git.ClonePath == "" {
		return fmt.Errorf("git repository path not configured")
	}

	// Validate repository directory exists
	if err := cheatsheetservice.ValidateCheatsheetsDirectory(cfg.Git.ClonePath); err != nil {
		return fmt.Errorf("git repository not found: %w\nRun 'chtsht repo clone' to clone the repository", err)
	}

	// Auto-switch to working branch if enabled and on main
	if cfg.Git.AutoBranch {
		currentBranch, err := reposervices.GetCurrentBranch(cfg.Git.ClonePath)
		if err == nil && (currentBranch == "main" || currentBranch == "master") {
			workingBranch := cfg.Git.WorkingBranch
			if workingBranch == "" {
				workingBranch = "working"
			}
			_, _ = reposervices.EnsureWorkingBranch(cfg.Git.ClonePath, workingBranch)
		}
	}

	// Handle different scenarios
	if len(args) == 0 && typeFilter == "" {
		// No arguments - show interactive selector
		return editWithSelector(cfg.Git.ClonePath, "")
	} else if len(args) == 0 && typeFilter != "" {
		// Only type provided - show selector filtered by type
		return editWithSelector(cfg.Git.ClonePath, typeFilter)
	} else if len(args) == 1 && typeFilter != "" {
		// Both type and name provided
		return editCheatsheet(cfg.Git.ClonePath, typeFilter, args[0])
	} else {
		// Only name provided - search across all types
		return editCheatsheetByName(cfg.Git.ClonePath, args[0])
	}
}

// editCheatsheet opens a specific cheatsheet file in the default editor
func editCheatsheet(repoPath, typeDir, name string) error {
	filePath := cheatsheetservice.GetCheatsheetPath(repoPath, typeDir, name)

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("cheatsheet '%s' not found in type '%s'", name, typeDir)
	}

	return openInEditor(filePath)
}

// editCheatsheetByName finds and edits a cheatsheet by name across all types
func editCheatsheetByName(repoPath, name string) error {
	availableTypes, err := cheatsheetservice.GetAvailableTypes(repoPath)
	if err != nil {
		return fmt.Errorf("error getting available types: %w", err)
	}

	// Search for the cheatsheet in all types
	var foundPaths []struct {
		Type string
		Path string
	}

	for _, t := range availableTypes {
		filePath := cheatsheetservice.GetCheatsheetPath(repoPath, t, name)
		if _, err := os.Stat(filePath); err == nil {
			foundPaths = append(foundPaths, struct {
				Type string
				Path string
			}{Type: t, Path: filePath})
		}
	}

	if len(foundPaths) == 0 {
		return fmt.Errorf("cheatsheet '%s' not found in any type", name)
	}

	// If only one match, open it directly
	if len(foundPaths) == 1 {
		fmt.Printf("Editing [%s] %s\n", foundPaths[0].Type, name)
		return openInEditor(foundPaths[0].Path)
	}

	// Multiple matches - show selector
	fmt.Printf("Multiple cheatsheets named '%s' found. Please specify type with -t flag:\n", name)
	for _, match := range foundPaths {
		fmt.Printf("  chtsht edit -t %s %s\n", match.Type, name)
	}
	return nil
}

// editWithSelector shows an interactive selector for choosing a cheatsheet to edit
func editWithSelector(repoPath, typeFilter string) error {
	// Implementation would be similar to show command's selector
	// For now, return an error suggesting to specify a name
	if typeFilter != "" {
		return fmt.Errorf("please specify a cheatsheet name: chtsht edit -t %s <name>", typeFilter)
	}
	return fmt.Errorf("please specify a cheatsheet name: chtsht edit <name>")
}

// openInEditor opens a file in the user's default editor using a temp file for safety
func openInEditor(originalPath string) error {
	editor, err := getEditor()
	if err != nil {
		return err
	}

	// Create cheatsheets temp directory for easier cleanup
	tempDir := filepath.Join(os.TempDir(), "cheatsheets")
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}

	// Create a temporary file for editing
	tempFile, err := os.CreateTemp(tempDir, "chtsht-edit-*.md")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	tempPath := tempFile.Name()

	// Ensure temp file is cleaned up on exit
	defer func() {
		tempFile.Close()
		os.Remove(tempPath)
	}()

	// Copy original file to temp file
	originalFile, err := os.Open(originalPath)
	if err != nil {
		return fmt.Errorf("failed to open original file: %w", err)
	}

	_, err = io.Copy(tempFile, originalFile)
	originalFile.Close()
	tempFile.Close() // Close before opening in editor

	if err != nil {
		return fmt.Errorf("failed to copy file to temp location: %w", err)
	}

	// Get original file info for permission preservation
	originalInfo, err := os.Stat(originalPath)
	if err != nil {
		return fmt.Errorf("failed to stat original file: %w", err)
	}

	fmt.Printf("Opening %s in %s\n", filepath.Base(originalPath), editor)

	// Create command to open editor with temp file
	cmd := exec.Command(editor, tempPath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run editor and wait for it to exit
	if err := cmd.Run(); err != nil {
		// Check if this is just a non-zero exit code
		if exitErr, ok := err.(*exec.ExitError); ok {
			return fmt.Errorf("editor exited with code %d, changes not saved", exitErr.ExitCode())
		}
		return fmt.Errorf("error running editor: %w", err)
	}

	// Editor exited successfully - check if temp file was modified
	tempInfo, err := os.Stat(tempPath)
	if err != nil {
		return fmt.Errorf("failed to stat temp file after editing: %w", err)
	}

	// Compare modification times
	if !tempInfo.ModTime().After(originalInfo.ModTime()) {
		fmt.Println("No changes made.")
		return nil
	}

	// Read the edited content from temp file
	editedContent, err := os.ReadFile(tempPath)
	if err != nil {
		return fmt.Errorf("failed to read edited content: %w", err)
	}

	// Update the last_updated field in the frontmatter
	editedContent = updateLastUpdatedDate(editedContent)

	// Write to original file atomically using a backup
	backupPath := originalPath + ".backup"

	// Create backup of original
	if err := copyFile(originalPath, backupPath); err != nil {
		return fmt.Errorf("failed to create backup: %w", err)
	}

	// Ensure backup is cleaned up after successful write
	defer os.Remove(backupPath)

	// Write new content to original file
	if err := os.WriteFile(originalPath, editedContent, originalInfo.Mode()); err != nil {
		// Restore from backup on write failure
		if restoreErr := copyFile(backupPath, originalPath); restoreErr != nil {
			return fmt.Errorf("failed to write changes AND failed to restore backup: write error: %w, restore error: %v", err, restoreErr)
		}
		return fmt.Errorf("failed to write changes (backup restored): %w", err)
	}

	fmt.Printf("Successfully saved changes to %s\n", originalPath)
	return nil
}

// copyFile copies a file from src to dst
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	if _, err := io.Copy(destFile, sourceFile); err != nil {
		return err
	}

	// Sync to ensure data is written to disk
	return destFile.Sync()
}

// getEditor determines the default editor to use
func getEditor() (string, error) {
	// Try environment variables first (Unix/Linux/Mac)
	if editor := os.Getenv("VISUAL"); editor != "" {
		return editor, nil
	}
	if editor := os.Getenv("EDITOR"); editor != "" {
		return editor, nil
	}

	// Platform-specific defaults
	switch runtime.GOOS {
	case "windows":
		// Try notepad as fallback on Windows
		return "notepad", nil
	case "darwin":
		// Try nano on macOS
		if _, err := exec.LookPath("nano"); err == nil {
			return "nano", nil
		}
		if _, err := exec.LookPath("vim"); err == nil {
			return "vim", nil
		}
		return "vi", nil
	default:
		// Linux and other Unix-like systems
		if _, err := exec.LookPath("nano"); err == nil {
			return "nano", nil
		}
		if _, err := exec.LookPath("vim"); err == nil {
			return "vim", nil
		}
		return "vi", nil
	}
}

// updateLastUpdatedDate updates the last_updated field in the YAML frontmatter
func updateLastUpdatedDate(content []byte) []byte {
	contentStr := string(content)
	currentDate := time.Now().Format("2006-01-02")

	// Pattern to match last_updated field in YAML frontmatter
	// Matches: last_updated: "YYYY-MM-DD" or last_updated: 'YYYY-MM-DD' or last_updated: YYYY-MM-DD
	pattern := regexp.MustCompile(`(?m)^last_updated:\s*["']?[0-9]{4}-[0-9]{2}-[0-9]{2}["']?`)

	// Check if last_updated exists in frontmatter
	if pattern.MatchString(contentStr) {
		// Replace existing last_updated
		contentStr = pattern.ReplaceAllString(contentStr, fmt.Sprintf(`last_updated: "%s"`, currentDate))
	} else {
		// If last_updated doesn't exist but frontmatter does, add it
		frontmatterPattern := regexp.MustCompile(`(?s)^---\n(.*?)\n---`)
		if frontmatterPattern.MatchString(contentStr) {
			// Add last_updated to the frontmatter
			contentStr = frontmatterPattern.ReplaceAllStringFunc(contentStr, func(match string) string {
				// Insert last_updated before the closing ---
				lines := strings.Split(match, "\n")
				if len(lines) >= 2 {
					// Insert before the last line (which is ---)
					result := strings.Join(lines[:len(lines)-1], "\n")
					result += fmt.Sprintf("\nlast_updated: \"%s\"\n---", currentDate)
					return result
				}
				return match
			})
		}
	}

	return []byte(contentStr)
}
