package main

import (
	"context"
	"io"
	"log"
	"os"
	"path/filepath"
)

var database = Database{Path: "./metransfer.db", Type: "sqlite3", UploadDir: filepath.Join(".", "upload")}

func main() {
	setupApplication(os.Stdout, false)
}

func setupApplication(w io.Writer, isTest bool) {

	log.Println("::> Start Server")

	setupLog(isTest, w)
	database.Setup()
	database.Populate()

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
