# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What this is

`uncmnt` is a small Go CLI that strips comment/blank lines from one or more
text files and prints the result to stdout. It self-updates from GitHub
releases via `go-selfupdate`.

## Commands

This project uses [Task](https://taskfile.dev) (`Taskfile.yml`) as its command
runner. `project.ini` is loaded as dotenv and defines `BINARY_NAME=uncmnt` and
`VERSION_PKG=github.com/sig9org/uncmnt/internal/version`.

- `task build` (alias `b`) — build a binary for the current platform into
  `dist/`, with `-X $VERSION_PKG.Version=$VERSION` set via ldflags (`VERSION`
  comes from `git describe --tags --always`, falling back to `dev`).
- `task build-all` (alias `ba`) — cross-compile release binaries for all
  platforms in `PLATFORMS` into `dist/`, named
  `{BINARY_NAME}_{RELEASE_VERSION}_{goos}_{goarch}` (`RELEASE_VERSION` is
  `git describe --tags --abbrev=0`, i.e. the current tag without a commit
  suffix). This naming must stay in sync with what `go-selfupdate` expects to
  find on GitHub releases (see `internal/update/update.go`).
- `task test` (alias `t`) — runs `go vet ./...` then `go test ./...`.
- `task cleanup` (alias `c`) — removes stray editor/OS/terraform files
  (`.idea`, `.vscode`, `.DS_Store`, etc.) from the working directory.
- Plain Go also works directly: `go build .`, `go vet ./...`,
  `go test ./...`, and `go test ./internal/uncomment/ -run TestIsIgnorable`
  for a single test.

There is no linter configured beyond `go vet` (run as part of `task test`).

## Architecture

- `main.go` — all CLI wiring lives here. `run(args, stdout, stderr) int` is
  the testable entry point (`main()` just calls it with `os.Args`/real
  stdio); tests in `main_test.go` call `run` directly rather than exec'ing
  the binary. Flags are hand-registered pairs (`-h`/`-help`,
  `-v`/`-version`, `-u`/`-update`) via a fresh `flag.FlagSet` per run rather
  than the package-level `flag` API. No args, `-h`, or an empty file list
  all print the same help text and exit 0.
- `internal/uncomment/` — the actual line-filtering logic, decoupled from
  I/O concerns: `IsIgnorable(line string) bool` is the single source of
  truth for what counts as a comment (blank/whitespace-only, or first
  non-whitespace char is `!`, `#`, or `;`), and `ProcessFile(w io.Writer,
  path string) error` streams a file's surviving lines to a writer. When
  multiple files are passed on the CLI, `main.go` calls `ProcessFile` for
  each in order against the same buffered stdout writer — output is
  concatenated per-file, not merged/interleaved.
- `internal/version/` — a single mutable `Version` var, overwritten at
  build time by ldflags (see Commands above). Defaults to `"dev"` in
  unbuilt/`go run` contexts.
- `internal/update/` — wraps `github.com/creativeprojects/go-selfupdate`
  to let the binary replace itself from the `sig9org/uncmnt` GitHub repo's
  releases. Version comparisons use the injected `version.Version`, so
  self-update correctness depends on binaries actually being built through
  `task build`/`task build-all` (not bare `go build`) so that var is
  populated with a real semver tag.

## Release/versioning notes

Version strings come from `git tag` (via `git describe`), not a hardcoded
constant — there is no version to bump in source when cutting a release.
