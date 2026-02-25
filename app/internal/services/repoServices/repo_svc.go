package reposervices

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
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

// =====================================================
// BRANCH OPERATIONS
// =====================================================

// GetCurrentBranch returns the name of the current branch.
func GetCurrentBranch(path string) (string, error) {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return "", fmt.Errorf("failed to open repository: %w", err)
	}

	head, err := repo.Head()
	if err != nil {
		return "", fmt.Errorf("failed to get HEAD: %w", err)
	}

	return head.Name().Short(), nil
}

// CreateAndCheckoutBranch creates a new branch at the current HEAD and checks it out.
func CreateAndCheckoutBranch(path, branchName string) error {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return fmt.Errorf("failed to open repository: %w", err)
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("failed to get worktree: %w", err)
	}

	// Get current HEAD
	head, err := repo.Head()
	if err != nil {
		return fmt.Errorf("failed to get HEAD: %w", err)
	}

	// Create branch reference
	branchRef := plumbing.NewBranchReferenceName(branchName)
	err = worktree.Checkout(&git.CheckoutOptions{
		Branch: branchRef,
		Create: true,
		Hash:   head.Hash(),
	})
	if err != nil {
		return fmt.Errorf("failed to create and checkout branch: %w", err)
	}

	return nil
}

// BranchExists checks if a branch exists locally.
func BranchExists(path, branchName string) (bool, error) {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return false, fmt.Errorf("failed to open repository: %w", err)
	}

	branchRef := plumbing.NewBranchReferenceName(branchName)
	_, err = repo.Reference(branchRef, true)
	if err == plumbing.ErrReferenceNotFound {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("failed to check branch: %w", err)
	}

	return true, nil
}

// CheckoutBranch switches to an existing branch.
func CheckoutBranch(path, branchName string) error {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return fmt.Errorf("failed to open repository: %w", err)
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("failed to get worktree: %w", err)
	}

	branchRef := plumbing.NewBranchReferenceName(branchName)
	err = worktree.Checkout(&git.CheckoutOptions{
		Branch: branchRef,
	})
	if err != nil {
		return fmt.Errorf("failed to checkout branch: %w", err)
	}

	return nil
}

// EnsureMachineBranch ensures the machine-specific branch exists and is checked out.
// Returns the branch name that was ensured.
func EnsureMachineBranch(path, prefix string) (string, error) {
	// Get hostname
	hostname, err := os.Hostname()
	if err != nil {
		return "", fmt.Errorf("failed to get hostname: %w", err)
	}

	// Clean hostname (remove domain, special chars)
	hostname = strings.Split(hostname, ".")[0]
	hostname = strings.ToLower(hostname)
	branchName := fmt.Sprintf("%s-%s", prefix, hostname)

	// Get current branch
	currentBranch, err := GetCurrentBranch(path)
	if err != nil {
		return "", err
	}

	// If already on the machine branch, we're done
	if currentBranch == branchName {
		return branchName, nil
	}

	// Check if machine branch exists
	exists, err := BranchExists(path, branchName)
	if err != nil {
		return "", err
	}

	if exists {
		// Branch exists, just switch to it
		err = CheckoutBranch(path, branchName)
		if err != nil {
			return "", fmt.Errorf("failed to switch to machine branch: %w", err)
		}
		fmt.Printf("Switched to existing branch '%s'\n", branchName)
	} else {
		// Create and switch to new branch
		err = CreateAndCheckoutBranch(path, branchName)
		if err != nil {
			return "", fmt.Errorf("failed to create machine branch: %w", err)
		}
		fmt.Printf("Created and switched to new branch '%s'\n", branchName)
	}

	return branchName, nil
}

// ListBranches returns a list of local branch names.
func ListBranches(path string) ([]string, error) {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open repository: %w", err)
	}

	branches := []string{}
	refs, err := repo.Branches()
	if err != nil {
		return nil, fmt.Errorf("failed to list branches: %w", err)
	}

	err = refs.ForEach(func(ref *plumbing.Reference) error {
		branches = append(branches, ref.Name().Short())
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to iterate branches: %w", err)
	}

	return branches, nil
}

// =====================================================
// STAGING OPERATIONS
// =====================================================

// StageFiles stages specific files for commit.
func StageFiles(path string, files []string) error {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return fmt.Errorf("failed to open repository: %w", err)
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("failed to get worktree: %w", err)
	}

	for _, file := range files {
		_, err = worktree.Add(file)
		if err != nil {
			return fmt.Errorf("failed to stage file %s: %w", file, err)
		}
	}

	return nil
}

// StageAll stages all changes (modified, new, deleted).
func StageAll(path string) error {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return fmt.Errorf("failed to open repository: %w", err)
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("failed to get worktree: %w", err)
	}

	err = worktree.AddWithOptions(&git.AddOptions{
		All: true,
	})
	if err != nil {
		return fmt.Errorf("failed to stage all changes: %w", err)
	}

	return nil
}

// FileStatus represents the status of a file in the repository.
type FileStatus struct {
	Path     string
	Staging  string // Status in staging area (M, A, D, etc.)
	Worktree string // Status in worktree (M, A, D, etc.)
}

// GetDetailedFileStatus returns detailed status of all files.
func GetDetailedFileStatus(path string) ([]FileStatus, error) {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open repository: %w", err)
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return nil, fmt.Errorf("failed to get worktree: %w", err)
	}

	status, err := worktree.Status()
	if err != nil {
		return nil, fmt.Errorf("failed to get status: %w", err)
	}

	fileStatuses := []FileStatus{}
	for file, stat := range status {
		fileStatuses = append(fileStatuses, FileStatus{
			Path:     file,
			Staging:  string(stat.Staging),
			Worktree: string(stat.Worktree),
		})
	}

	return fileStatuses, nil
}

// GetStagedFiles returns list of staged files.
func GetStagedFiles(path string) ([]string, error) {
	fileStatuses, err := GetDetailedFileStatus(path)
	if err != nil {
		return nil, err
	}

	staged := []string{}
	for _, fs := range fileStatuses {
		if fs.Staging != " " && fs.Staging != "?" {
			staged = append(staged, fs.Path)
		}
	}

	return staged, nil
}

// GetUnstagedFiles returns list of unstaged modified files.
func GetUnstagedFiles(path string) ([]string, error) {
	fileStatuses, err := GetDetailedFileStatus(path)
	if err != nil {
		return nil, err
	}

	unstaged := []string{}
	for _, fs := range fileStatuses {
		if fs.Worktree == "M" || fs.Worktree == "D" {
			unstaged = append(unstaged, fs.Path)
		}
	}

	return unstaged, nil
}

// GetUntrackedFiles returns list of untracked files.
func GetUntrackedFiles(path string) ([]string, error) {
	fileStatuses, err := GetDetailedFileStatus(path)
	if err != nil {
		return nil, err
	}

	untracked := []string{}
	for _, fs := range fileStatuses {
		if fs.Worktree == "?" && fs.Staging == "?" {
			untracked = append(untracked, fs.Path)
		}
	}

	return untracked, nil
}

// =====================================================
// COMMIT OPERATIONS
// =====================================================

// CommitInfo holds information about a commit.
type CommitInfo struct {
	Hash    string
	Author  string
	Date    time.Time
	Message string
}

// CommitChanges creates a commit with staged changes.
func CommitChanges(path, message, authorName, authorEmail string) (string, error) {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return "", fmt.Errorf("failed to open repository: %w", err)
	}

	worktree, err := repo.Worktree()
	if err != nil {
		return "", fmt.Errorf("failed to get worktree: %w", err)
	}

	// Create commit
	hash, err := worktree.Commit(message, &git.CommitOptions{
		Author: &object.Signature{
			Name:  authorName,
			Email: authorEmail,
			When:  time.Now(),
		},
	})
	if err != nil {
		return "", fmt.Errorf("failed to create commit: %w", err)
	}

	return hash.String()[:7], nil
}

// GetLastCommit returns information about the last commit.
func GetLastCommit(path string) (*CommitInfo, error) {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return nil, fmt.Errorf("failed to open repository: %w", err)
	}

	ref, err := repo.Head()
	if err != nil {
		return nil, fmt.Errorf("failed to get HEAD: %w", err)
	}

	commit, err := repo.CommitObject(ref.Hash())
	if err != nil {
		return nil, fmt.Errorf("failed to get commit: %w", err)
	}

	return &CommitInfo{
		Hash:    commit.Hash.String()[:7],
		Author:  commit.Author.Name,
		Date:    commit.Author.When,
		Message: strings.Split(commit.Message, "\n")[0], // First line only
	}, nil
}

// =====================================================
// PUSH/PULL OPERATIONS
// =====================================================

// PushBranch pushes the current branch to remote.
func PushBranch(path, branchName, token string, setUpstream bool) error {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return fmt.Errorf("failed to open repository: %w", err)
	}

	pushOpts := &git.PushOptions{
		Progress: os.Stdout,
	}

	// If token is provided, configure authentication
	if token != "" {
		pushOpts.Auth = &http.BasicAuth{
			Username: token,
			Password: "",
		}
	}

	// Set upstream if requested
	if setUpstream {
		pushOpts.RefSpecs = []config.RefSpec{
			config.RefSpec(fmt.Sprintf("refs/heads/%s:refs/heads/%s", branchName, branchName)),
		}
	}

	err = repo.Push(pushOpts)
	if err == git.NoErrAlreadyUpToDate {
		fmt.Println("Branch is already up to date with remote")
		return nil
	}
	if err != nil {
		return fmt.Errorf("failed to push: %w", err)
	}

	return nil
}

// GetBranchDivergence returns commits ahead/behind remote.
func GetBranchDivergence(path, branchName, token string) (ahead, behind int, err error) {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to open repository: %w", err)
	}

	// Get local branch reference
	localRef, err := repo.Reference(plumbing.NewBranchReferenceName(branchName), true)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get local branch: %w", err)
	}

	// Get remote reference
	remote, err := repo.Remote("origin")
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get remote: %w", err)
	}

	listOpts := &git.ListOptions{}
	if token != "" {
		listOpts.Auth = &http.BasicAuth{
			Username: token,
			Password: "",
		}
	}

	refs, err := remote.List(listOpts)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to list remote references: %w", err)
	}

	var remoteHash plumbing.Hash
	remoteRefName := plumbing.NewRemoteReferenceName("origin", branchName)
	for _, ref := range refs {
		if ref.Name() == remoteRefName {
			remoteHash = ref.Hash()
			break
		}
	}

	// If remote branch doesn't exist, we're only ahead
	if remoteHash.IsZero() {
		// Count local commits
		commits, err := repo.Log(&git.LogOptions{
			From: localRef.Hash(),
		})
		if err != nil {
			return 0, 0, fmt.Errorf("failed to get log: %w", err)
		}

		count := 0
		err = commits.ForEach(func(c *object.Commit) error {
			count++
			return nil
		})
		return count, 0, err
	}

	// Count commits between local and remote
	localCommit, err := repo.CommitObject(localRef.Hash())
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get local commit: %w", err)
	}

	remoteCommit, err := repo.CommitObject(remoteHash)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get remote commit: %w", err)
	}

	// Simple comparison - if hashes match, no divergence
	if localCommit.Hash == remoteCommit.Hash {
		return 0, 0, nil
	}

	// Count commits ahead (from local that aren't in remote)
	commits, err := repo.Log(&git.LogOptions{
		From: localCommit.Hash,
	})
	if err != nil {
		return 0, 0, fmt.Errorf("failed to get log: %w", err)
	}

	ahead = 0
	behind = 0
	foundRemote := false

	err = commits.ForEach(func(c *object.Commit) error {
		if c.Hash == remoteCommit.Hash {
			foundRemote = true
			return fmt.Errorf("stop") // Use error to break
		}
		if !foundRemote {
			ahead++
		}
		return nil
	})

	// Count commits behind (from remote that aren't in local)
	if !foundRemote {
		commits, err = repo.Log(&git.LogOptions{
			From: remoteCommit.Hash,
		})
		if err != nil {
			return 0, 0, fmt.Errorf("failed to get remote log: %w", err)
		}

		err = commits.ForEach(func(c *object.Commit) error {
			if c.Hash == localCommit.Hash {
				return fmt.Errorf("stop")
			}
			behind++
			return nil
		})
	}

	return ahead, behind, nil
}

// =====================================================
// STATUS & DIAGNOSTICS
// =====================================================

// DetailedStatus extends RepositoryStatus with more details.
type DetailedStatus struct {
	RepositoryStatus             // Embed existing status
	StagedFiles      []string    `json:"staged_files"`
	UnstagedFiles    []string    `json:"unstaged_files"`
	UntrackedFiles   []string    `json:"untracked_files"`
	AheadBy          int         `json:"ahead_by"`
	BehindBy         int         `json:"behind_by"`
	LastCommit       *CommitInfo `json:"last_commit,omitempty"`
	IsOnMainBranch   bool        `json:"is_on_main_branch"`
}

// GetDetailedStatus returns comprehensive status information.
func GetDetailedStatus(path, token string) (*DetailedStatus, error) {
	status := &DetailedStatus{}

	// Get basic status
	basicStatus := ValidateRepository(path, token)
	status.RepositoryStatus = basicStatus

	if !basicStatus.IsCloned {
		return status, nil
	}

	if basicStatus.Error != nil {
		return status, basicStatus.Error
	}

	// Check if on main branch
	status.IsOnMainBranch = (basicStatus.CurrentBranch == "main" || basicStatus.CurrentBranch == "master")

	// Get file statuses
	staged, err := GetStagedFiles(path)
	if err != nil {
		return status, fmt.Errorf("failed to get staged files: %w", err)
	}
	status.StagedFiles = staged

	unstaged, err := GetUnstagedFiles(path)
	if err != nil {
		return status, fmt.Errorf("failed to get unstaged files: %w", err)
	}
	status.UnstagedFiles = unstaged

	untracked, err := GetUntrackedFiles(path)
	if err != nil {
		return status, fmt.Errorf("failed to get untracked files: %w", err)
	}
	status.UntrackedFiles = untracked

	// Get last commit
	lastCommit, err := GetLastCommit(path)
	if err == nil {
		status.LastCommit = lastCommit
	}

	// Get branch divergence
	ahead, behind, err := GetBranchDivergence(path, basicStatus.CurrentBranch, token)
	if err == nil {
		status.AheadBy = ahead
		status.BehindBy = behind
	}

	return status, nil
}

// GetGitConfig reads git configuration values.
func GetGitConfig(key string) (string, error) {
	cmd := exec.Command("git", "config", "--get", key)
	output, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// GetGitAuthor returns the configured git author name and email.
// Falls back to system git config if not provided.
func GetGitAuthor(configName, configEmail string) (name, email string, err error) {
	name = configName
	email = configEmail

	// If name not provided, try to get from git config
	if name == "" {
		name, err = GetGitConfig("user.name")
		if err != nil || name == "" {
			return "", "", fmt.Errorf("git author name not configured. Set it with: git config --global user.name 'Your Name'")
		}
	}

	// If email not provided, try to get from git config
	if email == "" {
		email, err = GetGitConfig("user.email")
		if err != nil || email == "" {
			return "", "", fmt.Errorf("git author email not configured. Set it with: git config --global user.email 'you@example.com'")
		}
	}

	return name, email, nil
}
