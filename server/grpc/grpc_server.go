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

type SimpleGrpcServer struct {
	cfg       *config.Config
	logger    logger.Logger
	isRunning bool
	cache     cache.Cache[string, string]
	server    *grpc.Server
	UnimplementedCacheServer
}

func NewGrpcServer(cfg *config.Config, logger logger.Logger, cache cache.Cache[string, string]) server.Server {
	s := &SimpleGrpcServer{
		cfg:    cfg,
		logger: logger,
		cache:  cache,
		server: grpc.NewServer(),
	}
	return s
}

func (s *SimpleGrpcServer) Put(ctx context.Context, item *Item) (*Response, error) {
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

func (s *SimpleGrpcServer) Get(ctx context.Context, obj *Object) (*Response, error) {
	var value, ok = s.cache.Get(obj.Key)
	if !ok {
		return nil, http.ResourceNotFoundError
	}
	s.logger.Info("item: " + obj.Key + " - successfully retrieved")
	return &Response{
		Message: "success",
		Code:    1,
		Data:    value,
	}, nil
}

func (s *SimpleGrpcServer) Start() {
	if s.Running() {
		s.logger.Fatal("grpc server is already running")
	}
	go func() {
		listener, _ := net.Listen("tcp", fmt.Sprintf("localhost:%s", s.cfg.Server.Grpc.Port))
		RegisterCacheServer(s.server, s)
		err := s.server.Serve(listener)
		if err != nil {
			s.logger.Fatalf("failed to start grpc server: %v", err)
		}
	}()
	s.isRunning = true
	s.logger.Infof("grpc server is listening on port: %s", s.cfg.Server.Grpc.Port)
}

func (s *SimpleGrpcServer) Stop() {
	if !s.Running() {
		s.logger.Fatal("grpc server is not running")
	}
	s.logger.Info("grpc server exited properly")
	s.server.GracefulStop()
	s.isRunning = false
}

func (s *SimpleGrpcServer) Running() bool {
	return s.isRunning
}
