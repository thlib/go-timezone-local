package tzlocal

import (
	"errors"
	"testing"
	"time"
)

func TestLocalTZInfoFromSources(t *testing.T) {
	apiErr := errors.New("API unavailable")
	regErr := errors.New("registry unavailable")
	for _, tc := range []struct {
		name        string
		apiErr      error
		regErr      error
		dstOff      bool
		wantRegRead bool
		wantErr     string
	}{
		{name: "API with DST enabled"},
		{name: "API with DST disabled", dstOff: true},
		{name: "API succeeds despite registry failure", regErr: regErr},
		{name: "registry with DST enabled", apiErr: apiErr, wantRegRead: true},
		{name: "registry with DST disabled", apiErr: apiErr, dstOff: true, wantRegRead: true},
		{
			name: "both sources fail", apiErr: apiErr, regErr: regErr, wantRegRead: true,
			wantErr: "failed to read timezone from Windows API (API unavailable) and registry (registry unavailable)",
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var calls []string
			name, dstOff, err := localTZInfoFromSources(
				func() (string, bool, error) {
					calls = append(calls, "API")
					return "FLE Standard Time", tc.dstOff, tc.apiErr
				},
				func() (string, bool, error) {
					calls = append(calls, "registry")
					return "India Standard Time", tc.dstOff, tc.regErr
				},
			)
			if len(calls) == 0 || calls[0] != "API" {
				t.Fatalf("source calls = %v, want API first", calls)
			}
			wantCalls := 1
			if tc.wantRegRead {
				wantCalls = 2
			}
			if len(calls) != wantCalls {
				t.Errorf("source calls = %v, want %d calls", calls, wantCalls)
			}
			if tc.wantErr != "" {
				if err == nil || err.Error() != tc.wantErr {
					t.Errorf("error = %v, want %q", err, tc.wantErr)
				}
				if name != "" || dstOff {
					t.Errorf("failed lookup = (%q, %t), want empty name and false", name, dstOff)
				}
				return
			}
			wantName := "FLE Standard Time"
			if tc.wantRegRead {
				wantName = "India Standard Time"
			}
			if err != nil || name != wantName || dstOff != tc.dstOff {
				t.Errorf("lookup = (%q, %t, %v), want (%q, %t, nil)", name, dstOff, err, wantName, tc.dstOff)
			}
		})
	}
}

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
