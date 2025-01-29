package utils

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/charmbracelet/log"
)

// return true if the directory already exists, false if it was created
func EnsureDir(dirName string) bool {
	if _, err := os.Stat(dirName); !os.IsNotExist(err) {
		return true
	}

	err := os.MkdirAll(dirName, os.ModePerm)
	HandleErr(err)
	return false
}

// Runs command on directory
func RunCommandOnDir(dir string, command string, args ...string) error {
	fmt.Printf("Running command: %s", command)
	for _, arg := range args {
		fmt.Printf(" %s", arg)
	}
	fmt.Println()

	cmd := exec.Command(command, args...)
	cmd.Dir = dir
	return cmd.Run()
}

func RunCommandOnDirOutput(dir string, command string, args ...string) ([]byte, error) {
	fmt.Printf("Running command: %s", command)
	for _, arg := range args {
		fmt.Printf(" %s", arg)
	}
	fmt.Println()

	cmd := exec.Command(command, args...)
	cmd.Dir = dir
	return cmd.Output()
}

func HandleErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
