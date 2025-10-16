package http

import (
	"net/http"
	"synapsis/inventory/app/service"

	"github.com/labstack/echo/v4"
)

type ProductHandler struct {
	createService   service.CreateProductService
	addStockService service.AddStockService
}

func NewProductHandler(
	createService service.CreateProductService,
	addStockService service.AddStockService,
) *ProductHandler {
	return &ProductHandler{
		createService:   createService,
		addStockService: addStockService,
	}
}

func (p *ProductHandler) CreateProduct(ctx echo.Context) error {
	var req service.CreateProductRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{
			"code":    http.StatusBadRequest,
			"message": err.Error(),
		})
	}

	output, err := p.createService.Execute(ctx.Request().Context(), req)
	if err != nil {
		return ctx.JSON(http.StatusUnprocessableEntity, echo.Map{
			"code":    http.StatusUnprocessableEntity,
			"message": err.Error(),
		})
	}

	return ctx.JSON(http.StatusCreated, echo.Map{
		"code": http.StatusCreated,
		"data": output,
	})
}

func (p *ProductHandler) AddStock(ctx echo.Context) error {
	var req service.AddStockRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{
			"code":    http.StatusBadRequest,
			"message": err.Error(),
		})
	}

	output, err := p.addStockService.Execute(ctx.Request().Context(), req)
	if err != nil {
		return ctx.JSON(http.StatusUnprocessableEntity, echo.Map{
			"code":    http.StatusUnprocessableEntity,
			"message": err.Error(),
		})
	}

	return ctx.JSON(http.StatusOK, echo.Map{
		"code": http.StatusOK,
		"data": output,
	})
}
