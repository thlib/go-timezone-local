package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestMappingComparison(t *testing.T) {
	old := map[string]string{}
	for i := 0; i < 10; i++ {
		old[fmt.Sprint(i)] = "old"
	}
	for _, tc := range []struct {
		name    string
		edit    func(map[string]string)
		wantErr bool
	}{
		{"unchanged", func(m map[string]string) {}, false},
		{"exactly 90 percent", func(m map[string]string) { m["0"] = "new" }, false},
		{"below 90 percent", func(m map[string]string) { m["0"], m["1"] = "new", "new" }, true},
		{"one removal at 90 percent", func(m map[string]string) { delete(m, "0") }, true},
		{"addition", func(m map[string]string) { m["extra"] = "new" }, false},
		{"large unexpected expansion", func(m map[string]string) { m["a"], m["b"] = "new", "new" }, true},
		{"empty", func(m map[string]string) {
			for k := range m {
				delete(m, k)
			}
		}, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			next := map[string]string{}
			for k, v := range old {
				next[k] = v
			}
			tc.edit(next)
			_, err := compareTable("test", old, next)
			if (err != nil) != tc.wantErr {
				t.Fatalf("error = %v, wantErr %v", err, tc.wantErr)
			}
		})
	}
}

const validMapping = `package tzlocal
var WinTZtoIANA = map[string]string{"UTC": "Etc/UTC"}
var IANAtoWinTZ = map[string]string{"Etc/UTC": "UTC"}
`

func TestUpdatePreservesFile(t *testing.T) {
	for _, tc := range []struct {
		name, candidate string
		wantErr         bool
	}{
		{"date only", "// a new date\n" + validMapping, false},
		{"syntax error", "broken", true},
		{"wrong table type", strings.Replace(validMapping, "map[string]string", "map[int]string", 1), true},
		{"lost map", "package tzlocal", true},
		{"empty table", strings.Replace(validMapping, `"UTC": "Etc/UTC"`, "", 1), true},
		{"duplicate key", strings.Replace(validMapping, `"UTC": "Etc/UTC"`, `"UTC": "Etc/UTC", "UTC": "Etc/UTC"`, 1), true},
		{"invalid zone", strings.ReplaceAll(validMapping, "Etc/UTC", "Invalid/Zone"), true},
		{"broken reverse reference", strings.Replace(validMapping, `"Etc/UTC": "UTC"`, `"Etc/UTC": "Missing"`, 1), true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			path := filepath.Join(t.TempDir(), "tzmapping.go")
			original := []byte("// old date\n" + validMapping)
			if err := os.WriteFile(path, original, 0644); err != nil {
				t.Fatal(err)
			}
			err := update(path, func(w io.Writer) error {
				_, err := io.WriteString(w, strings.Replace(tc.candidate, "package tzlocal\n", "", 1))
				return err
			})
			if (err != nil) != tc.wantErr {
				t.Fatalf("error = %v, wantErr %v", err, tc.wantErr)
			}
			got, err := os.ReadFile(path)
			if err != nil {
				t.Fatal(err)
			}
			if !bytes.Equal(got, original) {
				t.Fatal("existing file changed")
			}
		})
	}
}

func TestUpdateDownloadFailurePreservesFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "tzmapping.go")
	if err := os.WriteFile(path, []byte(validMapping), 0644); err != nil {
		t.Fatal(err)
	}
	err := update(path, func(w io.Writer) error {
		io.WriteString(w, "partial download")
		return errors.New("upstream unavailable")
	})
	if err == nil {
		t.Fatal("expected generation failure")
	}
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != validMapping {
		t.Fatal("failed generation modified existing file")
	}
}

func TestUpdateWritesSafeAddition(t *testing.T) {
	path := filepath.Join(t.TempDir(), "tzmapping.go")
	old := "package tzlocal\nvar WinTZtoIANA = map[string]string{\n"
	for i := 0; i < 10; i++ {
		old += fmt.Sprintf("%q: %q,\n", fmt.Sprint(i), "Etc/UTC")
	}
	old += "}\nvar IANAtoWinTZ = map[string]string{\n"
	for i := 0; i < 10; i++ {
		old += fmt.Sprintf("%q: %q,\n", fmt.Sprintf("Etc/GMT+%d", i), "0")
	}
	old += "\"Etc/UTC\": \"0\",\n}\n"
	if err := os.WriteFile(path, []byte(old), 0644); err != nil {
		t.Fatal(err)
	}
	next := strings.Replace(old, `"Etc/UTC": "0",`, "\"Etc/UTC\": \"0\",\n\"UTC\": \"0\",", 1)
	if err := update(path, func(w io.Writer) error {
		_, err := io.WriteString(w, strings.TrimPrefix(next, "package tzlocal\n"))
		return err
	}); err != nil {
		t.Fatal(err)
	}
	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Contains(got, []byte(`"UTC"`)) {
		t.Fatal("safe addition not written")
	}
}
