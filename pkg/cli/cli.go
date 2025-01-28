package cli

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"syncedpz/config"
)

func tryParseCommand(cmd *flag.FlagSet) {
	err := cmd.Parse(os.Args[2:])
	if err != nil {
		cmd.Usage()
		runtime.Goexit()
	}
}

func validateArgs() {
	if len(os.Args) < 2 {
		printUsage()
		runtime.Goexit()
	}
}

func Run() {
	validateArgs()

	configCmd := flag.NewFlagSet("config", flag.ExitOnError)
	listCmd := flag.NewFlagSet("list", flag.ExitOnError)
	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	cloneCmd := flag.NewFlagSet("clone", flag.ExitOnError)
	runCmd := flag.NewFlagSet("run", flag.ExitOnError)

	listType := listCmd.String("type", "local", "Type of servers to list")

	switch os.Args[1] {
	case "config":
		tryParseCommand(configCmd)
	case "run":
		tryParseCommand(runCmd)
	case "list":
		tryParseCommand(listCmd)
	case "add":
		tryParseCommand(addCmd)
	case "clone":
		tryParseCommand(cloneCmd)
	default:
		printUsage()
		runtime.Goexit()
	}

	if configCmd.Parsed() {
		if configCmd.Arg(0) == "setup" {
			config.FirstTimeSetup = true
			setup()
		} else if configCmd.Arg(0) == "list" {
			setup()
			listConfig()
		} else {
			fmt.Println("No argument for config")
			runtime.Goexit()
		}
		return
	}

	setup()

	if listCmd.Parsed() {
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
	} else if cloneCmd.Parsed() {
		cloneServer()
	}
}
