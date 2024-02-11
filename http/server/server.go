package server

import (
	"github.com/labstack/echo/v4"
	"golang.org/x/net/context"
	"ka-cache/config"
	"ka-cache/pkg/logger"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	echo   *echo.Echo
	cfg    *config.Config
	logger logger.Logger
}

func NewServer(cfg *config.Config, logger logger.Logger) *Server {
	return &Server{echo: echo.New(), cfg: cfg, logger: logger}
}

func (s *Server) Run() error {
	server := &http.Server{
		Addr:           ":" + s.cfg.Server.Default.Port,
		ReadTimeout:    time.Second * s.cfg.Server.Default.ReadTimeout,
		WriteTimeout:   time.Second * s.cfg.Server.Default.WriteTimeout,
		MaxHeaderBytes: s.cfg.Server.Default.MaxHeaderBytes,
	}

	go func() {
		s.logger.Infof("Server is listening on PORT: %s", s.cfg.Server.Default.Port)
		if err := s.echo.StartServer(server); err != nil {
			s.logger.Errorf("Error starting Server: ", err)
			os.Exit(1)
		}
	}()

	if err := s.MapHandlers(s.echo); err != nil {
		return err
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	s.logger.Info("Server Exited Properly")
	return s.echo.Server.Shutdown(ctx)
}
