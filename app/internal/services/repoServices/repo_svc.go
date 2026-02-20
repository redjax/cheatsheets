package reposervices

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

// IsRepositoryCloned checks if the repository is already cloned at the given path.
func IsRepositoryCloned(path string) (bool, error) {
	// Check if dir already exists
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return false, fmt.Errorf("error checking path: %w", err)
	}

	// Check if .git dir exists inside the path
	gitDir := filepath.Join(path, ".git")
	if _, err := os.Stat(gitDir); os.IsNotExist(err) {
		return false, nil
	}

	// Verify it's a valid git repository
	_, err := git.PlainOpen(path)
	if err != nil {
		return false, fmt.Errorf("directory exists but is not a valid git repository: %w", err)
	}

	return true, nil
}

// CloneRepository clones the repository from the given URL to the specified path.
// If a token is provided, it will be used for authentication.
func CloneRepository(url, path, token string) error {
	fmt.Printf("Cloning repository from %s to %s\n", url, path)

	cloneOpts := &git.CloneOptions{
		URL:      url,
		Progress: os.Stdout,
	}

	// If token is provided, configure authentication
	if token != "" {
		cloneOpts.Auth = &http.BasicAuth{
			Username: token, // GitHub uses token as username
			Password: "",    // Password should be empty
		}
	}

	_, err := git.PlainClone(path, false, cloneOpts)

	if err != nil {
		return fmt.Errorf("failed to clone repository: %w", err)
	}

	fmt.Println("Repository cloned successfully")
	return nil
}

// EnsureRepository checks if the repository is already cloned at the given path,
// and if not, tests that the repo URL is valid, then clones it to the path.
// If a token is provided, it will be used for authentication during cloning.
func EnsureRepository(url, path, token string) error {
	cloned, err := IsRepositoryCloned(path)
	if err != nil {
		return fmt.Errorf("error checking repository: %w", err)
	}

	if cloned {
		fmt.Printf("Repository already exists at %s\n", path)
		return nil
	}

	// Repository doesn't exist, clone it
	return CloneRepository(url, path, token)
}
