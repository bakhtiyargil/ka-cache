package http

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/net/context"
	"ka-cache/config"
	"ka-cache/logger"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type SimpleHttpServer struct {
	server    *http.Server
	handler   Handler
	echo      *echo.Echo
	logger    logger.Logger
	isRunning bool
}

func NewHttpServer(cfg *config.Config, logger logger.Logger, handler Handler) *SimpleHttpServer {
	server := &http.Server{
		Addr:           ":" + cfg.Server.Default.Port,
		ReadTimeout:    time.Second * cfg.Server.Default.ReadTimeout,
		WriteTimeout:   time.Second * cfg.Server.Default.WriteTimeout,
		MaxHeaderBytes: cfg.Server.Default.MaxHeaderBytes,
	}
	return &SimpleHttpServer{
		server:  server,
		handler: handler,
		echo:    echo.New(),
		logger:  logger,
	}
}

func (s *SimpleHttpServer) Start() {
	if s.Running() {
		s.logger.Fatal("http server is already running")
	}
	go func() {
		s.logger.Infof("http server is listening on port: %s", s.server.Addr)
		if err := s.echo.StartServer(s.server); err != nil {
			s.logger.Fatalf("error starting http server: %v", err)
		}
	}()

	amw := NewApiMiddlewareManager([]string{"*"}, s.logger)
	s.appendMiddleware(s.echo, amw)
	s.appendRoutes(s.echo)
	s.isRunning = true

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

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
	e.Use(manager.RequestLoggerMiddleware)
	e.Use(manager.CorsMiddleware)
	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize:         1 << 10,
		DisablePrintStack: true,
		DisableStackAll:   true,
	}))
	e.Use(middleware.RequestID())
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
	}))
	e.Use(middleware.Secure())
	e.Use(middleware.BodyLimit("2M"))
}

func (s *SimpleHttpServer) appendRoutes(e *echo.Echo) {
	base := e.Group("/cache")
	s.handler.mapBaseRouteHandlers(base)
	health := base.Group("/health")
	s.handler.mapHealthRouteHandlers(health)
}
