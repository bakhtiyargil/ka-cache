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
		Addr:         s.cfg.Server.Port,
		ReadTimeout:  time.Second * s.cfg.Server.ReadTimeout,
		WriteTimeout: time.Second * s.cfg.Server.WriteTimeout,
	}

	go func() {
		s.logger.Infof("Server is listening on PORT: %s", s.cfg.Server.Port)
		if err := s.echo.StartServer(server); err != nil {
			s.logger.Errorf("Error starting Server: ", err)
		}
	}()

	//todo refactor handle part above
	v1 := s.echo.Group("/api/cache")
	v1.GET("", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "OK"})
	})

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), 11*time.Second)
	defer shutdown()

	s.logger.Info("Server Exited Properly")
	return s.echo.Server.Shutdown(ctx)
}
