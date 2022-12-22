package main

import "testing"

func assertEqual(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("got %q want %q", got, want)
	}
}

func assertUnequal(t testing.TB, got, want string) {
	t.Helper()
	if got == want {
		t.Errorf("got %q which should not be equal to want %q", got, want)
	}
}

func TestCreateHash(t *testing.T) {

	want := "41289770591ec7768a9bddc3a282f04fb36db5c6"

	t.Run("basic hash", func(t *testing.T) {
		got := CreateHash("Hello World", 24)
		assertEqual(t, got, want)
	})

	t.Run("changes based on filename", func(t *testing.T) {
		got := CreateHash("Hello Other World", 24)
		assertUnequal(t, got, want)
	})

	t.Run("changes based on filesize", func(t *testing.T) {
		got := CreateHash("Hello World", 48)
		assertUnequal(t, got, want)
	})

}
