package main

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
