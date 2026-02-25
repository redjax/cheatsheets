package validatecommand

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/redjax/cheatsheets/internal/config"
	"github.com/redjax/cheatsheets/internal/constants"
	reposervices "github.com/redjax/cheatsheets/internal/services/repoServices"
	"github.com/spf13/cobra"
)

// ValidateCmd represents the validate command
var ValidateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate the cheatsheets repository",
	Long:  `Checks if the repository is cloned, has uncommitted changes, and if there are remote updates available.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get the config file path from the persistent flag
		configFile, err := cmd.Flags().GetString("config-file")
		if err != nil {
			return fmt.Errorf("error getting config-file flag: %w", err)
		}

		// Load config
		var cfg *config.Config
		if configFile == "" {
			// Use FindConfigFile for .local fallback
			configFile = config.FindConfigFile("")
			cfg, err = config.LoadConfig(nil, configFile)
		} else {
			// Use explicit file
			cfg, err = config.LoadConfig(nil, configFile)
		}

		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		fmt.Println("=== Repository Validation ===")
		fmt.Printf("Repository Path: %s\n", cfg.Git.ClonePath)
		fmt.Println()

		// Perform validation
		status := reposervices.ValidateRepository(cfg.Git.ClonePath, cfg.Git.Token)

		// Check for errors
		if status.Error != nil {
			return fmt.Errorf("validation error: %w", status.Error)
		}

		// Display results
		if status.IsCloned {
			fmt.Println("Repository is cloned")
			fmt.Printf("  Current branch: %s\n", status.CurrentBranch)
		} else {
			fmt.Println("✗ Repository is NOT cloned")
			fmt.Println()

			// Prompt user to clone
			fmt.Print("Would you like to clone the repository now? (y/n): ")
			reader := bufio.NewReader(os.Stdin)
			response, err := reader.ReadString('\n')
			if err != nil {
				return fmt.Errorf("error reading input: %w", err)
			}

			response = strings.TrimSpace(strings.ToLower(response))
			if response == "y" || response == "yes" {
				fmt.Println()
				fmt.Println("=== Cloning Repository ===")
				fmt.Printf("Repository URL: %s\n", constants.RepoURL)
				fmt.Printf("Clone Path: %s\n", cfg.Git.ClonePath)
				fmt.Println()

				err = reposervices.CloneRepository(constants.RepoURL, cfg.Git.ClonePath, cfg.Git.Token)
				if err != nil {
					return fmt.Errorf("failed to clone repository: %w", err)
				}

				fmt.Println()
				fmt.Println("Repository cloned successfully")
			} else {
				fmt.Println("\nSkipping clone. Run 'repo clone' when ready.")
			}

			return nil
		}

		if status.IsClean {
			fmt.Println("Working tree is clean")
		} else {
			fmt.Println("✗ Working tree has uncommitted changes")
		}

		if status.HasRemoteUpdates {
			fmt.Println("⚠ Remote has updates available")
			fmt.Println("  Run 'repo update' to pull the latest changes")
		} else {
			fmt.Println("Repository is up to date with remote")
		}

		fmt.Println()

		// Summary
		if status.IsClean && !status.HasRemoteUpdates {
			fmt.Println("Repository is in good state")
		} else {
			fmt.Println("⚠ Repository needs attention")
		}

		return nil
	},
}
