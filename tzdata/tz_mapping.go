/*
 * Based on a python package to update windows mappings
 * @see https://github.com/regebro/tzlocal/blob/master/update_windows_mappings.py
 */

package tzdata

import (
	"bytes"
	"fmt"
	"io"
	"sort"
	"strings"
	"time"
)

func UpdateWindowsTZMapping(target io.Writer) error {

	backward, err := DownloadOldNames()
	if err != nil {
		return err
	}

	data, err := DownloadWindowsZones()
	if err != nil {
		return err
	}

	win_tz := make(map[string]string)
	tz_win := make(map[string]string)

	// # UTC is a common but non-standard alias for Etc/UTC:
	tz_win["Etc/UTC"] = "UTC"

	for _, element := range data.WindowsZones.MapTimezones {
		// if element.Type == "windows" {
		// 	break
		// }

		// Making windows mapping
		for _, m := range element.MapZone {
			t := strings.Split(m.Type, " ")
			if m.Territory == "001" {
				win_tz[m.Other] = t[0]
			}
			for _, tz_name := range t {
				tz_win[tz_name] = m.Other
			}
		}

		// Map in the backwards compatible zone names
		for backward_compat_name, standard_name := range backward {
			if win_zone, ok := tz_win[standard_name]; ok {
				tz_win[backward_compat_name] = win_zone
			}
		}
	}

	// sort the keys
	win_tz_keys := make([]string, 0, len(win_tz))
	for k := range win_tz {
		win_tz_keys = append(win_tz_keys, k)
	}
	sort.Strings(win_tz_keys)

	tz_win_keys := make([]string, 0, len(tz_win))
	for k := range tz_win {
		tz_win_keys = append(tz_win_keys, k)
	}
	sort.Strings(tz_win_keys)

	// Generate the code
	out := bytes.Buffer{}
	out.WriteString("// A lookup table, mapping Windows time zone names to IANA time zone names and vice versa.\n\n")
	out.WriteString(fmt.Sprintf("// Last created %v\n\n", time.Now().UTC().Format(time.RFC3339)))
	out.WriteString("// WinTZtoIANA maps time zone names used by Windows to those used by IANA\n")
	out.WriteString("var WinTZtoIANA = map[string]string{\n")

	for _, k := range win_tz_keys {
		out.WriteString(fmt.Sprintf("\t\"%v\": \"%v\",\n", k, win_tz[k]))
	}
	out.WriteString("}\n\n")
	out.WriteString("// IANAtoWinTZ maps time zone names used by IANA to those used by Windows\n")
	out.WriteString("var IANAtoWinTZ = map[string]string{\n")

	for _, k := range tz_win_keys {
		out.WriteString(fmt.Sprintf("    \"%v\":\"%v\",\n", k, tz_win[k]))
	}
	out.WriteString("}\n")

	// Write buffered code to target writer
	target.Write(out.Bytes())

	return nil
}
