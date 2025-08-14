package main

import (
	"ka-cache/bootstrap"
	"ka-cache/server/grpc"
	"ka-cache/server/http"
	"os"
)

func main() {
	go startHttpServer()
	startGrpcServer()
}

func startHttpServer() {
	eServer := http.NewHttpServer(bootstrap.App.Config, bootstrap.App.Logger)
	eServer.Start()
}

func startGrpcServer() {
	err := grpc.NewGrpcServer(bootstrap.App.Config).Start()
	if err != nil {
		os.Exit(1)
	}
}
