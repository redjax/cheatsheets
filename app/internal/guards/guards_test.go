package guards

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/object"
	appconfig "github.com/redjax/cheatsheets/internal/config"
)

// Helper function to create a test repository
func createTestRepo(t *testing.T) string {
	t.Helper()
	tempDir := t.TempDir() // Auto-cleanup after test

	repo, err := git.PlainInit(tempDir, false)
	if err != nil {
		t.Fatalf("failed to init test repo: %v", err)
	}

	// Create initial commit so we have a proper repository
	worktree, err := repo.Worktree()
	if err != nil {
		t.Fatalf("failed to get worktree: %v", err)
	}

	// Create a test file
	testFile := filepath.Join(tempDir, "test.txt")
	if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
		t.Fatalf("failed to create test file: %v", err)
	}

	// Stage and commit
	if _, err := worktree.Add("test.txt"); err != nil {
		t.Fatalf("failed to add file: %v", err)
	}

	_, err = worktree.Commit("Initial commit", &git.CommitOptions{
		Author: &object.Signature{
			Name:  "Test User",
			Email: "test@example.com",
			When:  time.Now(),
		},
	})
	if err != nil {
		t.Fatalf("failed to create initial commit: %v", err)
	}

	return tempDir
}

// Helper to create test context
func createTestContext(repoPath string) *GuardContext {
	return &GuardContext{
		Config: &appconfig.Config{
			Git: appconfig.GitConfig{
				ClonePath: repoPath,
			},
		},
		RepoPath: repoPath,
	}
}

// TestCheckType tests that all check types are defined
func TestCheckType(t *testing.T) {
	checkTypes := []CheckType{
		RepoCloned,
		CleanWorkingTree,
		OnWorkingBranch,
		NotOnWorkingBranch,
		NoMergeInProgress,
		HasUpstream,
		ValidBranchName,
		RemoteReachable,
	}

	for _, ct := range checkTypes {
		if ct == "" {
			t.Errorf("check type should not be empty")
		}
	}
}

// TestCheckRepoCloned tests repository cloned verification
func TestCheckRepoCloned(t *testing.T) {
	tests := []struct {
		name       string
		setupFunc  func() string
		wantPassed bool
		wantFix    string
	}{
		{
			name: "valid repository",
			setupFunc: func() string {
				return createTestRepo(t)
			},
			wantPassed: true,
			wantFix:    "",
		},
		{
			name: "nonexistent repository",
			setupFunc: func() string {
				return filepath.Join(t.TempDir(), "nonexistent")
			},
			wantPassed: false,
			wantFix:    "chtsht repo clone",
		},
		{
			name: "empty directory",
			setupFunc: func() string {
				return t.TempDir()
			},
			wantPassed: false,
			wantFix:    "chtsht repo clone",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoPath := tt.setupFunc()
			ctx := createTestContext(repoPath)

			result := Check(ctx, RepoCloned)

			if result.Passed != tt.wantPassed {
				t.Errorf("Check(RepoCloned) passed = %v, want %v", result.Passed, tt.wantPassed)
			}

			if result.Fix != tt.wantFix {
				t.Errorf("Check(RepoCloned) fix = %q, want %q", result.Fix, tt.wantFix)
			}

			if result.Type != RepoCloned {
				t.Errorf("Check(RepoCloned) type = %v, want %v", result.Type, RepoCloned)
			}
		})
	}
}

// TestCheckCleanWorkingTree tests working tree cleanliness
func TestCheckCleanWorkingTree(t *testing.T) {
	tests := []struct {
		name       string
		setupFunc  func() string
		wantPassed bool
	}{
		{
			name: "clean working tree",
			setupFunc: func() string {
				return createTestRepo(t)
			},
			wantPassed: true,
		},
		{
			name: "dirty working tree",
			setupFunc: func() string {
				repoPath := createTestRepo(t)
				// Create an uncommitted file
				testFile := filepath.Join(repoPath, "dirty.txt")
				if err := os.WriteFile(testFile, []byte("uncommitted"), 0644); err != nil {
					t.Fatalf("failed to create dirty file: %v", err)
				}
				return repoPath
			},
			wantPassed: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoPath := tt.setupFunc()
			ctx := createTestContext(repoPath)

			result := Check(ctx, CleanWorkingTree)

			if result.Passed != tt.wantPassed {
				t.Errorf("Check(CleanWorkingTree) passed = %v, want %v", result.Passed, tt.wantPassed)
			}

			if !result.Passed && result.Fix == "" {
				t.Error("Check(CleanWorkingTree) should provide fix suggestion on failure")
			}
		})
	}
}

// TestCheckOnWorkingBranch tests branch type verification
func TestCheckOnWorkingBranch(t *testing.T) {
	tests := []struct {
		name       string
		branchName string
		wantPassed bool
	}{
		{
			name:       "on working branch",
			branchName: "working",
			wantPassed: true,
		},
		{
			name:       "on feature branch",
			branchName: "feat/new-feature",
			wantPassed: true,
		},
		{
			name:       "on main branch",
			branchName: "main",
			wantPassed: false,
		},
		{
			name:       "on master branch",
			branchName: "master",
			wantPassed: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoPath := createTestRepo(t)

			// Create and checkout the test branch if not main
			if tt.branchName != "main" && tt.branchName != "master" {
				repo, err := git.PlainOpen(repoPath)
				if err != nil {
					t.Fatalf("failed to open repo: %v", err)
				}

				worktree, err := repo.Worktree()
				if err != nil {
					t.Fatalf("failed to get worktree: %v", err)
				}

				// Create new branch
				headRef, err := repo.Head()
				if err != nil {
					t.Fatalf("failed to get head: %v", err)
				}

				branchRef := plumbing.NewBranchReferenceName(tt.branchName)
				ref := plumbing.NewHashReference(branchRef, headRef.Hash())
				if err := repo.Storer.SetReference(ref); err != nil {
					t.Fatalf("failed to create branch: %v", err)
				}

				// Checkout the branch
				if err := worktree.Checkout(&git.CheckoutOptions{
					Branch: branchRef,
				}); err != nil {
					t.Fatalf("failed to checkout branch: %v", err)
				}
			}

			ctx := createTestContext(repoPath)
			result := Check(ctx, OnWorkingBranch)

			if result.Passed != tt.wantPassed {
				t.Errorf("Check(OnWorkingBranch) on branch %q: passed = %v, want %v",
					tt.branchName, result.Passed, tt.wantPassed)
			}
		})
	}
}

// TestCheckNotOnWorkingBranch tests inverse branch check
func TestCheckNotOnWorkingBranch(t *testing.T) {
	tests := []struct {
		name       string
		branchName string
		wantPassed bool
	}{
		{
			name:       "on main branch",
			branchName: "main",
			wantPassed: true,
		},
		{
			name:       "on working branch",
			branchName: "working",
			wantPassed: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoPath := createTestRepo(t)

			if tt.branchName != "main" {
				repo, err := git.PlainOpen(repoPath)
				if err != nil {
					t.Fatalf("failed to open repo: %v", err)
				}

				worktree, err := repo.Worktree()
				if err != nil {
					t.Fatalf("failed to get worktree: %v", err)
				}

				headRef, err := repo.Head()
				if err != nil {
					t.Fatalf("failed to get head: %v", err)
				}

				branchRef := plumbing.NewBranchReferenceName(tt.branchName)
				ref := plumbing.NewHashReference(branchRef, headRef.Hash())
				if err := repo.Storer.SetReference(ref); err != nil {
					t.Fatalf("failed to create branch: %v", err)
				}

				if err := worktree.Checkout(&git.CheckoutOptions{
					Branch: branchRef,
				}); err != nil {
					t.Fatalf("failed to checkout branch: %v", err)
				}
			}

			ctx := createTestContext(repoPath)
			result := Check(ctx, NotOnWorkingBranch)

			if result.Passed != tt.wantPassed {
				t.Errorf("Check(NotOnWorkingBranch) on branch %q: passed = %v, want %v",
					tt.branchName, result.Passed, tt.wantPassed)
			}
		})
	}
}

// TestCheckAll tests running multiple checks
func TestCheckAll(t *testing.T) {
	tests := []struct {
		name      string
		setupFunc func() string
		checks    []CheckType
		wantError bool
	}{
		{
			name: "all checks pass",
			setupFunc: func() string {
				repoPath := createTestRepo(t)
				// Create working branch
				repo, _ := git.PlainOpen(repoPath)
				worktree, _ := repo.Worktree()
				headRef, _ := repo.Head()
				branchRef := plumbing.NewBranchReferenceName("working")
				ref := plumbing.NewHashReference(branchRef, headRef.Hash())
				repo.Storer.SetReference(ref)
				worktree.Checkout(&git.CheckoutOptions{Branch: branchRef})
				return repoPath
			},
			checks: []CheckType{
				RepoCloned,
				CleanWorkingTree,
				OnWorkingBranch,
			},
			wantError: false,
		},
		{
			name: "first check fails",
			setupFunc: func() string {
				return filepath.Join(t.TempDir(), "nonexistent")
			},
			checks: []CheckType{
				RepoCloned,
				CleanWorkingTree,
			},
			wantError: true,
		},
		{
			name: "second check fails",
			setupFunc: func() string {
				repoPath := createTestRepo(t)
				// Make working tree dirty
				testFile := filepath.Join(repoPath, "dirty.txt")
				os.WriteFile(testFile, []byte("dirty"), 0644)
				return repoPath
			},
			checks: []CheckType{
				RepoCloned,
				CleanWorkingTree,
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			repoPath := tt.setupFunc()
			ctx := createTestContext(repoPath)

			err := CheckAll(ctx, tt.checks...)

			if (err != nil) != tt.wantError {
				t.Errorf("CheckAll() error = %v, wantError %v", err, tt.wantError)
			}
		})
	}
}

// TestCheckAllWithResults tests collecting all check results
func TestCheckAllWithResults(t *testing.T) {
	repoPath := createTestRepo(t)
	ctx := createTestContext(repoPath)

	checks := []CheckType{
		RepoCloned,       // Should pass
		CleanWorkingTree, // Should pass
		OnWorkingBranch,  // Should fail (on main)
	}

	results := CheckAllWithResults(ctx, checks...)

	if len(results) != len(checks) {
		t.Errorf("CheckAllWithResults() returned %d results, want %d", len(results), len(checks))
	}

	// Verify we got results for all checks
	for i, result := range results {
		if result.Type != checks[i] {
			t.Errorf("Result[%d] type = %v, want %v", i, result.Type, checks[i])
		}
	}

	// Verify at least one check failed (OnWorkingBranch on main)
	failedCount := 0
	for _, result := range results {
		if !result.Passed {
			failedCount++
		}
	}

	if failedCount == 0 {
		t.Error("Expected at least one check to fail, but all passed")
	}
}

// TestUnknownCheckType tests handling of invalid check types
func TestUnknownCheckType(t *testing.T) {
	repoPath := createTestRepo(t)
	ctx := createTestContext(repoPath)

	unknownCheck := CheckType("unknown_check")
	result := Check(ctx, unknownCheck)

	if result.Passed {
		t.Error("Check(unknown) should fail")
	}

	if result.Type != unknownCheck {
		t.Errorf("Check(unknown) type = %v, want %v", result.Type, unknownCheck)
	}
}

// TestNewGuardContext tests creating guard context
func TestNewGuardContext(t *testing.T) {
	cfg := &appconfig.Config{
		Git: appconfig.GitConfig{
			ClonePath: "/test/path",
		},
	}

	ctx := NewGuardContext(cfg)

	if ctx.Config != cfg {
		t.Error("NewGuardContext() config not set correctly")
	}

	if ctx.RepoPath != cfg.Git.ClonePath {
		t.Errorf("NewGuardContext() repoPath = %q, want %q", ctx.RepoPath, cfg.Git.ClonePath)
	}
}

// TestFormatCheckError tests error message formatting
func TestFormatCheckError(t *testing.T) {
	tests := []struct {
		name    string
		result  *CheckResult
		wantMsg string
	}{
		{
			name: "error with fix suggestion",
			result: &CheckResult{
				Type:    RepoCloned,
				Passed:  false,
				Message: "repository not found",
				Fix:     "chtsht repo clone",
			},
			wantMsg: "repository not found",
		},
		{
			name: "error without fix suggestion",
			result: &CheckResult{
				Type:    CleanWorkingTree,
				Passed:  false,
				Message: "uncommitted changes",
				Fix:     "",
			},
			wantMsg: "uncommitted changes",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := formatCheckError(tt.result)

			if err == nil {
				t.Fatal("formatCheckError() should return error")
			}

			errMsg := err.Error()
			if !contains(errMsg, tt.wantMsg) {
				t.Errorf("formatCheckError() message = %q, should contain %q", errMsg, tt.wantMsg)
			}

			if tt.result.Fix != "" && !contains(errMsg, tt.result.Fix) {
				t.Errorf("formatCheckError() message should contain fix suggestion %q", tt.result.Fix)
			}
		})
	}
}

// Helper function to check if string contains substring
func contains(s, substr string) bool {
	return len(s) >= len(substr) &&
		(s == substr || len(s) > len(substr) &&
			indexString(s, substr) >= 0)
}

func indexString(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
