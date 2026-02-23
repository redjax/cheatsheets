package cdcommand

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/redjax/cheatsheets/internal/config"
	"github.com/spf13/cobra"
)

var CdCmd = &cobra.Command{
	Use:   "cd",
	Short: "Open a shell in the cheatsheets repository directory",
	Long: `Open a new shell session in the cheatsheets repository directory.
This works similarly to 'chezmoi cd' - it spawns a subshell in the repository location.
Use 'exit' to return to your original directory.`,
	RunE: runCd,
}

func runCd(cmd *cobra.Command, args []string) error {
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

	repoPath := cfg.Git.ClonePath
	if repoPath == "" {
		return fmt.Errorf("git.ClonePath not configured")
	}

	// Verify the repository path exists
	if _, err := os.Stat(repoPath); os.IsNotExist(err) {
		return fmt.Errorf("repository path does not exist: %s\nRun 'chtsht repo clone' first", repoPath)
	}

	return openShellInDirectory(repoPath)
}

// openShellInDirectory spawns a new shell in the specified directory
func openShellInDirectory(path string) error {
	var shell string
	var args []string

	// Determine the shell to use based on the operating system
	if runtime.GOOS == "windows" {
		// On Windows, use PowerShell if available, otherwise cmd.exe
		shell = os.Getenv("ComSpec")
		if shell == "" {
			shell = "cmd.exe"
		}
		// Check if PowerShell is available
		if _, err := exec.LookPath("pwsh.exe"); err == nil {
			shell = "pwsh.exe"
			args = []string{"-NoLogo"}
		} else if _, err := exec.LookPath("powershell.exe"); err == nil {
			shell = "powershell.exe"
			args = []string{"-NoLogo"}
		}
	} else {
		// On Unix-like systems, use the user's shell
		shell = os.Getenv("SHELL")
		if shell == "" {
			shell = "/bin/sh"
		}
	}

	fmt.Printf("Opening shell in: %s\n", path)
	fmt.Println("Type 'exit' to return to your original directory.")
	fmt.Println()

	// Create the command
	shellCmd := exec.Command(shell, args...)
	shellCmd.Dir = path
	shellCmd.Stdin = os.Stdin
	shellCmd.Stdout = os.Stdout
	shellCmd.Stderr = os.Stderr

	// Run the shell and wait for it to exit
	if err := shellCmd.Run(); err != nil {
		return fmt.Errorf("error running shell: %w", err)
	}

	return nil
}
