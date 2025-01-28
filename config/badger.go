package config

import (
	"log"
	"syncedpz/pkg/utils"

	"github.com/dgraph-io/badger"
)

const (
	badgerDBPath = DataPath + "/badger"
)

var (
	FirstTimeSetup bool
	DB             *badger.DB
)

func init() {
	var err error

	FirstTimeSetup = !utils.EnsureDir(DataPath)

	opts := badger.DefaultOptions(badgerDBPath)
	opts.Logger = nil
	opts.Truncate = true

	DB, err = badger.Open(opts)
	if err != nil {
		log.Fatal(err)
	}
}
