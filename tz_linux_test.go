package timezone

import "testing"

func TestInferFromPathSuccess(t *testing.T) {
	tz, err := inferFromPath("/usr/share/zoneinfo/Asia/Tokyo")
	if err != nil {
		t.Errorf("got err=%d; want: nil", err)
	}
	want := "Asia/Tokyo"
	if tz != want {
		t.Errorf("got tz=%s; want: %s", tz, want)
	}
}
