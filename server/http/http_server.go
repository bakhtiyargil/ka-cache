package http

import (
	"context"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"ka-cache/config"
	"ka-cache/logger"
	"ka-cache/server"
	"net/http"
	"os"
	"time"
)

type SimpleHttpServer struct {
	server    *http.Server
	handler   Handler
	cfg       *config.Config
	echo      *echo.Echo
	logger    logger.Logger
	stopChan  chan os.Signal
	isRunning bool
}

func NewHttpServer(cfg *config.Config, logger logger.Logger, handler Handler) server.Server {
	s := &http.Server{
		Addr:           ":" + cfg.Server.Default.Port,
		ReadTimeout:    time.Second * cfg.Server.Default.ReadTimeout,
		WriteTimeout:   time.Second * cfg.Server.Default.WriteTimeout,
		MaxHeaderBytes: cfg.Server.Default.MaxHeaderBytes,
	}
	return &SimpleHttpServer{
		server:   s,
		handler:  handler,
		cfg:      cfg,
		echo:     echo.New(),
		logger:   logger,
		stopChan: make(chan os.Signal, 1),
	}
}

func (s *SimpleHttpServer) Start() {
	if s.Running() {
		s.logger.Fatal("http server is already running")
	}
	go func() {
		s.logger.Infof("http server is listening on port: %s", s.cfg.Server.Default.Port)
		s.echo.HideBanner = true
		s.echo.HidePort = true
		if err := s.echo.StartServer(s.server); err != nil {
			s.logger.Fatalf("failed to start http server: %v", err)
		}
	}()

	amw := NewApiMiddlewareManager(s.cfg.Server.Default.AllowOrigins, s.logger)
	s.appendMiddleware(s.echo, amw)
	s.appendRoutes(s.echo)
	s.isRunning = true
}

func (s *SimpleHttpServer) Stop() {
	if !s.Running() {
		s.logger.Fatal("http server is not running")
	}
	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	s.logger.Info("http server exited properly")
	err := s.echo.Server.Shutdown(ctx)
	if err != nil {
		s.logger.Fatalf("error shutting down http server: %v", err)
	}
}

func (s *SimpleHttpServer) Running() bool {
	return s.isRunning
}

func (s *SimpleHttpServer) appendMiddleware(e *echo.Echo, manager MiddlewareManager) {
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

	e.Use(manager.RequestLoggerMiddleware)
	e.Use(manager.CorsMiddleware)
}

func (s *SimpleHttpServer) appendRoutes(e *echo.Echo) {
	base := e.Group("/cache")
	s.handler.mapBaseRouteHandlers(base)
	health := base.Group("/health")
	s.handler.mapHealthRouteHandlers(health)
}
