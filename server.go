package main

import (
	"context"
	"database/sql"
	"errors"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
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

	db := database.Connect()
	defer db.Close()

	exists, upload := database.GetRecord(h)

	if exists {

		log.Println("!!! Error: File Already Exists")

		writeJson(w, map[string]string{
			"status": "301",
			"url":    "http://127.0.0.1:1337/api/download/" + upload.hash},
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

		database.InsertRecord(h, handler.Filename)

		log.Printf("Uploaded File: %v", handler.Filename)
		log.Printf("File Size:\t%v", handler.Size)
		log.Printf("File Hash:\t%s", h)
		log.Printf("File URL:\thttp://127.0.0.1:1337/api/download/%s", h)

		writeJson(w, map[string]string{
			"status": "200",
			"url":    "http://127.0.0.1:2080/api/download/" + h,
		}, http.StatusOK)

	}
}

func downloadFile(w http.ResponseWriter, r *http.Request) {

	log.Println("::> New Download Request")
	h := strings.Replace(r.URL.String(), "/api/download/", "", -1)

	db, _ := sql.Open("sqlite3", dbFile)
	defer db.Close()

	exists, upload := database.GetRecord(h)
	if !exists {

		log.Printf("!!! Error: File Not Found")
		writeJson(w, map[string]string{
			"status": "404",
		}, http.StatusNotFound)

	} else {

		log.Printf("~~~ Serving File: %s", upload.name)
		w.Header().Set("Content-Disposition", "attachment; filename="+upload.name)
		http.ServeFile(w, r, "./upload/"+upload.name)

	}

}

// Main Server

func createChannel() (chan os.Signal, func()) {
	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	return stopCh, func() {
		close(stopCh)
	}
}

func start(s *http.Server) {
	log.Println("::> Awaiting Server Connections")
	if err := s.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		panic(err)
	} else {
		log.Println("::> Server Stopped Gracefully")
	}
}

func shutdown(ctx context.Context, s *http.Server) {
	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		panic(err)
	} else {
		log.Println("::> Application Shut Down")
	}
}

func setupRoutes() *http.Server {

	http.HandleFunc("/api/upload", uploadFile)
	http.HandleFunc("/api/download/", downloadFile)

	s := &http.Server{
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
		Addr:         ":2080",
	}

	return s

}
