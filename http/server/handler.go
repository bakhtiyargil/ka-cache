package server

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"ka-cache/cache"
	"ka-cache/model"
	"net/http"
	"strings"
)

var cac = cache.NewCache(5)

func (s *Server) MapHandlers(e *echo.Echo) error {
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderXRequestID},
	}))
	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize:         1 << 10, // 1 KB
		DisablePrintStack: true,
		DisableStackAll:   true,
	}))
	e.Use(middleware.RequestID())
	e.Use(middleware.GzipWithConfig(middleware.GzipConfig{
		Level: 5,
		Skipper: func(c echo.Context) bool {
			return strings.Contains(c.Request().URL.Path, "swagger")
		},
	}))
	e.Use(middleware.Secure())
	e.Use(middleware.BodyLimit("2M"))

	v1 := e.Group("/cache/v1")
	health := v1.Group("/health")

	MapCacheRoutes(v1)

	health.GET("", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"status": "OK"})
	})

	return nil
}

func GetHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		itemKey := c.Param("key")
		item := cac.Get(itemKey)
		return c.JSON(http.StatusOK, item)
	}
}

func PutHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		i := &model.Item{}
		c.Bind(i)
		cac.Set(i.Key, i.Value)
		return c.JSON(http.StatusOK, map[string]string{"status": "OK"})
	}
}
