package main

import (
	"io"
	"log"
	"net/http"
	"os"
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
