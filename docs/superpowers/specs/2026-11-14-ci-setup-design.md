# CI Setup — Design

**Date:** 2026-11-14

## Intent

`baseline` has no CI today — only `.github/dependabot.yml`. Dependabot PRs
(and everything else) get merged on manual judgment alone; PR #7 was merged
after the author ran `go build`/`go test` locally by hand because there was
nothing else to check. This adds automated build/vet/test/lint signal to
every PR and push to `main`.

## Goals

- `.github/workflows/ci.yml` with three jobs:
  - `test`: matrix (`ubuntu-latest`, `macos-latest`) — `go build ./...`,
    `go vet ./...`, `go test ./... -race`.
  - `lint`: `golangci-lint` with default rules, `ubuntu-latest` only.
  - `markdownlint`: lints `*.md` using the existing `.markdownlint.json`,
    `ubuntu-latest` only.
- Triggers: `pull_request` and `push` to `main`.
- Go version read from `go.mod` via `actions/setup-go`'s `go-version-file`,
  not hardcoded, so it tracks future `go.mod` bumps automatically.
- Add a `github-actions` entry to `.github/dependabot.yml` (now that
  workflow files exist to keep current), using a two-group pattern —
  `major-updates` and `minor-and-patch`, each batching its own update
  types, 7-day cooldown — unlike `gomod`'s single `minor-and-patch` group
  (majors ungrouped there). GitHub Actions gets majors grouped too because
  coordinated action releases (e.g. `actions/upload-artifact` and
  `actions/download-artifact`) need their major bumps to land together.
- CI must be green on its first run. Local dry-runs of each tool surfaced
  pre-existing findings unrelated to CI itself:
  - `go vet`: two `fmt.Println` calls in
    `internal/sources/bitbucket/bitbucket.go` end their format string with
    a redundant `\n` (lines 125, 149).
  - `golangci-lint` (default rules): 5 `errcheck` findings — unchecked
    `resp.Body.Close()` / permission-restore error returns in
    `internal/git/git.go`, `internal/git/git_test.go`,
    `internal/sources/bitbucket/bitbucket.go`, and
    `internal/sources/github/github.go` (two spots).
  - `markdownlint-cli2`: a missing trailing newline in
    `.github/copilot-instructions.md`, and a misaligned table (`MD060`) in
    `README.md`'s "SSH vs HTTPS" section.
  All of the above get fixed in a PR that lands *before* the CI workflow
  PR, so the workflow itself never has to tolerate known-red findings.

## Non-goals

- Release automation (goreleaser, cross-compiled binaries, automated GitHub
  Releases) — tracked separately in
  [issue #8](https://github.com/jonasbn/baseline/issues/8), deliberately
  deferred while the project is alpha.
- A Windows runner — nothing in the tool targets Windows-specific paths or
  permission semantics.
- Code coverage reporting/badges.
- Enabling GitHub branch-protection "required status checks" — a repo
  setting, not code; can be done separately once the workflow has run a
  few times successfully.

## Anti-goal

CI heavy or slow enough to discourage fast iteration while still in alpha
testing (hence: no coverage gating, no extra matrix dimensions beyond OS,
default linter rules rather than a strict custom `.golangci.yml`).

## Constraints

- Go pinned at `1.25.1` in `go.mod`.
- `cmd`, `internal/sources/github`, `internal/sources/bitbucket` have no
  tests today. Adding tests for them is out of scope — CI runs what exists.
- `.markdownlint.json` already exists and disables `line-length`,
  `first-line-heading`, `no-emphasis-as-heading`, `link-fragments`,
  `no-bare-urls`, `no-duplicate-heading`. The CI markdownlint job must
  respect it, not a default ruleset.

## Approaches considered (CI workflow shape)

- **A — minimal, single OS.** One job, `ubuntu-latest` only. Cheapest, but
  misses platform-specific bugs.
- **B — matrix across OS.** Same steps, matrix over `ubuntu-latest` and
  `macos-latest`. Worth it because `internal/git`'s `GitOps` chmod's repos
  read-only/writable, and permission semantics genuinely differ between
  Linux and macOS.
- **C — B plus a separate lint job (chosen).** Adds `golangci-lint` as its
  own job so a style/lint finding doesn't get conflated with a real
  regression in the status checks. Confirmed worthwhile: a lint job would
  have caught the `bitbucket.go` vet issue before it reached `main`.

Chosen: **C**, plus a `markdownlint` job (agreed in conversation), plus
fixing all pre-existing findings in their own PR first (agreed in
conversation, after discovering the errcheck and markdownlint findings
during local dry-runs).

## Decisions made in conversation

- Scope is CI only; release automation is tracked in issue #8, not built now.
- The 5 `errcheck` findings get fixed now, not deferred and not
  lint-config'd away.
- Pre-existing fixes (vet, errcheck, markdownlint) land in their own
  branch/PR, separate from the PR that introduces the CI workflow itself,
  so the workflow PR's diff is purely additive and its first run is green
  without having also changed application behavior.

## Left to the implementer

- Exact action versions, SHA-pinned per GitHub's security hardening
  guidance (with the resolved version as a trailing comment) rather than
  pinned to a mutable tag — see plan for the full YAML and resolved SHAs.
- Exact wording of commit messages and PR descriptions.
- Exact whitespace used to realign the `README.md` table (mechanical,
  verified by running `markdownlint-cli2` until clean rather than
  hand-specified here).
