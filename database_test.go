package main

import (
	"bytes"
	"log"
	"strings"
	"testing"
)

func TestInitDb(t *testing.T) {
	if got, want := db[hash], "testUpload"; got != want {
		t.Errorf("initDb() not working as expected")
		compareTest(got, want)
	}
}

func TestPrintDb(t *testing.T) {

	var buf bytes.Buffer
	log.SetOutput(&buf)
	printDb()
	output := buf.String()

	split := strings.Split(output, " ")

	hashBufString := "\"38623461383531323230343930633438356233363333626330393133343739343464383665313963\","
	testFileName := "\"testUpload\""

	if got, want := split[8], hashBufString; got != want {
		t.Errorf("printDb() Not working as expected")
		compareTest(got, want)
	}

	if got, want := split[10], testFileName; got != want {
		t.Errorf("printDb() Not working as expected")
		compareTest(got, want)
	}

}
