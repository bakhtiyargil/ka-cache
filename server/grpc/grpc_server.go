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
	cache     cache.Cache[string, string]
	UnimplementedCacheServer
}

func NewGrpcServer(cfg *config.Config, logger logger.Logger, cache cache.Cache[string, string]) server.Server {
	s := &GrpcServer{
		cfg:    cfg,
		logger: logger,
		cache:  cache,
	}
	return s
}

func (s *GrpcServer) Put(ctx context.Context, item *Item) (*Response, error) {
	err := s.cache.Put(item.Key, item.Value, item.Ttl)
	if err != nil {
		return nil, http.InternalServerError
	}
	log.Print("item: " + item.Key + " - successfully set")
	return &Response{
		Message: "success",
		Code:    1,
		Data:    "",
	}, nil
}

func (s *GrpcServer) Get(ctx context.Context, obj *Object) (*Response, error) {
	var entry, ok = s.cache.Get(obj.Key)
	if !ok {
		return nil, http.ResourceNotFoundError
	}
	s.logger.Info("item: " + obj.Key + " - successfully retrieved")
	return &Response{
		Message: "success",
		Code:    1,
		Data:    entry.Value,
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
