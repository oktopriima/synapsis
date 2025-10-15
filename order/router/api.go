package router

import (
	"net/http"

	"synapsis/order/app/controller"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func NewRoute(
	e *echo.Echo,
	createOrder *controller.CreateController,
) {
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middleware.RequestID())

	route := e.Group("/api")

	// ping
	{
		route.GET("/ping", func(c echo.Context) error {
			return c.JSON(http.StatusOK, map[string]string{"message": "pong"})
		})
	}

	// order route group
	{
		order := route.Group("/orders")
		order.POST("", createOrder.Serve)
	}
}
