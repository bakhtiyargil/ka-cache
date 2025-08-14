package main

import (
	"ka-cache/bootstrap"
	"ka-cache/server/grpc"
	"ka-cache/server/http"
	"os"
)

func main() {
	go startDefaultServer()
	startGrpcServer()
}

func startDefaultServer() {
	eServer := http.NewHttpServer(bootstrap.App.Config, bootstrap.App.Logger)
	err := eServer.Run()
	if err != nil {
		os.Exit(1)
	}
}

func startGrpcServer() {
	err := grpc.NewGrpcServer(bootstrap.App.Config).Run()
	if err != nil {
		os.Exit(1)
	}
}
