package main

import (
	"ka-cache/bootstrap"
	"ka-cache/server/grpc"
	"ka-cache/server/http"
)

func main() {
	go startHttpServer()
	startGrpcServer()
}

func startHttpServer() {
	hServer := http.NewHttpServer(bootstrap.App.Config, bootstrap.App.Logger, http.NewCacheHandler())
	hServer.Start()
}

func startGrpcServer() {
	gServer := grpc.NewGrpcServer(bootstrap.App.Config, bootstrap.App.Logger)
	gServer.Start()
}
