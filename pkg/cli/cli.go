package cli

import (
	"fmt"
	"os"
	"runtime"
	"syncedpz/config"
	"syncedpz/pkg/syncedpz"
)

func printUsage() {
	fmt.Println("Usage: ")
	fmt.Println("  syncedpz setup")
	fmt.Println("  syncedpz run")
}

func validateArgs() {
	if len(os.Args) < 2 {
		printUsage()
		runtime.Goexit()
	}
}

func setup() {
	handleErr := func(err error) {
		if err != nil {
			config.FirstTimeSetup = true
			setup()
		}
	}

	if config.FirstTimeSetup {
		exe_path := askForInput("Enter the path to the pz executable: ")
		data_path := askForInput("Enter the path to the pz data directory: ")
		syncedpz.SetupPzDirs(exe_path, data_path)

		steamID := askForInput("Enter your steam id: ")
		syncedpz.SetupSteamId(steamID)
		config.PZ_SteamID = steamID
	} else {
		handleErr(syncedpz.LoadSteamID())
		handleErr(syncedpz.LoadPzDirs())
	}
}

func Run() {
	//validateArgs()

	setup()

	fmt.Printf("Hello, World!\n\n")
	fmt.Printf("PZ exe path: %s\n", config.PZ_ExePath)
	fmt.Printf("PZ data path: %s\n", config.PZ_DataPath)
	fmt.Printf("Steam ID: %s\n", config.PZ_SteamID)
}
