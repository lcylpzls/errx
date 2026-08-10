package errx

import "testing"

func TestVersion(t *testing.T) {
	if Version != "v1.5.4" {
		t.Errorf("Version = %s,want v1.5.4", Version)
	}
}
