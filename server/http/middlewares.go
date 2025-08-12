package http

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"ka-cache/logger"
	"time"
)

type MiddlewareManager interface {
	RequestLoggerMiddleware(next echo.HandlerFunc) echo.HandlerFunc
	CorsMiddleware(next echo.HandlerFunc) echo.HandlerFunc
}

type ApiMiddlewareManager struct {
	origins []string
	logger  logger.Logger
}

func NewApiMiddlewareManager(origins []string, logger logger.Logger) MiddlewareManager {
	return &ApiMiddlewareManager{origins: origins, logger: logger}
}

func (mw *ApiMiddlewareManager) RequestLoggerMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		start := time.Now()
		err := next(ctx)

		req := ctx.Request()
		res := ctx.Response()
		status := res.Status
		size := res.Size
		s := time.Since(start).String()
		requestID := GetRequestID(ctx)

		mw.logger.Infof("RequestID: %s, Method: %s, URI: %s, Status: %v, Size: %v, Time: %s",
			requestID, req.Method, req.URL, status, size, s,
		)
		return err
	}
}

func (mw *ApiMiddlewareManager) CorsMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: mw.origins,
			AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderXRequestID},
		})
		return next(ctx)
	}
}

func GetRequestID(c echo.Context) string {
	return c.Response().Header().Get(echo.HeaderXRequestID)
}
