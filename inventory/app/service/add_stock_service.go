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
	AddStockRequest struct {
		ProductID int64 `json:"product_id"`
		Stock     int64 `json:"stock"`
	}

	AddStockResponse struct {
		Stock *model.Stock `json:"stock"`
	}

	AddStockService interface {
		Execute(ctx context.Context, req AddStockRequest) (*AddStockResponse, error)
	}

	addStockService struct {
		db                      *gorm.DB
		productRepository       repository.ProductRepository
		stockRepository         repository.StockRepository
		stockMovementRepository repository.StockMovementRepository
	}
)

func NewAddStockService(
	instance connection.DBInstance,
	productRepository repository.ProductRepository,
	stockRepository repository.StockRepository,
	stockMovementRepository repository.StockMovementRepository,
) AddStockService {
	return &addStockService{
		db:                      instance.Database(),
		productRepository:       productRepository,
		stockRepository:         stockRepository,
		stockMovementRepository: stockMovementRepository,
	}
}

func (a *addStockService) Execute(ctx context.Context, req AddStockRequest) (*AddStockResponse, error) {
	output := new(AddStockResponse)
	tx := a.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	stock, err := a.stockRepository.FindByProductTransaction(req.ProductID, tx, ctx)
	if err != nil {
		return nil, err
	}

	stock.AvailableStock += req.Stock
	stock.TotalStock += req.Stock

	// update stock
	if err = a.stockRepository.Update(stock, tx, ctx); err != nil {
		return nil, err
	}

	// add stock movement
	if err = a.stockMovementRepository.Create(&model.StockMovement{
		ProductId:   stock.ProductId,
		ChangeType:  "ADD",
		Quantity:    req.Stock,
		ReferenceId: stock.ProductId,
		Note:        fmt.Sprintf("Update stock for product %d ", req.ProductID),
	}, tx, ctx); err != nil {
		return nil, err
	}

	tx.Commit()

	output.Stock = stock
	return output, nil
}
