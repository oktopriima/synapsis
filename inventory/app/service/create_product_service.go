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
	CreateProductRequest struct {
		Name        string  `json:"name"`
		Sku         string  `json:"sku"`
		Description string  `json:"description"`
		Price       float64 `json:"price"`
		Stock       struct {
			AvailableStock int64 `json:"available_stock"`
		} `json:"stock"`
	}

	CreateProductResponse struct {
		Product *model.Product `json:"product"`
		Stock   *model.Stock   `json:"stock"`
	}

	createProductService struct {
		db                      *gorm.DB
		productRepository       repository.ProductRepository
		stockRepository         repository.StockRepository
		stockMovementRepository repository.StockMovementRepository
	}

	CreateProductService interface {
		Execute(ctx context.Context, req CreateProductRequest) (*CreateProductResponse, error)
	}
)

func NewCreateProductService(
	instance connection.DBInstance,
	productRepository repository.ProductRepository,
	stockRepository repository.StockRepository,
	stockMovementRepository repository.StockMovementRepository) CreateProductService {
	return &createProductService{
		db:                      instance.Database(),
		productRepository:       productRepository,
		stockRepository:         stockRepository,
		stockMovementRepository: stockMovementRepository,
	}
}

func (c *createProductService) Execute(ctx context.Context, req CreateProductRequest) (*CreateProductResponse, error) {
	output := new(CreateProductResponse)
	tx := c.db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// create product
	product, err := c.productRepository.Create(&model.Product{
		Name:        req.Name,
		Sku:         req.Sku,
		Description: req.Description,
		Price:       req.Price,
	}, tx, ctx)
	if err != nil {
		return nil, err
	}

	// create stock
	stock, err := c.stockRepository.Create(&model.Stock{
		ProductId:      product.ID,
		TotalStock:     req.Stock.AvailableStock,
		AvailableStock: req.Stock.AvailableStock,
	}, tx, ctx)
	if err != nil {
		return nil, err
	}

	// create initial stock movement
	if err = c.stockMovementRepository.Create(&model.StockMovement{
		ProductId:   product.ID,
		ChangeType:  "ADD",
		Quantity:    req.Stock.AvailableStock,
		ReferenceId: product.ID,
		Note:        fmt.Sprintf("Initial stock for product %d", product.ID),
	}, tx, ctx); err != nil {
		return nil, err
	}

	tx.Commit()

	output.Product = product
	output.Stock = stock
	return output, nil
}
