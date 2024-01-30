package main

import (
	"ka-cache/api/server"
	"ka-cache/config"
)

func main() {
	lConfig, _ := config.LoadConfig("./config/config-local")
	eServer := server.NewServer(lConfig)
	err := eServer.Run()
	if err != nil {
		return
	}
}
