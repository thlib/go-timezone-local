//go:build !windows
// +build !windows

package tzlocal

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func inferFromPath(p string) (string, error) {
	var name string
	var err error

	parts := strings.Split(p, string(filepath.Separator))
	for i := range parts {
		if parts[i] == "zoneinfo" || parts[i] == "zoneinfo.default" {
			parts = parts[i+1:]
			break
		}
	}

	if len(parts) < 1 {
		err = fmt.Errorf("cannot infer timezone name from path: %q", p)
		return name, err
	}

	return filepath.Join(parts...), nil
}

// LocalTZ gets the timezone name by resolving symlink /etc/localtime that points to a timezone
// file.
func LocalTZ() (string, error) {
	const localZoneFile = "/etc/localtime"

	p, err := filepath.EvalSymlinks(localZoneFile)
	if err != nil {
		// Specific error if we can't access localZoneFile
		if _, err := os.Lstat(localZoneFile); err != nil {
			return "", fmt.Errorf("failed to stat %q: %w", localZoneFile, err)
		}
		// Generic EvalSymlink error
		return "", err
	}

	if p == localZoneFile {
		return "", fmt.Errorf("%q is not a symlink - cannot infer name", localZoneFile)
	}

	return inferFromPath(p)
}
