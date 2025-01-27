package main

import (
	"os"
	"syncedpz/config"
	"syncedpz/pkg/cli"
)

func main() {
	defer os.Exit(0)
	defer config.DB.Close()

	cli.Run()
}
