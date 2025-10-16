package http

import (
	"net/http"
	"synapsis/order/app/service"

	"github.com/labstack/echo/v4"
)

type OrderHandler struct {
	createService service.CreateOrderService
	cancelService service.CancelOrderService
}

func NewOrderHandler(
	createService service.CreateOrderService,
	cancelService service.CancelOrderService,
) *OrderHandler {
	return &OrderHandler{
		createService: createService,
		cancelService: cancelService,
	}
}

func (c *OrderHandler) CreateOrderHandler(ctx echo.Context) error {
	var req service.CreateOrderRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{
			"code":    http.StatusBadRequest,
			"message": err.Error(),
			"status":  "rejected",
		})
	}

	output, err := c.createService.Serve(&req, ctx.Request().Context())
	if err != nil {
		return ctx.JSON(http.StatusUnprocessableEntity, echo.Map{
			"code":    http.StatusUnprocessableEntity,
			"message": err.Error(),
			"status":  "rejected",
		})
	}

	return ctx.JSON(http.StatusOK, echo.Map{
		"code":   http.StatusOK,
		"status": "confirmed",
		"data": map[string]interface{}{
			"order":   output.Order,
			"product": output.Product,
		},
	})
}

func (c *OrderHandler) CancelOrderHandler(ctx echo.Context) error {
	var req service.CancelOrderRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{
			"code":    http.StatusBadRequest,
			"message": err.Error(),
			"status":  "failed",
		})
	}
	//
	output, err := c.cancelService.Execute(ctx.Request().Context(), req)
	if err != nil {
		return ctx.JSON(http.StatusUnprocessableEntity, echo.Map{
			"code":    http.StatusUnprocessableEntity,
			"message": err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, echo.Map{
		"code": http.StatusOK,
		"data": output.Order,
	})
}
