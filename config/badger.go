package config

import (
	"syncedpz/pkg/utils"

	"github.com/charmbracelet/log"
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
	utils.HandleErr(err)

	log.Info("BadgerDB initialized")
}
