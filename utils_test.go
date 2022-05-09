package main

import (
	"testing"
)

func TestMakeHash(t *testing.T) {
	if got, want := makeHash("testUpload", 2500), "4401a81d8f8905d0f82461fd46637e2f984e10b4"; got != want {
		t.Errorf("makeHash() not functioning as expected.")
	}
}
