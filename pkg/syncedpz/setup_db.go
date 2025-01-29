package syncedpz

import (
	"encoding/binary"
	"fmt"
	"os"
	"syncedpz/config"
	"syncedpz/pkg/utils"

	"github.com/charmbracelet/log"
	"github.com/dgraph-io/badger"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
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
	if err != nil {
		return err
	}
	if _, err := os.Stat(config.PZ_ExePath); os.IsNotExist(err) {
		return fmt.Errorf("%s dir does not exists\n", config.PZ_ExePath)
	}
	if _, err := os.Stat(config.PZ_DataPath); os.IsNotExist(err) {
		return fmt.Errorf("%s dir does not exists\n", config.PZ_DataPath)
	}
	return nil
}

func SetupPzDirs(PZ_ExePath, PZ_DataPath string) {
	if _, err := os.Stat(PZ_ExePath); os.IsNotExist(err) {
		log.Fatalf("%s dir does not exists\n", PZ_ExePath)
	}
	if _, err := os.Stat(PZ_DataPath); os.IsNotExist(err) {
		log.Fatalf("%s dir does not exists\n", PZ_DataPath)
	}

	err := config.DB.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte("pz_exe_path"), []byte(PZ_ExePath))
		if err != nil {
			return err
		}
		err = txn.Set([]byte("pz_data_path"), []byte(PZ_DataPath))
		return err
	})
	utils.HandleErr(err)

	utils.HandleErr(LoadPzDirs())
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
	if err != nil {
		return err
	}

	if config.PZ_SteamID == "" {
		return fmt.Errorf("Steam ID cannot be empty")
	}
	return nil
}

func SetupSteamId(steamID string) {
	if steamID == "" {
		log.Fatal("Steam ID cannot be empty")
	}

	err := config.DB.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte("steam_id"), []byte(steamID))
		return err
	})
	utils.HandleErr(err)

	utils.HandleErr(LoadSteamID())
}

func LoadGitAuth() error {
	var gitUser, gitPass string
	err := config.DB.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("git_username"))
		if err != nil {
			return err
		}
		err = item.Value(func(val []byte) error {
			gitUser = string(val)
			return nil
		})
		if err != nil {
			return err
		}

		item, err = txn.Get([]byte("git_password"))
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			gitPass = string(val)
			return nil
		})
	})
	if err != nil {
		return err
	}

	if gitUser == "" || gitPass == "" {
		return fmt.Errorf("Git username or password cannot be empty")
	}
	config.GitAuth = &http.BasicAuth{
		Username: gitUser,
		Password: gitPass,
	}
	return nil
}

func SetupGitAuth(username, password string) {
	if username == "" || password == "" {
		log.Fatal("Git username or password cannot be empty")
	}

	err := config.DB.Update(func(txn *badger.Txn) error {
		err := txn.Set([]byte("git_username"), []byte(username))
		if err != nil {
			return err
		}
		err = txn.Set([]byte("git_password"), []byte(password))
		return err
	})
	utils.HandleErr(err)

	utils.HandleErr(LoadGitAuth())
}

func LoadLanguage() error {
	err := config.DB.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("language"))
		if err != nil {
			return err
		}

		return item.Value(func(val []byte) error {
			config.Launguage = int(binary.BigEndian.Uint32(val))
			return nil
		})
	})

	if err != nil {
		return err
	}
	if !config.IsLanguageValid(config.Launguage) {
		return fmt.Errorf("Invalid language")
	}
	return nil
}

func SetupLanguage(lang int) {
	if !config.IsLanguageValid(lang) {
		log.Fatal("Invalid language")
	}

	err := config.DB.Update(func(txn *badger.Txn) error {
		val := make([]byte, 4)
		binary.BigEndian.PutUint32(val, uint32(lang))
		err := txn.Set([]byte("language"), val)
		return err
	})
	utils.HandleErr(err)

	utils.HandleErr(LoadLanguage())
}
