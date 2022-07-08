package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
)

func TestUploadFile(t *testing.T) {

	file, _ := os.Open("./upload/testUpload")
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("fileUpload", filepath.Base(file.Name()))
	io.Copy(part, file)
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/api/upload", body)
	req.Header.Add("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()
	Server.UploadFile(w, req)
	res := w.Result()

	if got, want := res.Status, "200 OK"; got != want {
		t.Error("Uploading a new file does not return a 200")
		compareTest(got, want)
	}

}

func TestUploadExistingFile(t *testing.T) {

	file, _ := os.Open("./upload/testUpload")
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	part, _ := writer.CreateFormFile("fileUpload", filepath.Base(file.Name()))
	io.Copy(part, file)
	writer.Close()

	req := httptest.NewRequest(http.MethodPost, "/api/upload", body)
	req.Header.Add("Content-Type", writer.FormDataContentType())
	w := httptest.NewRecorder()
	Server.UploadFile(w, req)
	res := w.Result()

	if got, want := res.Status, "301 Moved Permanently"; got != want {
		t.Error("Uploading an existing file doesn't return a 301")
		compareTest(got, want)
	}

}

func TestDownloadFile(t *testing.T) {

	req := httptest.NewRequest(http.MethodGet, "/api/download/"+hash, nil)
	w := httptest.NewRecorder()
	Server.DownloadFile(w, req)
	res := w.Result()
	defer res.Body.Close()

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("Error hitting download endpoint")
	}

	if got, want := string(data), "testing"; got != want {
		t.Errorf("Incorrect data in test file download")
		compareTest(got, want)
	}

}

func TestDownloadBadFile(t *testing.T) {

	req := httptest.NewRequest(http.MethodGet, "/api/download/badhash", nil)
	w := httptest.NewRecorder()
	Server.DownloadFile(w, req)
	res := w.Result()

	if got, want := res.Status, "404 Not Found"; got != want {
		t.Errorf("Bad download isn't returning 404")
		compareTest(got, want)
	}

}

func TestLogRecords(t *testing.T) {

	AppDatabase.LogRecords()

}
