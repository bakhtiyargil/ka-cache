package http

import (
	"github.com/labstack/echo/v4"
	"ka-cache/cache"
	"ka-cache/model"
	"net/http"
)

type Handler interface {
	mapHealthRouteHandlers(health *echo.Group)
	mapBaseRouteHandlers(base *echo.Group)
}

type CacheHandler struct {
	cache cache.Cache[string, string]
}

func NewCacheHandler(cache cache.Cache[string, string]) Handler {
	return &CacheHandler{
		cache: cache,
	}
}

func (h *CacheHandler) mapHealthRouteHandlers(health *echo.Group) {
	health.GET("", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})
}

func (h *CacheHandler) mapBaseRouteHandlers(base *echo.Group) {
	base.GET("/:key", func(c echo.Context) error {
		itemKey := c.Param("key")
		value, ok := h.cache.Get(itemKey)
		reqId := c.Response().Header().Get(echo.HeaderXRequestID)
		if !ok {
			err := NewResourceNotFound(reqId, nil)
			return err
		}
		data := model.DataResponse{
			Data: value,
		}
		return c.JSON(http.StatusOK, data)
	})

	base.PUT("", func(c echo.Context) error {
		i := &model.Item{}
		if err := c.Bind(i); err != nil {
			return err
		}
		err := h.cache.Put(i.Key, i.Value, i.Ttl)
		if err != nil {
			return err
		}
		return c.NoContent(http.StatusOK)
	})
}
