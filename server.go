package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

func uploadFile(w http.ResponseWriter, r *http.Request) {

	log.Println("::> New Upload Started")

	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("fileUpload")
	if err != nil {
		log.Println("!!! Error Retrieving File")
		log.Println(err)
		return
	}
	defer file.Close()

	h := makeHash(handler.Filename, handler.Size)

	if len(db[h]) != 0 {

		log.Println("!!! Error: File Already Exists")

		writeJson(w, map[string]string{
			"status": "301",
			"url":    "http://127.0.0.1:1337/api/download/" + h},
			http.StatusMovedPermanently)

	} else {

		f, err := os.OpenFile("./upload/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			log.Println("!!! Error Saving File")
			log.Println(err)
			return
		}
		defer f.Close()
		io.Copy(f, file)

		updateDb(h, handler.Filename)

		log.Printf("Uploaded File: %v", handler.Filename)
		log.Printf("File Size: %v", handler.Size)
		log.Printf("File Hash: %s", h)
		log.Printf("URL: http://127.0.0.1:1337/api/download/%s", h)

		writeJson(w, map[string]string{
			"status": "200",
			"url":    "http://127.0.0.1:1337/api/download/" + h,
		}, http.StatusOK)

	}
}

func downloadFile(w http.ResponseWriter, r *http.Request) {

	log.Println("::> New Download Request")
	h := strings.Replace(r.URL.String(), "/api/download/", "", -1)
	file, exists := db[h]

	if !exists {

		log.Printf("!!! Error: File Not Found")
		writeJson(w, map[string]string{
			"status": "404",
		}, http.StatusNotFound)

	} else {

		log.Printf("~~~ Serving File: %s", file)
		w.Header().Set("Content-Disposition", "attachment; filename="+file)
		http.ServeFile(w, r, "./upload/"+file)

	}

}

// Main Server

func setupRoutes() {
	log.Println("::> Awaiting Server Connections")

	http.HandleFunc("/api/upload", uploadFile)
	http.HandleFunc("/api/download/", downloadFile)

	s := &http.Server{
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 10 * time.Second,
		Addr:         ":1337",
	}

	if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		fmt.Println(err)
	}
}
