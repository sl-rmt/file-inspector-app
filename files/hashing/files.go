// Package hashing provides methods for hashing files and data
package hashing

import (
	"crypto/sha256"
	"fmt"
	"io"
	"log"
	"os"
)

// GetFileSHA256HashString reads the file, sha256 hashes it and return the hash as a hex string
func GetFileSHA256HashString(fileString string) string {
	f, err := os.Open(fileString)

	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
	}

	return fmt.Sprintf("%x", h.Sum(nil))
}

// GetFileSHA256Hash reads the file, sha256 hashes it and return the hash as a byte array
func GetFileSHA256Hash(fileString string) []byte {
	f, err := os.Open(fileString)

	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		log.Fatal(err)
	}

	return h.Sum(nil)
}
