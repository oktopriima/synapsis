package service

import (
	"context"
	"synapsis/inventory/app/model"
	"synapsis/inventory/app/repository"
)

type (
	CheckStockRequest struct {
		ProductId int64 `json:"product_id"`
		Quantity  int64 `json:"quantity"`
	}

	CheckStockResponse struct {
		IsAvailable bool           `json:"is_available"`
		Product     *model.Product `json:"product"`
		Stock       *model.Stock   `json:"stock"`
	}

	checkStockService struct {
		productRepository repository.ProductRepository
		stockRepository   repository.StockRepository
	}

	CheckStockService interface {
		Execute(ctx context.Context, req CheckStockRequest) (*CheckStockResponse, error)
	}
)

func NewCheckStockService(
	productRepository repository.ProductRepository,
	stockRepository repository.StockRepository,
) CheckStockService {
	return &checkStockService{productRepository: productRepository, stockRepository: stockRepository}
}

func (c *checkStockService) Execute(ctx context.Context, req CheckStockRequest) (*CheckStockResponse, error) {
	output := new(CheckStockResponse)

	product, err := c.productRepository.Find(req.ProductId, ctx)
	if err != nil {
		return nil, err
	}

	stock, err := c.stockRepository.FindByProduct(product.ID, ctx)
	if err != nil {
		return nil, err
	}

	output.IsAvailable = stock.AvailableStock >= req.Quantity
	output.Product = product
	output.Stock = stock

	return output, nil
}
