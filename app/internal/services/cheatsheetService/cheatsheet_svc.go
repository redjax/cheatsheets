package cheatsheetservice

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// CheatsheetInfo represents a cheatsheet with its type and name
type CheatsheetInfo struct {
	Type string
	Name string
}

// GetCheatsheetsPath returns the path to the cheatsheets directory
func GetCheatsheetsPath(repoPath string) string {
	return filepath.Join(repoPath, "cheatsheets")
}

// ValidateCheatsheetsDirectory checks if the cheatsheets directory exists
func ValidateCheatsheetsDirectory(repoPath string) error {
	cheatsheetsPath := GetCheatsheetsPath(repoPath)
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
