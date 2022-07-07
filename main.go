package main

import (
	"context"
	"path/filepath"
)

var database = Database{
	Path:      "./metransfer.db",
	Type:      "sqlite3",
	UploadDir: filepath.Join(".", "upload"),
}

var server = FileServer{
	Port: ":2080",
}

var logger = Log{
	LogPath: "/tmp/metransfer.log",
	Multi:   nil,
}

func main() {

	logger.Setup()
	database.Setup()
	server.Setup()

	go server.Start()

	stopCh, closeCh := server.CreateChannel()
	defer closeCh()
	<-stopCh

	server.Stop(context.Background())

	// Test Start Server and Graceful Shutdown
	//go server.Start()
	//_, closeCh := server.CreateChannel()
	//closeCh()
	//server.Stop(context.Background())

}
