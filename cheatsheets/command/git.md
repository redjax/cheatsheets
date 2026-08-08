---
description: "Version control system."
last_updated: "2026-07-25"
tags: ["git", "command", "cli"]
---

# Git <!-- omit in toc -->

## Table of Contents <!-- omit in toc -->

- [Usage](#usage)
- [Branches](#branches)
  - [Filter git branch command](#filter-git-branch-command)
    - [Show branches that are fully merged](#show-branches-that-are-fully-merged)
    - [Show only unmerged branches](#show-only-unmerged-branches)
    - [Compare against remote](#compare-against-remote)
- [Push empty commit](#push-empty-commit)
- [Git Log](#git-log)
  - [Git Log Cheatsheet](#git-log-cheatsheet)
  - [Pipe commit history into less](#pipe-commit-history-into-less)
  - [Show log with ASCII graph](#show-log-with-ascii-graph)
  - [Show shortened commit hashes](#show-shortened-commit-hashes)
  - [Show change count statistics](#show-change-count-statistics)
  - [Show code diffs in log](#show-code-diffs-in-log)
  - [Filter commits](#filter-commits)
  - [Search history](#search-history)
  - [Commit ranges](#commit-ranges)
  - [Merge commit filtering](#merge-commit-filtering)
  - [Group commits by author](#group-commits-by-author)
  - [Pretty formatting](#pretty-formatting)
    - [Pretty log formatting tokens](#pretty-log-formatting-tokens)
    - [Pretty log formatting examples](#pretty-log-formatting-examples)
- [Git config](#git-config)
  - [Set preserve colors in less](#set-preserve-colors-in-less)
- [Git aliases](#git-aliases)
- [Troubleshooting](#troubleshooting)
  - [Rewrite Git commit history](#rewrite-git-commit-history)
    - [Bash script to rewrite history](#bash-script-to-rewrite-history)

## Usage

## Branches

### Filter git branch command

Running `git branch` or `git branch -a` prints all of the branches in your repository. Sometimes you may want to filter the list to see branches that only exist locally, or have been deleted from the remote.

> [!NOTE]
> These command operates on the branch you currently have checked out.
> To see branches that have been merged into `main`, make sure to `git switch main` first.

#### Show branches that are fully merged

```shell
git branch --merged
```

#### Show only unmerged branches

```shell
git branch --no-merged
```

#### Compare against remote

  ```shell
  git fetch --prune && git branch -vv
  ```

Example output:

```shell
feature/login   a1b2c3d [origin/feature/login: gone] WIP login
old-feature     d4e5f6g [origin/old-feature: gone] fix
```

## Push empty commit

You can create and push an empty commit to trigger CI/CD pipelines using:

```shell
git commit --allow-empty -m "Trigger pipeline" ; git push
```

## Git Log

*[Git log documentation](https://git-scm.com/docs/git-log)*

Use the `git log` command to show/explore the repository's history.

Syntax:

```shell
git log [--oneline] [--decorate] [--color=always]
```

There are many other options for the `git log` command, detailed below in specific examples.

### Git Log Cheatsheet

| Use Case                                   | Command                                                                                                                       |
| ------------------------------------------ | ----------------------------------------------------------------------------------------------------------------------------- |
| Basic log                                  | `git log`                                                                                                                     |
| Compact one-line log                       | `git log --oneline`                                                                                                           |
| Show branch/tag refs                       | `git log --decorate`                                                                                                          |
| ASCII graph view                           | `git log --graph --oneline --decorate`                                                                                        |
| Preserve colors through pager              | `git log --oneline --decorate --color=always \| less -R`                                                                      |
| Show shortened hashes                      | `git log --abbrev-commit`                                                                                                     |
| Show changed files stats                   | `git log --stat`                                                                                                              |
| Show full patches/diffs                    | `git log -p`                                                                                                                  |
| Show stats + patches                       | `git log --stat -p`                                                                                                           |
| Show last 5 commits                        | `git log -5`                                                                                                                  |
| Filter commits after date                  | `git log --after="2025-01-01"`                                                                                                |
| Filter commits before date                 | `git log --before="2025-01-31"`                                                                                               |
| Filter commits between dates               | `git log --after="2025-01-01" --before="2025-01-31"`                                                                          |
| Filter by author                           | `git log --author="John"`                                                                                                     |
| Filter by commit message                   | `git log --grep="fix"`                                                                                                        |
| Search for added/removed text              | `git log -S"TODO"`                                                                                                            |
| Show commits affecting file                | `git log -- src/main.c`                                                                                                       |
| Show commits in branch range               | `git log main..feature/login`                                                                                                 |
| Hide merge commits                         | `git log --no-merges`                                                                                                         |
| Show only merge commits                    | `git log --merges`                                                                                                            |
| Group commits by author                    | `git shortlog`                                                                                                                |
| Sort authors by commit count               | `git shortlog -n`                                                                                                             |
| Pretty oneline format                      | `git log --pretty=oneline`                                                                                                    |
| Pretty short format                        | `git log --pretty=short`                                                                                                      |
| Pretty medium format                       | `git log --pretty=medium`                                                                                                     |
| Pretty full format                         | `git log --pretty=full`                                                                                                       |
| Pretty fuller format                       | `git log --pretty=fuller`                                                                                                     |
| Pretty email format                        | `git log --pretty=email`                                                                                                      |
| Pretty raw format                          | `git log --pretty=raw`                                                                                                        |
| Custom pretty format                       | `git log --pretty=format:"%h - %an, %ar : %s"`                                                                                |
| Compact colored format                     | `git log --pretty=format:"%C(yellow)%h%C(reset) %C(green)%ar%C(reset) %C(blue)%an%C(reset) %s"`                               |
| Multi-line pretty format                   | `git log --pretty=format:"Commit: %h%nAuthor: %an%nDate: %ar%n%n%s%n"`                                                        |
| Graph + custom pretty format               | `git log --graph --pretty=format:"%C(yellow)%h%C(reset) - %C(green)(%cr)%C(reset) %s %C(blue)<%an>%C(reset)" --abbrev-commit` |
| Show refs without patches                  | `git log --name-only`                                                                                                         |
| Show refs with file status                 | `git log --name-status`                                                                                                       |
| Show commits from all branches             | `git log --all --decorate --oneline --graph`                                                                                  |
| Show commit timestamps in ISO format       | `git log --date=iso`                                                                                                          |
| Show relative timestamps                   | `git log --date=relative`                                                                                                     |
| Follow file renames                        | `git log --follow -- path/to/file`                                                                                            |
| Show commits by committer                  | `git log --committer="John"`                                                                                                  |
| Limit to one branch only                   | `git log --first-parent`                                                                                                      |
| Show commits touching a directory          | `git log -- src/components/`                                                                                                  |
| Show one commit per author summary         | `git shortlog -sne`                                                                                                           |
| Show history & changes for a specific file | `git log -p -- path/to/filename.ext`                                                                                          |

### Pipe commit history into less

```shell
git log --oneline --decorate --color=always | less -R
```

Flags:

- `--oneline`: Compact one-line output
- `--decorate`: Show branch/tag references
- `--color=always`: Force color output through pipes
- `less -R`: Preserve ANSI colors in pager

### Show log with ASCII graph

```shell
git log --graph --oneline --decorate --color=always
```

This will show commits in an ASCII "graph" like:

```shell
* 8f3c1a2 (HEAD -> main) Add README updates
* 7b2d9e1 Fix parsing bug
|\
| * a19cd34 Add feature branch work
|/
* 3f29aa0 Initial commit
```

### Show shortened commit hashes

```shell
git log --abbrev-commit
```

Shows a shortened commit hash, i.e. `commit 8f3c1a2` instead of `commit 8f3c1a2f15d6f61f0db95f2d8f8f7f0c79d1abc`.

### Show change count statistics

```shell
git log --stat
```

Shows an overview of changes with `+` and `-`, i.e.:

```shell
README.md | 12 +++++++++---
main.c    |  4 ++--
```

### Show code diffs in log

```shell
git log [-p/--patch]
```

Show full diff with stats:

```shell
git log --stat -p --color=always | less -R
```

### Filter commits

Show last N commits:

```shell
git log -n 3
```

You can also omit the `-n` and just use `-#`, for example:

```shell
git log -3
```

- Filter by date:
  - Commits after a specific date:

    ```shell
    git log --after="2014-07-01"
    ```

  - Relative dates:

    ```shell
    git log --after="yesterday"
    git log --after="2 weeks ago"
    ```

  - Between 2 dates:

    ```shell
    git log --after="2014-07-01" --before="2014-07-04"
    ```

  - Commits up to a specific date:

    ```shell
    git log --until="2024-07-04"
    ```

- Filter by author:
  - Single author by name/email:

    ```shell
    git log --author="John"
    git log --author="john@example.com"
    ```

  - Multiple authors:

    ```shell
    git log --author="John\|Mary"
    ```

### Search history

Search by commit message:

```shell
git log --grep="text-to-search"
```

Search for changes to specific files:

```shell
git log -- filename.ext subdir/filename2.ext
```

Search for when text was added/removed:

```shell
git log -S"Hello, World!"
```

### Commit ranges

Show commits in a range. The general syntax is:

```shell
git log <since>..<until>
```

Show commits between 2 branches:

```shell
git log main..feat/branch-name
```

### Merge commit filtering

Hide merge commits:

```shell
git log --no-merges
```

Show *only* merge commits:

```shell
git log --merges
```

### Group commits by author

Show commits grouped by author with:

```shell
git shortlog
```

Sort authors by commit count:

```shell
git shortlog -n
```

### Pretty formatting

Git supports multiple display formats, which you can customize with the `--pretty` flag.

- Show each commit on a single line

  ```shell
  git log --pretty=oneline
  ```

- Include commit hash, author, commit title:

  ```shell
  git log --pretty=short
  ```

- Include commit hash, author, date, full commit message

  ```shell
  git log --pretty=medium
  ```

- Include commit hash, full author (name + email), date, full commit message:

  ```shell
  git log --pretty=full
  ```

- Include commit hash, full author (name + email), date (author and commit date included), and full metadata

  ```shell
  git log --pretty=fuller
  ```

- Format commits like email patches:

  ```shell
  git log --pretty=email
  ```

- Show complete internal commit object data:

  ```shell
  git log --pretty=raw
  ```

- Customize pretty formatting (see [Pretty log formatting tokens section](https://git-scm.com/docs/pretty-formats#Documentation/pretty-formats.txt-formatformat-string) for more):

  ```shell
  git log --pretty=format:"%h - %an, %ar : %s"
  ```

  - Example output:

    ```shell
    8f3c1a2 - John Smith, 2 days ago : Fix parsing bug
    ```

#### Pretty log formatting tokens

These tables are cheatsheets for the `git log --pretty=format:"..."` functionality.

- Commit Information

| Token | Meaning                   |
| ----- | ------------------------- |
| `%H`  | Full commit hash          |
| `%h`  | Abbreviated commit hash   |
| `%T`  | Full tree hash            |
| `%t`  | Abbreviated tree hash     |
| `%P`  | Parent hashes             |
| `%p`  | Abbreviated parent hashes |

- Author Information

| Token | Meaning              |
| ----- | -------------------- |
| `%an` | Author name          |
| `%ae` | Author email         |
| `%ad` | Author date          |
| `%ar` | Relative author date |
| `%ai` | ISO-8601 author date |

- Committer Information

| Token | Meaning                 |
| ----- | ----------------------- |
| `%cn` | Committer name          |
| `%ce` | Committer email         |
| `%cd` | Committer date          |
| `%cr` | Relative committer date |

- Commit Message Information

| Token | Meaning            |
| ----- | ------------------ |
| `%s`  | Subject            |
| `%f`  | Sanitized subject  |
| `%b`  | Body               |
| `%B`  | Raw body + subject |

- Reference Information

| Token | Meaning                       |
| ----- | ----------------------------- |
| `%d`  | Ref names                     |
| `%D`  | Ref names without parentheses |

- Formatting/Color Tokens

| Token        | Meaning     |
| ------------ | ----------- |
| `%n`         | Newline     |
| `%%`         | Literal `%` |
| `%C(red)`    | Red text    |
| `%C(green)`  | Green text  |
| `%C(yellow)` | Yellow text |
| `%C(reset)`  | Reset color |

#### Pretty log formatting examples

- Compact colored log:

  ```shell
  git log --pretty=format:"%C(yellow)%h%C(reset) %C(green)%ar%C(reset) %C(blue)%an%C(reset) %s"
  ```

  - Example:

    ```shell
    8f3c1a2 2 days ago John Smith Fix parsing bug
    ```

- Multi-line format

  ```shell
  git log --pretty=format:"Commit: %h%nAuthor: %an%nDate: %ar%n%n%s%n"
  ```

- Graph + pretty format

  ```shell
  git log --graph \
  --pretty=format:"%C(yellow)%h%C(reset) - %C(green)(%cr)%C(reset) %s %C(blue)<%an>%C(reset)" \
  --abbrev-commit
  ```

## Git config

### Set preserve colors in less

```shell
git config --global core.pager "less -R"
```

## Git aliases

Set an alias with `git config [--global] alias.<alias> "cmd"`.

- Decorated git log graph:

  ```shell
  git config --global alias.lg "log --graph --oneline --decorate"
  ```

- Advanced git log:

  ```shell
  git config --global alias.lga \
    "log --graph --pretty=format:'%C(yellow)%h%C(reset) %C(green)%ar%C(reset) %C(blue)%an%C(reset)%C(auto)%d%C(reset) %s' --abbrev-commit"
  ```

## Troubleshooting

### Rewrite Git commit history

Say you accidentally commit code from a work `git config user.name`/`git config user.email`, or using a username/email for the wrong forge (i.e. committing code to a repository hosted on Github using your Gitlab username & email). You can rewrite a git repository's history, changing all commits authored by a user to a different desired user.

First, install the [`git-filter-repo` plugin for Git](https://github.com/newren/git-filter-repo/blob/main/INSTALL.md) with Python:

```shell
pip install git-filter-repo
```

Then, follow the steps below to rewrite history. Wherever you see `Old Name` and `old.email@example.com`, use the current author's information (the one you want to rewrite), and use the new/desired username/email for `New Name` and `new.email@example.com`.

- Clone the repository using the `--bare` flag:

  ```shell
  git clone --bare git@github.com:user/repo.git
  ```

- Run the following command to replace all instances of the old/wrong name & email with a different Git user:

  ```shell
  git filter-branch --env-filter '
  if [ "$GIT_COMMITTER_NAME" = "Old Name" ] && [ "$GIT_COMMITTER_EMAIL" = "old.email@example.com" ]; then
      GIT_COMMITTER_NAME="New Name"
      GIT_COMMITTER_EMAIL="new.email@example.com"
      GIT_AUTHOR_NAME="New Name"
      GIT_AUTHOR_EMAIL="new.email@example.com"
  fi
  ' --tag-name-filter cat -- --branches --tags
  ```

- Clean the `reflog` and run Git garbage collection to remove any cached history with the old Git user:

  ```shell
  git reflog expire --expire=now --all
  git gc --prune=now
  ```

- Force push the changes back to the remote (note: if you have branch protection rules that prevent pushing to `main` directly, you will need to temporarily disable that to complete this step):

  ```shell
  git push --force --all
  git push --force --tags
  ```

- Verify the rewrite succeeded by cloning the repository to a new path and search the log for the old user:

  ```shell
  mkdir ~/tmp
  git clone git@github.com:user/repo.git ~/tmp/repo
  cd ~/tmp/repo
  git log --author="Old Name"
  ```

  - You should not see any results for the old user. If you do, look back through the history of your commands and make sure there were no errors during the process.

#### Bash script to rewrite history

On a Linux or Mac system, you can use this Bash script to automate the steps above. Run the script with `--help` to see the usage menu.

```shell
#!/usr/bin/env bash
set -euo pipefail

## Find a Python interpreter for pip installs
PYTHON_BIN=""
for bin in python3 python py py3 python; do
  if command -v "$bin" >/dev/null 2>&1; then
    PYTHON_BIN=$bin
    break
  fi
done

## Ensure git-filter-repo is installed (try uv first, then pip)
if ! command -v git-filter-repo >/dev/null 2>&1; then
  echo "git-filter-repo not found."

  ## Install with uv, if available
  if command -v uv >/dev/null 2>&1; then
    echo "uv found. Installing git-filter-repo as a tool"
    uv tool install git-filter-repo
    export PATH="$HOME/.local/bin:$PATH"
  fi

  if [[ -z "$PYTHON_BIN" ]]; then
    echo "No Python interpreter found. Please install Python or uv."
    exit 1
  fi

  ## Fallback to Python
  echo "Using $PYTHON_BIN to install git-filter-repo via pip"
  "$PYTHON_BIN" -m pip install --user git-filter-repo
  export PATH="$HOME/.local/bin:$PATH"
fi

## Test git-filter-repo was installed correctly
if ! command -v git-filter-repo >/dev/null 2>&1; then
  echo "git-filter-repo still not found after installation attempts."
  exit 1
fi

## Function to print help menu/usage
usage() {
  echo ""
  echo "Usage: $0 [--force] \\"
  echo "          --repo-url git@github.com:user/repo.git \\"
  echo "          --source-email 'old.email@example.com' \\"
  echo "          --target-email 'new.email@example.com' \\"
  echo "          [--source-name 'Old Name'] \\"
  echo "          [--target-name 'New Name']"
  echo ""

  exit 1
}

## Default vars
REPO_URL=""
SRC_EMAIL=""
TGT_EMAIL=""
SRC_NAME=""
TGT_NAME=""
FORCE_PUSH=""

## Parse arguments
while [[ $# -gt 0 ]]; do
  case $1 in
  --repo-url)
    REPO_URL="$2"
    shift 2
    ;;
  --source-email)
    SRC_EMAIL="$2"
    shift 2
    ;;
  --target-email)
    TGT_EMAIL="$2"
    shift 2
    ;;
  --source-name)
    SRC_NAME="$2"
    shift 2
    ;;
  --target-name)
    TGT_NAME="$2"
    shift 2
    ;;
  --force)
    FORCE_PUSH=1
    shift
    ;;
  -h | --help)
    usage
    ;;
  *)
    echo "Invalid argument: $1"
    usage
    ;;
  esac
done

if [[ -z "$REPO_URL" || -z "$SRC_EMAIL" || -z "$TGT_EMAIL" ]]; then
  echo "Missing required arguments."
  usage
fi

echo "Repo URL: $REPO_URL"
echo "Replacing source email <$SRC_EMAIL> with target email <$TGT_EMAIL>"

if [[ -n "$SRC_NAME" ]]; then
  echo "Source name: $SRC_NAME"
fi
if [[ -n "$TGT_NAME" ]]; then
  echo "Target name: $TGT_NAME"
fi

## Create temporary directory to clone repo into
TMP_DIR=$(mktemp -d)
echo "Mirror cloning repository into temporary directory: $TMP_DIR"
git clone --mirror "$REPO_URL" "$TMP_DIR/repo"
cd "$TMP_DIR/repo"

echo "Rewriting commit history emails with git-filter-repo"
## Read history with source username/email, replace with target
COMMIT_CALLBACK="
if commit.author_email.decode('utf-8') == '$SRC_EMAIL':
    commit.author_email = b'$TGT_EMAIL'
"
if [[ -n "$TGT_NAME" ]]; then
    COMMIT_CALLBACK+="
    commit.author_name = b'$TGT_NAME'
"
fi

COMMIT_CALLBACK+="
if commit.committer_email.decode('utf-8') == '$SRC_EMAIL':
    commit.committer_email = b'$TGT_EMAIL'
"
if [[ -n "$TGT_NAME" ]]; then
    COMMIT_CALLBACK+="
    commit.committer_name = b'$TGT_NAME'
"
fi

echo "Generated callback:"
echo "$COMMIT_CALLBACK"

## Rewrite history
git filter-repo --force --commit-callback "$COMMIT_CALLBACK"

echo "git filter-repo completed successfully"
echo "Verifying first few commits after rewrite"
git log --all --pretty=format:"%h %ad %an <%ae>" --date=iso -5

echo "Removing backup refs"
##  Remove all refs with old data
git for-each-ref --format='%(refname)' refs/original | xargs -r git update-ref -d
git for-each-ref --format='%(refname)' refs/backup | xargs -r git update-ref -d

echo "Expiring reflogs and pruning unreachable objects"
## Expire old data locally before pushing
git reflog expire --expire=now --all
git gc --prune=now --aggressive

echo "Verifying no lingering commits with source email anywhere"

## Get list of commits with source email
SOURCE_COMMITS=$(git for-each-ref --format='%(refname)' | while read -r ref; do
  git log "$ref" --pretty=format:"%H%x09%ad%x09%an%x09%ae%x09%cN%x09%cE" --date=iso |
    awk -v src_email="$SRC_EMAIL" '
        $4 == src_email || $6 == src_email {
            print FILENAME "\t" $0
        }' FILENAME="$ref"
done)

## If any commits remain with source email, exit
if [[ -n "$SOURCE_COMMITS" ]]; then
  echo ""
  echo "[ERROR] Some commits still contain the source email <$SRC_EMAIL>:"
  echo ""
  echo -e "Ref\tCommit\tDate\tAuthorName\tAuthorEmail\tCommitterName\tCommitterEmail"
  echo "$SOURCE_COMMITS"
  echo ""
  echo "Aborting push."

  exit 2
fi

echo "Removing local refs under refs/merge-requests/ to avoid push errors"
## Remove all refs under refs/merge-requests/
git for-each-ref --format='%(refname)' refs/merge-requests | xargs -r -n 1 git update-ref -d || true

echo "Adding remote origin after filter-repo cleanup"
## Re-add origin (git-filter-repo removes it)
git remote add origin "$REPO_URL"

echo "Pushing all branches and tags to origin forcibly"
## Push rewritten histories back up
if [[ -n "${FORCE_PUSH:-}" ]]; then
  if ! git push origin --force --all; then
    echo ""
    echo "[ERROR] Failed to push rewritten commits to origin."

    exit 1
  fi

  if ! git push origin --force --tags; then
    echo ""
    echo "[ERROR] Failed to push rewritten commits to origin."

    exit 1
  fi

else
  if ! git push origin --force --all; then
    echo ""
    echo "[ERROR] Failed to push rewritten commits to origin."

    exit 1
  fi

  if ! git push origin --force --tags; then
    echo ""
    echo "[ERROR] Failed to push rewritten commits to origin."

    exit 1
  fi
fi

echo "Successfully rewrote all commits and pushed to remote."
echo "Temporary repo location: $TMP_DIR/repo"

exit 0
```
