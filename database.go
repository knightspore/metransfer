package main

import (
	"database/sql"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"

	_ "github.com/mattn/go-sqlite3"
)

const dbFile string = "./metransfer.db"

func setupDb() {

	log.Printf("::> Setting up Database...")

	os.Remove(dbFile)

	_, err := os.Create(dbFile)
	if err != nil {
		log.Fatalf("!!! Error creating database file: %s", dbFile)
	}

	db, _ := sql.Open("sqlite3", dbFile)
	defer db.Close()

	log.Printf("~~~ Sqlite3 Database Instantiated")

	tableSql := `CREATE TABLE upload (
		hash TEXT PRIMARY KEY,
		name TEXT
	);`

	stmt, err := db.Prepare(tableSql)
	if err != nil {
		log.Println("!!! Error preparing database table")
		log.Fatalf(err.Error())
	}
	stmt.Exec()

	log.Printf("~~~ Database Table Created")

}

func populateDb() {

	db, _ := sql.Open("sqlite3", dbFile)
	defer db.Close()

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
		log.Printf("~~~ Adding %v Existing Files to Database", len(files))
		for _, file := range files {
			h := makeHash(file.Name(), file.Size())
			insertRecord(db, h, file.Name())
		}

		log.Printf("+++ Database Populated")
		printDb()
	}

}

func insertRecord(db *sql.DB, h string, n string) {

	uploadSql := `INSERT INTO upload(hash, name) VALUES (?, ?)`
	stmt, err := db.Prepare(uploadSql)
	if err != nil {
		log.Fatalf("Error creating insert SQL string")
		log.Fatalln(err.Error())
	}

	_, err = stmt.Exec(h, n)
	if err != nil {
		log.Fatalf("Error inserting upload into database")
		log.Fatalln(err.Error())
	}

	log.Printf("+++ New Database Entry: { hash: \"%s\", name: \"%s\" }", h, n)

}

func getRecord(db *sql.DB, h string) (bool, string, string) {

	getSql := `SELECT hash, name FROM upload WHERE hash=$1`
	var hash, name string

	row := db.QueryRow(getSql, h)
	switch err := row.Scan(&hash, &name); err {
	case sql.ErrNoRows:
		log.Printf("No rows were found for %s", h)
		return false, "", ""
	case nil:
		return true, hash, name
	default:
		return false, "", ""
	}

}

func printDb() {
	log.Printf("::> Database Printout")

	db, _ := sql.Open("sqlite3", dbFile)
	defer db.Close()

	row, err := db.Query("SELECT * FROM upload")
	if err != nil {
		log.Fatalf("Error querying uploads in Database")
	}

	defer row.Close()

	for row.Next() {
		var hash, name string
		err = row.Scan(&hash, &name)
		log.Printf("{ hash: \"%s\", name: \"%s\" }\n", hash, name)
	}
}
