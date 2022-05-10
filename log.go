package main

import (
	"io"
	"log"
	"os"
)

func setupLog(t bool, w io.Writer) {

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

}
