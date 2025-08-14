package http

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"ka-cache/cache"
	"ka-cache/model"
	"net/http"
	"strings"
)

func (s *SimpleHttpServer) MapHandlers(e *echo.Echo) error {
	amw := NewApiMiddlewareManager([]string{"*"}, s.logger)
	e.Use(amw.RequestLoggerMiddleware)
	e.Use(amw.CorsMiddleware)
	e.Use(middleware.RecoverWithConfig(middleware.RecoverConfig{
		StackSize:         1 << 10,
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
		return c.JSON(http.StatusOK, model.NewSuccessResponse())
	})

	return nil
}

func GetHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		itemKey := c.Param("key")
		item := cache.SimpleCache.Get(itemKey)
		if item == "" {
			err := NewResourceNotFound("")
			return c.JSON(err.Status(), err)
		}
		data := model.DataResponse{
			Data: item,
		}
		return c.JSON(http.StatusOK, data)
	}
}

func PutHandler() echo.HandlerFunc {
	return func(c echo.Context) error {
		i := &model.Item{}
		if err := c.Bind(i); err != nil {
			internalError := NewInternalServerError(err.Error())
			return c.JSON(internalError.Status(), internalError)
		}
		cache.SimpleCache.Set(i.Key, i.Value)
		return c.JSON(http.StatusOK, model.NewSuccessResponse())
	}
}
