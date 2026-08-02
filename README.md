<p align="center">
  <img src="https://raw.githubusercontent.com/sig9org/uncmnt/main/assets/logo.webp" alt="uncmnt">
</p>

# uncmnt

[![Go Reference](https://pkg.go.dev/badge/github.com/sig9org/uncmnt.svg)](https://pkg.go.dev/github.com/sig9org/uncmnt)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)

A CLI tool written in Go that removes comments from a specified file and outputs the clean code to stdout.

## Usage

```
uncmnt [options] <file> [file...]
```

Lines are dropped if they are blank, contain only whitespace, or start with `!`, `#`, or `;` (leading whitespace ignored). Multiple files are each processed and their results printed to stdout in order.

### Options

| Flag              | Description                                |
| ----------------- | ------------------------------------------ |
| `-h`, `-help`      | Show help                                  |
| `-v`, `-version`   | Show version information (tag plus commit hash, when available) |
| `-update`          | Update uncmnt to the latest release        |
| `-debug`           | Print detailed, timestamped debug output to stdout |
| `-c`, `-config`    | Force uncmnt to use this config file path  |

### Output

Messages are color-coded by level: normal output is uncolored, warnings are orange, errors are red, and (with `-debug`) debug messages are gray and timestamped.

## Configuration

By default, a line is a comment if it starts with `!`, `#`, or `;`. To override which prefixes count as comments, add a `config.yml` (or `config.yaml`) with a `comment_prefixes` list, at either location below. The first one found is used, in this priority order:

1. `config.yml`/`config.yaml` in the directory you run `uncmnt` from
2. `~/.config/uncmnt/config.yml`/`config.yaml`

Pass `-c`/`-config` with a path to use that config file instead, bypassing the search above entirely. Unlike the automatic search, a missing `-c`/`-config` path is an error.

```yaml
comment_prefixes:
  - "#"
  - "//"
  - "--"
```

When a config file is found and its `comment_prefixes` list is non-empty, it replaces the built-in prefixes entirely for that run. See [config.yml.example](config.yml.example).

## Install

```
go install github.com/sig9org/uncmnt@latest
```

## License

MIT — see [LICENSE](LICENSE).
