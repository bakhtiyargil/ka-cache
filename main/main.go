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
	//startGrpcServer()
	startDefaultServer()
}

func startDefaultServer() {
	lConfig, _ := config.LoadConfig("./config/config-local")
	loggr := logger.NewCustomLogger(lConfig)
	loggr.InitLogger()
	eServer := hs.NewServer(lConfig, loggr)
	err := eServer.Run()
	if err != nil {
		os.Exit(1)
	}
}

func startGrpcServer() {
	lis, _ := net.Listen("tcp", fmt.Sprintf("localhost:%d", 5000))
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterCacheServer(grpcServer, gs.NewServer())
	err := grpcServer.Serve(lis)
	if err != nil {
		os.Exit(1)
	}
}
