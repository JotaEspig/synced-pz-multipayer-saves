package syncedpz

import (
	"syncedpz/pkg/utils"

	"github.com/go-git/go-git/v5"
)

func removeUnstagedFiles(wt *git.Worktree) {
	status, err := wt.Status()
	utils.HandleErr(err)

	for file, statusEntry := range status {
		if statusEntry.Worktree == git.Untracked {
			_, err := wt.Remove(file)
			utils.HandleErr(err)
		}
	}
}
