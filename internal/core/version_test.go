package core

import "testing"

func TestVersion(t *testing.T) {
	if Version != "v1.6.0" {
		t.Errorf("Version = %s,want v1.6.0", Version)
	}
}
