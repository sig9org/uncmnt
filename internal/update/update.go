// Package update implements self-updating uncmnt from its GitHub releases.
package update

import (
	"context"
	"fmt"

	selfupdate "github.com/sig9org/selfupdate-go"
)

const repoSlug = "sig9org/uncmnt"

// SelfUpdate replaces the running binary with the latest GitHub release after
// verifying its SHA-256 checksum against checksums.txt. If currentVersion is
// already up to date, no download or replacement is performed. debugf, if
// non-nil, is called with progress details for each step of the update.
func SelfUpdate(currentVersion string, debugf func(format string, args ...any)) error {
	if debugf == nil {
		debugf = func(string, ...any) {}
	}
	ctx := context.Background()

	debugf("checking latest release for %s", repoSlug)
	updater, err := selfupdate.New(selfupdate.Config{
		Repository: repoSlug,
		Validator:  selfupdate.SHA256Validator{AssetName: "checksums.txt"},
	})
	if err != nil {
		return fmt.Errorf("configure self-update: %w", err)
	}

	result, err := updater.Update(ctx, currentVersion)
	if err != nil {
		return fmt.Errorf("update binary: %w", err)
	}

	debugf("latest release detected: %s", result.LatestVersion)
	if result.Updated {
		fmt.Printf("updated to version %s\n", result.LatestVersion)
	} else {
		fmt.Printf("current version (%s) is already the latest\n", currentVersion)
	}
	return nil
}
