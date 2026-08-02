# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What this is

`uncmnt` is a small Go CLI that strips comment/blank lines from one or more
text files and prints the result to stdout. Which lines count as comments
can be overridden by a YAML config file. It self-updates from GitHub
releases via `go-selfupdate`.

## Commands

This project uses [Task](https://taskfile.dev) (`Taskfile.yml`) as its command
runner. `project.ini` is loaded as dotenv and defines `BINARY_NAME=uncmnt`,
`GITHUB_REPO=uncmnt`, and `GITHUB_USER=sig9org`. `Taskfile.yml` itself only
declares thin wrappers (with `_time:start`/`_time:end` timing hooks) around
shared tasks pulled in from `https://raw.githubusercontent.com/sig9org/tasks`
(`_cleanup.yml`, `_go.yml`, `_time.yml`) — the real task bodies live there,
not in this repo.

- `task go-build` (alias `gb`) — build a binary for the current platform into
  `dist/`, with `-X $VERSION_PKG.Version=$VERSION` set via ldflags, where
  `VERSION_PKG=github.com/sig9org/uncmnt/internal/version` and `VERSION`
  comes from `git describe --tags --always`, falling back to `dev`.
- `task go-all-build` (alias `ga`) — cross-compile release binaries for all
  platforms into `dist/`, named `{BINARY_NAME}_{RELEASE_VERSION}_{goos}_{goarch}`
  (`RELEASE_VERSION` is `git describe --tags --abbrev=0`, i.e. the current
  tag without a commit suffix). This naming must stay in sync with what
  `go-selfupdate` expects to find on GitHub releases (see
  `internal/update/update.go`).
- `task go-clean` (alias `gc`) — empties `dist/`.
- `task go-test` (alias `gt`) — runs `go vet ./...` then `go test ./...`.
- `task go-register` (alias `gr`) — hits the pkg.go.dev/proxy.golang.org
  endpoints for the current tag to request indexing after a release.
- `task go-version` (alias `gv`) — runs the CLI's own `-v` via `go run` with
  the same ldflags as `go-build`.
- `task cleanup` (alias `c`) — removes stray editor/OS/terraform files
  (`.idea`, `.vscode`, `.DS_Store`, etc.) from the working directory.
- Plain Go also works directly: `go build .`, `go vet ./...`,
  `go test ./...`, and `go test ./internal/uncomment/ -run TestIsIgnorable`
  for a single test.

There is no linter configured beyond `go vet` (run as part of `task go-test`).

## Architecture

- `main.go` — all CLI wiring lives here. `run(args, stdout, stderr) int` is
  the testable entry point (`main()` just calls it with `os.Args`/real
  stdio); tests in `main_test.go` call `run` directly rather than exec'ing
  the binary. Flags are hand-registered pairs (`-h`/`-help`, `-v`/`-version`,
  `-c`/`-config`, plus standalone `-update` and `-debug`) via a fresh
  `flag.FlagSet` per run rather than the package-level `flag` API. No args,
  `-h`, or an empty file list all print the same help text and exit 0.
  Output is colorized by message level (errors red, warnings orange, debug
  gray via `newDebugLogger`, normal output uncolored) using raw ANSI escape
  codes defined in `main.go`. `-debug` timestamps and prints internal
  progress to stdout regardless of which stream a given message would
  otherwise use.
- `internal/uncomment/` — the actual line-filtering logic, decoupled from
  I/O concerns. `Rules{Prefixes []string}` defines what counts as a comment
  (always blank/whitespace-only, plus any line whose trimmed content starts
  with one of `Prefixes`); `DefaultRules()` returns the built-in `!`/`#`/`;`
  set, and the package-level `IsIgnorable(line string) bool` is a
  convenience wrapper over `DefaultRules().IsIgnorable`. `ProcessFile(w
  io.Writer, path string, rules Rules) error` streams a file's surviving
  lines to a writer under the given rules. When multiple files are passed
  on the CLI, `main.go` calls `ProcessFile` for each in order against the
  same buffered stdout writer — output is concatenated per-file, not
  merged/interleaved.
- `internal/config/` — loads optional comment-rule overrides from a
  `config.yml`/`config.yaml`. `Load()` checks, in priority order, the
  current working directory then `~/.config/uncmnt/`, and returns the first
  file found (nil, nil path if none exist). `LoadFrom(path)` parses an
  exact path instead, for `main.go`'s `-c`/`-config` flag — unlike `Load()`,
  a missing file here is an error, since the caller explicitly named it.
  `main.go` only swaps in the config's `comment_prefixes` over
  `uncomment.DefaultRules()` when that list is non-empty — an empty or
  key-less config file is a no-op rather than a footgun that disables
  comment stripping. See `config.yml.example` for the schema; `config.yml`
  in the repo root is a git-ignored local scratch file for manual testing.
- `internal/version/` — a single mutable `Version` var, overwritten at
  build time by ldflags (see Commands above). Defaults to `"dev"` in
  unbuilt/`go run` contexts. `Commit()` reads the short VCS revision from
  `runtime/debug.ReadBuildInfo()` (Go's automatic build stamping, not an
  ldflag), and `Full()` combines the two for display in `-v`/`-h` output.
- `internal/update/` — wraps `github.com/creativeprojects/go-selfupdate`
  to let the binary replace itself from the `sig9org/uncmnt` GitHub repo's
  releases. Version comparisons use the injected `version.Version`, so
  self-update correctness depends on binaries actually being built through
  `task go-build`/`task go-all-build` (not bare `go build`) so that var is
  populated with a real semver tag.

## Release/versioning notes

Version strings come from `git tag` (via `git describe`), not a hardcoded
constant — there is no version to bump in source when cutting a release.
