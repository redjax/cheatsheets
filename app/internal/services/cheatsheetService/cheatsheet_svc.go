package cheatsheetservice

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

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

	// Search for the cheatsheet in all types and collect all matches
	var foundTypes []string
	for _, t := range availableTypes {
		filePath := GetCheatsheetPath(repoPath, t, name)
		if _, err := os.Stat(filePath); err == nil {
			foundTypes = append(foundTypes, t)
		}
	}

	if len(foundTypes) == 0 {
		return fmt.Errorf("cheatsheet '%s' not found in any type", name)
	}

	// If only one match, show it directly
	if len(foundTypes) == 1 {
		return ShowCheatsheet(repoPath, foundTypes[0], name)
	}

	// Multiple matches - let user choose
	fmt.Printf("Multiple cheatsheets named '%s' found:\n\n", name)

	type cheatsheetOption struct {
		Display string
		Type    string
	}

	var options []cheatsheetOption
	for _, t := range foundTypes {
		options = append(options, cheatsheetOption{
			Display: fmt.Sprintf("[%s] %s", t, name),
			Type:    t,
		})
	}

	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}",
		Active:   "▸ {{ .Display | cyan }}",
		Inactive: "  {{ .Display }}",
		Selected: "✓ {{ .Display | green }}",
	}

	prompt := promptui.Select{
		Label:     "Select which cheatsheet to view",
		Items:     options,
		Templates: templates,
		Size:      len(options),
	}

	idx, _, err := prompt.Run()
	if err != nil {
		return fmt.Errorf("selection cancelled")
	}

	selected := options[idx]
	fmt.Println() // Add blank line for spacing
	return ShowCheatsheet(repoPath, selected.Type, name)
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

// CreateCheatsheetOptions contains options for creating a new cheatsheet
type CreateCheatsheetOptions struct {
	Type        string
	Name        string
	Title       string
	Description string
	Tags        string
}

// CreateCheatsheet creates a new cheatsheet from a template
func CreateCheatsheet(repoPath string, opts CreateCheatsheetOptions) (string, error) {
	// Validate type
	validTypes := []string{"app", "command", "language", "system"}
	isValidType := false
	for _, vt := range validTypes {
		if opts.Type == vt {
			isValidType = true
			break
		}
	}
	if !isValidType {
		return "", fmt.Errorf("invalid type '%s'. Must be one of: %s", opts.Type, strings.Join(validTypes, ", "))
	}

	// Validate name
	if opts.Name == "" {
		return "", fmt.Errorf("cheatsheet name cannot be empty")
	}

	// Build paths
	templatePath := filepath.Join(repoPath, ".templates", opts.Type+".md")
	cheatsheetsPath := GetCheatsheetsPath(repoPath)
	targetDir := filepath.Join(cheatsheetsPath, opts.Type)
	targetFile := filepath.Join(targetDir, opts.Name+".md")

	// Check if template exists
	if _, err := os.Stat(templatePath); os.IsNotExist(err) {
		return "", fmt.Errorf("template not found: %s", templatePath)
	}

	// Check if target file already exists
	if _, err := os.Stat(targetFile); err == nil {
		return "", fmt.Errorf("cheatsheet already exists: %s", targetFile)
	}

	// Read template
	templateContent, err := os.ReadFile(templatePath)
	if err != nil {
		return "", fmt.Errorf("failed to read template: %w", err)
	}

	// Prepare replacements
	content := string(templateContent)
	replacements := map[string]string{
		"{{title}}":        opts.Title,
		"{{description}}":  opts.Description,
		"{{last_updated}}": time.Now().Format("2006-01-02"),
	}

	// Format tags
	var formattedTags string
	if opts.Tags != "" {
		tagList := strings.Split(opts.Tags, ",")
		quotedTags := make([]string, 0, len(tagList))
		for _, tag := range tagList {
			tag = strings.TrimSpace(tag)
			if tag != "" {
				quotedTags = append(quotedTags, fmt.Sprintf(`"%s"`, tag))
			}
		}
		formattedTags = strings.Join(quotedTags, ", ")
	}
	replacements["{{tags}}"] = formattedTags

	// Apply replacements
	for placeholder, value := range replacements {
		content = strings.ReplaceAll(content, placeholder, value)
	}

	// Create target directory if it doesn't exist
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create target directory: %w", err)
	}

	// Write the new file
	if err := os.WriteFile(targetFile, []byte(content), 0644); err != nil {
		return "", fmt.Errorf("failed to write cheatsheet file: %w", err)
	}

	return targetFile, nil
}

// DeleteCheatsheet deletes a cheatsheet file
func DeleteCheatsheet(repoPath, cheatsheetType, name string) error {
	filePath := GetCheatsheetPath(repoPath, cheatsheetType, name)

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return fmt.Errorf("cheatsheet not found: %s", filePath)
	}

	// Delete the file
	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("failed to delete cheatsheet: %w", err)
	}

	return nil
}
