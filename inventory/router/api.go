package router

import (
	"net/http"
	handler "synapsis/inventory/app/handler/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func NewApiRouter(
	e *echo.Echo,
	productHandler *handler.ProductHandler,
) {
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middleware.RequestID())

	route := e.Group("/api")

	{
		route.GET("/ping", func(c echo.Context) error {
			return c.JSON(http.StatusOK, map[string]string{"message": "pong"})
		})
	}

	{
		route.POST("/product", productHandler.CreateProduct)
		route.POST("/product/stock", productHandler.AddStock)
	}
}
