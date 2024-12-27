//go:build !windows
// +build !windows

package tzlocal

import (
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

func TestRuntimeTZ(t *testing.T) {
	got, err := RuntimeTZ()
	if err != nil {
		t.Errorf("RuntimeTZ failed with %v", err.Error())
	}
	if got == "" {
		t.Error("RuntimeTZ returned empty timezone string")
	}
}
