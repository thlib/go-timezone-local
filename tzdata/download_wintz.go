package tzdata

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"time"
)

const winZonesURL = `https://raw.githubusercontent.com/unicode-org/cldr/master/common/supplemental/windowsZones.xml`

const maxDownloadSize = 10 << 20

type SupplementalData struct {
	WindowsZones WindowsZones `xml:"windowsZones"`
}
type Version struct {
	Number string `xml:"number,attr"`
}
type WindowsZones struct {
	MapTimezones []MapTimezones `xml:"mapTimezones"`
}
type MapTimezones struct {
	MapZone []MapZone `xml:"mapZone"`
	Type    string    `xml:"type,attr"`
}
type MapZone struct {
	Other     string `xml:"other,attr"`
	Territory string `xml:"territory,attr"`
	Type      string `xml:"type,attr"`
}

// DownloadWindowsZones fetches Windows mapping info from unicode.org
func DownloadWindowsZones() (SupplementalData, error) {
	return downloadWindowsZones(winZonesURL)
}

func downloadWindowsZones(target string) (SupplementalData, error) {
	buf, err := downloadHTTP(target)
	if err != nil {
		return SupplementalData{}, err
	}
	return parseWindowsZones(buf.Bytes(), target)
}

func parseWindowsZones(contents []byte, source string) (SupplementalData, error) {
	var data SupplementalData

	if err := xml.Unmarshal(contents, &data); err != nil {
		return data, fmt.Errorf("failed to parse %q: %w", source, err)
	}
	foundMappings := false
	for _, mappings := range data.WindowsZones.MapTimezones {
		if len(mappings.MapZone) > 0 {
			foundMappings = true
			break
		}
	}
	if !foundMappings {
		return data, fmt.Errorf("failed to parse %q: no time zone mappings found", source)
	}

	return data, nil
}

func downloadHTTP(target string) (bytes.Buffer, error) {
	client := http.Client{Timeout: 30 * time.Second}
	return downloadHTTPWithClient(&client, target)
}

func downloadHTTPWithClient(client *http.Client, target string) (bytes.Buffer, error) {
	var buf bytes.Buffer
	resp, err := client.Get(target)
	if err != nil {
		return buf, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return buf, fmt.Errorf("failed to download %q, HTTP status: %s", target, resp.Status)
	}

	if _, err := io.Copy(&buf, io.LimitReader(resp.Body, maxDownloadSize+1)); err != nil {
		return buf, err
	}
	if buf.Len() > maxDownloadSize {
		return bytes.Buffer{}, fmt.Errorf("failed to download %q: response exceeds %d bytes", target, maxDownloadSize)
	}

	return buf, nil
}
