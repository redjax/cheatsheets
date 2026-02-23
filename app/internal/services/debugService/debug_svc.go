package debugservice

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/redjax/cheatsheets/internal/config"
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
		Token:     maskToken(cfg.Git.Token),
	}

	gitJSON, _ := json.MarshalIndent(displayGit, "", "  ")
	fmt.Printf("Git Config:\n%s\n", string(gitJSON))
	fmt.Println("==========================")
}

// maskToken masks a git token while showing the prefix and first 7 characters of the secret
func maskToken(token string) string {
	if token == "" {
		return "<empty>"
	}

	// Common git forge token prefixes
	prefixes := []string{
		"github_pat_", // GitHub fine-grained PAT
		"ghp_",        // GitHub personal access token
		"gho_",        // GitHub OAuth token
		"ghu_",        // GitHub user-to-server token
		"ghs_",        // GitHub server-to-server token
		"ghr_",        // GitHub refresh token
		"glpat-",      // GitLab personal access token
		"gloas-",      // GitLab OAuth application secret
		"glptt-",      // GitLab project access token
	}

	// Find matching prefix
	var prefix string
	secretStart := 0

	for _, p := range prefixes {
		if strings.HasPrefix(token, p) {
			prefix = p
			secretStart = len(p)
			break
		}
	}

	// If no known prefix, treat entire token as secret
	if prefix == "" {
		if len(token) <= 7 {
			return "***"
		}

		return token[:7] + "***"
	}

	// Show prefix + first 7 chars of secret
	secretPart := token[secretStart:]
	if len(secretPart) <= 7 {
		return prefix + "***"
	}

	return prefix + secretPart[:7] + "***"
}
