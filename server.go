package main

import (
	"log"
	"net/http"
)

func setupRoutes() {
	log.Println("::> Awaiting Server Connections")

	http.HandleFunc("/api/upload", uploadFile)
	http.HandleFunc("/api/download/", downloadFile)

	http.ListenAndServe(":1337", nil)
}
