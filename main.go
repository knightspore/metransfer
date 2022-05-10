package main

import (
	"context"
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
}

func setupApplication(w io.Writer, isTest bool) {

	log.Println("")
	log.Println("::> Start Server")

	setupLog(isTest, w)
	setupDb()
	populateDb()

	s := setupRoutes()

	if isTest == false {
		// Start Server with Graceful Shutdown
		go start(s)

		stopCh, closeCh := createChannel()
		defer closeCh()
		log.Println("::> Notified:", <-stopCh)

		shutdown(context.Background(), s)
	} else {
		// Test Start Server and Graceful Shutdown
		go start(s)

		_, closeCh := createChannel()

		closeCh()
		shutdown(context.Background(), s)

	}

}
