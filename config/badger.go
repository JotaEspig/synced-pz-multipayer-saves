package config

import (
	"log"
	"os"

	"github.com/dgraph-io/badger"
)

const (
	badgerDBPath = "data"
)

var (
	FirstTimeSetup bool
	DB             *badger.DB
)

// return true if the directory already exists, false if it was created
func ensureDir(dirName string) bool {
	if _, err := os.Stat(dirName); !os.IsNotExist(err) {
		return true
	}

	err := os.MkdirAll(dirName, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	return false
}

func init() {
	var err error

	FirstTimeSetup = !ensureDir(badgerDBPath)

	opts := badger.DefaultOptions(badgerDBPath)
	opts.Logger = nil
	opts.Truncate = true

	DB, err = badger.Open(opts)
	if err != nil {
		log.Fatal(err)
	}
}
