package main

import (
	"fmt"
	"google.golang.org/grpc"
	"ka-cache/config"
	pb "ka-cache/grpc"
	gs "ka-cache/grpc/server"
	hs "ka-cache/http/server"
	"ka-cache/pkg/logger"
	"net"
	"os"
)

func main() {
	lConfig, _ := config.LoadConfig("./config/config-local")
	ch := make(chan string)

	go startGrpcServer(ch, lConfig)
	go startDefaultServer(ch, lConfig)

	<-ch
	<-ch
}

func startDefaultServer(ch chan string, cnfg *config.Config) {
	loggr := logger.NewCustomLogger(cnfg)
	loggr.InitLogger()
	eServer := hs.NewServer(cnfg, loggr)
	err := eServer.Run()
	if err != nil {
		os.Exit(1)
	}
}

func startGrpcServer(ch chan string, cnfg *config.Config) {
	listener, _ := net.Listen("tcp", fmt.Sprintf("localhost:%s", cnfg.Server.Grpc.Port))
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterCacheServer(grpcServer, gs.NewServer())
	err := grpcServer.Serve(listener)
	if err != nil {
		os.Exit(1)
	}
}
