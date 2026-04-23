# clean-checkout

clean-checkout is a small set of tools for safely refreshing a local git branch — check out the target branch, sync with a remote, and delete local branches that are not protected. The primary CLI in this repo is git-cc (implemented in `cmd/git-cc`) which is intended to be used directly or integrated into IDE extensions (future folders: `intellij-clean-checkout`, `vs-clean-checkout`).

## Quick overview

- CLI entrypoint: [`main.main`](cmd/git-cc/main.go)
- Core planning logic: [`refresh.PlanRefresh`](internal/refresh/plan.go)
- Execution logic: [`refresh.Execute`](internal/refresh/execute.go)
- Remote validation: [`refresh.validateRemote`](internal/refresh/remote.go)
- Git helpers: [`git.ListBranches`](internal/git/branches.go), [`git.ListRemotes`](internal/git/remotes.go), and git command model/runner in [internal/git/command.go](internal/git/command.go) and [internal/git/runner.go](internal/git/runner.go)

## Usage

Build the CLI and run:

```sh
go build -o git-cc ./cmd/git-cc
./git-cc [options] <branch>
```

Or run directly with `go run`:

```sh
go run ./cmd/git-cc -- <options> <branch>
```

Flags (see implementation in [cmd/git-cc/main.go](cmd/git-cc/main.go)):
- `-reset` : Hard reset to origin/<branch>
- `-dry-run` : Show what would be executed without running
- `-json` : Output structured JSON (for IDE integration)
- `-only-destructive` : Run only destructive commands
- `-confirm` : Ask for confirmation before execution
- `-remote` : Git remote to use (default: "origin")

Example:

```sh
# preview actions on branch "feature" without making changes
./git-cc -dry-run feature

# run reset+dry-run and get JSON output
./git-cc -json -reset -dry-run feature
```

## How it works (high level)

1. The CLI parses flags and builds an `refresh.Options` object.
2. `refresh.PlanRefresh` builds a command plan:
   - Check out requested branch (`git.Checkout`).
   - Fetch from the configured remote (`git.Fetch`).
   - Either reset (`git.ResetHard`) or pull (`git.PullFFOnly`).
   - Enumerate local branches (`git.ListBranches`) and append `git.DeleteBranch` commands for non-protected branches.
   - Optionally filter to only destructive commands.
3. `refresh.Execute` runs the plan (or simulates it on dry-run). Failures are returned (and printed as JSON when `-json` is set).

Protected branches are defined in [internal/refresh/plan.go](internal/refresh/plan.go).

## Development & tests

- Build the CLI: `go build ./cmd/git-cc`
- Run tests for the CLI package: `go test ./cmd/git-cc` — see tests in [cmd/git-cc/main_test.go](cmd/git-cc/main_test.go)
  - The tests build the CLI binary and exercise flag handling, JSON error output, and run end-to-end dry-run cases against temporary git repositories.
- Core packages:
  - [internal/refresh](internal/refresh)
  - [internal/git](internal/git)

## Notes & future integrations

- This repository is designed to be extended with IDE/plugin integrations:
  - `intellij-clean-checkout` — an IntelliJ plugin that would call the CLI or embed similar logic
  - `vs-clean-checkout` — a VS Code extension that would call the CLI or integrate the logic
- The CLI supports JSON output to make integrating into IDEs and tooling straightforward.