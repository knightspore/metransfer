package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"
)

const hash string = "8b4a851220490c485b3633bc091347944d86e19c"

func compareTest(g string, w string) {
	fmt.Printf("\nGot: %s\nWant: %s\n", g, w)
}

func setup() {

	fmt.Println("Setup Tests")

	newpath := filepath.Join(".", "upload")
	err := os.MkdirAll(newpath, os.ModePerm)
	if err != nil {
		log.Println("Error Creating Upload Folder")
		log.Fatal(err)
	}

	f, err := os.Create("./upload/testUpload")
	if err != nil {
		log.Println("Error creating 'testUpload' File")
		log.Fatal(err)
	}

	defer f.Close()

	_, err2 := f.WriteString("testing")
	if err2 != nil {
		log.Fatal(err2)
	}

	var buf bytes.Buffer
	setupApplication(&buf, true)

}

func teardown() {
	fmt.Println("Test Shutdown")
	err := os.RemoveAll("./upload")
	if err != nil {
		log.Fatalf("Could not remove test directory")
	}
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	teardown()
	os.Exit(code)
}
