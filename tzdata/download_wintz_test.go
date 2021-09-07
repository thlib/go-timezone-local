package tzdata

import "testing"

func TestDownloadWindowsZones(t *testing.T) {
	_, err := DownloadWindowsZones()
	if err != nil {
		t.Errorf("error: %v", err)
	}
}
