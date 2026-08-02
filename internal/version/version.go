// Package version holds uncmnt's build-time version metadata.
package version

import "runtime/debug"

// Version is set at build time via -ldflags, based on the current git tag.
var Version = "dev"

// Commit returns the short VCS revision the running binary was built from,
// as recorded by the Go toolchain's automatic build info stamping. It
// returns "" if that information isn't available (e.g. built outside a git
// checkout).
func Commit() string {
	info, ok := debug.ReadBuildInfo()
	if !ok {
		return ""
	}
	for _, s := range info.Settings {
		if s.Key == "vcs.revision" {
			if len(s.Value) > 7 {
				return s.Value[:7]
			}
			return s.Value
		}
	}
	return ""
}

// Full returns Version, plus the short commit hash in parentheses when
// available.
func Full() string {
	if c := Commit(); c != "" {
		return Version + " (" + c + ")"
	}
	return Version
}
