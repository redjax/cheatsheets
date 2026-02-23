package cleanupcommand

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var dryRun bool

// CleanupCmd represents the cleanup command
var CleanupCmd = &cobra.Command{
	Use:     "cleanup",
	Short:   "Clean up temporary editing files. Alias: gc",
	Long:    `Remove the temporary directory used for editing cheatsheets ($TMPDIR/cheatsheets/).`,
	Aliases: []string{"gc"},
	RunE:    runCleanup,
}

func init() {
	CleanupCmd.Flags().BoolVarP(&dryRun, "dry-run", "n", false, "Show what would be deleted without actually deleting")
}

func runCleanup(cmd *cobra.Command, args []string) error {
	tempDir := filepath.Join(os.TempDir(), "cheatsheets")

	// Check if the directory exists
	info, err := os.Stat(tempDir)
	if os.IsNotExist(err) {
		fmt.Printf("Temp directory does not exist: %s\n", tempDir)
		fmt.Println("Nothing to clean up.")
		return nil
	}
	if err != nil {
		return fmt.Errorf("error checking temp directory: %w", err)
	}

	if !info.IsDir() {
		return fmt.Errorf("expected directory but found file: %s", tempDir)
	}

	// Count files in the directory
	entries, err := os.ReadDir(tempDir)
	if err != nil {
		return fmt.Errorf("error reading temp directory: %w", err)
	}

	if len(entries) == 0 {
		fmt.Printf("Temp directory is already empty: %s\n", tempDir)
		return nil
	}

	if dryRun {
		fmt.Printf("Would remove %d file(s) from: %s\n", len(entries), tempDir)
		for _, entry := range entries {
			fmt.Printf("  - %s\n", entry.Name())
		}
		fmt.Println("\nRun without --dry-run to actually delete these files.")
		return nil
	}

	// Remove the entire directory
	if err := os.RemoveAll(tempDir); err != nil {
		return fmt.Errorf("error removing temp directory: %w", err)
	}

	fmt.Printf("Successfully cleaned up %d file(s) from: %s\n", len(entries), tempDir)
	return nil
}
