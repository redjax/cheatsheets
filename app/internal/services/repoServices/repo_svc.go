package reposervices

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
)

// GetRepositoryStatus returns detailed status information about the repository.
type RepositoryStatus struct {
	IsCloned         bool
	IsClean          bool
	HasRemoteUpdates bool
	CurrentBranch    string
	Error            error
}

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

// UpdateRepository pulls the latest changes from the remote repository.
// If a token is provided, it will be used for authentication.
func UpdateRepository(path, token string) error {
	fmt.Printf("Updating repository at %s\n", path)

	// Open the existing repository
	repo, err := git.PlainOpen(path)
	if err != nil {
		return fmt.Errorf("failed to open repository: %w", err)
	}

	// Get remote to check URL format
	remote, err := repo.Remote("origin")
	if err != nil {
		return fmt.Errorf("failed to get remote: %w", err)
	}

	// Get the remote config
	config := remote.Config()
	if len(config.URLs) > 0 {
		fmt.Printf("Remote URL: %s\n", config.URLs[0])
	}

	// Get the working tree
	worktree, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("failed to get worktree: %w", err)
	}

	// Configure pull options
	pullOpts := &git.PullOptions{
		RemoteName: "origin",
		Progress:   os.Stdout,
		// Force HTTPS-based pull, don't use remote URL
		Force: false,
	}

	// If token is provided, configure authentication
	if token != "" {
		pullOpts.Auth = &http.BasicAuth{
			Username: token, // GitHub uses token as username
			Password: "",    // Password should be empty
		}
	}

	// Pull latest changes
	err = worktree.Pull(pullOpts)
	if err == git.NoErrAlreadyUpToDate {
		fmt.Println("Repository is already up to date")
		return nil
	}
	if err != nil {
		return fmt.Errorf("failed to pull updates: %w", err)
	}

	fmt.Println("Repository updated successfully")
	return nil
}

// IsWorkingTreeClean checks if the repository has uncommitted changes.
// Returns true if the working tree is clean (no changes), false if there are uncommitted changes.
func IsWorkingTreeClean(path string) (bool, error) {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return false, fmt.Errorf("failed to open repository: %w", err)
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return false, fmt.Errorf("failed to get worktree: %w", err)
	}

	status, err := worktree.Status()
	if err != nil {
		return false, fmt.Errorf("failed to get working tree status: %w", err)
	}

	return status.IsClean(), nil
}

// CheckForRemoteUpdates checks if there are updates available from the remote repository.
// Returns true if remote has new commits, false if local is up to date.
func CheckForRemoteUpdates(path, token string) (bool, error) {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return false, fmt.Errorf("failed to open repository: %w", err)
	}

	// Get the remote
	remote, err := repo.Remote("origin")
	if err != nil {
		return false, fmt.Errorf("failed to get remote: %w", err)
	}

	// Configure list options
	listOpts := &git.ListOptions{}
	if token != "" {
		listOpts.Auth = &http.BasicAuth{
			Username: token,
			Password: "",
		}
	}

	// Fetch remote references without downloading objects
	refs, err := remote.List(listOpts)
	if err != nil {
		return false, fmt.Errorf("failed to list remote references: %w", err)
	}

	// Get local HEAD reference
	head, err := repo.Head()
	if err != nil {
		return false, fmt.Errorf("failed to get HEAD reference: %w", err)
	}

	// Find the remote HEAD for the current branch
	localHash := head.Hash()
	branchName := head.Name()

	var remoteHash string
	for _, ref := range refs {
		if ref.Name() == branchName {
			remoteHash = ref.Hash().String()
			break
		}
	}

	if remoteHash == "" {
		return false, fmt.Errorf("could not find remote reference for branch %s", branchName)
	}

	// Compare hashes
	return localHash.String() != remoteHash, nil
}

// ValidateRepository performs a comprehensive validation of the repository.
func ValidateRepository(path, token string) RepositoryStatus {
	status := RepositoryStatus{}

	// Check if cloned
	cloned, err := IsRepositoryCloned(path)
	if err != nil {
		status.Error = err
		return status
	}
	status.IsCloned = cloned

	if !cloned {
		// If not cloned, nothing else to check
		return status
	}

	// Get current branch
	repo, err := git.PlainOpen(path)
	if err != nil {
		status.Error = fmt.Errorf("failed to open repository: %w", err)
		return status
	}

	head, err := repo.Head()
	if err != nil {
		status.Error = fmt.Errorf("failed to get HEAD: %w", err)
		return status
	}
	status.CurrentBranch = head.Name().Short()

	// Check if working tree is clean
	clean, err := IsWorkingTreeClean(path)
	if err != nil {
		status.Error = err
		return status
	}
	status.IsClean = clean

	// Check for remote updates
	hasUpdates, err := CheckForRemoteUpdates(path, token)
	if err != nil {
		status.Error = err
		return status
	}
	status.HasRemoteUpdates = hasUpdates

	return status
}
