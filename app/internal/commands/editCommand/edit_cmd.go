package editcommand

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/redjax/cheatsheets/internal/config"
	cheatsheetservice "github.com/redjax/cheatsheets/internal/services/cheatsheetService"
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

// openInEditor opens a file in the user's default editor
func openInEditor(filePath string) error {
	editor, err := getEditor()
	if err != nil {
		return err
	}

	fmt.Printf("Opening %s in %s...\n", filePath, editor)

	// Create command to open editor
	cmd := exec.Command(editor, filePath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Run editor and wait for it to exit
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error running editor: %w", err)
	}

	fmt.Println("Edit complete.")
	return nil
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
