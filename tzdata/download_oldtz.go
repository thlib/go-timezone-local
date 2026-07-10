package tzdata

import (
	"bytes"
	"strings"
)

const tzdataURL = `https://data.iana.org/time-zones/tzdata-latest.tar.gz`

// DownloadOldNames fetches the list of old tz names and returns a mapping
func DownloadOldNames() (map[string]string, error) {
	var err error

	buf, err := downloadHTTP(tzdataURL)
	if err != nil {
		return nil, err
	}

	var b bytes.Buffer
	err = ExtractTarGz(&buf, &b, "backward")
	if err != nil {
		return nil, err
	}

	backward := make(map[string]string)
	for _, line := range strings.Split(b.String(), "\n") {
		if len(strings.TrimSpace(line)) == 0 {
			continue
		}
		if line[0:1] == "#" {
			continue
		}
		if !strings.HasPrefix(line, "Link") {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) != 3 {
			continue
		}
		backward[parts[2]] = parts[1]
	}

	return backward, nil
}
