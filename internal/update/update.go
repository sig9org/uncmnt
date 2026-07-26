// Package update implements self-updating uncmnt from its GitHub releases.
package update

import (
	"context"
	"errors"
	"fmt"

	selfupdate "github.com/creativeprojects/go-selfupdate"
)

const repoSlug = "sig9org/uncmnt"

// SelfUpdate replaces the running binary with the latest GitHub release,
// unless currentVersion is already up to date.
func SelfUpdate(currentVersion string) error {
	ctx := context.Background()

	latest, found, err := selfupdate.DetectLatest(ctx, selfupdate.ParseSlug(repoSlug))
	if err != nil {
		return fmt.Errorf("detect latest version: %w", err)
	}
	if !found {
		return errors.New("no release found for this platform")
	}

	if latest.LessOrEqual(currentVersion) {
		fmt.Printf("current version (%s) is already the latest\n", currentVersion)
		return nil
	}

	exe, err := selfupdate.ExecutablePath()
	if err != nil {
		return fmt.Errorf("locate executable: %w", err)
	}

	if err := selfupdate.UpdateTo(ctx, latest.AssetURL, latest.AssetName, exe); err != nil {
		return fmt.Errorf("update binary: %w", err)
	}

	fmt.Printf("updated to version %s\n", latest.Version())
	return nil
}
