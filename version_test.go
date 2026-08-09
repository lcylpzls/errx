package errx

import "testing"

func TestVersion(t *testing.T) {
	if Version != "v1.3.1" {
		t.Errorf("Version = %s,want v1.3.1", Version)
	}
}
