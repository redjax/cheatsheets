package upgradecommand

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/redjax/cheatsheets/internal/version"
	"github.com/spf13/cobra"
)

var (
	checkOnly bool
	force     bool
)

// GitHubRelease represents a GitHub release
type GitHubRelease struct {
	TagName string `json:"tag_name"`
	Name    string `json:"name"`
	Assets  []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

// UpgradeCmd represents the upgrade command
var UpgradeCmd = &cobra.Command{
	Use:     "upgrade",
	Aliases: []string{"update"},
	Short:   "Upgrade chtsht to the latest version",
	Long:    `Check for and install the latest version of chtsht from GitHub releases.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		currentVersion := version.GetVersion()
		fmt.Printf("Current version: %s (commit: %s, built: %s)\n",
			currentVersion, version.GetCommit(), version.GetDate())

		// Check if this is a dev build
		if currentVersion == "dev" || strings.Contains(currentVersion, "dev") {
			fmt.Println("⚠️  Running a development build")
			if !force {
				fmt.Println("   Use --force to upgrade a dev build, or rebuild from source")
				return nil
			}
			fmt.Println("   --force specified, proceeding with upgrade check...")
		}

		// Fetch latest release from GitHub
		release, err := getLatestRelease()
		if err != nil {
			return fmt.Errorf("failed to check for updates: %w", err)
		}

		fmt.Printf("Latest version:  %s\n", release.TagName)

		if !force {
			// Use semantic version comparison
			cmp := version.CompareVersions(currentVersion, release.TagName)

			switch cmp {
			case 0:
				fmt.Println("✓ You are already running the latest version")
				return nil
			case 1:
				fmt.Printf("🤯 You're ahead of the latest release: current=%s, latest=%s\n",
					currentVersion, release.TagName)
				return nil
			case -1:
				fmt.Printf("\n🚀 Upgrade available: %s → %s\n", currentVersion, release.TagName)
			}
		} else {
			fmt.Println("\n⚡ --force specified, skipping version comparison")
		}

		if checkOnly {
			fmt.Println("\n✅ Run without --check to upgrade")
			return nil
		}

		// Find the appropriate asset for the current platform
		assetURL, assetName, err := findAsset(release)
		if err != nil {
			// Print available assets to help debug
			fmt.Fprintln(os.Stderr, "\nAvailable assets:")
			for _, asset := range release.Assets {
				fmt.Fprintf(os.Stderr, "  - %s\n", asset.Name)
			}
			return err
		}

		// Download and install
		fmt.Printf("\nDownloading %s...\n", assetName)
		if err := downloadAndInstall(assetURL); err != nil {
			return fmt.Errorf("failed to upgrade: %w", err)
		}

		fmt.Printf("\n✓ Successfully upgraded to %s\n", release.TagName)
		return nil
	},
}

func init() {
	UpgradeCmd.Flags().BoolVar(&checkOnly, "check", false, "Check for updates without installing")
	UpgradeCmd.Flags().BoolVar(&force, "force", false, "Force upgrade even if versions match or on dev builds")
}

// getLatestRelease fetches the latest release from GitHub
func getLatestRelease() (*GitHubRelease, error) {
	apiURL := "https://api.github.com/repos/redjax/cheatsheets/releases/latest"

	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	var release GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, err
	}

	return &release, nil
}

// normalizeOS maps runtime.GOOS to the naming used in goreleaser archives.
func normalizeOS(goos string) string {
	switch goos {
	case "darwin":
		return "macOS"
	default:
		return goos // "linux", "windows"
	}
}

// findAsset locates the correct zip asset for the current platform from a release.
// Goreleaser produces archives named: chtsht-{os}-{arch}-{version}.zip
// e.g. chtsht-linux-amd64-1.2.3.zip, chtsht-macOS-arm64-1.2.3.zip
func findAsset(release *GitHubRelease) (downloadURL, assetName string, err error) {
	osName := normalizeOS(runtime.GOOS)
	arch := runtime.GOARCH

	// Build the expected prefix: "chtsht-linux-amd64-" or "chtsht-macOS-arm64-"
	expectedPrefix := fmt.Sprintf("chtsht-%s-%s-", osName, arch)

	for _, asset := range release.Assets {
		if asset.Name == "" {
			continue
		}

		// Case-insensitive prefix match + .zip suffix
		nameLower := strings.ToLower(asset.Name)
		prefixLower := strings.ToLower(expectedPrefix)

		if strings.HasPrefix(nameLower, prefixLower) && strings.HasSuffix(nameLower, ".zip") {
			return asset.BrowserDownloadURL, asset.Name, nil
		}
	}

	return "", "", fmt.Errorf("no release asset found for %s/%s (expected prefix: %s*.zip)",
		runtime.GOOS, runtime.GOARCH, expectedPrefix)
}

// downloadAndInstall downloads a zip archive, extracts the binary, verifies it,
// and replaces the current executable in-place.
func downloadAndInstall(url string) error {
	// Download zip to temporary file
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("download failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed with status %d", resp.StatusCode)
	}

	zipTmp, err := os.CreateTemp("", "chtsht-upgrade-*.zip")
	if err != nil {
		return fmt.Errorf("failed to create temp zip file: %w", err)
	}
	defer os.Remove(zipTmp.Name())

	if _, err := io.Copy(zipTmp, resp.Body); err != nil {
		zipTmp.Close()
		return fmt.Errorf("failed to write zip file: %w", err)
	}
	zipTmp.Close()

	// Extract binary from zip
	fmt.Println("Extracting binary from archive...")
	binaryTmpPath, err := extractBinaryFromZip(zipTmp.Name())
	if err != nil {
		return fmt.Errorf("failed to extract binary: %w", err)
	}
	defer os.Remove(binaryTmpPath)

	// Get current executable path (resolve symlinks)
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("failed to get executable path: %w", err)
	}
	exePath, err = filepath.EvalSymlinks(exePath)
	if err != nil {
		return fmt.Errorf("failed to resolve executable path: %w", err)
	}

	// Create backup of current binary
	backupPath := exePath + ".bak"
	fmt.Printf("Backing up current binary to %s\n", backupPath)
	if err := copyFile(exePath, backupPath); err != nil {
		return fmt.Errorf("failed to create backup: %w", err)
	}

	// Replace the binary (platform-specific)
	if runtime.GOOS == "windows" {
		if err := replaceWindows(exePath, binaryTmpPath); err != nil {
			fmt.Fprintln(os.Stderr, "Restoring backup after failed install...")
			restoreErr := os.Rename(backupPath, exePath)
			if restoreErr != nil {
				fmt.Fprintf(os.Stderr, "⚠️  Failed to restore backup: %v\n", restoreErr)
			}
			return fmt.Errorf("failed to install new binary: %w", err)
		}
	} else {
		// Unix: os.Rename works because the old inode stays alive while the process runs
		if err := os.Rename(binaryTmpPath, exePath); err != nil {
			return fmt.Errorf("failed to install new binary: %w", err)
		}
	}

	// Verify the new binary actually works
	fmt.Println("Verifying new binary...")
	if err := verifyBinary(exePath); err != nil {
		fmt.Fprintf(os.Stderr, "⚠️  Verification failed: %v\n", err)
		fmt.Fprintln(os.Stderr, "Rolling back to previous version...")

		rollbackErr := os.Rename(backupPath, exePath)
		if rollbackErr != nil {
			return fmt.Errorf("rollback also failed: %w (original error: %v)", rollbackErr, err)
		}

		fmt.Fprintln(os.Stderr, "✓ Rolled back successfully")
		return fmt.Errorf("upgrade aborted: new binary failed verification: %w", err)
	}

	// Clean up backup after successful verification
	os.Remove(backupPath)

	return nil
}

// replaceWindows handles binary replacement on Windows where the running exe is locked.
// It moves the old binary out of the way, then copies the new one in.
func replaceWindows(exePath, newBinaryPath string) error {
	oldPath := exePath + ".old"

	// Remove any stale .old file from a previous upgrade
	os.Remove(oldPath)

	// Move current exe to .old (Windows allows renaming a running exe)
	if err := os.Rename(exePath, oldPath); err != nil {
		return fmt.Errorf("failed to move old binary: %w", err)
	}

	// Copy new binary into place
	if err := copyFile(newBinaryPath, exePath); err != nil {
		// Try to restore the old binary
		os.Rename(oldPath, exePath)
		return fmt.Errorf("failed to copy new binary: %w", err)
	}

	// Schedule cleanup of .old — best effort
	defer os.Remove(oldPath)

	return nil
}

// extractBinaryFromZip extracts the chtsht binary from a zip archive.
func extractBinaryFromZip(zipPath string) (string, error) {
	r, err := zip.OpenReader(zipPath)
	if err != nil {
		return "", fmt.Errorf("failed to open zip: %w", err)
	}
	defer r.Close()

	expectedName := "chtsht"
	if runtime.GOOS == "windows" {
		expectedName = "chtsht.exe"
	}

	var binaryFile *zip.File
	for _, f := range r.File {
		if f.FileInfo().IsDir() {
			continue
		}
		// Match by filename (may be at root or in a subdirectory)
		fileName := strings.ToLower(filepath.Base(f.Name))
		if fileName == expectedName {
			binaryFile = f
			break
		}
	}

	if binaryFile == nil {
		return "", fmt.Errorf("binary '%s' not found in zip archive", expectedName)
	}

	rc, err := binaryFile.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open binary in zip: %w", err)
	}
	defer rc.Close()

	tmpBin, err := os.CreateTemp("", "chtsht-bin-*")
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}

	// Limit extraction size to 500MB to prevent decompression bombs
	limitedReader := io.LimitReader(rc, 500*1024*1024)
	if _, err := io.Copy(tmpBin, limitedReader); err != nil {
		tmpBin.Close()
		return "", fmt.Errorf("failed to extract binary: %w", err)
	}
	tmpBin.Close()

	// Make executable (no-op on Windows, needed for Unix)
	if err := os.Chmod(tmpBin.Name(), 0755); err != nil {
		return "", fmt.Errorf("failed to chmod binary: %w", err)
	}

	return tmpBin.Name(), nil
}

// copyFile copies a file from src to dst, preserving permissions.
func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	if _, err := io.Copy(destFile, sourceFile); err != nil {
		return err
	}

	sourceInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	return os.Chmod(dst, sourceInfo.Mode())
}

// verifyBinary runs `<binary> self version` to confirm the new binary is functional.
func verifyBinary(path string) error {
	cmd := exec.Command(path, "self", "version")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("binary at %s failed to run: %w (output: %s)", path, err, strings.TrimSpace(string(output)))
	}
	return nil
}
