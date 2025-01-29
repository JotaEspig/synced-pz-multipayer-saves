package cli

import (
	"fmt"
	"strconv"
	"strings"
	"syncedpz/config"
	"syncedpz/pkg/syncedpz"
	"time"
)

func printUsage() {
	fmt.Println(config.GTM("Usage: "))
	fmt.Println(config.GTM("  syncedpz help = shows this message"))
	fmt.Println(config.GTM("  syncedpz menu = use menu mode"))
	fmt.Println(config.GTM("  syncedpz config [setup | list] = sets up or list the syncedpz configuration"))
	fmt.Println(config.GTM("  syncedpz list -type [local | synced] = list servers according to its type (default is local))"))
	fmt.Println(config.GTM("  syncedpz add = adds a new synced PZ server from your local files"))
	fmt.Println(config.GTM("  syncedpz delete = deletes a synced PZ server from the database only"))
	fmt.Println(config.GTM("  syncedpz clone = adds a new synced PZ server from a git repository"))
	fmt.Println(config.GTM("  syncedpz sync = syncs all servers"))
	fmt.Println(config.GTM("  syncedpz play = syncs all servers at the start and at the end. And starts Project Zomboid"))
	fmt.Println(config.GTM("  syncedpz language = sets the language of the application"))
}

func menu() {
	var err error

	functions := []func(){printUsage, setup, listConfig, listLocalServers, listSyncedServers, addServer, deleteServer, cloneServer, syncServers, setLanguage}

	choice := -1
	// Runs the menu until the user chooses to exit
	for true {
		// use > len(functions) and not >= len(functions) to include the exit option
		for choice < 0 || choice > len(functions) {
			fmt.Println(config.GTM("Menu:"))
			fmt.Println(config.GTM("  [0] Help"))
			fmt.Println(config.GTM("  [1] Setup config"))
			fmt.Println(config.GTM("  [2] List config"))
			fmt.Println(config.GTM("  [3] List local servers"))
			fmt.Println(config.GTM("  [4] List synced servers"))
			fmt.Println(config.GTM("  [5] Add synced server"))
			fmt.Println(config.GTM("  [6] Delete synced server"))
			fmt.Println(config.GTM("  [7] Clone synced server"))
			fmt.Println(config.GTM("  [8] Sync servers"))
			fmt.Println(config.GTM("  [9] Set language"))
			fmt.Println(config.GTM("  [10] Exit"))

			choiceStr := askForInput(config.GTM("Enter the number of the option you want to choose: "))
			choice, err = strconv.Atoi(choiceStr)
			if err != nil || choice < 0 || choice > len(functions) {
				fmt.Println(config.GTM("Invalid choice"))
				time.Sleep(1 * time.Second)
			}
		}

		if choice == len(functions) {
			break
		}
		if choice == 1 {
			config.FirstTimeSetup = true
		}

		fn := functions[choice]
		fn()

		choice = -1 // Reset choice to run the menu again
		time.Sleep(1 * time.Second)
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
		fmt.Println(config.GTM("Leave the field empty to use the previous value (if it exists)"))
		exe_path := askForInput(config.GTM("Enter the path to the pz executable: "))
		if exe_path == "" {
			exe_path = config.PZ_ExePath
		}
		data_path := askForInput(config.GTM("Enter the path to the pz data directory: "))
		if data_path == "" {
			data_path = config.PZ_DataPath
		}
		syncedpz.SetupPzDirs(exe_path, data_path)

		steamID := askForInput(config.GTM("Enter your steam id: "))
		if steamID == "" {
			steamID = config.PZ_SteamID
		}
		syncedpz.SetupSteamId(steamID)

		gitUsername := askForInput(config.GTM("Enter your git username: "))
		gitPassword := askForInput(config.GTM("Enter your git password (or your github token)): "))
		if gitUsername != "" || gitPassword != "" {
			syncedpz.SetupGitAuth(gitUsername, gitPassword)
		}
	} else {
		handleErr(syncedpz.LoadSteamID())
		handleErr(syncedpz.LoadPzDirs())
		handleErr(syncedpz.LoadGitAuth())
	}
}

func listConfig() {
	fmt.Println(config.GTM("PZ Exe Path: "), config.PZ_ExePath)
	fmt.Println(config.GTM("PZ Data Path: "), config.PZ_DataPath)
	fmt.Println(config.GTM("Steam ID: "), config.PZ_SteamID)
}

func listLocalServers() {
	servers := syncedpz.GetLocalServers()
	for i, s := range servers {
		fmt.Printf("[%d] - %s\n", i, s.Name)
	}
}

func listSyncedServers() {
	servers := syncedpz.GetSyncedServers()
	i := 0
	for _, ss := range servers {
		fmt.Printf("[%d] - %s\n", i, ss.Name)
		i++
	}
}

func addServer() {
	var err error

	localServers := syncedpz.GetLocalServers()
	choiceInt := -1
	for err != nil || choiceInt < 0 || choiceInt > len(localServers) {
		fmt.Println(config.GTM("Local Servers:"))
		for i, server := range localServers {
			fmt.Printf("[%d] %s\n", i, server.Name)
		}

		choice := askForInput(config.GTM("Enter the number of the server you want to add: "))
		choiceInt, err = strconv.Atoi(choice)
		if err != nil || choiceInt < 0 || choiceInt > len(localServers) {
			fmt.Println(config.GTM("Invalid choice"))
			time.Sleep(1 * time.Second)
		}
	}

	server := localServers[choiceInt]
	gitURL := askForInput(config.GTM("Enter the git repository link to the server: "))
	ss := syncedpz.NewSyncedServer(server.Name, gitURL)

	ss.InitGit()
	changes := ss.Pull()
	if changes {
		fmt.Println(config.GTM("Warning! Apparently a server using this git repository already exists"))
		fmt.Println(config.GTM("and it already has some content."))
		fmt.Println(config.GTM("Do you want to continue copying your local content to it?"))
		choice := askForInput(config.GTM("Enter y/N: "))
		choice = strings.ToLower(choice)
		choice = strings.TrimSpace(choice)
		if choice != "y" {
			fmt.Println(config.GTM("Aborting."))
			time.Sleep(1 * time.Second)
			return
		}
	}

	ss.CopyLocalServerToSynced()
	ss.UpdatePlayersFile()
	ss.CommitAndPush()
	ss.Save()

	fmt.Println(config.GTM("Server added successfully"))
}

func deleteServer() {
	servers := syncedpz.GetSyncedServers()
	if (len(servers)) == 0 {
		fmt.Println(config.GTM("No servers to delete"))
		return
	}

	serversArray := make([]*syncedpz.SyncedServer, len(servers))
	i := 0
	for _, ss := range servers {
		fmt.Printf("[%d] %s\n", i, ss.Name)
		serversArray[i] = ss
		i++
	}
	choice := askForInput(config.GTM("Enter the number of the server you want to delete: "))
	choiceInt, err := strconv.Atoi(choice)
	if err != nil || choiceInt < 0 || choiceInt >= len(servers) {
		fmt.Println(config.GTM("Invalid choice"))
		deleteServer()
		return
	}

	server := serversArray[choiceInt]
	server.Delete()
}

func cloneServer() {
	gitURL := askForInput(config.GTM("Enter the git repository link to the server: "))
	ss := syncedpz.SyncedServer{GitURL: gitURL}

	ss.Clone()
	ss.UpdatePlayersFile()
	ss.CommitAndPush()
	ss.Save()

	fmt.Println(config.GTM("Server cloned successfully"))
}

func syncServers() {
	servers := syncedpz.GetSyncedServers()
	for _, ss := range servers {
		// If there are changes in the server, prioritize pulling
		if ss.Pull() {
			ss.CopySyncedServerToLocal()
		} else { // If there are no changes, prioritize pushing
			ss.CopyLocalServerToSynced()
			ss.CommitAndPush()
		}
		ss.EnsureUpdatedPlayerSaveFolders()
	}
}

func setLanguage() {
	var err error

	choice := -1
	for !config.IsLanguageValid(choice) {
		fmt.Printf(" [%d] English\n", config.LANG_EN)
		fmt.Printf(" [%d] Português\n", config.LANG_PTBR)
		choiceStr := askForInput(config.GTM("Enter the number of the language you want to choose: "))
		choice, err = strconv.Atoi(choiceStr)
		if err != nil || !config.IsLanguageValid(choice) {
			fmt.Println(config.GTM("Invalid choice"))
			time.Sleep(1 * time.Second)
		}
	}

	syncedpz.SetupLanguage(choice)
}
