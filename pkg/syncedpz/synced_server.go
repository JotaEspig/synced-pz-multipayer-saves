package syncedpz

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"syncedpz/config"
	"syncedpz/pkg/utils"

	"github.com/dgraph-io/badger"
	cp "github.com/otiai10/copy"
)

type SyncedServer struct {
	Server
	GitURL  string
	Players []string
}

func NewSyncedServer(name, gitURL string) *SyncedServer {
	return &SyncedServer{
		Server:  Server{Name: name},
		GitURL:  gitURL,
		Players: []string{},
	}
}

func LoadSyncedServer(name string) *SyncedServer {
	ss := &SyncedServer{}
	err := config.DB.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte("server_" + name))
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			dec := gob.NewDecoder(bytes.NewReader(val))
			return dec.Decode(&ss)
		})
	})
	if err != nil {
		log.Fatal(err)
	}
	return ss
}

func (ss SyncedServer) GetKey() []byte {
	return []byte("server_" + ss.Name)
}

func (ss SyncedServer) Serialize() []byte {
	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)
	err := enc.Encode(ss)
	if err != nil {
		log.Fatal(err)
	}
	return buff.Bytes()
}

func (ss SyncedServer) Save() {
	err := config.DB.Update(func(txn *badger.Txn) error {
		return txn.Set(ss.GetKey(), ss.Serialize())
	})
	if err != nil {
		log.Fatal(err)
	}
}

func (ss SyncedServer) EnsureDirs() {
	serverPath := ss.GetServerPath()
	utils.EnsureDir(serverPath)
	utils.EnsureDir(filepath.Join(serverPath, "config"))
	utils.EnsureDir(filepath.Join(serverPath, "save"))
}

func (ss SyncedServer) GetServerPath() string {
	return filepath.Join(config.ServersPath, ss.Name)
}

func (ss SyncedServer) CopyPZServerToDir() {
	ss.EnsureDirs()

	// copy config files
	configPath := filepath.Join(ss.GetServerPath(), "config")
	pzConfigFilesPath := filepath.Join(config.PZ_DataPath, "Server")
	err := filepath.Walk(pzConfigFilesPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Checks if the filename starts with the server name
		if strings.HasPrefix(filepath.Base(path), ss.Name) {
			newConfigPath := filepath.Join(configPath, filepath.Base(path))
			return cp.Copy(path, newConfigPath)
		}
		return nil
	})
	if err != nil {
		log.Fatal(err)
	}

	// copy save files
	savePath := filepath.Join(ss.GetServerPath(), "save")
	pzSaveFilesPath := filepath.Join(config.PZ_DataPath, "Saves", "Multiplayer")
	fmt.Println(pzSaveFilesPath)

	ssNameWithUnderScore := strings.ReplaceAll(ss.Name, " ", "_")
	entries, err := os.ReadDir(pzSaveFilesPath)
	if err != nil {
		log.Fatal(err)
	}
	for _, entry := range entries {
		path := filepath.Join(pzSaveFilesPath, entry.Name())
		fmt.Println(path)
		// Checks if the filename starts with the server name and doesnt end with _player
		if strings.HasPrefix(filepath.Base(path), ssNameWithUnderScore) && !strings.HasSuffix(filepath.Base(path), "_player") {
			fmt.Println("ups")
			if err := cp.Copy(path, savePath); err != nil {
				log.Fatal(err)
			}
		}
	}

}

func GetSyncedServers() map[string]*SyncedServer {
	servers := map[string]*SyncedServer{}
	err := config.DB.View(func(txn *badger.Txn) error {
		// Get all servers (it starts with prefix : "server_")
		opts := badger.DefaultIteratorOptions
		opts.PrefetchValues = false
		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Seek([]byte("server_")); it.ValidForPrefix([]byte("server_")); it.Next() {
			item := it.Item()
			err := item.Value(func(v []byte) error {
				dec := gob.NewDecoder(bytes.NewReader(v))
				ss := &SyncedServer{}
				err := dec.Decode(&ss)
				if err != nil {
					return err
				}
				servers[ss.Name] = ss
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return map[string]*SyncedServer{}
	}

	return servers
}
