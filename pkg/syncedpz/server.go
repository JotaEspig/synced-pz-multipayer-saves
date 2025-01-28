package syncedpz

import (
	"io/fs"
	"path/filepath"
	"syncedpz/config"
)

type Server struct {
	Name string
}

func GetLocalServers() []*Server {
	serversConfigFilesPath := filepath.Join(config.PZ_DataPath, "Server")
	servers := []*Server{}

	filepath.Walk(serversConfigFilesPath, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if filepath.Ext(path) == ".ini" {
			server := &Server{
				Name: filepath.Base(path[:len(path)-len(filepath.Ext(path))]),
			}
			servers = append(servers, server)
		}

		return nil
	})

	return servers
}
