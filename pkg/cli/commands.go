package cli

import (
	"fmt"
	"strconv"
	"syncedpz/config"
	"syncedpz/pkg/syncedpz"
)

func printUsage() {
	fmt.Println("Usage: ")
	fmt.Println("  syncedpz help = shows this message")
	fmt.Println("  syncedpz config [setup | list] = sets up or list the syncedpz configuration")
	fmt.Println("  syncedpz list -type [local | synced] = list servers according to its type (default is local)")
	fmt.Println("  syncedpz add = adds a new synced PZ server from your local files")
	fmt.Println("  syncedpz clone = adds a new synced PZ server from a git repository")
	fmt.Println("  syncedpz sync = syncs all servers")
	fmt.Println("  syncedpz play = syncs a server at the start and at the end. And starts Project Zomboid")
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

func listConfig() {
	fmt.Println("PZ Exe Path: ", config.PZ_ExePath)
	fmt.Println("PZ Data Path: ", config.PZ_DataPath)
	fmt.Println("Steam ID: ", config.PZ_SteamID)
}

func listLocalServers() {
	servers := syncedpz.GetLocalServers()
	for _, server := range servers {
		fmt.Println(server.Name)
	}
}

func listSyncedServers() {
	servers := syncedpz.GetSyncedServers()
	for _, ss := range servers {
		fmt.Println(ss.Name)
	}
}

func addServer() {
	localServers := syncedpz.GetLocalServers()
	fmt.Println("Local Servers:")
	for i, server := range localServers {
		fmt.Printf("[%d] %s\n", i, server.Name)
	}
	choice := askForInput("Enter the number of the server you want to add: ")
	choiceInt, err := strconv.Atoi(choice)
	if err != nil || choiceInt < 0 || choiceInt >= len(localServers) {
		fmt.Println("Invalid choice")
		addServer()
		return
	}

	server := localServers[choiceInt]
	gitURL := askForInput("Enter the github link to the server: ")
	ss := syncedpz.NewSyncedServer(server.Name, gitURL)
	ss.Save()
	ss.CopyPZServerToDir()
}
