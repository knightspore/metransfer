package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
)

const hash string = "8b4a851220490c485b3633bc091347944d86e19c"

func compareTest(g string, w string) {
	fmt.Printf("\nGot: %s\nWant: %s\n", g, w)
}

func setup() {

	fmt.Println("Setup Tests")

	Logger.Setup()
	AppDatabase.Setup()
	Server.Setup()

	go Server.Start()

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

}

func teardown() {
	fmt.Println("Test Shutdown")
	err := os.RemoveAll("./upload")
	if err != nil {
		log.Fatalf("Could not remove test upload file")
	}
}

func TestMain(m *testing.M) {
	_, closeCh := Server.CreateChannel()
	defer closeCh()
	setup()
	code := m.Run()
	teardown()
	closeCh()
	Server.Stop(context.Background())
	Logger.Warn(code)
	os.Exit(0)

}
