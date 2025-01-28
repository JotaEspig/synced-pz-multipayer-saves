package syncedpz

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"errors"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"syncedpz/config"
	"syncedpz/pkg/utils"

	"github.com/dgraph-io/badger"
	"github.com/go-git/go-git/v5"
	gitconfig "github.com/go-git/go-git/v5/config"
	cp "github.com/otiai10/copy"
)

type SyncedServer struct {
	Server
	GitURL string
	repo   *git.Repository
}

// NewSyncedServer creates a new synced server object
func NewSyncedServer(name, gitURL string) *SyncedServer {
	return &SyncedServer{
		Server: Server{Name: name},
		GitURL: gitURL,
	}
}

// LoadSyncedServer loads a synced server from the database
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
	utils.HandleErr(err)
	return ss
}

// GetKey returns the key of the server used in the database
func (ss SyncedServer) GetKey() []byte {
	return []byte("server_" + ss.Name)
}

// Serialize serializes the server object
func (ss SyncedServer) Serialize() []byte {
	var buff bytes.Buffer
	enc := gob.NewEncoder(&buff)
	err := enc.Encode(ss)
	utils.HandleErr(err)
	return buff.Bytes()
}

// Save saves the server object to the database
func (ss SyncedServer) Save() {
	err := config.DB.Update(func(txn *badger.Txn) error {
		return txn.Set(ss.GetKey(), ss.Serialize())
	})
	utils.HandleErr(err)
}

// Delete deletes the server object from the database
func (ss SyncedServer) Delete() {
	err := config.DB.Update(func(txn *badger.Txn) error {
		return txn.Delete(ss.GetKey())
	})
	utils.HandleErr(err)
}

// Delete deletes the server object from the database
func (ss SyncedServer) EnsureDirs() {
	serverPath := ss.GetServerPath()
	utils.EnsureDir(serverPath)
	utils.EnsureDir(filepath.Join(serverPath, "config"))
	utils.EnsureDir(filepath.Join(serverPath, "save"))
}

// GetServerPath returns the path of the server repository
func (ss SyncedServer) GetServerPath() string {
	return filepath.Join(config.ServersPath, ss.Name)
}

// CopyLocalServerToSynced copies the local server files to the synced server repository
func (ss SyncedServer) CopyLocalServerToSynced() {
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
	utils.HandleErr(err)

	// copy save files
	savePath := filepath.Join(ss.GetServerPath(), "save")
	pzSaveFilesPath := filepath.Join(config.PZ_DataPath, "Saves", "Multiplayer")

	ssNameWithUnderScore := strings.ReplaceAll(ss.Name, " ", "_")

	entries, err := os.ReadDir(pzSaveFilesPath)
	utils.HandleErr(err)
	for _, entry := range entries {
		path := filepath.Join(pzSaveFilesPath, entry.Name())
		// Checks if the filename starts with the server name and doesnt end with _player
		if strings.HasPrefix(filepath.Base(path), ssNameWithUnderScore) && !strings.HasSuffix(filepath.Base(path), "_player") {
			if err := cp.Copy(path, savePath); err != nil {
				log.Fatal(err)
			}
		}
	}
}

// UpdatePlayersFile updates the players file of the server.
// It creates the file if it doesn't exist, otherwise ensures your steam id is in the file
func (ss SyncedServer) UpdatePlayersFile() {
	playersFilePath := filepath.Join(ss.GetServerPath(), "players.txt")
	// if the file exists, ensure your steam id is in the file
	if _, err := os.Stat(playersFilePath); !os.IsNotExist(err) {
		file, err := os.Open(playersFilePath)
		utils.HandleErr(err)
		defer file.Close()

		scanner := bufio.NewScanner(file)
		found := false
		for scanner.Scan() {
			line := scanner.Text()
			line = strings.TrimSpace(line)
			line = strings.Trim(line, "\n")
			line = strings.Trim(line, "\r")
			if line == config.PZ_SteamID {
				found = true
				break
			}
		}
		utils.HandleErr(scanner.Err())

		if !found {
			file, err := os.OpenFile(playersFilePath, os.O_APPEND|os.O_WRONLY, 0600)
			utils.HandleErr(err)
			defer file.Close()

			_, err = file.WriteString(config.PZ_SteamID + "\n")
			utils.HandleErr(err)
		}
	} else { // otherwise create the file and write your steam id
		file, err := os.Create(playersFilePath)
		utils.HandleErr(err)
		defer file.Close()

		_, err = file.WriteString(config.PZ_SteamID + "\n")
		utils.HandleErr(err)
	}
}

// InitGit initializes the git repository for the synced server
func (ss *SyncedServer) InitGit() {
	utils.EnsureDir(ss.GetServerPath())

	repo, err := git.PlainOpen(ss.GetServerPath())
	if errors.Is(err, git.ErrRepositoryNotExists) {
		fmt.Println("Creating local git repository")
		repo, err = git.PlainInit(ss.GetServerPath(), false)
		utils.HandleErr(err)
		_, err = repo.CreateRemote(&gitconfig.RemoteConfig{
			Name: "origin",
			URLs: []string{ss.GitURL},
		})
		utils.HandleErr(err)
	} else if err != nil {
		log.Fatal(err)
	}

	ss.repo = repo
}

// Pull pulls the latest changes from the git repository
// Returns true if there are new changes
func (ss SyncedServer) Pull() bool {
	if ss.repo == nil {
		ss.InitGit()
	}

	w, err := ss.repo.Worktree()
	utils.HandleErr(err)

	err = w.Pull(&git.PullOptions{
		RemoteName: "origin",
		Auth:       config.GitAuth,
	})
	if err == git.NoErrAlreadyUpToDate {
		return false
	} else {
		utils.HandleErr(err)
	}

	return true
}

func (ss SyncedServer) Commit() {
	if ss.repo == nil {
		ss.InitGit()
	}

	w, err := ss.repo.Worktree()
	utils.HandleErr(err)

	_, err = w.Add(".")
	utils.HandleErr(err)

	commitMsg := fmt.Sprintf("SyncedPZ: synced by %s", config.PZ_SteamID)
	_, err = w.Commit(commitMsg, &git.CommitOptions{})
	utils.HandleErr(err)
}

func (ss SyncedServer) Push() {
	if ss.repo == nil {
		ss.InitGit()
	}

	err := ss.repo.Push(&git.PushOptions{
		RemoteName: "origin",
		Auth:       config.GitAuth,
		Progress:   os.Stdout,
	})
	utils.HandleErr(err)
}

func (ss SyncedServer) CommitAndPush() {
	ss.Commit()
	ss.Push()
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
