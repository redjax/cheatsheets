package editcommand

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/redjax/cheatsheets/internal/config"
	"github.com/spf13/cobra"
)

// EditCmd represents the config edit command
var EditCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit configuration file",
	Long:  `Opens the configuration file in your default editor ($EDITOR, $VISUAL, or system default).`,
	RunE:  runEdit,
}

func runEdit(cmd *cobra.Command, args []string) error {
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

	// Create config file if it doesn't exist
	if _, err := os.Stat(configFile); os.IsNotExist(err) {
		if err := createDefaultConfig(configFile); err != nil {
			return fmt.Errorf("failed to create config file: %w", err)
		}
		fmt.Printf("Created new config file: %s\n", configFile)
	}

	// Get editor
	editor, err := getEditor()
	if err != nil {
		return err
	}

	fmt.Printf("Opening %s in %s...\n", configFile, editor)

	// Open editor
	editorCmd := exec.Command(editor, configFile)
	editorCmd.Stdin = os.Stdin
	editorCmd.Stdout = os.Stdout
	editorCmd.Stderr = os.Stderr

	if err := editorCmd.Run(); err != nil {
		return fmt.Errorf("failed to open editor: %w", err)
	}

	fmt.Println("Config file saved.")
	return nil
}

func getEditor() (string, error) {
	// Try environment variables first
	if editor := os.Getenv("EDITOR"); editor != "" {
		return editor, nil
	}
	if editor := os.Getenv("VISUAL"); editor != "" {
		return editor, nil
	}

	// Fall back to system defaults
	switch runtime.GOOS {
	case "windows":
		return "notepad", nil
	case "darwin":
		return "nano", nil
	default:
		// Try common Linux editors
		for _, editor := range []string{"nano", "vim", "vi"} {
			if _, err := exec.LookPath(editor); err == nil {
				return editor, nil
			}
		}
	}

	return "", fmt.Errorf("no suitable editor found. Set $EDITOR environment variable")
}

func createDefaultConfig(path string) error {
	defaultConfig := `# Cheatsheets Configuration
# See https://github.com/redjax/cheatsheets for details

# Path where cheatsheets are cloned (optional, defaults to ~/.local/share/cheatsheets)
# sheets_path: ~/.cheatsheets

# Git repository configuration
git:
  # Repository URL
  repo_url: "https://github.com/redjax/cheatsheets.git"
  
  # Local clone path (defaults to ~/.local/share/cheatsheets)
  # clone_path: ~/.local/share/cheatsheets
  
  # GitHub token for push access (optional, for private repos or pushing)
  # token: ""
  
  # Auto-switch to working branch when editing (recommended)
  auto_branch: true
  
  # Name of the shared working branch
  working_branch: "working"
  
  # Git author info (optional, falls back to git config)
  # author_name: "Your Name"
  # author_email: "you@example.com"

# Debug mode
# debug: false
`

	return os.WriteFile(path, []byte(defaultConfig), 0644)
}
