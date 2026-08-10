package errx

import "testing"

func TestVersion(t *testing.T) {
	if Version != "v1.4.0" {
		t.Errorf("Version = %s,want v1.4.0", Version)
	}
}
