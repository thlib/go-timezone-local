package tzdata

import (
	"os"
	"testing"
)

func TestDownloadOldNames(t *testing.T) {
	if os.Getenv("TZLOCAL_LIVE_TESTS") == "" {
		t.Skip("set TZLOCAL_LIVE_TESTS=1 to run live download tests")
	}
	_, err := DownloadOldNames()
	if err != nil {
		t.Errorf("error: %v", err)
	}
}
