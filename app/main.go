package main

import (
	"fmt"

	"github.com/redjax/cheatsheets/internal/config"
)

func main() {
	// Load configuration (checks for config.local.yml, falls back to config.yml)
	configFile := config.FindConfigFile("config.yml")
	cfg, err := config.LoadConfig(nil, configFile)
	if err != nil {
		panic(fmt.Sprintf("Failed to load config: %v", err))
	}

	fmt.Printf("Config loaded: %+v\n", cfg)
}
