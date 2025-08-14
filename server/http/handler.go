package http

import (
	"github.com/labstack/echo/v4"
	"ka-cache/cache"
	"ka-cache/model"
	"net/http"
)

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
