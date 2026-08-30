package main

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRunNoArgsShowsHelp(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := run(nil, &stdout, &stderr)
	if code != 0 {
		t.Errorf("exit code = %d, want 0", code)
	}
	if !strings.Contains(stdout.String(), "Usage:") {
		t.Errorf("stdout missing usage text: %q", stdout.String())
	}
}

func TestRunVersion(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := run([]string{"-version"}, &stdout, &stderr)
	if code != 0 {
		t.Errorf("exit code = %d, want 0", code)
	}
	if !strings.Contains(stdout.String(), toolName) {
		t.Errorf("stdout missing tool name: %q", stdout.String())
	}
	if strings.Contains(stdout.String(), "(") {
		t.Errorf("stdout must not include a git commit ID: %q", stdout.String())
	}
}

func TestRunHelp(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := run([]string{"-h"}, &stdout, &stderr)
	if code != 0 {
		t.Errorf("exit code = %d, want 0", code)
	}
	if !strings.Contains(stdout.String(), "Usage:") {
		t.Errorf("stdout missing usage text: %q", stdout.String())
	}
}

func TestRunHelpOptionsAreAlphabeticalAndAligned(t *testing.T) {
	var stdout, stderr bytes.Buffer
	if code := run([]string{"-help"}, &stdout, &stderr); code != 0 {
		t.Fatalf("exit code = %d, want 0", code)
	}

	help := stdout.String()
	options := []string{
		"  -c, -config <path>",
		"      -debug",
		"  -h, -help",
		"      -update",
		"  -v, -version",
	}
	previous := -1
	for _, option := range options {
		position := strings.Index(help, option)
		if position <= previous {
			t.Fatalf("option %q is out of order in help output", option)
		}
		previous = position
	}

	for _, description := range []string{
		"config file path",
		"show debug output",
		"show this help message",
		"update uncmnt",
		"show version information",
	} {
		lineStart := strings.LastIndex(help[:strings.Index(help, description)], "\n") + 1
		if got := strings.Index(help[lineStart:], description); got != 22 {
			t.Errorf("description %q starts at column %d, want 22", description, got)
		}
	}
}

func TestRunProcessesFiles(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "input.txt")
	if err := os.WriteFile(path, []byte("# comment\nkeep me\n"), 0o644); err != nil {
		t.Fatalf("write input: %v", err)
	}

	var stdout, stderr bytes.Buffer
	code := run([]string{path}, &stdout, &stderr)
	if code != 0 {
		t.Errorf("exit code = %d, want 0; stderr=%q", code, stderr.String())
	}
	if got, want := stdout.String(), "keep me\n"; got != want {
		t.Errorf("stdout = %q, want %q", got, want)
	}
}

func TestRunDebug(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "input.txt")
	if err := os.WriteFile(path, []byte("# comment\nkeep me\n"), 0o644); err != nil {
		t.Fatalf("write input: %v", err)
	}

	var stdout, stderr bytes.Buffer
	code := run([]string{"-debug", path}, &stdout, &stderr)
	if code != 0 {
		t.Errorf("exit code = %d, want 0; stderr=%q", code, stderr.String())
	}
	if !strings.Contains(stdout.String(), "DEBUG:") {
		t.Errorf("stdout missing debug output: %q", stdout.String())
	}
	if !strings.Contains(stdout.String(), "keep me") {
		t.Errorf("stdout missing file content: %q", stdout.String())
	}
}

func TestRunUsesConfigFileCommentPrefixes(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)
	t.Setenv("HOME", t.TempDir())

	if err := os.WriteFile(filepath.Join(dir, "config.yml"), []byte("comment_prefixes:\n  - \"//\"\n"), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	inputPath := filepath.Join(dir, "input.txt")
	if err := os.WriteFile(inputPath, []byte("# not a comment now\n// this is a comment\nkeep me\n"), 0o644); err != nil {
		t.Fatalf("write input: %v", err)
	}

	var stdout, stderr bytes.Buffer
	code := run([]string{inputPath}, &stdout, &stderr)
	if code != 0 {
		t.Errorf("exit code = %d, want 0; stderr=%q", code, stderr.String())
	}
	if got, want := stdout.String(), "# not a comment now\nkeep me\n"; got != want {
		t.Errorf("stdout = %q, want %q", got, want)
	}
}

func TestRunUsesExplicitConfigFlag(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)
	t.Setenv("HOME", t.TempDir())

	// A config.yml in the working directory would normally win, but an
	// explicit -config path must take priority over it.
	if err := os.WriteFile(filepath.Join(dir, "config.yml"), []byte("comment_prefixes:\n  - \"IGNORED\"\n"), 0o644); err != nil {
		t.Fatalf("write cwd config: %v", err)
	}
	altPath := filepath.Join(dir, "alt.yml")
	if err := os.WriteFile(altPath, []byte("comment_prefixes:\n  - \"//\"\n"), 0o644); err != nil {
		t.Fatalf("write alt config: %v", err)
	}

	inputPath := filepath.Join(dir, "input.txt")
	if err := os.WriteFile(inputPath, []byte("# not a comment now\n// this is a comment\nkeep me\n"), 0o644); err != nil {
		t.Fatalf("write input: %v", err)
	}

	var stdout, stderr bytes.Buffer
	code := run([]string{"-config", altPath, inputPath}, &stdout, &stderr)
	if code != 0 {
		t.Errorf("exit code = %d, want 0; stderr=%q", code, stderr.String())
	}
	if got, want := stdout.String(), "# not a comment now\nkeep me\n"; got != want {
		t.Errorf("stdout = %q, want %q", got, want)
	}
}

func TestRunExplicitConfigMissingIsError(t *testing.T) {
	dir := t.TempDir()
	inputPath := filepath.Join(dir, "input.txt")
	if err := os.WriteFile(inputPath, []byte("keep me\n"), 0o644); err != nil {
		t.Fatalf("write input: %v", err)
	}

	var stdout, stderr bytes.Buffer
	code := run([]string{"-c", filepath.Join(dir, "missing.yml"), inputPath}, &stdout, &stderr)
	if code != 1 {
		t.Errorf("exit code = %d, want 1", code)
	}
	if stderr.Len() == 0 {
		t.Error("expected error message on stderr for missing -c config file")
	}
}

func TestRunMissingFile(t *testing.T) {
	var stdout, stderr bytes.Buffer
	code := run([]string{filepath.Join(t.TempDir(), "missing.txt")}, &stdout, &stderr)
	if code != 1 {
		t.Errorf("exit code = %d, want 1", code)
	}
	if stderr.Len() == 0 {
		t.Error("expected error message on stderr")
	}
}
