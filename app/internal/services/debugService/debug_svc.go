package debugservice

import (
	"encoding/json"
	"fmt"

	"github.com/redjax/cheatsheets/internal/config"
	"github.com/redjax/cheatsheets/internal/utils"
)

// PrintConfig prints the configuration in a debug-friendly format
func PrintConfig(cfg *config.Config) {
	fmt.Println("=== Configuration Debug ===")
	fmt.Printf("Debug mode: %v\n", cfg.Debug)
	fmt.Printf("Sheets path: %v\n", cfg.SheetsPath)

	// Create a display version with masked token
	displayGit := struct {
		RepoUrl   string `json:"repo_url"`
		ClonePath string `json:"clone_path"`
		Token     string `json:"token"`
	}{
		RepoUrl:   cfg.Git.RepoUrl,
		ClonePath: cfg.Git.ClonePath,
		Token:     utils.MaskToken(cfg.Git.Token),
	}

	gitJSON, _ := json.MarshalIndent(displayGit, "", "  ")
	fmt.Printf("Git Config:\n%s\n", string(gitJSON))
	fmt.Println("==========================")
}
