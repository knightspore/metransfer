package main

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"

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
	path         string
	databaseName string
}

func NewMetransfer(path string, dbName string) *MeTransfer {
	s := MeTransfer{
		path,
		dbName,
	}
	return &s
}

func (m *MeTransfer) Path() string {
	return m.path
}

func (m *MeTransfer) DBPath() string {
	return fmt.Sprintf("%s/%s", m.path, m.databaseName)
}

func (m *MeTransfer) ReadUploadsFolder() []fs.FileInfo {
	files, err := ioutil.ReadDir(m.path)
	if err != nil {
		log.Fatal(err)
	}
	return files
}

func CreateHash(name string, size int64) string {
	bytes := sha1.New()
	bytes.Write([]byte(name + fmt.Sprint(size)))
	return hex.EncodeToString(bytes.Sum(nil))
}
