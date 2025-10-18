package http

import (
	"errors"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"ka-cache/logger"
	"time"
)

type MiddlewareManager interface {
	requestLoggerMiddleware(next echo.HandlerFunc) echo.HandlerFunc
	corsMiddleware(next echo.HandlerFunc) echo.HandlerFunc
	errorHandlerMiddleware(next echo.HandlerFunc) echo.HandlerFunc
}

type ApiMiddlewareManager struct {
	allowOrigins []string
	logger       logger.Logger
}

func NewApiMiddlewareManager(origins []string, logger logger.Logger) MiddlewareManager {
	return &ApiMiddlewareManager{allowOrigins: origins, logger: logger}
}

func (mw *ApiMiddlewareManager) requestLoggerMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		start := time.Now()

		err := next(ctx)
		req := ctx.Request()
		res := ctx.Response()
		status := res.Status
		requestID := getRequestID(ctx)
		elapsed := time.Since(start).String()
		mw.logger.Infof("RequestID: %s, Method: %s, URI: %s, Status: %v, Time: %s",
			requestID, req.Method, req.URL, status, elapsed,
		)
		return err
	}
}

func (mw *ApiMiddlewareManager) corsMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(ctx echo.Context) error {
		middleware.CORSWithConfig(middleware.CORSConfig{
			AllowOrigins: mw.allowOrigins,
			AllowHeaders: []string{
				echo.HeaderOrigin,
				echo.HeaderContentType,
				echo.HeaderAccept,
				echo.HeaderXRequestID,
			},
			ExposeHeaders: []string{echo.HeaderXRequestID},
		})
		return next(ctx)
	}
}

func (mw *ApiMiddlewareManager) errorHandlerMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		err := next(c)
		if err == nil {
			return nil
		}

		reqID := c.Response().Header().Get(echo.HeaderXRequestID)
		var restErr RestError
		if errors.As(err, &restErr) {
			mw.logger.Errorf("RequestID: %s, Error: %s", reqID, restErr.Causes())
			return c.JSON(restErr.Status(), restErr)
		}

		internalErr := NewInternalServerError(reqID, err.Error())
		mw.logger.Errorf("RequestID: %s, Error: %s", reqID, internalErr.Causes())
		return c.JSON(internalErr.Status(), internalErr)
	}
}

func getRequestID(c echo.Context) string {
	return c.Response().Header().Get(echo.HeaderXRequestID)
}
