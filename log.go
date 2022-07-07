package main

import (
	"io"
	"log"
	"os"
)

var Logger = Log{LogPath: "/tmp/metransfer.log", TestLogPath: "/tmp/metransfer_0.log", multi: nil}

func setupLog(t bool, w io.Writer) {

	var path string

	if t == true {
		path = Logger.TestLogPath
	} else {
		path = Logger.LogPath
	}

	logFile, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	Logger.multi = io.MultiWriter(logFile, w)
	log.SetOutput(Logger.multi)

}
