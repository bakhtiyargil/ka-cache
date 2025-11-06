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
		if !ok {
			reqId := c.Response().Header().Get(echo.HeaderXRequestID)
			return NewResourceNotFound(reqId, nil)
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

		if (len(i.Key) == 0) || (i.Ttl == 0) {
			reqId := c.Response().Header().Get(echo.HeaderXRequestID)
			return NewValidationError(reqId, "validation failed for object: [ttl] > 0; [key] not empty")
		}

		err := h.cache.Put(i.Key, i.Value, i.Ttl)
		if err != nil {
			return err
		}
		return c.NoContent(http.StatusOK)
	})

	base.DELETE("/:key", func(c echo.Context) error {
		itemKey := c.Param("key")
		h.cache.Delete(itemKey)
		return c.NoContent(http.StatusOK)
	})
}
