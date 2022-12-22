package main

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
)

func CreateHash(name string, size int64) string {
	bytes := sha1.New()
	bytes.Write([]byte(name + fmt.Sprint(size)))
	return hex.EncodeToString(bytes.Sum(nil))
}
