package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadNoConfig(t *testing.T) {
	t.Chdir(t.TempDir())
	t.Setenv("HOME", t.TempDir())

	cfg, path, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg != nil || path != "" {
		t.Errorf("Load() = %+v, %q, want nil config and empty path", cfg, path)
	}
}

func TestLoadFromWorkingDirectory(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)
	t.Setenv("HOME", t.TempDir())

	write(t, filepath.Join(dir, "config.yml"), "comment_prefixes:\n  - \"//\"\n")

	cfg, path, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg == nil {
		t.Fatal("Load() cfg = nil, want non-nil")
	}
	if want := []string{"//"}; len(cfg.CommentPrefixes) != 1 || cfg.CommentPrefixes[0] != want[0] {
		t.Errorf("cfg.CommentPrefixes = %v, want %v", cfg.CommentPrefixes, want)
	}
	if path != "config.yml" {
		t.Errorf("path = %q, want %q", path, "config.yml")
	}
}

func TestLoadFromHomeDirectory(t *testing.T) {
	t.Chdir(t.TempDir())
	home := t.TempDir()
	t.Setenv("HOME", home)

	confDir := filepath.Join(home, ".config", "uncmnt")
	if err := os.MkdirAll(confDir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	write(t, filepath.Join(confDir, "config.yaml"), "comment_prefixes:\n  - \"REM\"\n")

	cfg, path, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if cfg == nil {
		t.Fatal("Load() cfg = nil, want non-nil")
	}
	if want := "REM"; len(cfg.CommentPrefixes) != 1 || cfg.CommentPrefixes[0] != want {
		t.Errorf("cfg.CommentPrefixes = %v, want [%q]", cfg.CommentPrefixes, want)
	}
	if path != filepath.Join(confDir, "config.yaml") {
		t.Errorf("path = %q, want %q", path, filepath.Join(confDir, "config.yaml"))
	}
}

func TestLoadWorkingDirectoryTakesPriority(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)
	home := t.TempDir()
	t.Setenv("HOME", home)

	confDir := filepath.Join(home, ".config", "uncmnt")
	if err := os.MkdirAll(confDir, 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	write(t, filepath.Join(confDir, "config.yml"), "comment_prefixes:\n  - \"HOME\"\n")
	write(t, filepath.Join(dir, "config.yml"), "comment_prefixes:\n  - \"CWD\"\n")

	cfg, path, err := Load()
	if err != nil {
		t.Fatalf("Load: %v", err)
	}
	if want := "CWD"; cfg == nil || len(cfg.CommentPrefixes) != 1 || cfg.CommentPrefixes[0] != want {
		t.Errorf("cfg.CommentPrefixes = %v, want [%q]", cfg.CommentPrefixes, want)
	}
	if path != "config.yml" {
		t.Errorf("path = %q, want %q", path, "config.yml")
	}
}

func write(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}
