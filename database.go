package main

import (
	"database/sql"
	_ "github.com/mattn/go-sqlite3"
	"io/fs"
	"io/ioutil"
	"os"
	"strconv"
)

const dbFile string = "./metransfer.db"

// Connect returns a pointer to a sql.DB object
func (d *Database) Connect() *sql.DB {
	db, _ := sql.Open(d.Type, d.Path)
	return db
}

// Setup creates the database file and upload folder if it does not exist
func (d *Database) Setup() {

	err := os.Remove(d.Path)
	if err != nil {
		logger.Warn("No database file found on server start")
	}

	_, err = os.Create(d.Path)
	if err != nil {
		logger.Warn("Error creating database file: (d.Path)" + d.Path)
	}

	db := d.Connect()
	defer db.Close()

	stmt, err := db.Prepare(`CREATE TABLE upload (
		hash TEXT PRIMARY KEY,
		name TEXT
	);`)

	if err != nil {
		logger.Error("Error creating table", err.Error())
	}

	stmt.Exec()
	if err != nil {
		logger.Error("Error executing table creation", err.Error())
	}

	err = os.MkdirAll(d.UploadDir, os.ModePerm)
	if err != nil {
		logger.Error("Error creating upload directory", err.Error())
	}

	// Populate Upload Folder
	files := d.Files()
	if len(files) > 0 {
		logger.Info("Adding " + strconv.Itoa(len(files)) + " Existing Files to Database")
		for _, file := range files {
			h := makeHash(file.Name(), file.Size())
			d.InsertRecord(h, file.Name())
		}
		logger.Info("Database populated with existing files")
	} else {
		logger.Warn("No existing files found in upload directory")
	}

}

// Files returns a slice of files in the upload directory
func (d *Database) Files() []fs.FileInfo {

	files, err := ioutil.ReadDir(d.UploadDir)
	if err != nil {
		logger.Error("Error reading upload directory", err.Error())
	}

	return files
}

// InsertRecord inserts a record into the database (hash, name)
func (d *Database) InsertRecord(h string, n string) {
	db := d.Connect()
	defer db.Close()

	stmt, err := db.Prepare(`INSERT INTO upload(hash, name) VALUES (?, ?)`)
	if err != nil {
		logger.Error("Error preparing insert statement", err.Error())
	}

	_, err = stmt.Exec(h, n)
	if err != nil {
		logger.Error("Error executing insert statement", err.Error())
	}

	logger.Info("Inserted record into database")
	logger.Info("{ hash: \"" + h + "\", name: \"" + n + "\" }")
}

// GetRecord returns a record from the database for a file hash
func (d *Database) GetRecord(h string) (bool, Upload) {
	db := d.Connect()
	defer db.Close()

	var file Upload

	row := db.QueryRow(`SELECT hash, name FROM upload WHERE hash=$1`, h)
	switch err := row.Scan(&file.hash, &file.name); err {
	case sql.ErrNoRows:
		logger.Info("No record found for hash: " + h)
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
		logger.Error("Error querying database", err.Error())
	}

	defer row.Close()

	for row.Next() {
		var hash, name string
		err = row.Scan(&hash, &name)
		logger.Info("{ hash: \"" + hash + "\", name: \"" + name + "\" }")
	}

}
