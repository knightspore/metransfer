package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"
	"testing"
)

func TestInitDb(t *testing.T) {
	initDb()

	_, err := os.ReadDir("./upload")
	if err != nil {
		t.Errorf("%v", err)
	}
}

func TestUpdateDb(t *testing.T) {
	updateDb("12345", "testUpload")

	if got, want := db["12345"], "testUpload"; got != want {
		t.Errorf("Created record %v is not equal to expected %v", got, want)
	}
}

func TestPrintDb(t *testing.T) {

	var buf bytes.Buffer
	log.SetOutput(&buf)
	printDb()
	output := buf.String()
	log.SetOutput(os.Stderr)

	split := strings.Split(output, " ")

	if got, want := split[8], "\"3132333435\","; got != want {
		fmt.Println("printDb() Not working as expected")
	}

	if got, want := split[10], "\"testUpload\""; got != want {
		fmt.Println("printDb() Not working as expected")
	}

}
