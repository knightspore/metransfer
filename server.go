package main

import (
	"context"
	"errors"
	"io"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
	"time"
)

func (s *FileServer) Setup() {

	fs := http.FileServer(http.Dir("./"))
	http.Handle("/", fs)
	http.HandleFunc("/api/upload", s.UploadFile)
	http.HandleFunc("/api/download/", s.DownloadFile)

	s.Server = &http.Server{
		Addr: s.Port,
	}

}

func (s *FileServer) Start() {

	logger.Info("Starting Application")
	if err := s.Server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		panic(err)
	} else {
		logger.Info("Server Stopped Gracefully")
	}

}

func (s *FileServer) Stop(ctx context.Context) {

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := s.Server.Shutdown(ctx); err != nil {
		panic(err)
	} else {
		logger.Warn("File Server Stopped Gracefully")
	}

}

func (s *FileServer) CreateChannel() (chan os.Signal, func()) {
	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	return stopCh, func() {
		close(stopCh)
	}
}

func (s *FileServer) UploadFile(w http.ResponseWriter, r *http.Request) {

	logger.Info("New Upload Request")

	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("fileUpload")
	if err != nil {
		logger.Warn("Error Retrieving File", err)
		return
	}
	defer file.Close()

	h := makeHash(handler.Filename, handler.Size)

	exists, _ := database.GetRecord(h)

	if exists {

		logger.Warn("File Already Exists", h)

		writeJson(w, map[string]string{
			"status": "301",
			"url":    "http://127.0.0.1:" + server.Port + "/api/download/" + h},
			http.StatusMovedPermanently)

	} else {

		f, err := os.OpenFile("./upload/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			logger.Warn("Error Opening File to Save Upload", err)
			return
		}

		defer f.Close()
		io.Copy(f, file)

		database.InsertRecord(h, handler.Filename)

		logger.Info(
			"File Uploaded: "+handler.Filename,
			"File Size: "+strconv.FormatInt(handler.Size, 10),
			"File Hash: "+h,
			"File URL: http://127.0.0.1:1337/api/download/"+h,
		)

		writeJson(w, map[string]string{
			"status": "200",
			"url":    "http://127.0.0.1:2080/api/download/" + h,
		}, http.StatusOK)

	}
}

func (s *FileServer) DownloadFile(w http.ResponseWriter, r *http.Request) {

	logger.Info("New Download Request")
	h := strings.Replace(r.URL.String(), "/api/download/", "", -1)

	exists, upload := database.GetRecord(h)
	if !exists {

		logger.Warn("File Does Not Exist", h)
		writeJson(w, map[string]string{
			"status": "404",
		}, http.StatusNotFound)

	} else {

		logger.Info("Serving File: " + upload.name)
		w.Header().Set("Content-Disposition", "attachment; filename="+upload.name)
		http.ServeFile(w, r, "./upload/"+upload.name)

	}

}
