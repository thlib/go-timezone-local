package tzlocal

import (
	"testing"
	"time"
)

func TestMappings(t *testing.T) {
	// sanity check: can we find all IANA names from WinTZtoIANA in IANAtoWinTZ?
	for _, v := range WinTZtoIANA {
		if _, ok := IANAtoWinTZ[v]; !ok {
			t.Errorf("could not find IANA name %v in IANAtoWinTZ", v)
		}
	}
	// sanity checks: can we find all IANA names from IANAtoWinTZ in WinTZtoIANA?
	// can we successfully call time.LoadLocation(tzname) for all given IANA names?
	for k, v := range IANAtoWinTZ {
		if _, ok := WinTZtoIANA[v]; !ok {
			t.Errorf("could not find Win tz name %v in WinTZtoIANA", v)
		}
		if _, err := time.LoadLocation(k); err != nil {
			t.Errorf("time.LoadLocation failed for IANA tz name %v", k)
		}
	}
}

func TestLocalTZ(t *testing.T) {

	s, _, err := localTZfromReg()
	if err != nil {
		t.Errorf("got unexpected error %v", err)
	}
	if s == "" {
		t.Error("got unexpected empty result with no error")
	}

	s, err = LocalTZ()
	if err != nil {
		t.Errorf("got unexpected error %v", err)
	}
	if s == "" {
		t.Error("got unexpected empty result with no error")
	}

	tmp := WinTZtoIANA
	WinTZtoIANA = map[string]string{}
	s, err = LocalTZ()
	if err == nil {
		t.Error("expected error but got nil")
	}
	if s != "" {
		t.Errorf("expected empty result but got %v", s)
	}
	WinTZtoIANA = tmp
}
