package main

import (
	"fmt"
	"os"
	"testing"
)

const (
	testdir        = "/tmp/metransfer"
	testdb         = "metransfer.db"
	file1          = "foo.txt"
	file2          = "bar.txt"
	helloworldhash = "41289770591ec7768a9bddc3a282f04fb36db5c6"
	fakefile       = "fake.txt"
	fakehash       = "1234abcd"
)

func TestCreateHash(t *testing.T) {

	m := MeTransfer{}

	t.Run("creates basic hash", func(t *testing.T) {
		got := m.CreateHash("Hello World", 24)
		assertEqual(t, got, helloworldhash)
	})

	t.Run("changes based on file name and size", func(t *testing.T) {
		got := m.CreateHash("Hello Other World", 24)
		assertUnequal(t, got, helloworldhash)
		got2 := m.CreateHash("Hello World", 48)
		assertUnequal(t, got2, helloworldhash)
	})

}

func TestInit(t *testing.T) {

	m, teardown := setupTest(t)
	defer teardown(t)

	err := m.Init()
	if err != nil {
		t.Errorf("initialization failed: %+v", err)
	}

}

func TestDatabase(t *testing.T) {

	m, teardown := setupTest(t)
	defer teardown(t)

	t.Run("populated with existing files", func(t *testing.T) {
		got, err := m.ListFiles()
		if err != nil {
			t.Errorf("could not list existing files: %+v", err)
		}

		want := 2

		if len(got) != want {
			t.Errorf("got %d, want %d", len(got), want)
		}

	})

	t.Run("Insert", func(t *testing.T) {
		err := m.Insert(fakehash, fakefile)
		if err != nil {
			t.Error("could not write to database")
		}
	})

	t.Run("GetByHash", func(*testing.T) {
		got, err := m.GetByHash(fakehash)
		if err != nil {
			t.Errorf("could not read from database: %+v", err)
		}
		want := fakefile
		assertEqual(t, got, want)
	})

}

func setupTest(tb testing.TB) (MeTransfer, func(tb testing.TB)) {

	err := os.Mkdir(testdir, os.ModePerm)
	f1, err := os.Create(fmt.Sprintf("%s/%s", testdir, file1))
	f2, err := os.Create(fmt.Sprintf("%s/%s", testdir, file2))
	_, err = f1.WriteString("foo\n")
	_, err = f2.WriteString("foo bar\n")

	if err != nil {
		tb.Fatal(err)
	}

	return MeTransfer{
			testdir,
			testdb,
		}, func(tb testing.TB) {
			err := os.RemoveAll(testdir)
			if err != nil {
				tb.Fatalf("could not remove test directory %q", testdir)
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
