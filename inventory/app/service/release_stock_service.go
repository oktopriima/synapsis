package service

import (
	"context"
	"fmt"
	"synapsis/inventory/app/model"
	"synapsis/inventory/app/repository"
	"synapsis/inventory/database/connection"

	"gorm.io/gorm"
)

type (
	ReleaseStockRequest struct {
		ProductId int64 `json:"product_id"`
		Quantity  int64 `json:"quantity"`
		OrderId   int64 `json:"order_id"`
	}

	ReleaseStockResponse struct {
		Message string       `json:"message"`
		Stock   *model.Stock `json:"stock"`
	}
	releaseStockService struct {
		db                      *gorm.DB
		stockRepository         repository.StockRepository
		stockMovementRepository repository.StockMovementRepository
	}

	ReleaseStockService interface {
		Execute(req ReleaseStockRequest, ctx context.Context) (*ReleaseStockResponse, error)
	}
)

func NewReleaseStockService(
	instance connection.DBInstance,
	stockRepository repository.StockRepository,
	stockMovementRepository repository.StockMovementRepository) ReleaseStockService {
	return &releaseStockService{
		db:                      instance.Database(),
		stockRepository:         stockRepository,
		stockMovementRepository: stockMovementRepository,
	}
}

func (r *releaseStockService) Execute(req ReleaseStockRequest, ctx context.Context) (*ReleaseStockResponse, error) {
	output := new(ReleaseStockResponse)
	tx := r.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	stock, err := r.stockRepository.FindByProductTransaction(req.ProductId, tx, ctx)
	if err != nil {
		return nil, err
	}

	stock.AvailableStock += req.Quantity
	stock.ReservedStock = stock.ReservedStock - req.Quantity

	if err = r.stockRepository.Update(stock, tx, ctx); err != nil {
		return nil, err
	}

	if err = r.stockMovementRepository.Create(&model.StockMovement{
		ProductId:   req.ProductId,
		ChangeType:  "RELEASE",
		Quantity:    req.Quantity,
		ReferenceId: req.OrderId,
		Note:        fmt.Sprintf("release stock %d for order %d", req.Quantity, req.OrderId),
	}, tx, ctx); err != nil {
		return nil, err
	}

	tx.Commit()
	return output, nil
}
