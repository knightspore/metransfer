package main

import (
	"testing"
)

func TestMakeHash(t *testing.T) {
	if got, want := makeHash("testUpload", 7), hash; got != want {
		t.Errorf("makeHash() not functioning as expected.")
		compareTest(got, want)
	}
}
