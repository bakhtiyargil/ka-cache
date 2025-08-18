package grpc

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"ka-cache/cache"
	"ka-cache/config"
	"ka-cache/logger"
	"ka-cache/server"
	"ka-cache/server/http"
	"log"
	"net"
)

type GrpcServer struct {
	cfg       *config.Config
	logger    logger.Logger
	isRunning bool
	UnimplementedCacheServer
}

func NewGrpcServer(cfg *config.Config, logger logger.Logger) server.Server {
	s := &GrpcServer{
		cfg:    cfg,
		logger: logger,
	}
	return s
}

func (s *GrpcServer) Put(ctx context.Context, item *Item) (*Response, error) {
	cache.SimpleCache.Put(item.Key, item.Value)
	log.Print("item: " + item.Key + " - successfully set")
	return &Response{
		Message: "success",
		Code:    1,
		Data:    "",
	}, nil
}

func (s *GrpcServer) Get(ctx context.Context, obj *Object) (*Response, error) {
	var item = cache.SimpleCache.Get(obj.Key)
	if item == "" {
		return nil, http.ResourceNotFoundError
	}
	log.Print("item: " + item + " - successfully get")
	return &Response{
		Message: "success",
		Code:    1,
		Data:    item,
	}, nil
}

func (s *GrpcServer) Start() {
	listener, _ := net.Listen("tcp", fmt.Sprintf("localhost:%s", s.cfg.Server.Grpc.Port))
	var opts []grpc.ServerOption
	grpcServer := grpc.NewServer(opts...)
	RegisterCacheServer(grpcServer, s)
	err := grpcServer.Serve(listener)
	if err != nil {
		s.logger.Fatalf("failed to start grpc server: %v", err)
	}
}

func (s *GrpcServer) Running() bool {
	return s.isRunning
}
