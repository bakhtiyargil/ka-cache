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
	cache cache.Cache
}

func NewCacheHandler(cache cache.Cache) Handler {
	return &CacheHandler{
		cache: cache,
	}
}

func (h *CacheHandler) mapHealthRouteHandlers(health *echo.Group) {
	health.GET("", func(c echo.Context) error {
		return c.JSON(http.StatusOK, model.NewSuccessResponse())
	})
}

func (h *CacheHandler) mapBaseRouteHandlers(base *echo.Group) {
	base.GET("/:key", func(c echo.Context) error {
		itemKey := c.Param("key")
		item := h.cache.Get(itemKey)
		if item == "" {
			err := NewResourceNotFound("")
			return c.JSON(err.Status(), err)
		}
		data := model.DataResponse{
			Data: item,
		}
		return c.JSON(http.StatusOK, data)
	})

	base.PUT("", func(c echo.Context) error {
		i := &model.Item{}
		if err := c.Bind(i); err != nil {
			internalError := NewInternalServerError(err.Error())
			return c.JSON(internalError.Status(), internalError)
		}
		h.cache.Put(i.Key, i.Value)
		return c.JSON(http.StatusOK, model.NewSuccessResponse())
	})
}
