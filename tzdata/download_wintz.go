package tzdata

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"net/http"
)

const winZonesURL = `https://raw.githubusercontent.com/unicode-org/cldr/master/common/supplemental/windowsZones.xml`

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
	var data SupplementalData

	resp, err := http.Get(winZonesURL)
	if err != nil {
		return data, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return data, fmt.Errorf("failed to download \"%v\", http error: %v", winZonesURL, resp.StatusCode)
	}

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return data, err
	}

	xml.Unmarshal(bytes, &data)

	return data, nil
}
