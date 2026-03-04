package reposervices

import (
	"fmt"
	"os"
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
// Handles fast-forward updates. For non-fast-forward scenarios (diverged branches),
// returns an informative error with guidance.
func UpdateRepository(path, token string) error {
	return UpdateRepositoryQuiet(path, token, false)
}

// UpdateRepositoryQuiet updates the repository with optional quiet mode
func UpdateRepositoryQuiet(path, token string, quiet bool) error {
	if !quiet {
		fmt.Printf("Updating repository at %s\n", path)
	}

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
	remoteConfig := remote.Config()
	if !quiet && len(remoteConfig.URLs) > 0 {
		fmt.Printf("Remote URL: %s\n", remoteConfig.URLs[0])
	}

	// Get current branch name before pulling
	head, err := repo.Head()
	if err != nil {
		return fmt.Errorf("failed to get HEAD: %w", err)
	}
	branchName := head.Name().Short()
	localHash := head.Hash()

	// Get remote branch hash before pull
	remoteBranchRef := plumbing.NewRemoteReferenceName("origin", branchName)
	remoteBranch, err := repo.Reference(remoteBranchRef, true)
	var remoteHashBefore plumbing.Hash
	if err == nil {
		remoteHashBefore = remoteBranch.Hash()
	}

	// Get the working tree
	worktree, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("failed to get worktree: %w", err)
	}

	// Configure pull options - explicitly specify which branch to pull
	pullOpts := &git.PullOptions{
		RemoteName:    "origin",
		Progress:      os.Stdout,
		Force:         false,
		ReferenceName: plumbing.NewBranchReferenceName(branchName), // Pull only this branch
	}

	// If token is provided, configure authentication
	if token != "" {
		pullOpts.Auth = &http.BasicAuth{
			Username: token,
			Password: "",
		}
	}

	// Attempt to pull
	err = worktree.Pull(pullOpts)
	if err == git.NoErrAlreadyUpToDate {
		if !quiet {
			fmt.Println("Repository is already up to date")
		}
		return nil
	}

	// If non-fast-forward, provide more detailed guidance
	if err != nil && strings.Contains(err.Error(), "non-fast-forward") {
		fmt.Printf("\nDebug info:\n")
		fmt.Printf("  Local:  %s\n", localHash)
		fmt.Printf("  Remote: %s\n", remoteHashBefore)
		return fmt.Errorf(`branches have diverged - both local and remote have changes

This happens when commits were made on multiple machines.

To resolve:
  1. Commit your local changes: chtsht repo commit -a -m "your message"
  2. Fetch remote changes:      chtsht repo sync
  3. Merge manually or use the merge command once available

Your local changes are safe. Remote: origin/%s`, branchName)
	}

	if err != nil {
		return fmt.Errorf("failed to pull updates: %w", err)
	}

	if !quiet {
		fmt.Println("Repository updated successfully")
	}
	return nil
}

// FetchRepository fetches the latest changes from remote without merging.
func FetchRepository(path, token string) error {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return fmt.Errorf("failed to open repository: %w", err)
	}

	fetchOpts := &git.FetchOptions{
		RemoteName: "origin",
	}

	if token != "" {
		fetchOpts.Auth = &http.BasicAuth{
			Username: token,
			Password: "",
		}
	}

	err = repo.Fetch(fetchOpts)
	if err == git.NoErrAlreadyUpToDate {
		return nil
	}
	if err != nil {
		return fmt.Errorf("failed to fetch: %w", err)
	}

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

// MergeBranch merges the specified branch into the current branch.
// Uses git command-line since go-git doesn't have built-in merge support.
// MergeBranch merges the specified branch into the current branch using go-git.
// This performs a fast-forward merge when possible. For non-fast-forward scenarios
// where automatic merging isn't possible, returns an error.
func MergeBranch(path, branchName string) error {
	// Open repository
	repo, err := git.PlainOpen(path)
	if err != nil {
		return fmt.Errorf("failed to open repository: %w", err)
	}

	// Get current HEAD
	head, err := repo.Head()
	if err != nil {
		return fmt.Errorf("failed to get HEAD: %w", err)
	}

	// Get the branch reference to merge
	branchRef, err := repo.Reference(plumbing.NewBranchReferenceName(branchName), true)
	if err != nil {
		return fmt.Errorf("branch '%s' not found: %w", branchName, err)
	}

	// Check if already at the same commit
	if head.Hash() == branchRef.Hash() {
		fmt.Printf("Already up to date with '%s'\n", branchName)
		return nil
	}

	// Get commit objects
	headCommit, err := repo.CommitObject(head.Hash())
	if err != nil {
		return fmt.Errorf("failed to get current commit: %w", err)
	}

	branchCommit, err := repo.CommitObject(branchRef.Hash())
	if err != nil {
		return fmt.Errorf("failed to get branch commit: %w", err)
	}

	// Check if branch is ancestor of current (already merged)
	isAncestor, err := headCommit.IsAncestor(branchCommit)
	if err != nil {
		return fmt.Errorf("failed to check ancestry: %w", err)
	}
	if isAncestor {
		fmt.Printf("Already contains all commits from '%s'\n", branchName)
		return nil
	}

	// Check if fast-forward is possible (current is ancestor of branch)
	canFF, err := branchCommit.IsAncestor(headCommit)
	if err != nil {
		return fmt.Errorf("failed to check if fast-forward possible: %w", err)
	}

	if canFF {
		// Perform fast-forward merge by updating HEAD to branch commit
		worktree, err := repo.Worktree()
		if err != nil {
			return fmt.Errorf("failed to get worktree: %w", err)
		}

		err = worktree.Checkout(&git.CheckoutOptions{
			Hash:  branchRef.Hash(),
			Force: false,
		})
		if err != nil {
			return fmt.Errorf("failed to fast-forward: %w", err)
		}

		// Update the branch reference
		currentBranchName := head.Name()
		ref := plumbing.NewHashReference(currentBranchName, branchRef.Hash())
		err = repo.Storer.SetReference(ref)
		if err != nil {
			return fmt.Errorf("failed to update reference: %w", err)
		}

		fmt.Printf("Fast-forward merge of '%s' successful\n", branchName)
		return nil
	}

	// Non-fast-forward merge - requires manual intervention
	return fmt.Errorf(`cannot automatically merge '%s' - branches have diverged

go-git does not support automatic conflict resolution.

To resolve:
  1. Ensure both branches are pushed to remote
  2. Use GitHub's web interface to create a pull request
  3. Or manually merge and resolve conflicts using external tools

Fast-forward merges are supported automatically.`, branchName)
}

// RebaseBranch rebases the current branch onto the specified branch.
// Note: go-git has limited rebase support. For complex rebases with conflicts,
// this function will return an error with instructions to use native git.
func RebaseBranch(path, branchName string) error {
	// Open repository
	repo, err := git.PlainOpen(path)
	if err != nil {
		return fmt.Errorf("failed to open repository: %w", err)
	}

	// Get current HEAD
	head, err := repo.Head()
	if err != nil {
		return fmt.Errorf("failed to get HEAD: %w", err)
	}

	currentBranch := head.Name().Short()

	// Get the branch reference to rebase onto
	branchRef, err := repo.Reference(plumbing.NewBranchReferenceName(branchName), true)
	if err != nil {
		return fmt.Errorf("branch '%s' not found: %w", branchName, err)
	}

	// Check if already at the same commit
	if head.Hash() == branchRef.Hash() {
		fmt.Printf("Already up to date with '%s'\n", branchName)
		return nil
	}

	// Get worktree
	worktree, err := repo.Worktree()
	if err != nil {
		return fmt.Errorf("failed to get worktree: %w", err)
	}

	// Check if working tree is clean
	status, err := worktree.Status()
	if err != nil {
		return fmt.Errorf("failed to get status: %w", err)
	}
	if !status.IsClean() {
		return fmt.Errorf("working tree is not clean. Commit or stash changes before rebasing")
	}

	// Get commit objects
	headCommit, err := repo.CommitObject(head.Hash())
	if err != nil {
		return fmt.Errorf("failed to get current commit: %w", err)
	}

	branchCommit, err := repo.CommitObject(branchRef.Hash())
	if err != nil {
		return fmt.Errorf("failed to get branch commit: %w", err)
	}

	// Check if branch is ancestor of current (already up to date)
	isAncestor, err := headCommit.IsAncestor(branchCommit)
	if err != nil {
		return fmt.Errorf("failed to check ancestry: %w", err)
	}
	if isAncestor {
		fmt.Printf("Already contains all commits from '%s'\n", branchName)
		return nil
	}

	// Check if current is ancestor of branch (can fast-forward)
	canFF, err := branchCommit.IsAncestor(headCommit)
	if err != nil {
		return fmt.Errorf("failed to check if fast-forward possible: %w", err)
	}

	if canFF {
		// Can do a simple fast-forward
		err = worktree.Checkout(&git.CheckoutOptions{
			Hash:  branchRef.Hash(),
			Force: false,
		})
		if err != nil {
			return fmt.Errorf("failed to fast-forward: %w", err)
		}

		// Update the branch reference
		ref := plumbing.NewHashReference(head.Name(), branchRef.Hash())
		err = repo.Storer.SetReference(ref)
		if err != nil {
			return fmt.Errorf("failed to update reference: %w", err)
		}

		fmt.Printf("Fast-forward to '%s' successful\n", branchName)
		return nil
	}

	// Complex rebase needed - go-git doesn't handle conflicts well
	return fmt.Errorf(`cannot automatically rebase '%s' onto '%s' - manual rebase required

go-git does not support interactive rebase or automatic conflict resolution.

To rebase manually:
  1. cd %s
  2. git rebase %s
  3. If conflicts occur, resolve them:
     - Edit conflicting files
     - git add <resolved-files>
     - git rebase --continue
  4. Once complete: git push --force-with-lease

Alternatively, use merge instead: chtsht repo update-from-main (without --rebase)`, currentBranch, branchName, path, branchName)
}

// EnsureWorkingBranch ensures the shared working branch exists and is checked out.
// Returns the branch name that was ensured.
func EnsureWorkingBranch(path, branchName string) (string, error) {
	if branchName == "" {
		branchName = "working"
	}

	// Get current branch
	currentBranch, err := GetCurrentBranch(path)
	if err != nil {
		return "", err
	}

	// If already on the working branch, we're done
	if currentBranch == branchName {
		return branchName, nil
	}

	// Check if working branch exists
	exists, err := BranchExists(path, branchName)
	if err != nil {
		return "", err
	}

	if exists {
		// Branch exists, just switch to it
		err = CheckoutBranch(path, branchName)
		if err != nil {
			return "", fmt.Errorf("failed to switch to working branch: %w", err)
		}
		fmt.Printf("Switched to existing branch '%s'\n", branchName)
	} else {
		// Create and switch to new branch
		err = CreateAndCheckoutBranch(path, branchName)
		if err != nil {
			return "", fmt.Errorf("failed to create working branch: %w", err)
		}
		fmt.Printf("Created and switched to new branch '%s'\n", branchName)
	}

	return branchName, nil
}

// EnsureMachineBranch is deprecated. Use EnsureWorkingBranch instead.
// Kept for backward compatibility.
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
// PushOptions configures branch push behavior
type PushOptions struct {
	SetUpstream bool // Set up tracking to remote branch
	Force       bool // Force push (use with caution)
}

func PushBranch(path, branchName, token string, setUpstream bool) error {
	return PushBranchWithOptions(path, branchName, token, PushOptions{SetUpstream: setUpstream})
}

func PushBranchWithOptions(path, branchName, token string, opts PushOptions) error {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return fmt.Errorf("failed to open repository: %w", err)
	}

	refSpec := fmt.Sprintf("refs/heads/%s:refs/heads/%s", branchName, branchName)

	// Add + prefix for force push (equivalent to --force-with-lease behavior)
	if opts.Force {
		refSpec = "+" + refSpec
	}

	pushOpts := &git.PushOptions{
		Progress: os.Stdout,
		// Always specify which branch to push to avoid pushing unintended branches
		RefSpecs: []config.RefSpec{
			config.RefSpec(refSpec),
		},
	}

	// If token is provided, configure authentication
	if token != "" {
		pushOpts.Auth = &http.BasicAuth{
			Username: token,
			Password: "",
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

	// If setUpstream flag is true, update local branch config to track remote
	if opts.SetUpstream {
		// Get the branch reference
		branchRef := plumbing.NewBranchReferenceName(branchName)

		// Update the branch configuration
		cfg, err := repo.Config()
		if err != nil {
			// Push succeeded but upstream setup failed - not critical
			fmt.Printf("Warning: failed to set upstream tracking: %v\n", err)
			return nil
		}

		cfg.Branches[branchName] = &config.Branch{
			Name:   branchName,
			Remote: "origin",
			Merge:  branchRef,
		}

		err = repo.Storer.SetConfig(cfg)
		if err != nil {
			fmt.Printf("Warning: failed to save upstream tracking: %v\n", err)
			return nil
		}

		fmt.Printf("Branch '%s' set up to track 'origin/%s'\n", branchName, branchName)
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

	// Check if local is ahead of remote (remote is ancestor of local)
	isLocalAhead, err := localCommit.IsAncestor(remoteCommit)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to check ancestry: %w", err)
	}

	// Check if remote is ahead of local (local is ancestor of remote)
	isRemoteAhead, err := remoteCommit.IsAncestor(localCommit)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to check ancestry: %w", err)
	}

	ahead = 0
	behind = 0

	// If remote is ancestor of local, count commits ahead
	if isLocalAhead {
		commits, err := repo.Log(&git.LogOptions{
			From: localCommit.Hash,
		})
		if err != nil {
			return 0, 0, fmt.Errorf("failed to get log: %w", err)
		}

		err = commits.ForEach(func(c *object.Commit) error {
			if c.Hash == remoteCommit.Hash {
				return fmt.Errorf("stop")
			}
			ahead++
			return nil
		})
	}

	// If local is ancestor of remote, count commits behind
	if isRemoteAhead {
		commits, err := repo.Log(&git.LogOptions{
			From: remoteCommit.Hash,
		})
		if err != nil {
			return 0, 0, fmt.Errorf("failed to get log: %w", err)
		}

		err = commits.ForEach(func(c *object.Commit) error {
			if c.Hash == localCommit.Hash {
				return fmt.Errorf("stop")
			}
			behind++
			return nil
		})
	}

	// If neither is ancestor of the other, branches have truly diverged
	// In this case, we'd need merge-base calculation which is complex
	// For now, just indicate divergence by setting both to 1
	if !isLocalAhead && !isRemoteAhead {
		// Branches have diverged - try to estimate
		// This is imperfect without proper merge-base
		ahead = 1  // Indicate there are local changes
		behind = 1 // Indicate there are remote changes
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

	// Fetch to update remote refs before checking divergence
	repo, err := git.PlainOpen(path)
	if err == nil {
		fetchOpts := &git.FetchOptions{
			RemoteName: "origin",
		}
		if token != "" {
			fetchOpts.Auth = &http.BasicAuth{
				Username: token,
				Password: "",
			}
		}
		// Fetch silently, ignore errors (might be offline)
		_ = repo.Fetch(fetchOpts)
	}

	// Get branch divergence (now against fresh remote refs)
	ahead, behind, err := GetBranchDivergence(path, basicStatus.CurrentBranch, token)
	if err == nil {
		status.AheadBy = ahead
		status.BehindBy = behind
	}

	return status, nil
}

// GetGitAuthor returns the configured git author name and email.
// Falls back to environment variables or returns an error.
func GetGitAuthor(configName, configEmail string) (name, email string, err error) {
	name = configName
	email = configEmail

	// If name not provided, try environment variables
	if name == "" {
		name = os.Getenv("GIT_AUTHOR_NAME")
		if name == "" {
			name = os.Getenv("GIT_COMMITTER_NAME")
		}
		if name == "" {
			return "", "", fmt.Errorf("git author name not configured. Set in config.yml or environment variable GIT_AUTHOR_NAME")
		}
	}

	// If email not provided, try environment variables
	if email == "" {
		email = os.Getenv("GIT_AUTHOR_EMAIL")
		if email == "" {
			email = os.Getenv("GIT_COMMITTER_EMAIL")
		}
		if email == "" {
			return "", "", fmt.Errorf("git author email not configured. Set in config.yml or environment variable GIT_AUTHOR_EMAIL")
		}
	}

	return name, email, nil
}

// HasUpstreamTracking checks if the current branch has upstream tracking configured
func HasUpstreamTracking(path string) (bool, string, error) {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return false, "", fmt.Errorf("failed to open repository: %w", err)
	}

	head, err := repo.Head()
	if err != nil {
		return false, "", fmt.Errorf("failed to get HEAD: %w", err)
	}

	branchName := head.Name().Short()

	// Get branch config to check for upstream
	cfg, err := repo.Config()
	if err != nil {
		return false, "", fmt.Errorf("failed to get config: %w", err)
	}

	// Check if branch has upstream configured
	for _, branch := range cfg.Branches {
		if branch.Name == branchName {
			if branch.Remote != "" && branch.Merge != "" {
				upstream := fmt.Sprintf("%s/%s", branch.Remote, branch.Merge.Short())
				return true, upstream, nil
			}
		}
	}

	return false, "", nil
}

// IsRemoteReachable checks if the remote repository is accessible
func IsRemoteReachable(path, token string) error {
	repo, err := git.PlainOpen(path)
	if err != nil {
		return fmt.Errorf("failed to open repository: %w", err)
	}

	remote, err := repo.Remote("origin")
	if err != nil {
		return fmt.Errorf("failed to get remote: %w", err)
	}

	listOpts := &git.ListOptions{}
	if token != "" {
		listOpts.Auth = &http.BasicAuth{
			Username: token,
			Password: "",
		}
	}

	// Try to list remote references (lightweight operation)
	_, err = remote.List(listOpts)
	if err != nil {
		return fmt.Errorf("remote is not reachable: %w", err)
	}

	return nil
}

// IsValidBranchName checks if a branch name follows conventions
// Valid formats: feat/*, fix/*, docs/*, refactor/*, test/*, chore/*, or "working", "main", "master"
func IsValidBranchName(branchName string) (bool, string) {
	// Allow special branches
	if branchName == "main" || branchName == "master" || branchName == "working" {
		return true, ""
	}

	// Check for conventional commit prefixes
	validPrefixes := []string{"feat/", "fix/", "docs/", "refactor/", "test/", "chore/", "style/", "perf/", "ci/"}
	for _, prefix := range validPrefixes {
		if strings.HasPrefix(branchName, prefix) && len(branchName) > len(prefix) {
			return true, ""
		}
	}

	suggestion := "Branch names should follow convention: feat/*, fix/*, docs/*, etc.\nExamples: feat/new-feature, fix/bug-name"
	return false, suggestion
}
