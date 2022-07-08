package main

import (
	"io"
	"net/http"
)

// Log is a struct for the log
type Log struct {
	LogPath     string
	TestLogPath string
	Multi       io.Writer
}

// Event is a struct for the log items
type Event struct {
	Level   int
	Message []interface{}
}

// Upload is a struct for the uploads table
type Upload struct {
	hash string
	name string
}

// Database is a struct for the AppDatabase
type Database struct {
	Path      string
	Type      string
	UploadDir string
}

// FileServer is a struct for the fileserver
type FileServer struct {
	Port   string
	Server *http.Server
}
