package main

import (
	"io"
	"log"
	"os"
)

type Upload struct {
	hash string
	name string
}

func main() {
	setupApplication(os.Stdout, false)
	setupRoutes()
}

func setupApplication(w io.Writer, t bool) {

	log.Println("")
	log.Println("")
	log.Println("::> Start Server")

	logPath, testLogPath := "/tmp/metransfer.log", "/tmp/metransfer_0.log"
	var path string

	if t == true {
		path = testLogPath
	} else {
		path = logPath
	}

	logFile, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	multi := io.MultiWriter(logFile, w)
	log.SetOutput(multi)

	setupDb()
	populateDb()

}
