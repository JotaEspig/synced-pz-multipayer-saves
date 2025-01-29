package syncedpz

import (
	"bufio"
	"bytes"
	"encoding/gob"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"syncedpz/config"
	"syncedpz/pkg/utils"

	"github.com/charmbracelet/log"
	"github.com/dgraph-io/badger"
	"github.com/go-git/go-git/v5"
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

	log.Info("Server saved to database")
}

// Delete deletes the server object from the database
func (ss SyncedServer) Delete() {
	err := config.DB.Update(func(txn *badger.Txn) error {
		return txn.Delete(ss.GetKey())
	})
	utils.HandleErr(err)

	log.Info("Server deleted from database")
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
func (ss *SyncedServer) CopyLocalServerToSynced() {
	log.Info("Copying local server to synced server")

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
			newConfigFilename := filepath.Join(configPath, filepath.Base(path))
			err = os.Remove(newConfigFilename)
			if err != nil && !os.IsNotExist(err) {
				return err
			}
			return cp.Copy(path, newConfigFilename)
		}
		return nil
	})
	utils.HandleErr(err)

	// copy save files
	savePath := filepath.Join(ss.GetServerPath(), "save")
	pzSaveFilesPath := filepath.Join(config.PZ_DataPath, "Saves", "Multiplayer")
	ssNameWithUnderScore := strings.ReplaceAll(ss.Name, " ", "_")
	fullLocalServerPath := filepath.Join(pzSaveFilesPath, ssNameWithUnderScore)

	log.Info("Removing old save files at synced server")
	err = os.RemoveAll(savePath)
	utils.HandleErr(err)

	log.Info("Copying new save files to synced server")
	utils.EnsureDir(savePath)
	err = cp.Copy(fullLocalServerPath, savePath)
	utils.HandleErr(err)

	log.Info("Local server copied to synced server")
}

func (ss *SyncedServer) CopySyncedServerToLocal() {
	log.Info("Copying synced server to local server")

	// copy config files
	configPath := filepath.Join(ss.GetServerPath(), "config")
	pzConfigFilesPath := filepath.Join(config.PZ_DataPath, "Server")
	err := filepath.Walk(configPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Checks if the filename starts with the server name
		if strings.HasPrefix(filepath.Base(path), ss.Name) {
			newConfigFilename := filepath.Join(pzConfigFilesPath, filepath.Base(path))
			err = os.Remove(newConfigFilename)
			if err != nil && !os.IsNotExist(err) {
				return err
			}
			return cp.Copy(path, newConfigFilename)
		}
		return nil
	})
	utils.HandleErr(err)

	// copy save files
	savePath := filepath.Join(ss.GetServerPath(), "save")
	pzSaveFilesPath := filepath.Join(config.PZ_DataPath, "Saves", "Multiplayer")
	ssNameWithUnderScore := strings.ReplaceAll(ss.Name, " ", "_")
	fullLocalServerPath := filepath.Join(pzSaveFilesPath, ssNameWithUnderScore)

	log.Info("Removing old save files at local server")
	err = os.RemoveAll(fullLocalServerPath)
	utils.HandleErr(err)

	log.Info("Copying new save files to local server")
	utils.EnsureDir(fullLocalServerPath)
	err = cp.Copy(savePath, fullLocalServerPath)
	utils.HandleErr(err)

	log.Info("Synced server copied to local server")
}

// EnsureUpdatedPlayerSaveFolders ensures that the player save folders are updated for each possible host of the server.
// In Project Zomboid, when player X is hosting, player Y's game will create a new player save folder having <SteamIDPlayerX> as suffix.
// So this function ensures that every player save folders for this server are updated.
func (ss *SyncedServer) EnsureUpdatedPlayerSaveFolders() {
	log.Info("Ensuring updated player save folders")

	playerSavePath := filepath.Join(config.PZ_DataPath, "Saves", "Multiplayer")
	ssNameWithUnderScore := strings.ReplaceAll(ss.Name, " ", "_")

	entries, err := os.ReadDir(playerSavePath)
	utils.HandleErr(err)

	playerFolders := []os.FileInfo{}
	for _, entry := range entries {
		hasNameInIt := strings.Contains(entry.Name(), ssNameWithUnderScore)
		if hasNameInIt && strings.HasSuffix(entry.Name(), "_player") {
			info, err := entry.Info()
			utils.HandleErr(err)
			playerFolders = append(playerFolders, info)
		}
	}
	if len(playerFolders) == 0 {
		return
	}

	sort.Slice(playerFolders, func(i, j int) bool {
		return playerFolders[i].ModTime().After(playerFolders[j].ModTime())
	})
	mostRecentPlayerFolderName := playerFolders[0].Name()

	// Ensures that a folder exist for every possible host
	for _, player := range ss.GetPlayers() {
		var playerFolderName string
		if player == config.PZ_SteamID {
			playerFolderName = ssNameWithUnderScore + "_player"
		} else {
			playerFolderName = player + "_" + ssNameWithUnderScore + "_player"
		}
		playerFolderPath := filepath.Join(playerSavePath, playerFolderName)
		utils.EnsureDir(playerFolderPath)
	}

	// Ensures that the most recent player folder is the most recent for every possible host
	entries, err = os.ReadDir(playerSavePath) // read again to get the updated list (if a new player folder was created)
	utils.HandleErr(err)

	for _, entry := range entries {
		hasNameInIt := strings.Contains(entry.Name(), ssNameWithUnderScore)
		if hasNameInIt && strings.HasSuffix(entry.Name(), "_player") {
			if entry.Name() == mostRecentPlayerFolderName {
				continue
			}
			err := os.RemoveAll(filepath.Join(playerSavePath, entry.Name()))
			utils.HandleErr(err)

			fullPathMostRecent := filepath.Join(playerSavePath, mostRecentPlayerFolderName)
			fullPathCurrent := filepath.Join(playerSavePath, entry.Name())
			err = cp.Copy(fullPathMostRecent, fullPathCurrent)
			utils.HandleErr(err)
		}
	}

	log.Info("Player save folders updated")
}

// GetPlayers returns the list of players in the server
func (ss SyncedServer) GetPlayers() []string {
	playersFilePath := filepath.Join(ss.GetServerPath(), "players.txt")
	players := []string{}

	if _, err := os.Stat(playersFilePath); !os.IsNotExist(err) {
		file, err := os.Open(playersFilePath)
		utils.HandleErr(err)
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			line := scanner.Text()
			line = strings.TrimSpace(line)
			line = strings.Trim(line, "\n")
			line = strings.Trim(line, "\r")
			players = append(players, line)
		}
		utils.HandleErr(scanner.Err())
	}

	return players
}

// UpdatePlayersFile updates the players file of the server.
// It creates the file if it doesn't exist, otherwise ensures your steam id is in the file
func (ss SyncedServer) UpdatePlayersFile() {
	log.Info("Updating players file")

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

	log.Info("Players file updated")
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
