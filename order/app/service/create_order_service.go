package service

import (
	"context"
	"fmt"
	"synapsis/order/app/model"
	"synapsis/order/app/repository"
	grpcclient "synapsis/order/bootstrap/grpc-client"
	"synapsis/order/database/connection"
	inventoryProto "synapsis/proto-definitions/inventory"

	"gorm.io/gorm"
)

type (
	CreateOrderRequest struct {
		Product struct {
			Id       int64 `json:"id"`
			Quantity int64 `json:"quantity"`
		} `json:"products"`
	}

	CreateOrderResponse struct {
		Message string       `json:"message"`
		Order   *model.Order `json:"order"`
		Product interface{}  `json:"product"`
	}

	createOrderService struct {
		db              *gorm.DB
		inventoryClient inventoryProto.InventoryServiceClient
		orderRepo       repository.OrderRepository
	}

	CreateOrderService interface {
		Serve(req *CreateOrderRequest, ctx context.Context) (*CreateOrderResponse, error)
	}
)

func NewCreateOrderService(
	instance connection.DBInstance,
	client grpcclient.GrpcClientInstance,
	orderRepository repository.OrderRepository,
) CreateOrderService {
	return &createOrderService{
		db:              instance.Database(),
		inventoryClient: client.InventoryConnection(),
		orderRepo:       orderRepository,
	}
}

func (c *createOrderService) Serve(req *CreateOrderRequest, ctx context.Context) (*CreateOrderResponse, error) {
	stock, err := c.inventoryClient.CheckStock(ctx, &inventoryProto.CheckStockRequest{
		ProductId: req.Product.Id,
		Quantity:  req.Product.Quantity,
	})

	if err != nil {
		return nil, err
	}

	if !stock.GetIsAvailable() {
		return nil, fmt.Errorf("stock is not available")
	}

	tx := c.db.Begin()
	defer func() {
		tx.Rollback()
	}()

	order, err := c.orderRepo.Create(&model.Order{
		ProductID:   stock.Product.Id,
		ProductName: stock.Product.Name,
		Quantity:    req.Product.Quantity,
		Price:       float64(stock.Product.Price),
		FinalPrice:  float64(stock.Product.Price) * float64(req.Product.Quantity),
		Status:      "confirmed",
	}, tx, ctx)

	if err != nil {
		return nil, err
	}

	reserveResponse, err := c.inventoryClient.ReserveStock(ctx, &inventoryProto.ReserveStockRequest{
		ProductId: req.Product.Id,
		Quantity:  req.Product.Quantity,
		OrderId:   order.ID,
	})
	if err != nil {
		return nil, err
	}

	if !reserveResponse.Success {
		return nil, fmt.Errorf("reserve stock fail")
	}

	tx.Commit()

	output := new(CreateOrderResponse)
	output.Order = order
	output.Product = stock.Product
	output.Message = reserveResponse.Message

	return output, nil
}
