package main

import (
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
)

/*
	Database
*/

var db = map[string]string{}

func initDb() {

	newpath := filepath.Join(".", "upload")
	err := os.MkdirAll(newpath, os.ModePerm)
	if err != nil {
		log.Println("!!! Error Creating Upload Folder")
		log.Fatal(err)
	}

	files, err := ioutil.ReadDir("./upload/")
	if err != nil {
		log.Println("!!! No Upload Folder Found")
		log.Fatal(err)
	}

	if len(files) > 0 {
		log.Println("~~~ Adding Existing Files to Database")
		for _, file := range files {
			h := makeHash(file.Name(), file.Size())
			updateDb(h, file.Name())
		}
	}

	log.Println("::> DB Initialized")
}

func updateDb(h string, n string) {
	db[h] = n // Save to DB
	log.Printf("+++ New Database Entry: { hash: \"%s\", name: \"%s\" }", h, n)
}

func printDb() {
	log.Printf("::> Database Printout")
	for key, el := range db {
		log.Printf("{ hash: \"%x\", name: \"%s\" }\n", key, el)
	}
}
