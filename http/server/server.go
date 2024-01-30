package server

import (
	"github.com/labstack/echo/v4"
	"golang.org/x/net/context"
	"ka-cache/config"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type Server struct {
	echo *echo.Echo
	cfg  *config.Config
}

func NewServer(cfg *config.Config) *Server {
	return &Server{echo: echo.New(), cfg: cfg}
}

func (s *Server) Run() error {
	server := &http.Server{
		Addr:         s.cfg.Server.Port,
		ReadTimeout:  time.Second * s.cfg.Server.ReadTimeout,
		WriteTimeout: time.Second * s.cfg.Server.WriteTimeout,
	}

	go func() {
		if err := s.echo.StartServer(server); err != nil {
		}
	}()

	v1 := s.echo.Group("/api/cache")
	v1.GET("", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "OK"})
	})

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), 11*time.Second)
	defer shutdown()

	log.Print("Server Exited Properly")
	return s.echo.Server.Shutdown(ctx)
}
