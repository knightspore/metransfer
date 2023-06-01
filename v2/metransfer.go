package main

import (
	"crypto/sha1"
	"database/sql"
	"encoding/hex"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"os"

	_ "github.com/mattn/go-sqlite3"
)

const (
	CreateTable = `CREATE TABLE upload (
		hash TEXT PRIMARY KEY,
		name TEXT
	);`
	InsertFile      = "INSERT INTO upload(hash, name) VALUES (?, ?)"
	QueryRecords    = "SELECT * FROM upload"
	QueryFileByHash = "SELECT hash, name FROM upload WHERE hash=$1"
)

type MeTransfer struct {
	Path string
	DB   string
}

func (m *MeTransfer) Connect() (*sql.DB, error) {

	db, err := sql.Open("sqlite3", m.DBPath())
	if err != nil {
		return nil, err
	}

	return db, nil

}

func (m *MeTransfer) CreateHash(name string, size int64) string {

	bytes := sha1.New()
	bytes.Write([]byte(name + fmt.Sprint(size)))

	return hex.EncodeToString(bytes.Sum(nil))

}

func (m *MeTransfer) CreateTable() (sql.Result, error) {

	db, err := m.Connect()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	stmt, err := db.Prepare(CreateTable)
	if err != nil {
		return nil, err
	}

	result, err := stmt.Exec()
	if err != nil {
		return nil, err
	}

	return result, nil

}

func (m *MeTransfer) DBPath() string {
	return fmt.Sprintf("./%s", m.DB)
}

func (m *MeTransfer) Init() error {

	_, err := os.Create(m.DBPath())
	if err != nil {
		return err
	}

	_, err = m.CreateTable()
	if err != nil {
		return err
	}

	files := m.ReadUploadsFolder()

	for _, file := range files {
		hash := m.CreateHash(file.Name(), file.Size())
		m.Insert(hash, file.Name())
	}

	return nil

}

func (m *MeTransfer) GetByHash(hash string) (string, error) {

	db, err := m.Connect()
	if err != nil {
		return "", err
	}

	defer db.Close()

	var name string

	row := db.QueryRow("SELECT name FROM upload WHERE hash=$1", hash)
	switch err := row.Scan(&name); err {
	case sql.ErrNoRows:
		return "", err
	case nil:
		return name, nil
	default:
		return "", err
	}

}

func (m *MeTransfer) Insert(hash, name string) error {

	db, err := m.Connect()
	if err != nil {
		return err
	}

	defer db.Close()

	stmt, err := db.Prepare("INSERT INTO upload(hash, name) VALUES (?, ?)")
	if err != nil {
		return err
	}

	_, err = stmt.Exec(hash, name)
	if err != nil {
		return err
	}

	return nil

}

func (m *MeTransfer) ListFiles() ([]string, error) {

	var results []string

	db, err := m.Connect()
	if err != nil {
		return nil, err
	}

	defer db.Close()

	row, err := db.Query(QueryRecords)
	if err != nil {
		return nil, err
	}

	defer row.Close()

	for row.Next() {
		var h, n string
		err = row.Scan(&h, &n)
		results = append(results, fmt.Sprintf("%s, %s", h, n))
	}

	return results, nil

}

func (m *MeTransfer) ReadUploadsFolder() []fs.FileInfo {

	files, err := ioutil.ReadDir(m.Path)
	if err != nil {
		log.Fatal(err)
	}

	return files

}
