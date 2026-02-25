package newcommand

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/redjax/cheatsheets/internal/config"
	cheatsheetservice "github.com/redjax/cheatsheets/internal/services/cheatsheetService"
	reposervices "github.com/redjax/cheatsheets/internal/services/repoServices"
	"github.com/spf13/cobra"
)

var (
	cheatsheetType string
	cheatsheetName string
	title          string
	description    string
	tags           string
)

// NewCmd represents the new command
var NewCmd = &cobra.Command{
	Use:   "new",
	Short: "Create a new cheatsheet from a template",
	Long: `Create a new cheatsheet from a template in the repository.
Templates are located in .templates/ and the new cheatsheet will be created in cheatsheets/<type>/.`,
	Example: `  chtsht new -t app -n docker
  chtsht new -t command -n kubectl --title "Kubectl Commands"
  chtsht new  # interactive prompts`,
	Aliases: []string{"n", "create"},
	RunE:    runNew,
}

func init() {
	NewCmd.Flags().StringVarP(&cheatsheetType, "type", "t", "", "Cheatsheet type (app, command, language, system)")
	NewCmd.Flags().StringVarP(&cheatsheetName, "name", "n", "", "Name for the new cheatsheet (without .md extension)")
	NewCmd.Flags().StringVar(&title, "title", "", "Title for the cheatsheet (prompts if not provided)")
	NewCmd.Flags().StringVar(&description, "description", "", "Description for the cheatsheet (prompts if not provided)")
	NewCmd.Flags().StringVar(&tags, "tags", "", "Comma-separated additional tags")
}

func runNew(cmd *cobra.Command, args []string) error {
	// Get the config file path from the persistent flag
	configFile, err := cmd.Flags().GetString("config-file")
	if err != nil {
		return fmt.Errorf("error getting config-file flag: %w", err)
	}

	// Load config
	var cfg *config.Config
	if configFile == "" {
		configFile = config.FindConfigFile("")
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

	// Auto-switch to working branch if enabled and on main
	if cfg.Git.AutoBranch {
		currentBranch, err := reposervices.GetCurrentBranch(repoPath)
		if err == nil && (currentBranch == "main" || currentBranch == "master") {
			workingBranch := cfg.Git.WorkingBranch
			if workingBranch == "" {
				workingBranch = "working"
			}
			_, _ = reposervices.EnsureWorkingBranch(repoPath, workingBranch)
		}
	}

	// Get available types from the repository
	availableTypes, err := cheatsheetservice.GetAvailableTypes(repoPath)
	if err != nil {
		return fmt.Errorf("failed to get available types: %w", err)
	}

	// Prompt for missing required fields
	reader := bufio.NewReader(os.Stdin)

	if cheatsheetType == "" {
		typesList := strings.Join(availableTypes, "/")
		fmt.Printf("Enter cheatsheet type (%s): ", typesList)
		input, _ := reader.ReadString('\n')
		cheatsheetType = strings.TrimSpace(input)
	}

	if cheatsheetName == "" {
		fmt.Print("Enter cheatsheet name: ")
		input, _ := reader.ReadString('\n')
		cheatsheetName = strings.TrimSpace(input)
	}

	// Prompt for optional metadata if not provided
	fmt.Printf("\nCreating new cheatsheet: %s in %s/\n", cheatsheetName, cheatsheetType)

	// Generate a suggested title from the name (capitalize first letter)
	suggestedTitle := cheatsheetName
	if len(cheatsheetName) > 0 {
		if len(cheatsheetName) == 1 {
			suggestedTitle = strings.ToUpper(cheatsheetName)
		} else {
			suggestedTitle = strings.ToUpper(string(cheatsheetName[0])) + cheatsheetName[1:]
		}
	}

	if title == "" {
		fmt.Printf("Enter title [%s]: ", suggestedTitle)
		input, _ := reader.ReadString('\n')
		title = strings.TrimSpace(input)
		if title == "" {
			title = suggestedTitle
		}
	}

	if description == "" {
		fmt.Print("Enter description: ")
		input, _ := reader.ReadString('\n')
		description = strings.TrimSpace(input)
	}

	if tags == "" {
		fmt.Print("Enter additional tags (comma-separated, optional): ")
		input, _ := reader.ReadString('\n')
		tags = strings.TrimSpace(input)
	}

	// Create the cheatsheet using the service
	opts := cheatsheetservice.CreateCheatsheetOptions{
		Type:        cheatsheetType,
		Name:        cheatsheetName,
		Title:       title,
		Description: description,
		Tags:        tags,
	}

	targetFile, err := cheatsheetservice.CreateCheatsheet(repoPath, opts)
	if err != nil {
		return err
	}

	fmt.Printf("\nSuccessfully created %s\n", targetFile)
	return nil
}
