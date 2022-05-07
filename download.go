package main

import (
	"log"
	"net/http"
	"strings"
)

func downloadFile(w http.ResponseWriter, r *http.Request) {

	log.Println("::> New Download Request")
	h := strings.Replace(r.URL.String(), "/api/download/", "", -1)
	file, exists := db[h]

	if !exists {

		log.Printf("!!! Error: File Not Found")
		w.WriteHeader(http.StatusNotFound)
		writeJson(w, map[string]string{
			"status": "404",
		})

	} else {

		log.Printf("~~~ Serving File: %s", file)
		w.Header().Set("Content-Disposition", "attachment; filename="+file)
		http.ServeFile(w, r, "./upload/"+file)

	}

}
