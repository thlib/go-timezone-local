package tzdata

import "testing"

func TestDownloadOldNames(t *testing.T) {
	_, err := DownloadOldNames()
	if err != nil {
		t.Errorf("error: %v", err)
	}
}
