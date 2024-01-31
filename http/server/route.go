package server

import (
	"github.com/labstack/echo/v4"
)

func MapCacheRoutes(cacheGroup *echo.Group) {
	cacheGroup.GET("/:key", GetHandler())
	cacheGroup.PUT("/", PutHandler())
}
