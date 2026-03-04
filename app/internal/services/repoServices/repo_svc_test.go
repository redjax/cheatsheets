package reposervices

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
)

// Helper to create a test git repository
func createTestRepo(t *testing.T) string {
	t.Helper()

	repoPath := filepath.Join(t.TempDir(), "testrepo")

	// Initialize the repository
	repo, err := git.PlainInit(repoPath, false)
	if err != nil {
		t.Fatalf("failed to init test repo: %v", err)
	}

	// Create initial commit so repo has a HEAD
	worktree, err := repo.Worktree()
	if err != nil {
		t.Fatalf("failed to get worktree: %v", err)
	}

	// Create a test file
	testFile := filepath.Join(repoPath, "README.md")
	if err := os.WriteFile(testFile, []byte("# Test Repo\n"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	// Stage and commit
	if _, err := worktree.Add("README.md"); err != nil {
		t.Fatalf("failed to stage file: %v", err)
	}

	sig := &object.Signature{
		Name:  "Test User",
		Email: "test@example.com",
	}

	if _, err := worktree.Commit("Initial commit", &git.CommitOptions{Author: sig}); err != nil {
		t.Fatalf("failed to create initial commit: %v", err)
	}

	return repoPath
}

// Helper to create a bare repository (acts as remote)
func createBareRepo(t *testing.T) string {
	t.Helper()

	repoPath := filepath.Join(t.TempDir(), "bare.git")

	_, err := git.PlainInit(repoPath, true)
	if err != nil {
		t.Fatalf("failed to init bare repo: %v", err)
	}

	return repoPath
}

// TestIsRepositoryCloned tests repository detection
func TestIsRepositoryCloned(t *testing.T) {
	tests := []struct {
		name      string
		setupFunc func() string
		want      bool
		wantError bool
	}{
		{
			name: "valid repository",
			setupFunc: func() string {
				return createTestRepo(t)
			},
			want:      true,
			wantError: false,
		},
		{
			name: "nonexistent path",
			setupFunc: func() string {
				return filepath.Join(t.TempDir(), "nonexistent")
			},
			want:      false,
			wantError: false,
		},
		{
			name: "directory without .git",
			setupFunc: func() string {
				dir := t.TempDir()
				os.WriteFile(filepath.Join(dir, "file.txt"), []byte("test"), 0644)
				return dir
			},
			want:      false,
			wantError: false,
		},
		{
			name: "empty directory",
			setupFunc: func() string {
				return t.TempDir()
			},
			want:      false,
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := tt.setupFunc()

			got, err := IsRepositoryCloned(path)

			if (err != nil) != tt.wantError {
				t.Errorf("IsRepositoryCloned() error = %v, wantError %v", err, tt.wantError)
				return
			}

			if got != tt.want {
				t.Errorf("IsRepositoryCloned() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestGetCurrentBranch tests branch name retrieval
func TestGetCurrentBranch(t *testing.T) {
	tests := []struct {
		name      string
		setupFunc func() string
		want      string
		wantError bool
	}{
		{
			name: "default branch (master)",
			setupFunc: func() string {
				return createTestRepo(t)
			},
			want:      "master",
			wantError: false,
		},
		{
			name: "custom branch",
			setupFunc: func() string {
				repoPath := createTestRepo(t)
				repo, _ := git.PlainOpen(repoPath)
				w, _ := repo.Worktree()

				// Create and checkout new branch
				branchName := plumbing.NewBranchReferenceName("feature")
				headRef, _ := repo.Head()
				ref := plumbing.NewHashReference(branchName, headRef.Hash())
				repo.Storer.SetReference(ref)
				w.Checkout(&git.CheckoutOptions{Branch: branchName})

				return repoPath
			},
			want:      "feature",
			wantError: false,
		},
		{
			name: "invalid repository",
			setupFunc: func() string {
				return t.TempDir()
			},
			want:      "",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := tt.setupFunc()

			got, err := GetCurrentBranch(path)

			if (err != nil) != tt.wantError {
				t.Errorf("GetCurrentBranch() error = %v, wantError %v", err, tt.wantError)
				return
			}

			if got != tt.want {
				t.Errorf("GetCurrentBranch() = %q, want %q", got, tt.want)
			}
		})
	}
}

// TestIsWorkingTreeClean tests working tree status
func TestIsWorkingTreeClean(t *testing.T) {
	tests := []struct {
		name      string
		setupFunc func() string
		want      bool
		wantError bool
	}{
		{
			name: "clean working tree",
			setupFunc: func() string {
				return createTestRepo(t)
			},
			want:      true,
			wantError: false,
		},
		{
			name: "modified file",
			setupFunc: func() string {
				repoPath := createTestRepo(t)

				// Modify existing file
				testFile := filepath.Join(repoPath, "README.md")
				os.WriteFile(testFile, []byte("# Modified\n"), 0644)

				return repoPath
			},
			want:      false,
			wantError: false,
		},
		{
			name: "new untracked file",
			setupFunc: func() string {
				repoPath := createTestRepo(t)

				// Create new file
				newFile := filepath.Join(repoPath, "new.txt")
				os.WriteFile(newFile, []byte("new content"), 0644)

				return repoPath
			},
			want:      false,
			wantError: false,
		},
		{
			name: "staged changes",
			setupFunc: func() string {
				repoPath := createTestRepo(t)

				// Create and stage a new file
				repo, _ := git.PlainOpen(repoPath)
				w, _ := repo.Worktree()

				newFile := filepath.Join(repoPath, "staged.txt")
				os.WriteFile(newFile, []byte("staged content"), 0644)
				w.Add("staged.txt")

				return repoPath
			},
			want:      false,
			wantError: false,
		},
		{
			name: "invalid repository",
			setupFunc: func() string {
				return t.TempDir()
			},
			want:      false,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := tt.setupFunc()

			got, err := IsWorkingTreeClean(path)

			if (err != nil) != tt.wantError {
				t.Errorf("IsWorkingTreeClean() error = %v, wantError %v", err, tt.wantError)
				return
			}

			if got != tt.want {
				t.Errorf("IsWorkingTreeClean() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestBranchExists tests branch existence checking
func TestBranchExists(t *testing.T) {
	tests := []struct {
		name      string
		setupFunc func() (repoPath, branchName string)
		want      bool
		wantError bool
	}{
		{
			name: "existing branch (master)",
			setupFunc: func() (string, string) {
				return createTestRepo(t), "master"
			},
			want:      true,
			wantError: false,
		},
		{
			name: "non-existing branch",
			setupFunc: func() (string, string) {
				return createTestRepo(t), "nonexistent"
			},
			want:      false,
			wantError: false,
		},
		{
			name: "created branch",
			setupFunc: func() (string, string) {
				repoPath := createTestRepo(t)
				repo, _ := git.PlainOpen(repoPath)

				// Create new branch
				branchName := "feature"
				headRef, _ := repo.Head()
				ref := plumbing.NewHashReference(
					plumbing.NewBranchReferenceName(branchName),
					headRef.Hash(),
				)
				repo.Storer.SetReference(ref)

				return repoPath, branchName
			},
			want:      true,
			wantError: false,
		},
		{
			name: "invalid repository",
			setupFunc: func() (string, string) {
				return t.TempDir(), "main"
			},
			want:      false,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoPath, branchName := tt.setupFunc()

			got, err := BranchExists(repoPath, branchName)

			if (err != nil) != tt.wantError {
				t.Errorf("BranchExists() error = %v, wantError %v", err, tt.wantError)
				return
			}

			if got != tt.want {
				t.Errorf("BranchExists() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestCreateAndCheckoutBranch tests branch creation
func TestCreateAndCheckoutBranch(t *testing.T) {
	tests := []struct {
		name       string
		branchName string
		wantError  bool
	}{
		{
			name:       "valid branch name",
			branchName: "feature",
			wantError:  false,
		},
		{
			name:       "branch with slash",
			branchName: "feature/new-feature",
			wantError:  false,
		},
		{
			name:       "branch with dash",
			branchName: "feature-123",
			wantError:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoPath := createTestRepo(t)

			err := CreateAndCheckoutBranch(repoPath, tt.branchName)

			if (err != nil) != tt.wantError {
				t.Errorf("CreateAndCheckoutBranch() error = %v, wantError %v", err, tt.wantError)
				return
			}

			if !tt.wantError {
				// Verify branch was created and checked out
				currentBranch, err := GetCurrentBranch(repoPath)
				if err != nil {
					t.Fatalf("failed to get current branch: %v", err)
				}

				if currentBranch != tt.branchName {
					t.Errorf("current branch = %q, want %q", currentBranch, tt.branchName)
				}

				// Verify branch exists
				exists, err := BranchExists(repoPath, tt.branchName)
				if err != nil {
					t.Fatalf("BranchExists() error = %v", err)
				}
				if !exists {
					t.Errorf("branch %q should exist after creation", tt.branchName)
				}
			}
		})
	}
}

// TestCheckoutBranch tests branch checkout
func TestCheckoutBranch(t *testing.T) {
	tests := []struct {
		name      string
		setupFunc func() (repoPath, branchName string)
		wantError bool
	}{
		{
			name: "checkout existing branch",
			setupFunc: func() (string, string) {
				repoPath := createTestRepo(t)

				// Create feature branch
				CreateAndCheckoutBranch(repoPath, "feature")

				// Go back to main
				CheckoutBranch(repoPath, "main")

				return repoPath, "feature"
			},
			wantError: false,
		},
		{
			name: "checkout non-existing branch",
			setupFunc: func() (string, string) {
				return createTestRepo(t), "nonexistent"
			},
			wantError: true,
		},
		{
			name: "checkout master branch",
			setupFunc: func() (string, string) {
				repoPath := createTestRepo(t)
				CreateAndCheckoutBranch(repoPath, "feat/test")
				return repoPath, "master"
			},
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoPath, branchName := tt.setupFunc()

			err := CheckoutBranch(repoPath, branchName)

			if (err != nil) != tt.wantError {
				t.Errorf("CheckoutBranch() error = %v, wantError %v", err, tt.wantError)
				return
			}

			if !tt.wantError {
				// Verify we're on the correct branch
				currentBranch, err := GetCurrentBranch(repoPath)
				if err != nil {
					t.Fatalf("failed to get current branch: %v", err)
				}

				if currentBranch != branchName {
					t.Errorf("current branch = %q, want %q", currentBranch, branchName)
				}
			}
		})
	}
}

// TestListBranches tests branch listing
func TestListBranches(t *testing.T) {
	tests := []struct {
		name      string
		setupFunc func() string
		wantCount int
		wantError bool
	}{
		{
			name: "repository with one branch",
			setupFunc: func() string {
				return createTestRepo(t)
			},
			wantCount: 1,
			wantError: false,
		},
		{
			name: "repository with multiple branches",
			setupFunc: func() string {
				repoPath := createTestRepo(t)

				// Create additional branches
				CreateAndCheckoutBranch(repoPath, "feature1")
				CheckoutBranch(repoPath, "main")
				CreateAndCheckoutBranch(repoPath, "feature2")

				return repoPath
			},
			wantCount: 3, // main, feature1, feature2
			wantError: false,
		},
		{
			name: "invalid repository",
			setupFunc: func() string {
				return t.TempDir()
			},
			wantCount: 0,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoPath := tt.setupFunc()

			branches, err := ListBranches(repoPath)

			if (err != nil) != tt.wantError {
				t.Errorf("ListBranches() error = %v, wantError %v", err, tt.wantError)
				return
			}

			if len(branches) != tt.wantCount {
				t.Errorf("ListBranches() returned %d branches, want %d", len(branches), tt.wantCount)
			}

			if !tt.wantError && tt.wantCount > 0 {
				// Verify master branch exists in list
				found := false
				for _, branch := range branches {
					if branch == "master" {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("master branch not found in branch list: %v", branches)
				}
			}
		})
	}
}

// TestIsValidBranchName tests branch name validation
func TestIsValidBranchName(t *testing.T) {
	tests := []struct {
		name       string
		branchName string
		want       bool
		wantReason string
	}{
		{
			name:       "valid special branch main",
			branchName: "main",
			want:       true,
			wantReason: "",
		},
		{
			name:       "valid special branch master",
			branchName: "master",
			want:       true,
			wantReason: "",
		},
		{
			name:       "valid special branch working",
			branchName: "working",
			want:       true,
			wantReason: "",
		},
		{
			name:       "valid conventional commit feat/",
			branchName: "feat/new-feature",
			want:       true,
			wantReason: "",
		},
		{
			name:       "valid conventional commit fix/",
			branchName: "fix/bug-123",
			want:       true,
			wantReason: "",
		},
		{
			name:       "invalid simple name",
			branchName: "feature",
			want:       false,
			wantReason: "Branch names should follow convention",
		},
		{
			name:       "empty string",
			branchName: "",
			want:       false,
			wantReason: "Branch names should follow convention",
		},
		{
			name:       "starts with slash",
			branchName: "/feature",
			want:       false,
			wantReason: "Branch names should follow convention",
		},
		{
			name:       "ends with slash",
			branchName: "feat/",
			want:       false,
			wantReason: "Branch names should follow convention",
		},
		{
			name:       "starts with dot",
			branchName: ".feature",
			want:       false,
			wantReason: "Branch names should follow convention",
		},
		{
			name:       "contains double dot",
			branchName: "feature..name",
			want:       false,
			wantReason: "Branch names should follow convention",
		},
		{
			name:       "contains space",
			branchName: "feature name",
			want:       false,
			wantReason: "Branch names should follow convention",
		},
		{
			name:       "contains special chars",
			branchName: "feature@branch",
			want:       false,
			wantReason: "Branch names should follow convention",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, reason := IsValidBranchName(tt.branchName)

			if got != tt.want {
				t.Errorf("IsValidBranchName(%q) = %v, want %v", tt.branchName, got, tt.want)
			}

			if !tt.want {
				if reason == "" {
					t.Errorf("IsValidBranchName(%q) should return a reason for invalid name", tt.branchName)
				}
				if tt.wantReason != "" && !strings.Contains(reason, tt.wantReason) {
					t.Errorf("IsValidBranchName(%q) reason = %q, should contain %q", tt.branchName, reason, tt.wantReason)
				}
			}
		})
	}
}

// TestGetGitAuthor tests author information retrieval
func TestGetGitAuthor(t *testing.T) {
	tests := []struct {
		name        string
		configName  string
		configEmail string
		setEnv      bool
		envName     string
		envEmail    string
		wantName    string
		wantEmail   string
		wantError   bool
	}{
		{
			name:        "from config values",
			configName:  "Config User",
			configEmail: "config@example.com",
			setEnv:      false,
			wantName:    "Config User",
			wantEmail:   "config@example.com",
			wantError:   false,
		},
		{
			name:        "from environment variables",
			configName:  "",
			configEmail: "",
			setEnv:      true,
			envName:     "Env User",
			envEmail:    "env@example.com",
			wantName:    "Env User",
			wantEmail:   "env@example.com",
			wantError:   false,
		},
		{
			name:        "config overrides empty env",
			configName:  "Config User",
			configEmail: "config@example.com",
			setEnv:      true,
			envName:     "",
			envEmail:    "",
			wantName:    "Config User",
			wantEmail:   "config@example.com",
			wantError:   false,
		},
		{
			name:        "config takes priority when provided",
			configName:  "Config User",
			configEmail: "config@example.com",
			setEnv:      true,
			envName:     "Env User",
			envEmail:    "env@example.com",
			wantName:    "Config User",
			wantEmail:   "config@example.com",
			wantError:   false,
		},
		{
			name:        "missing both",
			configName:  "",
			configEmail: "",
			setEnv:      false,
			wantName:    "",
			wantEmail:   "",
			wantError:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup environment if needed
			if tt.setEnv {
				if tt.envName != "" {
					t.Setenv("GIT_AUTHOR_NAME", tt.envName)
				}
				if tt.envEmail != "" {
					t.Setenv("GIT_AUTHOR_EMAIL", tt.envEmail)
				}
			}

			name, email, err := GetGitAuthor(tt.configName, tt.configEmail)

			if (err != nil) != tt.wantError {
				t.Errorf("GetGitAuthor() error = %v, wantError %v", err, tt.wantError)
				return
			}

			if name != tt.wantName {
				t.Errorf("GetGitAuthor() name = %q, want %q", name, tt.wantName)
			}

			if email != tt.wantEmail {
				t.Errorf("GetGitAuthor() email = %q, want %q", email, tt.wantEmail)
			}
		})
	}
}

// TestHasUpstreamTracking tests upstream tracking detection
func TestHasUpstreamTracking(t *testing.T) {
	tests := []struct {
		name      string
		setupFunc func() string
		want      bool
		wantError bool
	}{
		{
			name: "branch without upstream",
			setupFunc: func() string {
				return createTestRepo(t)
			},
			want:      false,
			wantError: false,
		},
		{
			name: "invalid repository",
			setupFunc: func() string {
				return t.TempDir()
			},
			want:      false,
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoPath := tt.setupFunc()

			got, _, err := HasUpstreamTracking(repoPath)

			if (err != nil) != tt.wantError {
				t.Errorf("HasUpstreamTracking() error = %v, wantError %v", err, tt.wantError)
				return
			}

			if got != tt.want {
				t.Errorf("HasUpstreamTracking() = %v, want %v", got, tt.want)
			}
		})
	}
}

// TestStageFiles tests file staging
func TestStageFiles(t *testing.T) {
	tests := []struct {
		name      string
		setupFunc func() (repoPath string, files []string)
		wantError bool
	}{
		{
			name: "stage single file",
			setupFunc: func() (string, []string) {
				repoPath := createTestRepo(t)

				// Create a new file
				newFile := filepath.Join(repoPath, "new.txt")
				os.WriteFile(newFile, []byte("content"), 0644)

				return repoPath, []string{"new.txt"}
			},
			wantError: false,
		},
		{
			name: "stage multiple files",
			setupFunc: func() (string, []string) {
				repoPath := createTestRepo(t)

				// Create multiple files
				os.WriteFile(filepath.Join(repoPath, "file1.txt"), []byte("content1"), 0644)
				os.WriteFile(filepath.Join(repoPath, "file2.txt"), []byte("content2"), 0644)

				return repoPath, []string{"file1.txt", "file2.txt"}
			},
			wantError: false,
		},
		{
			name: "stage non-existent file",
			setupFunc: func() (string, []string) {
				repoPath := createTestRepo(t)
				return repoPath, []string{"nonexistent.txt"}
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoPath, files := tt.setupFunc()

			err := StageFiles(repoPath, files)

			if (err != nil) != tt.wantError {
				t.Errorf("StageFiles() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

// TestStageAll tests staging all changes
func TestStageAll(t *testing.T) {
	tests := []struct {
		name      string
		setupFunc func() string
		wantError bool
	}{
		{
			name: "stage all changes",
			setupFunc: func() string {
				repoPath := createTestRepo(t)

				// Create multiple files
				os.WriteFile(filepath.Join(repoPath, "file1.txt"), []byte("content1"), 0644)
				os.WriteFile(filepath.Join(repoPath, "file2.txt"), []byte("content2"), 0644)

				return repoPath
			},
			wantError: false,
		},
		{
			name: "clean working tree",
			setupFunc: func() string {
				return createTestRepo(t)
			},
			wantError: false,
		},
		{
			name: "invalid repository",
			setupFunc: func() string {
				return t.TempDir()
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoPath := tt.setupFunc()

			err := StageAll(repoPath)

			if (err != nil) != tt.wantError {
				t.Errorf("StageAll() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

// TestCommitChanges tests commit creation
func TestCommitChanges(t *testing.T) {
	tests := []struct {
		name      string
		setupFunc func() string
		message   string
		wantError bool
	}{
		{
			name: "commit staged changes",
			setupFunc: func() string {
				repoPath := createTestRepo(t)

				// Create and stage a file
				newFile := filepath.Join(repoPath, "new.txt")
				os.WriteFile(newFile, []byte("content"), 0644)
				StageFiles(repoPath, []string{"new.txt"})

				return repoPath
			},
			message:   "Add new file",
			wantError: false,
		},
		{
			name: "commit with nothing staged",
			setupFunc: func() string {
				return createTestRepo(t)
			},
			message:   "Empty commit",
			wantError: true,
		},
		{
			name: "invalid repository",
			setupFunc: func() string {
				return t.TempDir()
			},
			message:   "Test commit",
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoPath := tt.setupFunc()

			hash, err := CommitChanges(repoPath, tt.message, "Test User", "test@example.com")

			if (err != nil) != tt.wantError {
				t.Errorf("CommitChanges() error = %v, wantError %v", err, tt.wantError)
				return
			}

			if !tt.wantError && hash == "" {
				t.Error("CommitChanges() should return non-empty hash on success")
			}
		})
	}
}

// TestGetLastCommit tests retrieving last commit info
func TestGetLastCommit(t *testing.T) {
	tests := []struct {
		name      string
		setupFunc func() string
		wantError bool
	}{
		{
			name: "get initial commit",
			setupFunc: func() string {
				return createTestRepo(t)
			},
			wantError: false,
		},
		{
			name: "get latest commit after multiple",
			setupFunc: func() string {
				repoPath := createTestRepo(t)

				// Create another commit
				newFile := filepath.Join(repoPath, "second.txt")
				os.WriteFile(newFile, []byte("content"), 0644)
				StageFiles(repoPath, []string{"second.txt"})
				CommitChanges(repoPath, "Second commit", "Test User", "test@example.com")

				return repoPath
			},
			wantError: false,
		},
		{
			name: "invalid repository",
			setupFunc: func() string {
				return t.TempDir()
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoPath := tt.setupFunc()

			commit, err := GetLastCommit(repoPath)

			if (err != nil) != tt.wantError {
				t.Errorf("GetLastCommit() error = %v, wantError %v", err, tt.wantError)
				return
			}

			if !tt.wantError {
				if commit == nil {
					t.Error("GetLastCommit() returned nil commit")
					return
				}

				if commit.Hash == "" {
					t.Error("GetLastCommit() returned commit with empty hash")
				}

				if commit.Author == "" {
					t.Error("GetLastCommit() returned commit with empty author")
				}

				if commit.Message == "" {
					t.Error("GetLastCommit() returned commit with empty message")
				}
			}
		})
	}
}
