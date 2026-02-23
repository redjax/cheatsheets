package cheatsheetservice

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/manifoldco/promptui"
)

// CheatsheetInfo represents a cheatsheet with its type and name
type CheatsheetInfo struct {
	Type string
	Name string
}

// GetCheatsheetsPath returns the path to the cheatsheets directory
// If the path already contains type directories (app, command, language, system),
// it returns the path as-is. Otherwise, it appends "cheatsheets".
func GetCheatsheetsPath(basePath string) string {
	// Check if this is already a cheatsheets directory by looking for common subdirs
	commonDirs := []string{"app", "command", "language", "system"}
	for _, dir := range commonDirs {
		checkPath := filepath.Join(basePath, dir)
		if info, err := os.Stat(checkPath); err == nil && info.IsDir() {
			// Found a type directory, so basePath is already the cheatsheets path
			return basePath
		}
	}

	// Not found, assume this is a repo path and append "cheatsheets"
	return filepath.Join(basePath, "cheatsheets")
}

// ValidateCheatsheetsDirectory checks if the cheatsheets directory exists
func ValidateCheatsheetsDirectory(basePath string) error {
	cheatsheetsPath := GetCheatsheetsPath(basePath)
	if _, err := os.Stat(cheatsheetsPath); os.IsNotExist(err) {
		return fmt.Errorf("cheatsheets directory not found at %s", cheatsheetsPath)
	}
	return nil
}

// GetAvailableTypes returns a list of subdirectories in the cheatsheets directory
func GetAvailableTypes(repoPath string) ([]string, error) {
	cheatsheetsPath := GetCheatsheetsPath(repoPath)

	entries, err := os.ReadDir(cheatsheetsPath)
	if err != nil {
		return nil, err
	}

	var types []string
	for _, entry := range entries {
		if entry.IsDir() {
			types = append(types, entry.Name())
		}
	}

	sort.Strings(types)
	return types, nil
}

// GetCheatsheetsByType returns a list of cheatsheet names in a specific type directory
func GetCheatsheetsByType(repoPath, typeDir string) ([]string, error) {
	cheatsheetsPath := GetCheatsheetsPath(repoPath)
	typePath := filepath.Join(cheatsheetsPath, typeDir)

	entries, err := os.ReadDir(typePath)
	if err != nil {
		return nil, err
	}

	var cheatsheets []string
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(strings.ToLower(entry.Name()), ".md") {
			// Remove .md extension for display
			name := strings.TrimSuffix(entry.Name(), ".md")
			cheatsheets = append(cheatsheets, name)
		}
	}

	sort.Strings(cheatsheets)
	return cheatsheets, nil
}

// ValidateType checks if a type exists in the available types
func ValidateType(repoPath, typeFilter string) (bool, []string, error) {
	availableTypes, err := GetAvailableTypes(repoPath)
	if err != nil {
		return false, nil, err
	}

	if typeFilter == "" {
		return true, availableTypes, nil
	}

	for _, t := range availableTypes {
		if t == typeFilter {
			return true, availableTypes, nil
		}
	}

	return false, availableTypes, nil
}

// ListResult contains the results of listing cheatsheets
type ListResult struct {
	TypesWithSheets map[string][]string
	TotalCount      int
}

// ListCheatsheets returns all cheatsheets, optionally filtered by type
func ListCheatsheets(repoPath, typeFilter string) (*ListResult, error) {
	availableTypes, err := GetAvailableTypes(repoPath)
	if err != nil {
		return nil, fmt.Errorf("error getting available types: %w", err)
	}

	if len(availableTypes) == 0 {
		return &ListResult{
			TypesWithSheets: make(map[string][]string),
			TotalCount:      0,
		}, nil
	}

	// Determine which types to list
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
			return nil, fmt.Errorf("type '%s' not found in repository", typeFilter)
		}
		typesToList = []string{typeFilter}
	}

	result := &ListResult{
		TypesWithSheets: make(map[string][]string),
		TotalCount:      0,
	}

	for _, t := range typesToList {
		cheatsheets, err := GetCheatsheetsByType(repoPath, t)
		if err != nil {
			// Skip types that error, but don't fail completely
			continue
		}

		if len(cheatsheets) > 0 {
			result.TypesWithSheets[t] = cheatsheets
			result.TotalCount += len(cheatsheets)
		}
	}

	return result, nil
}

// GetCheatsheetPath returns the full path to a specific cheatsheet file
func GetCheatsheetPath(repoPath, typeDir, name string) string {
	cheatsheetsPath := GetCheatsheetsPath(repoPath)
	// Add .md extension if not present
	if !strings.HasSuffix(name, ".md") {
		name = name + ".md"
	}
	return filepath.Join(cheatsheetsPath, typeDir, name)
}

// stripFrontmatter removes YAML frontmatter from markdown content
func stripFrontmatter(content string) string {
	lines := strings.Split(content, "\n")

	// Check if the file starts with frontmatter delimiter
	if len(lines) < 3 || strings.TrimSpace(lines[0]) != "---" {
		return content
	}

	// Find the closing delimiter
	for i := 1; i < len(lines); i++ {
		if strings.TrimSpace(lines[i]) == "---" {
			// Return everything after the closing delimiter
			return strings.Join(lines[i+1:], "\n")
		}
	}

	// If no closing delimiter found, return original content
	return content
}

// ShowCheatsheet displays a specific cheatsheet with markdown rendering
func ShowCheatsheet(repoPath, typeDir, name string) error {
	filePath := GetCheatsheetPath(repoPath, typeDir, name)

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("cheatsheet '%s' not found in type '%s'", name, typeDir)
	}

	// Read file content
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading cheatsheet: %w", err)
	}

	// Strip frontmatter and display in viewer
	contentStr := stripFrontmatter(string(content))
	return ShowInViewer(contentStr)
}

// ShowCheatsheetByName finds and displays a cheatsheet by name across all types
func ShowCheatsheetByName(repoPath, name string) error {
	availableTypes, err := GetAvailableTypes(repoPath)
	if err != nil {
		return fmt.Errorf("error getting available types: %w", err)
	}

	// Search for the cheatsheet in all types
	var foundType string
	for _, t := range availableTypes {
		filePath := GetCheatsheetPath(repoPath, t, name)
		if _, err := os.Stat(filePath); err == nil {
			foundType = t
			break
		}
	}

	if foundType == "" {
		return fmt.Errorf("cheatsheet '%s' not found in any type", name)
	}

	return ShowCheatsheet(repoPath, foundType, name)
}

// ShowCheatsheetSelector displays an interactive selector for choosing a cheatsheet
func ShowCheatsheetSelector(repoPath, typeFilter string) error {
	// Get all cheatsheets, optionally filtered by type
	result, err := ListCheatsheets(repoPath, typeFilter)
	if err != nil {
		return err
	}

	if result.TotalCount == 0 {
		if typeFilter != "" {
			return fmt.Errorf("no cheatsheets found for type '%s'", typeFilter)
		}
		return fmt.Errorf("no cheatsheets found")
	}

	// Build a list of cheatsheet options with type prefixes
	type cheatsheetOption struct {
		Display string
		Type    string
		Name    string
	}

	var options []cheatsheetOption
	var types []string
	for t := range result.TypesWithSheets {
		types = append(types, t)
	}
	sort.Strings(types)

	for _, t := range types {
		sheets := result.TypesWithSheets[t]
		for _, sheet := range sheets {
			options = append(options, cheatsheetOption{
				Display: fmt.Sprintf("[%s] %s", t, sheet),
				Type:    t,
				Name:    sheet,
			})
		}
	}

	// Create promptui selector
	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "▸ {{ .Display | cyan }}",
		Inactive: "  {{ .Display }}",
		Selected: "✓ {{ .Display | green }}",
	}

	prompt := promptui.Select{
		Label:     "Select a cheatsheet",
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
	return ShowCheatsheet(repoPath, selected.Type, selected.Name)
}
