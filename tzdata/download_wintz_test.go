package tzdata

import (
	"io"
	"net/http"
	"os"
	"strings"
	"testing"
)

func TestDownloadWindowsZones(t *testing.T) {
	if os.Getenv("TZLOCAL_LIVE_TESTS") == "" {
		t.Skip("set TZLOCAL_LIVE_TESTS=1 to run live download tests")
	}
	_, err := DownloadWindowsZones()
	if err != nil {
		t.Errorf("error: %v", err)
	}
}

func TestDownloadWindowsZonesRejectsInvalidData(t *testing.T) {
	for _, body := range []string{
		"not XML",
		"<supplementalData/>",
		"<supplementalData><windowsZones><mapTimezones/></windowsZones></supplementalData>",
	} {
		t.Run(body, func(t *testing.T) {
			if _, err := parseWindowsZones([]byte(body), "test data"); err == nil {
				t.Fatal("parseWindowsZones() error = nil, want an error")
			}
		})
	}
}

func TestDownloadHTTP(t *testing.T) {
	client := &http.Client{Transport: roundTripperFunc(func(r *http.Request) (*http.Response, error) {
		status := http.StatusOK
		body := "contents"
		switch r.URL.Path {
		case "/ok":
		case "/large":
			body = strings.Repeat("x", maxDownloadSize+1)
		default:
			status = http.StatusNotFound
		}
		return &http.Response{
			StatusCode: status,
			Status:     http.StatusText(status),
			Body:       io.NopCloser(strings.NewReader(body)),
			Header:     make(http.Header),
		}, nil
	})}

	buf, err := downloadHTTPWithClient(client, "http://example.test/ok")
	if err != nil {
		t.Fatalf("downloadHTTP() failed: %v", err)
	}
	if got, want := buf.String(), "contents"; got != want {
		t.Fatalf("downloadHTTP() = %q, want %q", got, want)
	}
	if _, err := downloadHTTPWithClient(client, "http://example.test/missing"); err == nil {
		t.Fatal("downloadHTTP() accepted a non-success HTTP status")
	}
	if _, err := downloadHTTPWithClient(client, "http://example.test/large"); err == nil {
		t.Fatal("downloadHTTP() accepted an oversized response")
	}
}

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(request *http.Request) (*http.Response, error) {
	return f(request)
}
