package syncedpz

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"syncedpz/config"
	"syncedpz/pkg/utils"

	"github.com/charmbracelet/log"
	"github.com/go-git/go-git/v5"
	gitconfig "github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing/transport"
)

// InitGit initializes the git repository for the synced server
func (ss *SyncedServer) InitGit() {
	log.Info("Initializing git repository")

	utils.EnsureDir(ss.GetServerPath())

	repo, err := git.PlainOpen(ss.GetServerPath())
	if err == git.ErrRepositoryNotExists {
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

	log.Info("Git repository initialized")
}

func (ss *SyncedServer) Clone() {
	log.Info("Starting to clone server")

	gitRepoName := filepath.Base(ss.GitURL)
	tempDirName := filepath.Join(config.ServersPath, gitRepoName)

	utils.EnsureDir(config.ServersPath)
	utils.EnsureDir(tempDirName)
	repo, err := git.PlainClone(tempDirName, false, &git.CloneOptions{
		URL:  ss.GitURL,
		Auth: config.GitAuth,
	})
	utils.HandleErr(err)

	ss.repo = repo

	// Get server name
	configPath := filepath.Join(tempDirName, "config")
	err = filepath.Walk(configPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if strings.HasSuffix(filepath.Base(path), ".ini") {
			ss.Name = strings.TrimSuffix(filepath.Base(path), ".ini")
			return errors.New("stop")
		}
		return nil
	})
	if err != nil && err.Error() != "stop" {
		log.Fatal(err)
	}

	// Renames the directory to the server name
	newDirName := ss.GetServerPath()
	err = os.Rename(tempDirName, newDirName)
	utils.HandleErr(err)

	// Recreates repo with the new directory
	ss.repo, err = git.PlainOpen(newDirName)
	utils.HandleErr(err)

	ss.Save()

	log.Info("Server cloned successfully")
}

// Restore restores the server to the last commit, useful to undo changes in case of a syncronization during
// an IO operation like copying files
func (ss *SyncedServer) Restore() {
	if ss.repo == nil {
		ss.InitGit()
	}

	log.Info("Starting to restore server")

	w, err := ss.repo.Worktree()
	utils.HandleErr(err)

	err = w.Reset(&git.ResetOptions{
		Mode: git.HardReset,
	})
	utils.HandleErr(err)
	removeUnstagedFiles(w)

	log.Info("Server restored")
}

// Fetch fetches the latest changes from the git repository
// Useful to check if anything was pushed when doing IO operations like
// copying local files to the synced server
func (ss *SyncedServer) Fetch() bool {
	if ss.repo == nil {
		ss.InitGit()
	}

	log.Info("Trying to fetch changes")

	err := ss.repo.Fetch(&git.FetchOptions{
		RemoteName: "origin",
		Auth:       config.GitAuth,
	})
	if err == git.NoErrAlreadyUpToDate {
		log.Info("Already up to date")
		return false
	} else {
		utils.HandleErr(err)
	}
	return true
}

// Pull pulls the latest changes from the git repository
// Returns true if there are new changes
func (ss *SyncedServer) Pull() bool {
	if ss.repo == nil {
		ss.InitGit()
	}

	log.Info("Starting to pull changes")

	w, err := ss.repo.Worktree()
	utils.HandleErr(err)

	err = w.Pull(&git.PullOptions{
		RemoteName: "origin",
		Auth:       config.GitAuth,
	})
	if err == git.NoErrAlreadyUpToDate || err == transport.ErrEmptyRemoteRepository {
		log.Info("Already up to date")
		return false
	} else {
		utils.HandleErr(err)
	}

	log.Info("Pulled changes")
	return true
}

func (ss *SyncedServer) Commit() {
	if ss.repo == nil {
		ss.InitGit()
	}

	log.Info("Starting to commit changes")

	w, err := ss.repo.Worktree()
	utils.HandleErr(err)

	_, err = w.Add(".")
	utils.HandleErr(err)

	commitMsg := fmt.Sprintf("SyncedPZ: synced by %s", config.PZ_SteamID)
	_, err = w.Commit(commitMsg, &git.CommitOptions{})
	if err != git.ErrEmptyCommit {
		utils.HandleErr(err)
		log.Info("Changes committed")
	} else {
		log.Info("No changes to commit")
	}
}

func (ss *SyncedServer) Push() {
	if ss.repo == nil {
		ss.InitGit()
	}

	log.Info("Starting to push changes")

	err := ss.repo.Push(&git.PushOptions{
		RemoteName: "origin",
		Auth:       config.GitAuth,
		Progress:   os.Stdout,
	})
	if err != git.NoErrAlreadyUpToDate {
		utils.HandleErr(err)
		log.Info("Changes pushed")
	} else {
		log.Info("Already up to date")
	}
}

func (ss *SyncedServer) CommitAndPush() {
	ss.Commit()
	ss.Push()
}
