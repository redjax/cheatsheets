package main

import (
	"fmt"

	"github.com/redjax/cheatsheets/internal/config"
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
	if err := reposervices.EnsureRepository(cfg.Git.RepoUrl, cfg.Git.ClonePath, cfg.Git.Token); err != nil {
		panic(fmt.Sprintf("Failed to ensure repository: %v", err))
	}

	// Update repository (this is a test, remove later)
	if err := reposervices.UpdateRepository(cfg.Git.ClonePath, cfg.Git.Token); err != nil {
		panic(fmt.Sprintf("Failed to update repository: %v", err))
	}
}

func debugConfig(cfg *config.Config) {
	fmt.Printf("Debug mode: %v\n", cfg.Debug)
	fmt.Printf("Git Config: %+v\n", cfg.Git)
}
