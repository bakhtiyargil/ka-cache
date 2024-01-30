package main

import (
	"fmt"
	"google.golang.org/grpc"
	pb "ka-cache/grpc"
	"ka-cache/grpc/server"
	"log"
	"net"
)

func main() {
	/*	lConfig, _ := config.LoadConfig("./config/config-local")
		eServer := server.NewServer(lConfig)
		err := eServer.Run()
		if err != nil {
			return
		}*/

	lis, err := net.Listen("tcp", fmt.Sprintf("localhost:%d", 5000))
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterCacheServer(grpcServer, server.NewServer())
	grpcServer.Serve(lis)
}
