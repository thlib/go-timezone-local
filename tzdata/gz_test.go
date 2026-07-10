package tzdata

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"strings"
	"testing"
)

func TestExtractTarGz(t *testing.T) {
	archive := makeTarGz(t, []tarEntry{
		{name: "directory/", typeflag: tar.TypeDir},
		{name: "other", body: "not this file"},
		{name: "backward", body: "wanted contents"},
	})
	var dst bytes.Buffer

	if err := ExtractTarGz(bytes.NewReader(archive), &dst, "backward"); err != nil {
		t.Fatalf("ExtractTarGz() failed: %v", err)
	}
	if got, want := dst.String(), "wanted contents"; got != want {
		t.Fatalf("ExtractTarGz() wrote %q, want %q", got, want)
	}
}

func TestExtractTarGzMissingFile(t *testing.T) {
	archive := makeTarGz(t, []tarEntry{{name: "other", body: "contents"}})

	err := ExtractTarGz(bytes.NewReader(archive), &bytes.Buffer{}, "backward")
	if err == nil || !strings.Contains(err.Error(), "not found") {
		t.Fatalf("ExtractTarGz() error = %v, want a not-found error", err)
	}
}

func TestExtractTarGzInvalidGzip(t *testing.T) {
	err := ExtractTarGz(strings.NewReader("not gzip"), &bytes.Buffer{}, "backward")
	if err == nil {
		t.Fatal("ExtractTarGz() error = nil, want an invalid-gzip error")
	}
}

type tarEntry struct {
	name     string
	body     string
	typeflag byte
}

func makeTarGz(t *testing.T, entries []tarEntry) []byte {
	t.Helper()
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gz)
	for _, entry := range entries {
		typeflag := entry.typeflag
		if typeflag == 0 {
			typeflag = tar.TypeReg
		}
		header := &tar.Header{
			Name:     entry.name,
			Mode:     0o644,
			Size:     int64(len(entry.body)),
			Typeflag: typeflag,
		}
		if err := tw.WriteHeader(header); err != nil {
			t.Fatalf("WriteHeader() failed: %v", err)
		}
		if _, err := tw.Write([]byte(entry.body)); err != nil {
			t.Fatalf("Write() failed: %v", err)
		}
	}
	if err := tw.Close(); err != nil {
		t.Fatalf("tar Close() failed: %v", err)
	}
	if err := gz.Close(); err != nil {
		t.Fatalf("gzip Close() failed: %v", err)
	}
	return buf.Bytes()
}
