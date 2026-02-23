package debugservice

import (
	"fmt"

	"github.com/redjax/cheatsheets/internal/config"
)

// PrintConfig prints the configuration in a debug-friendly format
func PrintConfig(cfg *config.Config) {
	fmt.Println("=== Configuration Debug ===")
	fmt.Printf("Debug mode: %v\n", cfg.Debug)
	fmt.Printf("Sheets path: %v\n", cfg.SheetsPath)
	fmt.Printf("Git Config: %+v\n", cfg.Git)
	fmt.Println("==========================")
}
