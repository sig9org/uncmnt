// Package uncomment strips comment and blank lines from text files.
package uncomment

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// defaultCommentPrefixes lists the leading strings that mark a line as a
// comment under the tool's built-in rules.
var defaultCommentPrefixes = []string{"!", "#", ";"}

// Rules defines which lines ProcessFile treats as comments: blank or
// whitespace-only lines, plus any line whose trimmed content starts with
// one of Prefixes.
type Rules struct {
	Prefixes []string
}

// DefaultRules returns uncmnt's built-in comment rules.
func DefaultRules() Rules {
	return Rules{Prefixes: defaultCommentPrefixes}
}

// IsIgnorable reports whether line should be dropped under the default
// rules: it is blank, made of only whitespace, or starts with one of the
// default comment prefixes (after leading whitespace is trimmed).
func IsIgnorable(line string) bool {
	return DefaultRules().IsIgnorable(line)
}

// IsIgnorable reports whether line should be dropped under r: it is blank,
// made of only whitespace, or starts with one of r.Prefixes (after leading
// whitespace is trimmed).
func (r Rules) IsIgnorable(line string) bool {
	trimmed := strings.TrimSpace(line)
	if trimmed == "" {
		return true
	}
	for _, prefix := range r.Prefixes {
		if prefix != "" && strings.HasPrefix(trimmed, prefix) {
			return true
		}
	}
	return false
}

// ProcessFile reads path and writes every line not ignorable under rules to w.
func ProcessFile(w io.Writer, path string, rules Rules) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("cannot open file: %w", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	for scanner.Scan() {
		line := scanner.Text()
		if rules.IsIgnorable(line) {
			continue
		}
		if _, err := fmt.Fprintln(w, line); err != nil {
			return fmt.Errorf("write output: %w", err)
		}
	}
	if err := scanner.Err(); err != nil {
		return fmt.Errorf("read file: %w", err)
	}
	return nil
}
