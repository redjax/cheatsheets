package guards

import (
	"fmt"

	"github.com/redjax/cheatsheets/internal/config"
	reposervices "github.com/redjax/cheatsheets/internal/services/repoServices"
)

// CheckType represents different types of pre-flight checks
type CheckType string

const (
	// RepoCloned ensures the repository exists at the configured path
	RepoCloned CheckType = "repo_cloned"

	// CleanWorkingTree ensures there are no uncommitted changes
	CleanWorkingTree CheckType = "clean_working_tree"

	// OnWorkingBranch ensures we're not on main/master branch
	OnWorkingBranch CheckType = "on_working_branch"

	// NotOnWorkingBranch ensures we ARE on main/master (e.g., for cloning)
	NotOnWorkingBranch CheckType = "not_on_working_branch"

	// NoMergeInProgress ensures there's no active merge conflict
	NoMergeInProgress CheckType = "no_merge_in_progress"

	// HasUpstream ensures the current branch has upstream tracking configured
	HasUpstream CheckType = "has_upstream"

	// ValidBranchName ensures branch follows naming conventions
	ValidBranchName CheckType = "valid_branch_name"

	// RemoteReachable ensures the remote repository is accessible
	RemoteReachable CheckType = "remote_reachable"
)

// CheckResult holds the result of a guard check
type CheckResult struct {
	Type    CheckType
	Passed  bool
	Message string
	Fix     string // Suggested fix command/action
}

// GuardContext holds configuration and state for guards
type GuardContext struct {
	Config   *config.Config
	RepoPath string
}

// Check executes a single guard check
func Check(ctx *GuardContext, checkType CheckType) *CheckResult {
	switch checkType {
	case RepoCloned:
		return checkRepoCloned(ctx)
	case CleanWorkingTree:
		return checkCleanWorkingTree(ctx)
	case OnWorkingBranch:
		return checkOnWorkingBranch(ctx)
	case NotOnWorkingBranch:
		return checkNotOnWorkingBranch(ctx)
	case NoMergeInProgress:
		return checkNoMergeInProgress(ctx)
	case HasUpstream:
		return checkHasUpstream(ctx)
	case ValidBranchName:
		return checkValidBranchName(ctx)
	case RemoteReachable:
		return checkRemoteReachable(ctx)
	default:
		return &CheckResult{
			Type:    checkType,
			Passed:  false,
			Message: fmt.Sprintf("unknown check type: %s", checkType),
		}
	}
}

// CheckAll runs multiple checks and returns the first failure, or nil if all pass
func CheckAll(ctx *GuardContext, checks ...CheckType) error {
	for _, checkType := range checks {
		result := Check(ctx, checkType)
		if !result.Passed {
			return formatCheckError(result)
		}
	}
	return nil
}

// CheckAllWithResults runs all checks and returns all results (doesn't stop on first failure)
func CheckAllWithResults(ctx *GuardContext, checks ...CheckType) []*CheckResult {
	results := make([]*CheckResult, len(checks))
	for i, checkType := range checks {
		results[i] = Check(ctx, checkType)
	}
	return results
}

// formatCheckError converts a failed check into a user-friendly error
func formatCheckError(result *CheckResult) error {
	msg := fmt.Sprintf("Pre-flight check failed: %s", result.Message)
	if result.Fix != "" {
		msg += fmt.Sprintf("\n\nTo fix: %s", result.Fix)
	}
	return fmt.Errorf("%s", msg)
}

// Individual check implementations

func checkRepoCloned(ctx *GuardContext) *CheckResult {
	cloned, err := reposervices.IsRepositoryCloned(ctx.RepoPath)
	if err != nil {
		return &CheckResult{
			Type:    RepoCloned,
			Passed:  false,
			Message: fmt.Sprintf("failed to check repository: %v", err),
		}
	}

	if !cloned {
		return &CheckResult{
			Type:    RepoCloned,
			Passed:  false,
			Message: fmt.Sprintf("repository not found at %s", ctx.RepoPath),
			Fix:     "chtsht repo clone",
		}
	}

	return &CheckResult{
		Type:   RepoCloned,
		Passed: true,
	}
}

func checkCleanWorkingTree(ctx *GuardContext) *CheckResult {
	clean, err := reposervices.IsWorkingTreeClean(ctx.RepoPath)
	if err != nil {
		return &CheckResult{
			Type:    CleanWorkingTree,
			Passed:  false,
			Message: fmt.Sprintf("failed to check working tree: %v", err),
		}
	}

	if !clean {
		return &CheckResult{
			Type:    CleanWorkingTree,
			Passed:  false,
			Message: "you have uncommitted changes",
			Fix:     "chtsht repo status (to see changes) or chtsht repo commit -a (to commit them)",
		}
	}

	return &CheckResult{
		Type:   CleanWorkingTree,
		Passed: true,
	}
}

func checkOnWorkingBranch(ctx *GuardContext) *CheckResult {
	branch, err := reposervices.GetCurrentBranch(ctx.RepoPath)
	if err != nil {
		return &CheckResult{
			Type:    OnWorkingBranch,
			Passed:  false,
			Message: fmt.Sprintf("failed to get current branch: %v", err),
		}
	}

	if branch == "main" || branch == "master" {
		return &CheckResult{
			Type:    OnWorkingBranch,
			Passed:  false,
			Message: fmt.Sprintf("you are on the '%s' branch", branch),
			Fix:     "chtsht repo branch (to switch to working branch)",
		}
	}

	return &CheckResult{
		Type:   OnWorkingBranch,
		Passed: true,
	}
}

func checkNotOnWorkingBranch(ctx *GuardContext) *CheckResult {
	branch, err := reposervices.GetCurrentBranch(ctx.RepoPath)
	if err != nil {
		return &CheckResult{
			Type:    NotOnWorkingBranch,
			Passed:  false,
			Message: fmt.Sprintf("failed to get current branch: %v", err),
		}
	}

	if branch != "main" && branch != "master" {
		return &CheckResult{
			Type:    NotOnWorkingBranch,
			Passed:  false,
			Message: fmt.Sprintf("you are on branch '%s', but should be on main/master", branch),
			Fix:     "git checkout main",
		}
	}

	return &CheckResult{
		Type:   NotOnWorkingBranch,
		Passed: true,
	}
}

func checkNoMergeInProgress(ctx *GuardContext) *CheckResult {
	// Check if .git/MERGE_HEAD exists (indicates merge in progress)
	// This is a simple check - could be expanded based on go-git capabilities
	clean, err := reposervices.IsWorkingTreeClean(ctx.RepoPath)
	if err != nil {
		return &CheckResult{
			Type:    NoMergeInProgress,
			Passed:  false,
			Message: fmt.Sprintf("failed to check merge status: %v", err),
		}
	}

	// For now, we'll consider an unclean tree as potentially having merge conflicts
	// A more sophisticated check could parse git status for "both modified" etc.
	if !clean {
		return &CheckResult{
			Type:    NoMergeInProgress,
			Passed:  false,
			Message: "repository may have merge conflicts or uncommitted changes",
			Fix:     "chtsht repo status (to check for conflicts)",
		}
	}

	return &CheckResult{
		Type:   NoMergeInProgress,
		Passed: true,
	}
}

func checkHasUpstream(ctx *GuardContext) *CheckResult {
	hasUpstream, upstream, err := reposervices.HasUpstreamTracking(ctx.RepoPath)
	if err != nil {
		return &CheckResult{
			Type:    HasUpstream,
			Passed:  false,
			Message: fmt.Sprintf("failed to check upstream tracking: %v", err),
		}
	}

	if !hasUpstream {
		branch, _ := reposervices.GetCurrentBranch(ctx.RepoPath)
		return &CheckResult{
			Type:    HasUpstream,
			Passed:  false,
			Message: fmt.Sprintf("branch '%s' has no upstream tracking", branch),
			Fix:     "chtsht repo push --set-upstream (to set upstream and push)",
		}
	}

	return &CheckResult{
		Type:    HasUpstream,
		Passed:  true,
		Message: fmt.Sprintf("tracking %s", upstream),
	}
}

func checkValidBranchName(ctx *GuardContext) *CheckResult {
	branch, err := reposervices.GetCurrentBranch(ctx.RepoPath)
	if err != nil {
		return &CheckResult{
			Type:    ValidBranchName,
			Passed:  false,
			Message: fmt.Sprintf("failed to get current branch: %v", err),
		}
	}

	valid, suggestion := reposervices.IsValidBranchName(branch)
	if !valid {
		return &CheckResult{
			Type:    ValidBranchName,
			Passed:  false,
			Message: fmt.Sprintf("branch name '%s' doesn't follow conventions", branch),
			Fix:     suggestion,
		}
	}

	return &CheckResult{
		Type:   ValidBranchName,
		Passed: true,
	}
}

func checkRemoteReachable(ctx *GuardContext) *CheckResult {
	err := reposervices.IsRemoteReachable(ctx.RepoPath, ctx.Config.Git.Token)
	if err != nil {
		return &CheckResult{
			Type:    RemoteReachable,
			Passed:  false,
			Message: fmt.Sprintf("remote is not reachable: %v", err),
			Fix:     "Check your internet connection or verify repository URL in config",
		}
	}

	return &CheckResult{
		Type:   RemoteReachable,
		Passed: true,
	}
}

// NewGuardContext creates a guard context from config
func NewGuardContext(cfg *config.Config) *GuardContext {
	return &GuardContext{
		Config:   cfg,
		RepoPath: cfg.Git.ClonePath,
	}
}
