package config

import (
	"syncedpz/pkg/utils"

	"github.com/charmbracelet/log"
	"github.com/dgraph-io/badger"
)

var (
	FirstTimeSetup bool
	DB             *badger.DB
)

func InitDB(path string) {
	var err error

	FirstTimeSetup = !utils.EnsureDir(path)

	opts := badger.DefaultOptions(path + "/badger")
	opts.Logger = nil
	opts.Truncate = true

	DB, err = badger.Open(opts)
	utils.HandleErr(err)

	log.Info("BadgerDB initialized")
}

func init() {
	InitDB(DataPath)
}
