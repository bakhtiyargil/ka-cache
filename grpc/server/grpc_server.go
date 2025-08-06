package server

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"ka-cache/config"
	err "ka-cache/http/error"
	"log"
	"net"
)
import pb "ka-cache/grpc"

type GrpcServer struct {
	cfg       *config.Config
	isRunning bool
	pb.UnimplementedCacheServer
}

func NewServer(cfg *config.Config) *GrpcServer {
	s := &GrpcServer{
		cfg: cfg,
	}
	return s
}

func (s *GrpcServer) Put(ctx context.Context, item *pb.Item) (*pb.Response, error) {
	config.DefaultCache.Set(item.Key, item.Value)
	log.Print("item: " + item.Key + " - successfully set")
	return &pb.Response{
		Message: "success",
		Code:    1,
		Data:    "",
	}, nil
}

func (s *GrpcServer) Get(ctx context.Context, obj *pb.Object) (*pb.Response, error) {
	var item = config.DefaultCache.Get(obj.Key)
	if item == "" {
		return nil, err.ResourceNotFoundError
	}
	log.Print("item: " + item + " - successfully get")
	return &pb.Response{
		Message: "success",
		Code:    1,
		Data:    item,
	}, nil
}

func (s *GrpcServer) Run() error {
	listener, _ := net.Listen("tcp", fmt.Sprintf("localhost:%s", s.cfg.Server.Grpc.Port))
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	pb.RegisterCacheServer(grpcServer, s)
	err1 := grpcServer.Serve(listener)
	return err1
}

func (s *GrpcServer) IsRunning() bool {
	return s.isRunning
}
