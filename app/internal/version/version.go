package version

import (
	"fmt"
	"strings"
)

// Build-time variables injected via ldflags:
//
//	go build -ldflags "\
//	  -X 'github.com/redjax/cheatsheets/internal/version.Version=v1.0.0' \
//	  -X 'github.com/redjax/cheatsheets/internal/version.Commit=abc1234' \
//	  -X 'github.com/redjax/cheatsheets/internal/version.Date=2025-01-01T00:00:00Z'"
var (
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
)

// GetVersion returns the current version string
func GetVersion() string {
	if Version == "" {
		return "dev"
	}
	return Version
}

// GetCommit returns the build commit hash
func GetCommit() string {
	if Commit == "" {
		return "none"
	}
	return Commit
}

// GetDate returns the build date
func GetDate() string {
	if Date == "" {
		return "unknown"
	}
	return Date
}

// CompareVersions compares two semver-like version strings.
// Returns -1 if v1 < v2, 0 if equal, 1 if v1 > v2.
// Strips leading "v" prefixes and trailing suffixes after "-" (e.g. "-rc1").
func CompareVersions(version1, version2 string) int {
	// Strip 'v' prefix
	v1 := strings.TrimPrefix(version1, "v")
	v2 := strings.TrimPrefix(version2, "v")

	// Strip suffixes like -rc1, -beta, -abcdef123
	v1 = strings.SplitN(v1, "-", 2)[0]
	v2 = strings.SplitN(v2, "-", 2)[0]

	parts1 := strings.Split(v1, ".")
	parts2 := strings.Split(v2, ".")

	maxLen := len(parts1)
	if len(parts2) > maxLen {
		maxLen = len(parts2)
	}

	for i := 0; i < maxLen; i++ {
		var n1, n2 int
		if i < len(parts1) {
			fmt.Sscanf(parts1[i], "%d", &n1)
		}
		if i < len(parts2) {
			fmt.Sscanf(parts2[i], "%d", &n2)
		}
		if n1 > n2 {
			return 1
		}
		if n1 < n2 {
			return -1
		}
	}
	return 0
}
