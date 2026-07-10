package tzdata

import (
	"bytes"
	"os"
	"testing"
)

func TestUpdateWindowsTZMapping(t *testing.T) {
	if os.Getenv("TZLOCAL_LIVE_TESTS") == "" {
		t.Skip("set TZLOCAL_LIVE_TESTS=1 to run live download tests")
	}
	var buf bytes.Buffer
	err := UpdateWindowsTZMapping(&buf)
	if err != nil {
		t.Errorf("error: %v", err)
	}
}
