package controller

import (
	"net/http"
	"synapsis/order/app/model"
	"synapsis/order/app/repository"
	"synapsis/order/database/connection"
	grpcclient "synapsis/order/grpc-client"
	pb "synapsis/proto-definitions/inventory"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type CreateController struct {
	db              *gorm.DB
	orderRepository repository.OrderRepository
	inventoryClient pb.InventoryServiceClient
}

func NewCreateController(
	instance connection.DBInstance,
	orderRepository repository.OrderRepository,
	grpcInstance grpcclient.GrpcClientInstance,
) *CreateController {
	return &CreateController{
		db:              instance.Database(),
		orderRepository: orderRepository,
		inventoryClient: grpcInstance.InventoryConnection(),
	}
}

type CreateOrderRequest struct {
	Products struct {
		Id       int64 `json:"id"`
		Quantity int64 `json:"quantity"`
	} `json:"products"`
}

func (c *CreateController) Serve(ctx echo.Context) error {
	var req CreateOrderRequest
	if err := ctx.Bind(&req); err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{
			"code":    http.StatusBadRequest,
			"message": err.Error(),
		})
	}

	product := req.Products
	stock, err := c.inventoryClient.CheckStock(ctx.Request().Context(), &pb.CheckStockRequest{
		ProductId: product.Id,
		Quantity:  product.Quantity,
	})

	if err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{
			"code":    http.StatusBadRequest,
			"message": err.Error(),
		})
	}

	if !stock.IsAvailable {
		return ctx.JSON(http.StatusBadRequest, echo.Map{
			"code":    http.StatusBadRequest,
			"message": "stock is not available",
		})
	}

	tx := c.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	order, err := c.orderRepository.Create(&model.Order{
		ID:          0,
		ProductID:   0,
		ProductName: "",
		Quantity:    0,
		Price:       0,
		FinalPrice:  0,
	}, tx, ctx.Request().Context())
	if err != nil {
		return ctx.JSON(http.StatusUnprocessableEntity, echo.Map{
			"code":    http.StatusUnprocessableEntity,
			"message": err.Error(),
		})
	}

	_, err = c.inventoryClient.ReserveStock(ctx.Request().Context(), &pb.ReserveStockRequest{
		ProductId: product.Id,
		Quantity:  product.Quantity,
		OrderId:   order.ID,
	})
	if err != nil {
		return ctx.JSON(http.StatusUnprocessableEntity, echo.Map{
			"code":    http.StatusUnprocessableEntity,
			"message": err.Error(),
		})
	}

	tx.Commit()
	return ctx.JSON(http.StatusOK, echo.Map{
		"code":     http.StatusOK,
		"message":  "SUCCESS",
		"products": req.Products,
	})
}
