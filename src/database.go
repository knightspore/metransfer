package main

import (
	"database/sql"
	"io/fs"
	"io/ioutil"
	"os"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)

// Connect returns a pointer to a sql.DB object
func (d *Database) Connect() *sql.DB {
	db, _ := sql.Open(d.Type, d.Path)
	return db
}

// Setup creates the AppDatabase file and upload folder if it does not exist
func (d *Database) Setup() {

	err := os.Remove(d.Path)
	if err != nil {
		Logger.Warn("No AppDatabase file found on Server start")
	}

	_, err = os.Create(d.Path)
	if err != nil {
		Logger.Warn("Error creating AppDatabase file: (d.Path)" + d.Path)
	}

	db := d.Connect()
	defer db.Close()

	stmt, err := db.Prepare(`CREATE TABLE upload (
		hash TEXT PRIMARY KEY,
		name TEXT
	);`)

	if err != nil {
		Logger.Error("Error creating table", err.Error())
	}

	stmt.Exec()
	if err != nil {
		Logger.Error("Error executing table creation", err.Error())
	}

	err = os.MkdirAll(d.UploadDir, os.ModePerm)
	if err != nil {
		Logger.Error("Error creating upload directory", err.Error())
	}

	// Populate Upload Folder
	files := d.Files()
	if len(files) > 0 {
		Logger.Info("Adding " + strconv.Itoa(len(files)) + " Existing Files to Database")
		for _, file := range files {
			h := makeHash(file.Name(), file.Size())
			d.InsertRecord(h, file.Name())
		}
		Logger.Info("Database populated with existing files")
	} else {
		Logger.Warn("No existing files found in upload directory")
	}

}

// Files returns a slice of files in the upload directory
func (d *Database) Files() []fs.FileInfo {

	files, err := ioutil.ReadDir(d.UploadDir)
	if err != nil {
		Logger.Error("Error reading upload directory", err.Error())
	}

	return files
}

// InsertRecord inserts a record into the AppDatabase (hash, name)
func (d *Database) InsertRecord(h string, n string) {
	db := d.Connect()
	defer db.Close()

	stmt, err := db.Prepare(`INSERT INTO upload(hash, name) VALUES (?, ?)`)
	if err != nil {
		Logger.Error("Error preparing insert statement", err.Error())
	}

	_, err = stmt.Exec(h, n)
	if err != nil {
		Logger.Error("Error executing insert statement", err.Error())
	}

	Logger.Info("Inserted record into AppDatabase")
}

// GetRecord returns a record from the AppDatabase for a file hash
func (d *Database) GetRecord(h string) (bool, Upload) {
	db := d.Connect()
	defer db.Close()

	var file Upload

	row := db.QueryRow(`SELECT hash, name FROM upload WHERE hash=$1`, h)
	switch err := row.Scan(&file.hash, &file.name); err {
	case sql.ErrNoRows:
		Logger.Info("No record found for hash: " + h)
		return false, file
	case nil:
		return true, file
	default:
		return false, file
	}

}

// LogRecords logs all records in the AppDatabase
func (d *Database) LogRecords() {
	db := d.Connect()
	defer db.Close()

	row, err := db.Query("SELECT * FROM upload")
	if err != nil {
		Logger.Error("Error querying AppDatabase", err.Error())
	}

	defer row.Close()

	for row.Next() {
		var hash, name string
		err = row.Scan(&hash, &name)
		Logger.Info("{ hash: \"" + hash + "\", name: \"" + name + "\" }")
	}

}
