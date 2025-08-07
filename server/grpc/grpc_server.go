package grpc

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"ka-cache/config"
	"ka-cache/server"
	err "ka-cache/server/http/error"
	"log"
	"net"
)

type GrpcServer struct {
	cfg       *config.Config
	isRunning bool
	UnimplementedCacheServer
}

func NewServer(cfg *config.Config) server.Server {
	s := &GrpcServer{
		cfg: cfg,
	}
	return s
}

func (s *GrpcServer) Put(ctx context.Context, item *Item) (*Response, error) {
	config.SimpleCache.Set(item.Key, item.Value)
	log.Print("item: " + item.Key + " - successfully set")
	return &Response{
		Message: "success",
		Code:    1,
		Data:    "",
	}, nil
}

func (s *GrpcServer) Get(ctx context.Context, obj *Object) (*Response, error) {
	var item = config.SimpleCache.Get(obj.Key)
	if item == "" {
		return nil, err.ResourceNotFoundError
	}
	log.Print("item: " + item + " - successfully get")
	return &Response{
		Message: "success",
		Code:    1,
		Data:    item,
	}, nil
}

func (s *GrpcServer) Run() error {
	listener, _ := net.Listen("tcp", fmt.Sprintf("localhost:%s", s.cfg.Server.Grpc.Port))
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	RegisterCacheServer(grpcServer, s)
	err1 := grpcServer.Serve(listener)
	return err1
}

func (s *GrpcServer) IsRunning() bool {
	return s.isRunning
}
