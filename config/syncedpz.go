package config

import (
	"os/exec"
	"syncedpz/pkg/utils"

	"github.com/charmbracelet/log"
	"github.com/go-git/go-git/v5/plumbing/transport"
)

const (
	LANG_START = iota
	LANG_EN    = iota
	LANG_PTBR  = iota
	LANG_END   = iota
)

var (
	PZ_DataPath string
	PZ_ExePath  string
	PZ_SteamID  string
	GitAuth     transport.AuthMethod
	ServersPath = DataPath + "/servers"
	Launguage   int
)

func IsLanguageValid(lang int) bool {
	return lang > LANG_START && lang < LANG_END
}

func checkExistenceOfGit() {
	_, err := exec.LookPath("git")
	if err != nil {
		log.Fatal("Git is not installed. Install it")
	}
}

func init() {
	utils.EnsureDir(DataPath)
	utils.EnsureDir(ServersPath)
	checkExistenceOfGit()
}
