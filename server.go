package main

import (
	"log"
	"net/http"
)

/*
	Server Routes
*/

func setupRoutes() {
	log.Println("::> Awaiting Server Connections")

	// File Upload
	http.HandleFunc("/api/upload", uploadFile)

	// File Download
	http.HandleFunc("/api/download/", downloadFile)

	http.ListenAndServe(":1337", nil)
}
