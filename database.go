package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"io/fs"
	"io/ioutil"
	"log"
	"os"
)

const dbFile string = "./metransfer.db"

// Connect returns a pointer to a sql.DB object
func (d *Database) Connect() *sql.DB {
	db, _ := sql.Open(d.Type, d.Path)
	return db
}

// Setup creates the database file and upload folder if it does not exist
func (d *Database) Setup() {

	os.Remove(d.Path)

	_, err := os.Create(d.Path)
	if err != nil {
		log.Fatalf("!!! Error creating database file: %s", d.Path)
	}

	db := d.Connect()
	defer db.Close()

	stmt, err := db.Prepare(`CREATE TABLE upload (
		hash TEXT PRIMARY KEY,
		name TEXT
	);`)

	if err != nil {
		log.Println("!!! Error preparing database table")
		log.Fatalf(err.Error())
	}

	stmt.Exec()

	err = os.MkdirAll(d.UploadDir, os.ModePerm)
	if err != nil {
		log.Println("!!! Error Creating Upload Folder")
		log.Fatal(err)
	}

}

// Files returns a slice of files in the upload directory
func (d *Database) Files() []fs.FileInfo {

	files, err := ioutil.ReadDir(d.UploadDir)
	if err != nil {
		log.Println("!!! No Upload Folder Found")
		log.Fatal(err)
	}

	return files
}

// Populate inserts existing files from the upload folder into the database on restart
func (d *Database) Populate() {

	files := d.Files()
	db := d.Connect()
	defer db.Close()

	if len(files) > 0 {
		log.Printf("~~~ Adding %v Existing Files to Database", len(files))
		for _, file := range files {
			h := makeHash(file.Name(), file.Size())
			d.InsertRecord(h, file.Name())
		}
		log.Printf("+++ Database Populated")
		database.LogRecords()
	} else {
		log.Println("No Files to Populate")
	}

}

// InsertRecord inserts a record into the database (hash, name)
func (d *Database) InsertRecord(h string, n string) {
	db := d.Connect()
	defer db.Close()

	stmt, err := db.Prepare(`INSERT INTO upload(hash, name) VALUES (?, ?)`)
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

// GetRecord returns a record from the database for a file hash
func (d *Database) GetRecord(h string) (bool, Upload) {
	db := d.Connect()
	defer db.Close()

	var file Upload

	row := db.QueryRow(`SELECT hash, name FROM upload WHERE hash=$1`, h)
	switch err := row.Scan(&file.hash, &file.name); err {
	case sql.ErrNoRows:
		log.Printf("No rows were found for %s", h)
		return false, file
	case nil:
		return true, file
	default:
		return false, file
	}

}

// LogRecords logs all records in the database
func (d *Database) LogRecords() {
	db := d.Connect()
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
