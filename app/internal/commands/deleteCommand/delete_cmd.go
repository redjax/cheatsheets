package deletecommand

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/manifoldco/promptui"
	"github.com/redjax/cheatsheets/internal/config"
	"github.com/redjax/cheatsheets/internal/guards"
	cheatsheetservice "github.com/redjax/cheatsheets/internal/services/cheatsheetService"
	"github.com/spf13/cobra"
)

var (
	typeFilter string
	force      bool
)

// DeleteCmd represents the delete command
var DeleteCmd = &cobra.Command{
	Use:     "delete [name]",
	Short:   "Delete a cheatsheet",
	Long:    `Delete a cheatsheet from the repository. Prompts for confirmation unless --force is used.`,
	Example: "  chtsht delete git\n  chtsht delete -t command git\n  chtsht delete git --force\n  chtsht delete  # interactive selector",
	Aliases: []string{"del", "rm", "remove"},
	Args:    cobra.MaximumNArgs(1),
	RunE:    runDelete,
}

func init() {
	DeleteCmd.Flags().StringVarP(&typeFilter, "type", "t", "", "Filter by cheatsheet type (app, command, language, system)")
	DeleteCmd.Flags().BoolVarP(&force, "force", "f", false, "Skip confirmation prompt")
}

func runDelete(cmd *cobra.Command, args []string) error {
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

	// Pre-flight checks - ensure repo exists
	guardCtx := guards.NewGuardContext(cfg)
	if err := guards.CheckAll(guardCtx, guards.RepoCloned); err != nil {
		return err
	}

	repoPath := cfg.Git.ClonePath

	// Validate repository directory exists
	if err := cheatsheetservice.ValidateCheatsheetsDirectory(repoPath); err != nil {
		return fmt.Errorf("git repository not found: %w\nRun 'chtsht repo clone' to clone the repository", err)
	}

	// If no name provided, show selector
	if len(args) == 0 {
		return deleteWithSelector(cfg, typeFilter, force)
	}

	name := args[0]

	// If type is specified, delete directly
	if typeFilter != "" {
		return deleteCheatsheet(cfg, typeFilter, name, force)
	}

	// Otherwise, search for the cheatsheet by name
	return deleteCheatsheetByName(cfg, name, force)
}

// deleteCheatsheet deletes a specific cheatsheet by type and name
func deleteCheatsheet(cfg *config.Config, cheatsheetType, name string, skipConfirm bool) error {
	filePath := cheatsheetservice.GetCheatsheetPath(cfg.Git.ClonePath, cheatsheetType, name)

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("cheatsheet not found: [%s] %s", cheatsheetType, name)
	}

	// Confirm deletion unless force flag is set
	if !skipConfirm {
		if !confirmDeletion(cheatsheetType, name) {
			fmt.Println("Deletion cancelled.")
			return nil
		}
	}

	// Delete the file
	if err := cheatsheetservice.DeleteCheatsheet(cfg.Git.ClonePath, cheatsheetType, name); err != nil {
		return err
	}

	fmt.Printf("Deleted [%s] %s\n", cheatsheetType, name)
	return nil
}

// deleteCheatsheetByName finds and deletes a cheatsheet by name (searching all types)
func deleteCheatsheetByName(cfg *config.Config, name string, skipConfirm bool) error {
	availableTypes, err := cheatsheetservice.GetAvailableTypes(cfg.Git.ClonePath)
	if err != nil {
		return fmt.Errorf("error getting available types: %w", err)
	}

	// Search for cheatsheets with this name across all types
	var foundPaths []struct {
		Type string
		Path string
	}

	for _, t := range availableTypes {
		filePath := cheatsheetservice.GetCheatsheetPath(cfg.Git.ClonePath, t, name)
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

	// If only one match, delete it directly
	if len(foundPaths) == 1 {
		return deleteCheatsheet(cfg, foundPaths[0].Type, name, skipConfirm)
	}

	// Multiple matches - show selector
	fmt.Printf("Multiple cheatsheets named '%s' found:\n\n", name)

	// Create options for the selector
	type selectOption struct {
		Display string
		Type    string
	}

	options := make([]selectOption, len(foundPaths))
	for i, match := range foundPaths {
		options[i] = selectOption{
			Display: fmt.Sprintf("[%s] %s", match.Type, name),
			Type:    match.Type,
		}
	}

	// Create promptui selector
	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "▸ {{ .Display | red }}",
		Inactive: "  {{ .Display }}",
		Selected: "{{ .Display | red }}",
	}

	prompt := promptui.Select{
		Label:     "Select a cheatsheet to delete",
		Items:     options,
		Templates: templates,
		Size:      15,
	}

	idx, _, err := prompt.Run()
	if err != nil {
		return fmt.Errorf("selection cancelled")
	}

	selected := options[idx]
	fmt.Println() // Add blank line for spacing

	return deleteCheatsheet(cfg, selected.Type, name, skipConfirm)
}

// confirmDeletion prompts the user to confirm deletion
func confirmDeletion(cheatsheetType, name string) bool {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("Are you sure you want to delete [%s] %s? (y/N): ", cheatsheetType, name)
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(strings.ToLower(input))
	return input == "y" || input == "yes"
}

// deleteWithSelector shows an interactive selector for choosing a cheatsheet to delete
func deleteWithSelector(cfg *config.Config, typeFilter string, skipConfirm bool) error {
	availableTypes, err := cheatsheetservice.GetAvailableTypes(cfg.Git.ClonePath)
	if err != nil {
		return fmt.Errorf("error getting available types: %w", err)
	}

	// Filter types if specified
	typesToList := availableTypes
	if typeFilter != "" {
		// Validate the provided type
		valid := false
		for _, t := range availableTypes {
			if t == typeFilter {
				valid = true
				break
			}
		}
		if !valid {
			return fmt.Errorf("type '%s' not found in repository", typeFilter)
		}
		typesToList = []string{typeFilter}
	}

	// Collect all cheatsheets
	type selectOption struct {
		Display string
		Type    string
		Name    string
	}

	var options []selectOption
	for _, t := range typesToList {
		cheatsheets, err := cheatsheetservice.GetCheatsheetsByType(cfg.Git.ClonePath, t)
		if err != nil {
			return fmt.Errorf("error getting cheatsheets for type %s: %w", t, err)
		}

		for _, name := range cheatsheets {
			options = append(options, selectOption{
				Display: fmt.Sprintf("[%s] %s", t, name),
				Type:    t,
				Name:    name,
			})
		}
	}

	if len(options) == 0 {
		if typeFilter != "" {
			return fmt.Errorf("no cheatsheets found for type: %s", typeFilter)
		}
		return fmt.Errorf("no cheatsheets found")
	}

	// Create promptui selector
	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "▸ {{ .Display | red }}",
		Inactive: "  {{ .Display }}",
		Selected: "{{ .Display | red }}",
	}

	prompt := promptui.Select{
		Label:     "Select a cheatsheet to delete",
		Items:     options,
		Templates: templates,
		Size:      15,
	}

	idx, _, err := prompt.Run()
	if err != nil {
		return fmt.Errorf("selection cancelled")
	}

	selected := options[idx]
	fmt.Println() // Add blank line for spacing

	return deleteCheatsheet(cfg, selected.Type, selected.Name, skipConfirm)
}
