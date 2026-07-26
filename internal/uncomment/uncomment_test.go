package uncomment

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

func TestIsIgnorable(t *testing.T) {
	cases := []struct {
		name string
		line string
		want bool
	}{
		{"empty", "", true},
		{"whitespace only", "   \t  ", true},
		{"hash comment", "# comment", true},
		{"bang comment", "! comment", true},
		{"semicolon comment", "; comment", true},
		{"indented comment", "   # comment", true},
		{"code line", "line one", false},
		{"code with hash later", "x = 1 # not a leading comment", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := IsIgnorable(tc.line); got != tc.want {
				t.Errorf("IsIgnorable(%q) = %v, want %v", tc.line, got, tc.want)
			}
		})
	}
}

func TestProcessFile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "input.txt")
	input := "# comment\nline one\n\n! bang\n  ; semi\nline two\n   \n"
	if err := os.WriteFile(path, []byte(input), 0o644); err != nil {
		t.Fatalf("write input: %v", err)
	}

	var buf bytes.Buffer
	if err := ProcessFile(&buf, path); err != nil {
		t.Fatalf("ProcessFile: %v", err)
	}

	want := "line one\nline two\n"
	if got := buf.String(); got != want {
		t.Errorf("ProcessFile output = %q, want %q", got, want)
	}
}

func TestProcessFileMissing(t *testing.T) {
	var buf bytes.Buffer
	if err := ProcessFile(&buf, filepath.Join(t.TempDir(), "missing.txt")); err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}
