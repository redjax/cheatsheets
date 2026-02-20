package main

import (
	"fmt"

	"github.com/redjax/cheatsheets/internal/config"
	"github.com/redjax/cheatsheets/internal/constants"
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
}

func debugConfig(cfg *config.Config) {
	fmt.Printf("Repository URL: %v\n", constants.RepoURL)
	fmt.Printf("Debug mode: %v\n", cfg.Debug)
	fmt.Printf("App Config: %+v\n", cfg.App)
}
