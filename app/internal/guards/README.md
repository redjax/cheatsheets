# Guards System

The guards system provides pre-flight checks to protect repository operations from failing or causing corruption.

## Overview

Guards are lightweight checks that run before commands execute. They validate preconditions (like "working tree is clean" or "repository exists") and provide helpful error messages when checks fail.

## Available Checks

| Check Type | Purpose | When to Use |
|------------|---------|-------------|
| `RepoCloned` | Repository exists at configured path | Almost all repo commands |
| `CleanWorkingTree` | No uncommitted changes | Before merge, pull, branch switch |
| `OnWorkingBranch` | Not on main/master | Before making edits, committing new work |
| `NotOnWorkingBranch` | Currently on main/master | Before operations that require main |
| `NoMergeInProgress` | No active merge conflicts | Before starting new operations |
| `HasUpstream` | Branch has upstream tracking configured | Before push operations |
| `ValidBranchName` | Branch follows naming conventions | For enforcing branch naming standards |
| `RemoteReachable` | Remote repository is accessible | Before push/pull operations (optional) |

## Usage

### Basic Usage in Commands

```go
package mycommand

import (
    "github.com/redjax/cheatsheets/internal/guards"
    "github.com/redjax/cheatsheets/internal/config"
)

var MyCmd = &cobra.Command{
    Use: "mycommand",
    RunE: func(cmd *cobra.Command, args []string) error {
        // Load config first
        cfg, err := config.LoadConfig(nil, configFile)
        if err != nil {
            return err
        }

        // Run pre-flight checks
        guardCtx := guards.NewGuardContext(cfg)
        if err := guards.CheckAll(guardCtx, 
            guards.RepoCloned, 
            guards.CleanWorkingTree,
        ); err != nil {
            return err // Guard provides user-friendly error message
        }

        // Your command logic here - guards passed!
        // ...
    },
}
```

### Choosing Which Guards to Use

Read-only commands (show, list):

```go
guards.CheckAll(guardCtx, guards.RepoCloned)
```

Edit commands (edit, new):

```go
guards.CheckAll(guardCtx, 
    guards.RepoCloned,
    guards.OnWorkingBranch,
)
```

Merge/Branch operations (merge-from-main, merge-to-main):

```go
guards.CheckAll(guardCtx, 
    guards.RepoCloned,
    guards.CleanWorkingTree,
    guards.OnWorkingBranch,
)
```

Push/Pull operations:

```go
guards.CheckAll(guardCtx, 
    guards.RepoCloned,
    guards.CleanWorkingTree,
)
```

### Getting Detailed Results

If you need to inspect check results rather than just pass/fail:

```go
results := guards.CheckAllWithResults(guardCtx,
    guards.RepoCloned,
    guards.CleanWorkingTree,
)

for _, result := range results {
    if !result.Passed {
        fmt.Printf("Check %s failed: %s\n", result.Type, result.Message)
        if result.Fix != "" {
            fmt.Printf("  Fix: %s\n", result.Fix)
        }
    }
}
```

### Running a Single Check

```go
result := guards.Check(guardCtx, guards.CleanWorkingTree)
if !result.Passed {
    return fmt.Errorf(result.Message)
}
```

## Error Messages

When a guard fails, it returns a message describing what went wrong and a suggested fix.

Example output:

```shell
Pre-flight check failed: you have uncommitted changes

To fix: chtsht repo status (to see changes) or chtsht repo commit -a (to commit them)
```

## Adding New Guards

To add a new guard type:

- Define the constant in `guards.go`:

```go
const (
    MyNewCheck CheckType = "my_new_check"
)
```

- Add to switch statement in `Check()`:

```go
case MyNewCheck:
    return checkMyNewCheck(ctx)
```

- Implement check function:

```go
func checkMyNewCheck(ctx *GuardContext) *CheckResult {
    // Perform your check
    if somethingIsWrong {
        return &CheckResult{
            Type:    MyNewCheck,
            Passed:  false,
            Message: "clear description of problem",
            Fix:     "command to fix it",
        }
    }
    
    return &CheckResult{
        Type:   MyNewCheck,
        Passed: true,
    }
}
```

## Best Practices

Do:

- Use guards for all mutating operations (merge, commit, push)
- Keep checks fast (< 50ms per check)
- Provide actionable fix suggestions
- Only run checks that are necessary for the command

Don't:

- Add guards to simple read operations (unless checking repo exists)
- Use `RemoteReachable` for every command (it's a network call - use sparingly)
- Perform deep validation that belongs in the command itself
- Add guards that duplicate work the command will do anyway

## Examples

See these commands for guard usage examples:
- [`updateFromMainCommand/update_from_main_cmd.go`](../commands/repoCommand/updateFromMainCommand/update_from_main_cmd.go) - Full example with all typical guards
- [`mergeToMainCommand/merge_to_main_cmd.go`](../commands/repoCommand/mergeToMainCommand/merge_to_main_cmd.go) - Merge operation guards
