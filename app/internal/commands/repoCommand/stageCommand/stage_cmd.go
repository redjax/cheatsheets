package stagecommand

import (
	"fmt"

	"github.com/redjax/cheatsheets/internal/config"
	reposervices "github.com/redjax/cheatsheets/internal/services/repoServices"
	"github.com/spf13/cobra"
)

var (
	stageAll bool
)

// StageCmd represents the stage command
var StageCmd = &cobra.Command{
	Use:   "stage [files...]",
	Short: "Stage files for commit",
	Long: `Stage specific files or all changes for commit.
	
Examples:
  chtsht repo stage --all              # Stage all changes
  chtsht repo stage file1.md file2.md  # Stage specific files
  chtsht repo stage                    # Same as --all`,
	RunE: func(cmd *cobra.Command, args []string) error {
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

		// Check if repository is cloned
		cloned, err := reposervices.IsRepositoryCloned(cfg.Git.ClonePath)
		if err != nil {
			return fmt.Errorf("error checking repository: %w", err)
		}

		if !cloned {
			return fmt.Errorf("repository is not cloned at %s. Run 'chtsht repo clone' first", cfg.Git.ClonePath)
		}

		// Check if working tree is clean
		clean, err := reposervices.IsWorkingTreeClean(cfg.Git.ClonePath)
		if err != nil {
			return fmt.Errorf("failed to check working tree: %w", err)
		}

		if clean {
			fmt.Println("No changes to stage - working tree is clean")
			return nil
		}

		// If no files specified or --all flag, stage all
		if len(args) == 0 || stageAll {
			fmt.Println("Staging all changes")
			err = reposervices.StageAll(cfg.Git.ClonePath)
			if err != nil {
				return fmt.Errorf("failed to stage changes: %w", err)
			}

			// Show what was staged
			staged, err := reposervices.GetStagedFiles(cfg.Git.ClonePath)
			if err != nil {
				return fmt.Errorf("failed to get staged files: %w", err)
			}

			fmt.Printf("Staged %d file(s):\n", len(staged))
			for _, file := range staged {
				fmt.Printf("  %s\n", file)
			}
			fmt.Println("\nNext: chtsht repo commit -m \"your message\"")
		} else {
			// Stage specific files
			fmt.Printf("Staging %d file(s)\n", len(args))
			err = reposervices.StageFiles(cfg.Git.ClonePath, args)
			if err != nil {
				return fmt.Errorf("failed to stage files: %w", err)
			}

			for _, file := range args {
				fmt.Printf("  %s\n", file)
			}
			fmt.Println("\nNext: chtsht repo commit -m \"your message\"")
		}

		return nil
	},
}

func init() {
	StageCmd.Flags().BoolVarP(&stageAll, "all", "a", false, "Stage all changes")
}
