package service

import (
	"context"
	"errors"
	"fmt"
	"synapsis/inventory/app/model"
	"synapsis/inventory/app/repository"
	"synapsis/inventory/database/connection"

	"gorm.io/gorm"
)

type (
	ReserveStockRequest struct {
		OrderId   int64 `json:"order_id"`
		ProductId int64 `json:"product_id"`
		Quantity  int64 `json:"quantity"`
	}

	ReserveStockResponse struct {
		Product *model.Product `json:"product"`
	}

	reserveStockService struct {
		db                      *gorm.DB
		productRepository       repository.ProductRepository
		stockRepository         repository.StockRepository
		stockMovementRepository repository.StockMovementRepository
	}
	ReserveStockService interface {
		Execute(ctx context.Context, req ReserveStockRequest) (*ReserveStockResponse, error)
	}
)

func NewReserveStockService(
	instance connection.DBInstance,
	productRepository repository.ProductRepository,
	stockRepository repository.StockRepository,
	stockMovementRepository repository.StockMovementRepository,
) ReserveStockService {
	return &reserveStockService{
		db:                      instance.Database(),
		productRepository:       productRepository,
		stockRepository:         stockRepository,
		stockMovementRepository: stockMovementRepository,
	}
}

func (r *reserveStockService) Execute(ctx context.Context, req ReserveStockRequest) (*ReserveStockResponse, error) {
	output := new(ReserveStockResponse)

	product, err := r.productRepository.Find(req.ProductId, ctx)
	if err != nil {
		return nil, err
	}

	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	stock, err := r.stockRepository.FindByProductTransaction(product.ID, tx, ctx)
	if err != nil {
		return nil, err
	}

	if stock.AvailableStock < req.Quantity {
		return nil, errors.New("not enough stock")
	}

	stock.AvailableStock -= req.Quantity
	stock.ReservedStock += req.Quantity

	if err = r.stockRepository.Update(stock, tx, ctx); err != nil {
		return nil, err
	}

	// create stock movement
	if err = r.stockMovementRepository.Create(&model.StockMovement{
		ProductId:   product.ID,
		ChangeType:  "RESERVED",
		Quantity:    req.Quantity,
		ReferenceId: req.OrderId,
		Note:        fmt.Sprintf("reserved stock %d for order %d", req.Quantity, req.OrderId),
	}, tx, ctx); err != nil {
		return nil, err
	}

	tx.Commit()
	return output, nil
}
