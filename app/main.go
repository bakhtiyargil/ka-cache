package main

import (
	"ka-cache/bootstrap"
	"ka-cache/server/grpc"
	"ka-cache/server/http"
	"os"
)

func main() {
	ch := make(chan string)
	go startGrpcServer(ch)
	go startDefaultServer(ch)
	<-ch
	<-ch
}

func startDefaultServer(ch chan string) {
	eServer := http.NewHttpServer(bootstrap.App.Config, bootstrap.App.Logger)
	err := eServer.Run()
	if err != nil {
		os.Exit(1)
	}
}

func startGrpcServer(ch chan string) {
	err := grpc.NewGrpcServer(bootstrap.App.Config).Run()
	if err != nil {
		os.Exit(1)
	}
}
