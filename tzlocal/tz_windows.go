package tzlocal

//go:generate go run ./../tzlocal/cmd/update_tzmapping.go

import (
	"fmt"
	"strings"

	"golang.org/x/sys/windows/registry"
)

const tzKey = `SYSTEM\CurrentControlSet\Control\TimeZoneInformation`
const tzKeyVal = "TimeZoneKeyName"

// LocalTZ obtains the name of the time zone Windows is configured to use. Returns the corresponding IANA standard name
func LocalTZ() (string, error) {
	var winTZname string
	var errTzutil, errReg error

	// try tzutil command first - if that is not available, try to read from registry
	winTZname, errTzutil = localTZfromTzutil()
	if errTzutil != nil {
		winTZname, errReg = localTZfromReg()
		if errReg != nil { // both methods failed, return both errors
			return "", fmt.Errorf("failed to read time zone name with errors\n(1) %s\n(2) %s", errTzutil, errReg)
		}
	}

	if name, ok := WinTZtoIANA[winTZname]; ok {
		return name, nil
	}
	
	// If the timezone is not found in the WinTZtoIANA map, check if the timezone string contains "_dstoff",
	// which indicates that "Daylight Saving Time" adjustments for the timezone are disabled.
	// Remove the "_dstoff" part from the timezone string and check in the map again.
	// Refer: https://learn.microsoft.com/en-us/windows-server/administration/windows-commands/tzutil#syntax
	index := strings.Index(winTZname, "_dstoff")
	if index != -1 {
		// Remove "_dstoff" from the time zone string
		winTZname = winTZname[:index]

		// Get the IANA time zone for the updated time zone string.
		if name, ok := tzlocal.WinTZtoIANA[winTZname]; ok {
			return name, nil
		}
	}

	
	return "", fmt.Errorf("could not find IANA tz name for set time zone \"%s\"", winTZname)
}

// localTZfromReg obtains the time zone Windows is configured to use from registry.
func localTZfromReg() (string, error) {
	k, err := registry.OpenKey(registry.LOCAL_MACHINE, tzKey, registry.QUERY_VALUE)
	if err != nil {
		return "", err
	}
	defer k.Close()

	winTZname, _, err := k.GetStringValue(tzKeyVal)
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(winTZname), nil
}
