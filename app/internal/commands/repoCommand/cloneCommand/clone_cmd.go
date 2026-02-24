package clonecommand

import (
	"fmt"

	"github.com/redjax/cheatsheets/internal/config"
	"github.com/redjax/cheatsheets/internal/constants"
	reposervices "github.com/redjax/cheatsheets/internal/services/repoServices"
	"github.com/spf13/cobra"
)

// CloneCmd represents the clone command
var CloneCmd = &cobra.Command{
	Use:   "clone",
	Short: "Clone the cheatsheets repository",
	Long:  `Clones the cheatsheets repository from the configured URL to the configured path.`,
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
			configFile = config.FindConfigFile("config.yml")
			cfg, err = config.LoadConfig(nil, configFile)
		} else {
			// Use explicit file
			cfg, err = config.LoadConfig(nil, configFile)
		}

		if err != nil {
			return fmt.Errorf("failed to load config: %w", err)
		}

		// Check if already cloned
		cloned, err := reposervices.IsRepositoryCloned(cfg.Git.ClonePath)
		if err != nil {
			return fmt.Errorf("error checking repository: %w", err)
		}

		if cloned {
			fmt.Printf("Repository already cloned at %s\n", cfg.Git.ClonePath)
			return nil
		}

		// Clone the repository
		fmt.Println("=== Cloning Repository ===")
		fmt.Printf("Repository URL: %s\n", constants.RepoURL)
		fmt.Printf("Clone Path: %s\n", cfg.Git.ClonePath)
		fmt.Println()

		err = reposervices.CloneRepository(constants.RepoURL, cfg.Git.ClonePath, cfg.Git.Token)
		if err != nil {
			return fmt.Errorf("failed to clone repository: %w", err)
		}

		fmt.Println()
		fmt.Println("✓ Repository cloned successfully")
		return nil
	},
}
