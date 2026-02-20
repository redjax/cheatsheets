package main

import (
	"fmt"

	"github.com/redjax/cheatsheets/internal/config"
	"github.com/redjax/cheatsheets/internal/constants"
	reposervices "github.com/redjax/cheatsheets/internal/services/repoServices"
)

func main() {
	// Load configuration (checks for config.local.yml, falls back to config.yml)
	configFile := config.FindConfigFile("config.yml")
	cfg, err := config.LoadConfig(nil, configFile)
	if err != nil {
		panic(fmt.Sprintf("Failed to load config: %v", err))
	}

	if cfg.Debug {
		debugConfig(cfg)
	}

	// Ensure repository is cloned
	if err := reposervices.EnsureRepository(constants.RepoURL, cfg.Git.ClonePath); err != nil {
		panic(fmt.Sprintf("Failed to ensure repository: %v", err))
	}
}

func debugConfig(cfg *config.Config) {
	fmt.Printf("Repository URL: %v\n", constants.RepoURL)
	fmt.Printf("Debug mode: %v\n", cfg.Debug)
	fmt.Printf("App Config: %+v\n", cfg.Git)
}
