# uncmnt

A CLI tool written in Go that removes comments from a specified file and outputs the clean code to stdout.

## Usage

```
uncmnt [options] <file> [file...]
```

Lines are dropped if they are blank, contain only whitespace, or start with `!`, `#`, or `;` (leading whitespace ignored). Multiple files are each processed and their results printed to stdout in order.

### Options

| Flag             | Description                          |
| ---------------- | ------------------------------------ |
| `-h`, `-help`    | Show help                            |
| `-v`, `-version` | Show version information             |
| `-u`, `-update`  | Update uncmnt to the latest release  |

## Install

```
go install github.com/sig9org/uncmnt@latest
```

## License

MIT — see [LICENSE](LICENSE).
