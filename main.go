package main

import (
	"os"
	"os/signal"
	"syncedpz/config"
	"syncedpz/pkg/cli"
	"syscall"

	"github.com/charmbracelet/log"
)

func main() {
	defer os.Exit(0)
	defer func() {
		config.DB.Close()
		log.Info("BadgerDB closed")
	}()

	// Gracefully shutdown
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	go cli.Run(ch)

	<-ch
}
