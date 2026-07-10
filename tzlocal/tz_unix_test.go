//go:build !windows
// +build !windows

package tzlocal

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func Test_inferFromPath(t *testing.T) {
	tests := []struct {
		name    string
		file    string
		wantErr bool
	}{
		{
			name:    "Asia/Tokyo",
			file:    "/usr/share/zoneinfo/Asia/Tokyo",
			wantErr: false,
		},
		{
			name:    "America/Chicago",
			file:    "/usr/share/zoneinfo/America/Chicago",
			wantErr: false,
		},
		{
			name:    "America/Kentucky/Monticello",
			file:    "/usr/share/zoneinfo/America/Kentucky/Monticello",
			wantErr: false,
		},
		{
			name:    "America/Argentina/Buenos_Aires",
			file:    "/usr/share/zoneinfo/America/Argentina/Buenos_Aires",
			wantErr: false,
		},
		{
			name:    "UTC",
			file:    "/usr/share/zoneinfo/UTC",
			wantErr: false,
		},
		{ // MacOS - /var/db location
			name:    "Europe/Berlin",
			file:    "/var/db/timezone/zoneinfo/Europe/Berlin",
			wantErr: false,
		},
		{ // MacOS - Final symlink target for /usr/share/zoneinfo and /var/db/timezone/zoneinfo
			name:    "Europe/Berlin",
			file:    "/usr/share/zoneinfo.default/Europe/Berlin",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := inferFromPath(tt.file)
			if (err != nil) != tt.wantErr {
				t.Errorf("inferFromPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.name {
				t.Errorf("inferFromPath() = %v, want %v", got, tt.name)
			}
			_, err = time.LoadLocation(tt.name)
			if err != nil {
				t.Errorf("can't load timezone %s: %s", tt.name, err.Error())
			}
		})
	}
}

func TestLocalTZFromPath(t *testing.T) {
	t.Run("absolute target", func(t *testing.T) {
		root := t.TempDir()
		target := createZoneinfoFile(t, root, "Europe/Tallinn")
		localtime := filepath.Join(root, "etc", "localtime")
		mustMkdirAll(t, filepath.Dir(localtime))
		mustSymlink(t, target, localtime)

		assertLocalTZ(t, localtime, "Europe/Tallinn")
	})

	t.Run("relative target through symlinked parent", func(t *testing.T) {
		root := t.TempDir()
		realRoot := filepath.Join(root, "real")
		createZoneinfoFile(t, realRoot, "Europe/Berlin")
		realEtc := filepath.Join(realRoot, "etc")
		mustMkdirAll(t, realEtc)
		mustSymlink(t,
			filepath.Join("..", "usr", "share", "zoneinfo", "Europe", "Berlin"),
			filepath.Join(realEtc, "localtime"),
		)
		mustSymlink(t, realEtc, filepath.Join(root, "etc"))

		assertLocalTZ(t, filepath.Join(root, "etc", "localtime"), "Europe/Berlin")
	})

	t.Run("regular file through symlinked parent", func(t *testing.T) {
		root := t.TempDir()
		target := createZoneinfoFile(t, root, "localtime")
		mustSymlink(t, filepath.Dir(target), filepath.Join(root, "etc"))

		_, err := localTZ(filepath.Join(root, "etc", "localtime"))
		if err == nil || !strings.Contains(err.Error(), "is not a symlink") {
			t.Fatalf("localTZ() error = %v, want a not-a-symlink error", err)
		}
	})

	t.Run("missing localtime", func(t *testing.T) {
		localtime := filepath.Join(t.TempDir(), "etc", "localtime")

		_, err := localTZ(localtime)
		if err == nil || !strings.Contains(err.Error(), "failed to stat") {
			t.Fatalf("localTZ() error = %v, want a failed-to-stat error", err)
		}
	})

	t.Run("dangling symlink", func(t *testing.T) {
		root := t.TempDir()
		localtime := filepath.Join(root, "etc", "localtime")
		mustMkdirAll(t, filepath.Dir(localtime))
		mustSymlink(t, filepath.Join(root, "usr", "share", "zoneinfo", "Missing"), localtime)

		_, err := localTZ(localtime)
		if err == nil {
			t.Fatal("localTZ() error = nil, want an error for a dangling symlink")
		}
		if strings.Contains(err.Error(), "failed to stat") {
			t.Fatalf("localTZ() error = %v, final symlink exists and should not fail its stat check", err)
		}
	})
}

func createZoneinfoFile(t *testing.T, root, name string) string {
	t.Helper()
	p := filepath.Join(root, "usr", "share", "zoneinfo", filepath.FromSlash(name))
	mustMkdirAll(t, filepath.Dir(p))
	if err := os.WriteFile(p, nil, 0o644); err != nil {
		t.Fatalf("os.WriteFile(%q) failed: %v", p, err)
	}
	return p
}

func mustMkdirAll(t *testing.T, path string) {
	t.Helper()
	if err := os.MkdirAll(path, 0o755); err != nil {
		t.Fatalf("os.MkdirAll(%q) failed: %v", path, err)
	}
}

func mustSymlink(t *testing.T, oldname, newname string) {
	t.Helper()
	if err := os.Symlink(oldname, newname); err != nil {
		t.Fatalf("os.Symlink(%q, %q) failed: %v", oldname, newname, err)
	}
}

func assertLocalTZ(t *testing.T, localtime, want string) {
	t.Helper()
	got, err := localTZ(localtime)
	if err != nil {
		t.Fatalf("localTZ() failed: %v", err)
	}
	if got != want {
		t.Fatalf("localTZ() = %q, want %q", got, want)
	}
}

func TestRuntimeTZ(t *testing.T) {
	got, err := RuntimeTZ()
	if err != nil {
		t.Errorf("RuntimeTZ failed with %v", err.Error())
	}
	if got == "" {
		t.Error("RuntimeTZ returned empty timezone string")
	}
}
