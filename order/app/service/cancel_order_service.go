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
	CancelOrderRequest struct {
		OrderId int64 `json:"order_id" param:"order_id" query:"order_id"`
	}

	CancelOrderResponse struct {
		Message string       `json:"message"`
		Order   *model.Order `json:"order"`
	}

	cancelOrderService struct {
		db              *gorm.DB
		inventoryClient inventoryProto.InventoryServiceClient
		orderRepository repository.OrderRepository
	}

	CancelOrderService interface {
		Execute(ctx context.Context, req CancelOrderRequest) (*CancelOrderResponse, error)
	}
)

func NewCancelOrderService(
	instance connection.DBInstance,
	client grpcclient.GrpcClientInstance,
	orderRepository repository.OrderRepository) CancelOrderService {
	return &cancelOrderService{
		db:              instance.Database(),
		inventoryClient: client.InventoryConnection(),
		orderRepository: orderRepository,
	}
}

func (c *cancelOrderService) Execute(ctx context.Context, req CancelOrderRequest) (*CancelOrderResponse, error) {
	output := new(CancelOrderResponse)

	tx := c.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	order, err := c.orderRepository.Find(req.OrderId, ctx)
	if err != nil {
		return nil, err
	}

	if order.Status == "cancelled" {
		return nil, fmt.Errorf("order already cancelled")
	}

	// release stock first
	release, err := c.inventoryClient.ReleaseStock(ctx, &inventoryProto.ReleaseStockRequest{
		ProductId: order.ProductID,
		Quantity:  order.Quantity,
		OrderId:   order.ID,
	})
	if err != nil {
		return nil, err
	}

	if !release.Success {
		return nil, fmt.Errorf("failed to release stock")
	}

	// update order
	order.Status = "cancelled"
	if err = c.orderRepository.Update(order, tx, ctx); err != nil {
		return nil, err
	}

	tx.Commit()
	output.Order = order
	output.Message = release.Message
	return output, nil
}
