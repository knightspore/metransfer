package main

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

func makeHash(n string, s int64) string {
	b := sha1.New()
	b.Write([]byte(n + fmt.Sprint(s)))
	h := hex.EncodeToString(b.Sum(nil))
	return h
}

func writeJson(w http.ResponseWriter, d map[string]string, s int) {
	jd, err := json.Marshal(d)
	if err != nil {
		log.Println(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(s)
	_, err = w.Write(jd)
	if err != nil {
		Logger.Warn("Error Writing Response", err)
	}
}
