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

	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/api/upload", s.UploadFile)
	http.HandleFunc("/api/download/", s.DownloadFile)

	s.Server = &http.Server{
		Addr: s.Port,
	}

}

func (s *FileServer) Start() {

	Logger.Info("Starting Application")
	if err := s.Server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		panic(err)
	} else {
		Logger.Info("Server Stopped Gracefully")
	}

}

func (s *FileServer) Stop(ctx context.Context) {

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := s.Server.Shutdown(ctx); err != nil {
		panic(err)
	} else {
		Logger.Warn("File Server Stopped Gracefully")
	}

}

func (s *FileServer) CreateChannel() (chan os.Signal, func()) {
	stopCh := make(chan os.Signal, 1)
	signal.Notify(stopCh, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	return stopCh, func() {
		close(stopCh)
	}
}

func (s *FileServer) Index(w http.ResponseWriter, r *http.Request) {
	http.FileServer(http.Dir("./static"))
}

func (s *FileServer) UploadFile(w http.ResponseWriter, r *http.Request) {

	Logger.Info("New Upload Request")

	r.ParseMultipartForm(32 << 20)
	file, handler, err := r.FormFile("fileUpload")
	if err != nil {
		Logger.Warn("Error Retrieving File", err)
		return
	}
	defer file.Close()

	h := makeHash(handler.Filename, handler.Size)

	exists, _ := AppDatabase.GetRecord(h)

	if exists {

		Logger.Warn("File Already Exists", h)

		writeJson(w, map[string]string{
			"status":   "301",
			"url":      "/api/download/" + h,
			"filename": handler.Filename,
		},
			http.StatusMovedPermanently)

	} else {

		f, err := os.OpenFile("./upload/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			Logger.Warn("Error Opening File to Save Upload", err)
			return
		}

		defer f.Close()
		io.Copy(f, file)

		AppDatabase.InsertRecord(h, handler.Filename)

		Logger.Info(
			"File Uploaded: "+handler.Filename,
			"File Size: "+strconv.FormatInt(handler.Size, 10),
			"File Hash: "+h,
			"File URL: /api/download/"+h,
		)

		filename := handler.Filename

		writeJson(w, map[string]string{
			"status":   "200",
			"url":      "/api/download/" + h,
			"filename": string(filename),
		}, http.StatusOK)

	}
}

func (s *FileServer) DownloadFile(w http.ResponseWriter, r *http.Request) {

	Logger.Info("New Download Request")
	h := strings.Split(r.URL.String(), "/")[2]

	exists, upload := AppDatabase.GetRecord(h)
	if !exists {

		Logger.Warn("File Does Not Exist", h)
		writeJson(w, map[string]string{
			"status": "404",
		}, http.StatusNotFound)

	} else {

		Logger.Info("Serving File: " + upload.name)
		w.Header().Set("Content-Disposition", "attachment; filename=\""+upload.name+"\"")
		http.ServeFile(w, r, "./upload/"+upload.name)

	}

}
