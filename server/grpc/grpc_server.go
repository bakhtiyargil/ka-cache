package grpc

import (
	"context"
	"crypto/tls"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/status"
	"ka-cache/cache"
	"ka-cache/config"
	"ka-cache/logger"
	"ka-cache/server"
	"net"
)

const success = "success"

type SimpleGrpcServer struct {
	server          *grpc.Server
	serverSecure    *grpc.Server
	cfg             *config.Config
	logger          logger.Logger
	isRunning       bool
	isSecureRunning bool
	cache           cache.Cache[string, string]
	UnimplementedCacheServer
}

func NewGrpcServer(cfg *config.Config, logger logger.Logger, cache cache.Cache[string, string]) server.Server {
	s := &SimpleGrpcServer{
		cfg:    cfg,
		logger: logger,
		cache:  cache,
	}
	return s
}

func (s *SimpleGrpcServer) Put(ctx context.Context, item *Item) (*Response, error) {
	err := s.cache.Put(item.Key, item.Value, item.Ttl)
	if err != nil {
		s.logger.Error(err)
		return nil, status.Error(codes.Internal, err.Error())
	}
	r := &Response{
		Code:    0,
		Message: success,
	}
	return r, nil
}

func (s *SimpleGrpcServer) Get(ctx context.Context, obj *Object) (*Response, error) {
	var value, ok = s.cache.Get(obj.Key)
	if !ok {
		return nil, status.Error(codes.NotFound, "resource not found")
	}
	r := &Response{
		Code:    0,
		Message: success,
		Data:    value,
	}
	return r, nil
}

func (s *SimpleGrpcServer) Start() {
	if s.Running() {
		s.logger.Fatal("grpc server is already running")
	}

	s.server = grpc.NewServer()

	addr := fmt.Sprintf(":%s", s.cfg.Server.Grpc.Port)
	go func() {
		listener, err := net.Listen("tcp", addr)
		if err != nil {
			s.logger.Fatalf("failed to listen on port %s: %v", addr, err)
		}

		RegisterCacheServer(s.server, s)

		if err := s.server.Serve(listener); err != nil {
			s.logger.Fatalf("failed to start grpc server: %v", err)
		}
	}()
	s.logger.Infof("grpc server is listening on port: %s", addr)
	s.isRunning = true
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

func (s *SimpleGrpcServer) StartSecure() {
	if s.SecureRunning() {
		s.logger.Fatal("grpc secure server is already running")
	}

	cert, err := tls.LoadX509KeyPair(s.cfg.Server.Grpc.CertFile, s.cfg.Server.Grpc.KeyFile)
	if err != nil {
		s.logger.Fatalf("failed to load TLS certificate: %v", err)
	}
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12,
	}
	creds := credentials.NewTLS(tlsConfig)
	s.serverSecure = grpc.NewServer(grpc.Creds(creds))

	addr := fmt.Sprintf(":%s", s.cfg.Server.Grpc.SecurePort)
	go func() {
		listener, err := net.Listen("tcp", addr)
		if err != nil {
			s.logger.Fatalf("failed to listen on port %s: %v", addr, err)
		}

		RegisterCacheServer(s.serverSecure, s)

		if err := s.serverSecure.Serve(listener); err != nil {
			s.logger.Fatalf("failed to start grpc secure server: %v", err)
		}
	}()
	s.logger.Infof("grpc secure server is listening on port %s", addr)
	s.isSecureRunning = true
}

func (s *SimpleGrpcServer) StopSecure() {
	if !s.SecureRunning() {
		s.logger.Fatal("grpc secure server is not running")
	}
	s.serverSecure.GracefulStop()
	s.logger.Info("grpc secure server exited properly")
	s.isSecureRunning = false
}

func (s *SimpleGrpcServer) SecureRunning() bool {
	return s.isSecureRunning
}
