package main

import (
	"io"
	"log"
	"os"
)

func setupLog(w io.Writer) {
	logFile, err := os.OpenFile("/tmp/metransfer.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	multi := io.MultiWriter(logFile, os.Stdout)
	log.SetOutput(multi)
}

func main() {
	setupLog(os.Stdout)
	log.Println("::> Start Server")
	initDb()
	setupRoutes()
}
