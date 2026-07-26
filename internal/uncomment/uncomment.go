// Package uncomment strips comment and blank lines from text files.
package uncomment

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// commentPrefixes lists the leading characters that mark a line as a comment.
var commentPrefixes = []byte{'!', '#', ';'}

// IsIgnorable reports whether line should be dropped: it is blank, made of
// only whitespace, or starts with one of commentPrefixes (after leading
// whitespace is trimmed).
func IsIgnorable(line string) bool {
	trimmed := strings.TrimSpace(line)
	if trimmed == "" {
		return true
	}
	for _, prefix := range commentPrefixes {
		if trimmed[0] == prefix {
			return true
		}
	}
	return false
}

// ProcessFile reads path and writes every non-comment, non-blank line to w.
func ProcessFile(w io.Writer, path string) error {
	f, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("cannot open file: %w", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	scanner.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	for scanner.Scan() {
		line := scanner.Text()
		if IsIgnorable(line) {
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
