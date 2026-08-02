// Package config loads uncmnt's optional comment-rule overrides from a
// config.yml/config.yaml file.
package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

const dirName = ".config/uncmnt"

var configNames = []string{"config.yml", "config.yaml"}

// Config holds user-configurable overrides for uncmnt's default behavior.
type Config struct {
	// CommentPrefixes, when the config file that provided it exists,
	// replaces the tool's built-in set of line-comment markers.
	CommentPrefixes []string `yaml:"comment_prefixes"`
}

// Load searches, in priority order, for config.yml/config.yaml in the
// current working directory and then under ~/.config/uncmnt/, and parses
// the first one found. It returns a nil Config and an empty path if none of
// the candidate locations exist.
func Load() (*Config, string, error) {
	for _, path := range candidatePaths() {
		data, err := os.ReadFile(path)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return nil, "", fmt.Errorf("read config %s: %w", path, err)
		}

		cfg, err := parse(path, data)
		if err != nil {
			return nil, "", err
		}
		return cfg, path, nil
	}
	return nil, "", nil
}

// LoadFrom parses the config file at path exactly, without searching any
// other location. Unlike Load, a missing file is an error here since the
// caller explicitly requested this path (e.g. via -c/-config).
func LoadFrom(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config %s: %w", path, err)
	}
	return parse(path, data)
}

func parse(path string, data []byte) (*Config, error) {
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config %s: %w", path, err)
	}
	return &cfg, nil
}

func candidatePaths() []string {
	paths := make([]string, 0, len(configNames)*2)
	paths = append(paths, configNames...)

	if home, err := os.UserHomeDir(); err == nil {
		for _, name := range configNames {
			paths = append(paths, filepath.Join(home, dirName, name))
		}
	}
	return paths
}
