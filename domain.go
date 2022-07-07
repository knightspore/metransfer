package main

import "io"

// Log is a struct for the log
type Log struct {
	LogPath     string
	TestLogPath string
	multi       io.Writer
}

// LogItem is a struct for the log items
type LogItem struct {
	Message string
	Level   string
}

// Upload is a struct for the uploads table
type Upload struct {
	hash string
	name string
}

// Database is a struct for the database
type Database struct {
	Path      string
	Type      string
	UploadDir string
}
