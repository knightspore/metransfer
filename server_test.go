package main

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestDownloadFile404(t *testing.T) {

	req := httptest.NewRequest(http.MethodGet, "/api/download/badhash", nil)
	w := httptest.NewRecorder()
	downloadFile(w, req)
	res := w.Result()
	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("Error hitting download endpoint")
	}
	if string(data) != "{\"status\":\"404\"}" {
		t.Errorf("Bad download isn't returning 404")
	}

}

func TestDownloadFile(t *testing.T) {

	req := httptest.NewRequest(http.MethodGet, "/api/download/"+hash, nil)
	w := httptest.NewRecorder()
	downloadFile(w, req)
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
