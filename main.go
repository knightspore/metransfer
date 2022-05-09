package main

import (
	"io"
	"log"
	"os"
)

func main() {
	os.Exit(setupApplication(os.Stdout))
	setupRoutes()
}

func setupApplication(w io.Writer) int {
	logFile, err := os.OpenFile("/tmp/metransfer.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	multi := io.MultiWriter(logFile, os.Stdout)
	log.SetOutput(multi)

	log.Println("::> Start Server")

	initDb()
	return 0
}
