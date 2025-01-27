package syncedpz

import (
	"log"
	"os"
	"syncedpz/config"

	"github.com/dgraph-io/badger"
)

func LoadPzDirs() error {
	err := config.DB.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("pz_exe_path"))
		if err != nil {
			return err
		}
		err = item.Value(func(val []byte) error {
			config.PZ_ExePath = string(val)
			return nil
		})
		if err != nil {
			return err
		}

		item, err = txn.Get([]byte("pz_data_path"))
		if err != nil {
			return err
		}
		err = item.Value(func(val []byte) error {
			config.PZ_DataPath = string(val)
			return nil
		})
		return err
	})
	return err
}

func SetupPzDirs(pz_exe_path, pz_data_path string) {
	if _, err := os.Stat(pz_data_path); os.IsNotExist(err) {
		log.Fatalf("%s dir does not exists\n", pz_data_path)
	}
	if _, err := os.Stat(pz_exe_path); os.IsNotExist(err) {
		log.Fatalf("%s dir does not exists\n", pz_exe_path)
	}

	// Saves into badger
	err := config.DB.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte("pz_exe_path"), []byte(pz_exe_path))
		if err != nil {
			return err
		}
		err = txn.Set([]byte("pz_data_path"), []byte(pz_data_path))
		return err
	})
	if err != nil {
		log.Fatal(err)
	}

	config.PZ_ExePath = pz_exe_path
	config.PZ_DataPath = pz_data_path
}

func LoadSteamID() error {
	err := config.DB.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("steam_id"))
		if err != nil {
			return err
		}
		err = item.Value(func(val []byte) error {
			config.PZ_SteamID = string(val)
			return nil
		})
		return err
	})
	return err
}

func SetupSteamId(steamID string) {
	err := config.DB.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte("steam_id"), []byte(steamID))
		return err
	})
	if err != nil {
		log.Fatal(err)
	}

	config.PZ_SteamID = steamID
}
