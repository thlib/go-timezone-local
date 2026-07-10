package tzlocal

import (
	"os"
	"testing"
)

func TestEnvTZ(t *testing.T) {
	original, wasSet := os.LookupEnv("TZ")
	t.Cleanup(func() {
		if wasSet {
			_ = os.Setenv("TZ", original)
		} else {
			_ = os.Unsetenv("TZ")
		}
	})

	tests := []struct {
		value string
		want  string
	}{
		{value: "", want: "UTC"},
		{value: "UTC", want: "UTC"},
		{value: ":UTC", want: "UTC"},
		{value: "not/a/timezone", want: "UTC"},
	}
	for _, test := range tests {
		if err := os.Setenv("TZ", test.value); err != nil {
			t.Fatalf("Setenv() failed: %v", err)
		}
		got, ok := EnvTZ()
		if !ok || got != test.want {
			t.Errorf("EnvTZ() with TZ=%q = (%q, %t), want (%q, true)", test.value, got, ok, test.want)
		}
	}

	if err := os.Unsetenv("TZ"); err != nil {
		t.Fatalf("Unsetenv() failed: %v", err)
	}
	if got, ok := EnvTZ(); ok || got != "" {
		t.Errorf("EnvTZ() with TZ unset = (%q, %t), want (\"\", false)", got, ok)
	}
}
