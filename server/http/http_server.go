package http

import (
	"context"
	"crypto/tls"
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"ka-cache/config"
	"ka-cache/logger"
	"ka-cache/server"
	"net/http"
	"time"
)

type SimpleHttpServer struct {
	server            *http.Server
	serverSecure      *http.Server
	echo              *echo.Echo
	echoSecure        *echo.Echo
	handler           Handler
	middlewareManager MiddlewareManager
	cfg               *config.Config
	logger            logger.Logger
	isRunning         bool
	isSecureRunning   bool
}

func NewHttpServer(cfg *config.Config, logger logger.Logger, handler Handler) server.Server {
	s := &http.Server{
		Addr:           ":" + cfg.Server.Default.Port,
		ReadTimeout:    time.Second * cfg.Server.Default.ReadTimeout,
		WriteTimeout:   time.Second * cfg.Server.Default.WriteTimeout,
		MaxHeaderBytes: cfg.Server.Default.MaxHeaderBytes,
	}
	secureS := &http.Server{
		Addr:           ":" + cfg.Server.Default.SecurePort,
		ReadTimeout:    time.Second * cfg.Server.Default.ReadTimeout,
		WriteTimeout:   time.Second * cfg.Server.Default.WriteTimeout,
		MaxHeaderBytes: cfg.Server.Default.MaxHeaderBytes,
	}
	e := echo.New()
	eS := echo.New()
	amw := NewApiMiddlewareManager(cfg.Server.Default.AllowOrigins, logger)

	return &SimpleHttpServer{
		server:            s,
		serverSecure:      secureS,
		echo:              e,
		echoSecure:        eS,
		middlewareManager: amw,
		handler:           handler,
		cfg:               cfg,
		logger:            logger,
	}
}

func (s *SimpleHttpServer) Start() {
	if s.Running() {
		s.logger.Fatal("http server is already running")
	}

	s.setupEcho(s.echo)
	go func() {
		s.logger.Infof("http server is listening on port: %s", s.cfg.Server.Default.Port)
		if err := s.echo.StartServer(s.server); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.logger.Fatalf("failed to start http server: %v", err)
		}
	}()
	s.isRunning = true
}

func (s *SimpleHttpServer) Stop() {
	if !s.Running() {
		s.logger.Fatal("http server is not running")
	}
	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	err := s.server.Shutdown(ctx)
	if err != nil {
		s.logger.Fatalf("error shutting down http server: %v", err)
	} else {
		s.logger.Info("http server exited properly")
		s.isRunning = false
	}
}

func (s *SimpleHttpServer) Running() bool {
	return s.isRunning
}

func (s *SimpleHttpServer) StartSecure() {
	if s.SecureRunning() {
		s.logger.Fatal("https server is already running")
	}

	s.setupEcho(s.echoSecure)

	certFile := s.cfg.Server.Default.CertFile
	keyFile := s.cfg.Server.Default.KeyFile
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		s.logger.Fatalf("failed to load TLS certificate: %v", err)
	}
	s.serverSecure.TLSConfig = &tls.Config{
		Certificates: []tls.Certificate{cert},
		MinVersion:   tls.VersionTLS12,
	}

	go func() {
		s.logger.Infof("https server is listening on port: %s", s.cfg.Server.Default.SecurePort)
		if err := s.echoSecure.StartServer(s.serverSecure); err != nil && !errors.Is(err, http.ErrServerClosed) {
			s.logger.Fatalf("failed to start https server: %v", err)
		}
	}()
	s.isSecureRunning = true
}

func (s *SimpleHttpServer) StopSecure() {
	if !s.SecureRunning() {
		s.logger.Fatal("https server is not running")
	}
	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	err := s.serverSecure.Shutdown(ctx)
	if err != nil {
		s.logger.Fatalf("error shutting down https server: %v", err)
	} else {
		s.logger.Info("https server exited properly")
		s.isSecureRunning = false
	}
}

func (s *SimpleHttpServer) SecureRunning() bool {
	return s.isSecureRunning
}

func (s *SimpleHttpServer) setupEcho(e *echo.Echo) {
	s.appendMiddleware(e)
	s.appendRoutes(e)
	e.HideBanner = true
	e.HidePort = true
}

func (s *SimpleHttpServer) appendMiddleware(e *echo.Echo) {
	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize:         1 << 10,
		DisablePrintStack: true,
		DisableStackAll:   true,
	}))
	e.Use(middleware.SecureWithConfig(middleware.SecureConfig{
		XSSProtection:         middleware.DefaultSecureConfig.XSSProtection,
		ContentTypeNosniff:    middleware.DefaultSecureConfig.ContentTypeNosniff,
		XFrameOptions:         middleware.DefaultSecureConfig.XFrameOptions,
		HSTSMaxAge:            31536000,
		HSTSExcludeSubdomains: false,
		HSTSPreloadEnabled:    true,
		ContentSecurityPolicy: "default-src 'self'",
		ReferrerPolicy:        "strict-origin-when-cross-origin",
	}))

	e.Use(middleware.RequestID())
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{Level: 5}))
	e.Use(middleware.BodyLimit("2M"))

	e.Use(s.middlewareManager.RequestLoggerMiddleware)
	e.Use(s.middlewareManager.ErrorHandlerMiddleware)
	e.Use(s.middlewareManager.CorsMiddleware)
}

func (s *SimpleHttpServer) appendRoutes(e *echo.Echo) {
	base := e.Group("/cache")
	s.handler.mapBaseRouteHandlers(base)
	health := base.Group("/health")
	s.handler.mapHealthRouteHandlers(health)
}
