package errx

import "testing"

func TestVersion(t *testing.T) {
	if Version != "v1.5.0" {
		t.Errorf("Version = %s,want v1.5.0", Version)
	}
}
