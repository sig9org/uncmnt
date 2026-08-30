// Package version holds uncmnt's build-time version metadata.
package version

// Version is set at build time via -ldflags, based on the current git tag.
var Version = "dev"

// Full returns the version from the git tag. Commit IDs are intentionally not
// included in the displayed version.
func Full() string {
	return Version
}
