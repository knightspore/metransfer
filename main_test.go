package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"testing"
)

func setup() {
	fmt.Println("Setup Tests")
	var buf bytes.Buffer
	log.SetOutput(&buf)
}

func shutdown() {
	fmt.Println("Test Shutdown")
	err := os.RemoveAll("./upload")
	if err != nil {
		log.Fatalf("Could not remove test directory")
	}
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	shutdown()
	os.Exit(code)
}
