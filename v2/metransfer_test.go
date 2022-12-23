package main

import (
	"fmt"
	"os"
	"testing"
)

var (
	testdir        = "/tmp/metransfer"
	testdb         = "metransfer.db"
	file1          = "foo.txt"
	file2          = "bar.txt"
	helloworldhash = "41289770591ec7768a9bddc3a282f04fb36db5c6"
)

func TestServer(t *testing.T) {
	teardown := setupTest(t)
	defer teardown(t)

	m := NewMetransfer(testdir, testdb)

	t.Run("initialized correctly", func(t *testing.T) {
		t.Run("upload folder set", func(t *testing.T) {
			path := m.Path()
			assertEqual(t, path, testdir)
		})
		t.Run("db path set", func(t *testing.T) {
			dbPath := m.DBPath()
			assertEqual(t, dbPath, fmt.Sprintf("%s/%s", testdir, testdb))
		})
	})

	t.Run("can read uploads folder", func(t *testing.T) {
		files := m.ReadUploadsFolder()
		got := len(files)
		want := 2
		if got != want {
			t.Errorf("got %d want %d", got, want)
		}
	})

}

func TestCreateHash(t *testing.T) {
	t.Run("basic hash", func(t *testing.T) {
		got := CreateHash("Hello World", 24)
		assertEqual(t, got, helloworldhash)
	})
	t.Run("changes based on filename", func(t *testing.T) {
		got := CreateHash("Hello Other World", 24)
		assertUnequal(t, got, helloworldhash)
	})
	t.Run("changes based on filesize", func(t *testing.T) {
		got := CreateHash("Hello World", 48)
		assertUnequal(t, got, helloworldhash)
	})

}

func setupTest(tb testing.TB) func(tb testing.TB) {
	err := os.Mkdir(testdir, os.ModePerm)
	f1, err := os.Create(fmt.Sprintf("%s/%s", testdir, file1))
	f2, err := os.Create(fmt.Sprintf("%s/%s", testdir, file2))
	_, err = f1.WriteString("foo\n")
	_, err = f2.WriteString("foo bar\n")

	if err != nil {
		tb.Fatal(err)
	}

	return func(tb testing.TB) {
		err := os.RemoveAll(testdir)
		if err != nil {
			tb.Fatalf("Could not remove test directory %q", testdir)
		}
	}
}

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
