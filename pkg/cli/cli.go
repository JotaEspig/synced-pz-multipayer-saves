package cli

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"syncedpz/config"
	"syncedpz/pkg/syncedpz"

	"github.com/charmbracelet/log"
)

func tryParseCommand(cmd *flag.FlagSet) {
	err := cmd.Parse(os.Args[2:])
	if err != nil {
		cmd.Usage()
		runtime.Goexit()
	}
}

func Run() {
	defer func() {
		// Read a single byte (key press)
		fmt.Print(config.GTM("Press any key to exit..."))
		var b [1]byte
		_, _ = os.Stdin.Read(b[:])
	}()

	if err := syncedpz.LoadLanguage(); err != nil {
		setLanguage()
	}

	menuCmd := flag.NewFlagSet("menu", flag.ExitOnError)
	configCmd := flag.NewFlagSet("config", flag.ExitOnError)
	listCmd := flag.NewFlagSet("list", flag.ExitOnError)
	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	deleteCmd := flag.NewFlagSet("delete", flag.ExitOnError)
	cloneCmd := flag.NewFlagSet("clone", flag.ExitOnError)
	syncCmd := flag.NewFlagSet("sync", flag.ExitOnError)
	runCmd := flag.NewFlagSet("run", flag.ExitOnError)

	listType := listCmd.String("type", "local", config.GTM("Type of servers to list"))

	setup()

	if len(os.Args) < 2 {
		menu()
		return
	}

	switch os.Args[1] {
	case "menu":
		tryParseCommand(menuCmd)
	case "config":
		tryParseCommand(configCmd)
	case "run":
		tryParseCommand(runCmd)
	case "list":
		tryParseCommand(listCmd)
	case "add":
		tryParseCommand(addCmd)
	case "delete":
		tryParseCommand(deleteCmd)
	case "clone":
		tryParseCommand(cloneCmd)
	case "sync":
		tryParseCommand(syncCmd)
	default:
		printUsage()
		runtime.Goexit()
	}

	if menuCmd.Parsed() {
		menu()
	} else if configCmd.Parsed() {
		if configCmd.Arg(0) == "setup" {
			config.FirstTimeSetup = true
			setup()
		} else if configCmd.Arg(0) == "list" {
			listConfig()
		} else {
			log.Fatal(config.GTM("No argument for config"))
		}
	} else if listCmd.Parsed() {
		switch *listType {
		case "local":
			listLocalServers()
		case "synced":
			listSyncedServers()
		default:
			listCmd.Usage()
			runtime.Goexit()
		}
	} else if addCmd.Parsed() {
		addServer()
	} else if deleteCmd.Parsed() {
		deleteServer()
	} else if cloneCmd.Parsed() {
		cloneServer()
	} else if syncCmd.Parsed() {
		syncServers()
	}
}
