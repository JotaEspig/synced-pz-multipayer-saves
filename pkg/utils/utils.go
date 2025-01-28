package utils

import (
	"log"
	"os"
)

// return true if the directory already exists, false if it was created
func EnsureDir(dirName string) bool {
	if _, err := os.Stat(dirName); !os.IsNotExist(err) {
		return true
	}

	err := os.MkdirAll(dirName, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	return false
}
