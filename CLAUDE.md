# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

```bash
# Build
go build

# Run tests
go test ./...

# Run a single test
go test ./internal/git/... -run TestFunctionName

# Install from source
go install github.com/jonasbn/baseline@latest
```

## Architecture

`baseline` is a CLI tool (built with [cobra](https://github.com/spf13/cobra)) that clones Git repositories from GitHub or Bitbucket as bare, read-only repositories for code searching.

**Entry point:** `main.go` → `cmd.Execute()`

**`cmd/`** — cobra subcommands (`init`, `discover`, `clone`, `update`). Global flags (directory, tokens, org, source, threads) are defined in `cmd/root.go` as package-level vars shared across subcommands.

**`internal/types/`** — shared data model: `Repository`, `RepositorySource` interface, `CloneResult`, `UpdateResult`. All source clients must implement `RepositorySource`.

**`internal/sources/github/`** and **`internal/sources/bitbucket/`** — platform-specific API clients that implement `RepositorySource.GetRepositories()`. Adding a new source means implementing this interface and wiring it into the `switch source` blocks in each command.

**`internal/worker/`** — `WorkerPool` runs clone/update operations concurrently using goroutines. It reads from a repository channel and writes results to a buffered result channel.

**`internal/git/`** — `GitOps` wraps actual `git clone --bare` and `git fetch` shell calls, plus permission management (chmod to read-only after clone, write before update).

**Data flow:** `cmd` → `RepositorySource.GetRepositories()` → `WorkerPool` → `GitOps` → filesystem

## Key Design Decisions

- Repositories are cloned **bare** (`git clone --bare`) — no working tree, just object store. This is intentional for searching.
- Cloned repos are set **read-only recursively** to prevent accidental modifications.
- The `--ssh` flag on `clone`/`update` replaces `CloneURL` with `SSHURL` before passing to the worker pool; if no SSH URL exists it falls back to HTTPS with a warning.
- The `clone` command silently skips existing repos (reported as "skipped", not "failed").
- The `update` command silently skips repos that don't exist locally (reported as "skipped").
