package tzdata

import (
	"bytes"
	"testing"
)

func TestUpdateWindowsTZMapping(t *testing.T) {
	var buf bytes.Buffer
	err := UpdateWindowsTZMapping(&buf)
	if err != nil {
		t.Errorf("error: %v", err)
	}
}
