package http

import (
	"github.com/labstack/echo/v4"
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
	echo      *echo.Echo
	cfg       *config.Config
	logger    logger.Logger
	isRunning bool
}

func NewHttpServer(cfg *config.Config, logger logger.Logger) *SimpleHttpServer {
	return &SimpleHttpServer{echo: echo.New(), cfg: cfg, logger: logger}
}

func (s *SimpleHttpServer) Run() error {
	server := &http.Server{
		Addr:           ":" + s.cfg.Server.Default.Port,
		ReadTimeout:    time.Second * s.cfg.Server.Default.ReadTimeout,
		WriteTimeout:   time.Second * s.cfg.Server.Default.WriteTimeout,
		MaxHeaderBytes: s.cfg.Server.Default.MaxHeaderBytes,
	}

	go func() {
		s.logger.Infof("SimpleHttpServer is listening on PORT: %s", s.cfg.Server.Default.Port)
		if err := s.echo.StartServer(server); err != nil {
			s.logger.Errorf("Error starting SimpleHttpServer: ", err)
			os.Exit(1)
		}
	}()

	if err := s.MapHandlers(s.echo); err != nil {
		return err
	}
	s.isRunning = true

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	s.logger.Info("SimpleHttpServer Exited Properly")
	return s.echo.Server.Shutdown(ctx)
}

func (s *SimpleHttpServer) IsRunning() bool {
	return s.isRunning
}
