package main

import (
	"context"
	"path/filepath"
)

var AppDatabase = Database{
	Path:      "./metransfer.db",
	Type:      "sqlite3",
	UploadDir: filepath.Join(".", "upload"),
}

var Server = FileServer{
	Port: ":2080",
}

var Logger = Log{
	LogPath: "/tmp/metransfer.log",
	Multi:   nil,
}

func main() {

	Logger.Setup()
	AppDatabase.Setup()
	Server.Setup()

	go Server.Start()

	stopCh, closeCh := Server.CreateChannel()
	defer closeCh()
	<-stopCh

	Server.Stop(context.Background())

}
